/*
  Copyright (c) Evan Hazlett

  Permission is hereby granted, free of charge, to any person
  obtaining a copy of this software and associated documentation
  files (the "Software"), to deal in the Software without
  restriction, including without limitation the rights to use, copy,
  modify, merge, publish, distribute, sublicense, and/or sell copies
  of the Software, and to permit persons to whom the Software is
  furnished to do so, subject to the following conditions:
  The above copyright notice and this permission notice shall be
  included in all copies or substantial portions of the Software.

  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
  EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
  OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
  IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE
  OR OTHER DEALINGS IN THE SOFTWARE.
*/
package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/typeurl"
	api "github.com/ehazlett/circuit/api/circuit/v1"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	nats "github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

const (
	clusterQueueSubject = "circuit.cluster"
)

func (s *Server) Nodes(ctx context.Context, req *api.NodesRequest) (*api.NodesResponse, error) {
	cNodes := s.clusterNodes()
	nodes := []*api.NodeInfo{}
	for _, node := range cNodes {
		nodes = append(nodes, &api.NodeInfo{
			Name: node,
		})
	}

	return &api.NodesResponse{
		Nodes: nodes,
	}, nil
}

func (s *Server) clusterListener(ctx context.Context, errCh chan error) {
	// disable if not configured
	if !s.clusterEnabled() {
		return
	}

	logrus.Debug("starting cluster")
	nc, err := nats.Connect(s.config.NATSAddr)
	if err != nil {
		errCh <- err
		return
	}

	go s.clusterHeartbeat(nc, errCh)

	recvCh := make(chan *nats.Msg)
	if _, err := nc.ChanSubscribe(clusterQueueSubject, recvCh); err != nil {
		errCh <- err
		return
	}
	for {
		msg := <-recvCh
		cx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		if err := s.handleClusterEvent(cx, msg); err != nil {
			logrus.Error(err)
		}
		cancel()
	}
}

func (s *Server) clusterHeartbeat(nc *nats.Conn, errCh chan error) {
	t := time.NewTicker(heartbeatInterval)
	for range t.C {
		data, err := marshal(&api.NodeInfo{Name: s.config.NodeName})
		if err != nil {
			errCh <- err
			return
		}
		if err := nc.Publish(clusterQueueSubject, data); err != nil {
			errCh <- err
			return
		}
	}
}

func (s *Server) handleClusterEvent(ctx context.Context, msg *nats.Msg) error {
	i, err := unmarshal(msg.Data)
	if err != nil {
		return err
	}
	switch v := i.(type) {
	case *api.NodeInfo:
		if err := s.cache.Set(v.Name, v.Name); err != nil {
			return err
		}
	case *api.ContainerIPQuery:
		nc, err := nats.Connect(s.config.NATSAddr)
		if err != nil {
			return err
		}
		defer nc.Close()

		cIPs, err := s.getLocalContainerIPs(ctx, v.Container)
		if err != nil {
			// if the container is not found we want to send
			// a nil message to the caller to inform the node
			// has responded; only return an error here if
			// it's not a "NotFound" error
			if !errdefs.IsNotFound(err) {
				return err
			}
		}

		for _, cip := range cIPs {
			data, err := marshal(cip)
			if err != nil {
				return err
			}
			if err := nc.Publish(msg.Reply, data); err != nil {
				return err
			}
		}

		// send OpComplete
		data, err := marshal(&api.OpComplete{Node: s.config.NodeName})
		if err != nil {
			return err
		}

		if err := nc.Publish(msg.Reply, data); err != nil {
			return err
		}
	default:
		logrus.Warnf("unknown cluster event type %T", v)
	}

	return nil
}

func (s *Server) getClusterContainerIPs(ctx context.Context, containerID string) ([]*api.ContainerIP, error) {
	cIPs := []*api.ContainerIP{}

	doneCh := make(chan bool, 1)
	recvCh := make(chan interface{})
	go func() {
		for {
			i, ok := <-recvCh
			if !ok {
				doneCh <- true
				return
			}
			v, ok := i.(*api.ContainerIP)
			if !ok {
				logrus.Warnf("expected api.ContainerIP; received %T", i)
				continue
			}
			cIPs = append(cIPs, v)
		}
	}()

	if err := s.receiveClusterMessages(&api.ContainerIPQuery{
		Container: containerID,
	}, recvCh); err != nil {
		return nil, err
	}

	<-doneCh

	return cIPs, nil
}

func (s *Server) clusterNodes() []string {
	kvs := s.cache.GetAll()
	nodes := []string{}
	for _, kv := range kvs {
		nodes = append(nodes, kv.Key)
	}
	return nodes
}

func (s *Server) receiveClusterMessages(request interface{}, recvCh chan interface{}) error {
	replySubject := fmt.Sprintf("circuit.%d", time.Now().UnixNano())

	nc, err := nats.Connect(s.config.NATSAddr)
	if err != nil {
		return err
	}
	defer nc.Close()

	wg := &sync.WaitGroup{}
	wg.Add(len(s.clusterNodes()))
	sub, err := nc.Subscribe(replySubject, func(m *nats.Msg) {
		i, err := unmarshal(m.Data)
		if err != nil {
			logrus.Error(err)
			return
		}
		switch i.(type) {
		case *api.OpComplete:
			wg.Done()
		default:
			recvCh <- i
		}
	})
	if err != nil {
		return err
	}

	// send request and wait for response from all nodes
	data, err := marshal(request)
	if err != nil {
		return err
	}
	if err := nc.PublishRequest(clusterQueueSubject, replySubject, data); err != nil {
		return err
	}

	wg.Wait()

	if err := sub.Unsubscribe(); err != nil {
		return err
	}

	close(recvCh)

	return nil

}

// marshal marshals to protobuf for use with the queue
func marshal(i interface{}) ([]byte, error) {
	any, err := typeurl.MarshalAny(i)
	if err != nil {
		return nil, err
	}
	data, err := proto.Marshal(any)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// unmarshal deserializes data from a cluster event
func unmarshal(data []byte) (interface{}, error) {
	var any types.Any
	if err := proto.Unmarshal(data, &any); err != nil {
		return nil, err
	}

	i, err := typeurl.UnmarshalAny(&any)
	if err != nil {
		return nil, err
	}

	return i, nil
}
