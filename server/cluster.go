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
	"net/url"
	"time"

	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/typeurl"
	api "github.com/ehazlett/circuit/api/circuit/v1"
	"github.com/ehazlett/circuit/version"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	clusterHeartbeatExpire    = 10
	clusterReplyTimeout       = time.Millisecond * 5000
	clusterChannel            = "circuit.cluster"
	clusterReplyChannelPrefix = "circuit.reply"
)

var (
	clusterHeartbeatInterval = clusterHeartbeatExpire * time.Second
)

func (s *Server) Nodes(ctx context.Context, req *api.NodesRequest) (*api.NodesResponse, error) {
	nodes, err := s.clusterNodes(ctx)
	if err != nil {
		return nil, err
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
	pool, err := getPool(s.config.RedisURL)
	if err != nil {
		errCh <- err
		return
	}
	s.pool = pool

	go s.clusterHeartbeat(errCh)

	c := s.pool.Get()
	defer c.Close()

	psc := redis.PubSubConn{Conn: c}
	if err := psc.PSubscribe(clusterChannel); err != nil {
		errCh <- err
		return
	}
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			if err := s.handleClusterEvent(ctx, v); err != nil {
				errCh <- err
				return
			}
		case redis.Subscription:
		default:
			logrus.Warnf("unknown message type %T", v)
		}
	}
}

func (s *Server) clusterHeartbeat(errCh chan error) {
	t := time.NewTicker(clusterHeartbeatInterval)
	for range t.C {
		data, err := marshal(&api.NodeInfo{
			Name:    s.config.NodeName,
			Version: version.BuildVersion(),
		})
		if err != nil {
			errCh <- err
			return
		}

		key := s.clusterNodeKey(s.config.NodeName)
		ctx := context.Background()
		if _, err := s.do(ctx, "SET", key, data); err != nil {
			errCh <- err
			return
		}
		// set ttl
		if _, err := s.do(ctx, "EXPIRE", key, clusterHeartbeatExpire); err != nil {
			errCh <- err
			return
		}
	}
}

func (s *Server) handleClusterEvent(ctx context.Context, msg redis.Message) error {
	i, err := unmarshal(msg.Data)
	if err != nil {
		return err
	}
	switch v := i.(type) {
	case *api.ClusterRequest:
		if err := s.handleClusterRequest(ctx, v); err != nil {
			return err
		}
	default:
		logrus.Warnf("unknown cluster event type %T", v)
	}

	return nil
}

func (s *Server) handleClusterRequest(ctx context.Context, req *api.ClusterRequest) error {
	i, err := typeurl.UnmarshalAny(req.Request)
	if err != nil {
		return err
	}
	switch v := i.(type) {
	case *api.ContainerIPQuery:
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
			if _, err := s.do(ctx, "PUBLISH", req.Channel, data); err != nil {
				return err
			}
		}

		// send OpComplete
		data, err := marshal(&api.OpComplete{Node: s.config.NodeName})
		if err != nil {
			return err
		}
		if _, err := s.do(ctx, "PUBLISH", req.Channel, data); err != nil {
			return err
		}
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
				logrus.Warnf("expected *api.ContainerIP; received %T", i)
				continue
			}
			cIPs = append(cIPs, v)
		}
	}()

	if err := s.receiveClusterMessages(ctx, &api.ContainerIPQuery{
		Container: containerID,
	}, recvCh); err != nil {
		return nil, err
	}

	<-doneCh

	return cIPs, nil
}

func (s *Server) clusterNodes(ctx context.Context) ([]*api.NodeInfo, error) {
	nodes := []*api.NodeInfo{}
	keys, err := redis.Strings(s.do(ctx, "KEYS", s.clusterNodeKey("*")))
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		i, err := s.getData(ctx, key)
		if err != nil {
			return nil, err
		}

		v, ok := i.(*api.NodeInfo)
		if !ok {
			return nil, errors.Errorf("expected *api.NodeInfo; received %T", v)
		}
		nodes = append(nodes, v)
	}
	return nodes, nil
}

func (s *Server) receiveClusterMessages(ctx context.Context, req interface{}, recvCh chan interface{}) error {
	replyChannel := fmt.Sprintf("%s.%d", clusterReplyChannelPrefix, time.Now().UnixNano())

	doneCh := make(chan bool, 1)
	nodes, err := s.clusterNodes(ctx)
	if err != nil {
		return err
	}
	numNodes := len(nodes)

	c := s.pool.Get()
	defer c.Close()

	psc := redis.PubSubConn{Conn: c}
	if err := psc.Subscribe(replyChannel); err != nil {
		return err
	}
	go func() {
		replies := 0
		for {
			if replies >= numNodes {
				break
			}
			switch v := psc.Receive().(type) {
			case redis.Message:
				i, err := unmarshal(v.Data)
				if err != nil {
					logrus.Error(err)
					return
				}
				switch i.(type) {
				case *api.OpComplete:
					replies++
				default:
					recvCh <- i
				}
			}
		}
		doneCh <- true
	}()

	// send request and wait for response from all nodes
	any, err := typeurl.MarshalAny(req)
	if err != nil {
		return err
	}
	request := &api.ClusterRequest{
		Channel: replyChannel,
		Request: any,
	}
	data, err := marshal(request)
	if err != nil {
		return err
	}
	if _, err := s.do(ctx, "PUBLISH", clusterChannel, data); err != nil {
		return err
	}

	select {
	case <-doneCh:
	case <-time.After(clusterReplyTimeout):
		logrus.Warnf("timeout occured waiting on cluster reply (%s)", clusterReplyTimeout)
	}

	close(recvCh)

	if err := psc.Unsubscribe(replyChannel); err != nil {
		return err
	}

	return nil
}

func (s *Server) do(ctx context.Context, cmd string, args ...interface{}) (interface{}, error) {
	conn, err := s.pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	r, err := conn.Do(cmd, args...)
	return r, err
}

func (s *Server) clusterNodeKey(node string) string {
	return "nodes:" + node
}

func (s *Server) getData(ctx context.Context, key string) (interface{}, error) {
	data, err := redis.Bytes(s.do(ctx, "GET", key))
	if err != nil {
		return nil, err
	}

	return unmarshal(data)
}

func getPool(redisUrl string) (*redis.Pool, error) {
	pool := redis.NewPool(func() (redis.Conn, error) {
		conn, err := redis.DialURL(redisUrl)
		if err != nil {
			return nil, errors.Wrap(err, "unable to connect to redis")
		}

		u, err := url.Parse(redisUrl)
		if err != nil {
			return nil, err
		}

		auth, ok := u.User.Password()
		if ok {
			if _, err := conn.Do("CONFIG", "SET", "MASTERAUTH", auth); err != nil {
				return nil, errors.Wrap(err, "error authenticating to redis")
			}
		}
		return conn, nil
	}, 10)

	return pool, nil
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
