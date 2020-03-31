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
	"time"

	api "github.com/ehazlett/circuit/api/circuit/v1"
	nats "github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

const (
	clusterQueueSubject = "circuit.cluster"
)

func (s *Server) clusterListener(ctx context.Context, errCh chan error) {
	if s.config.NATSAddr == "" {
		return
	}
	nc, err := nats.Connect(s.config.NATSAddr)
	if err != nil {
		errCh <- err
		return
	}

	recvCh := make(chan *nats.Msg)
	if _, err := nc.ChanSubscribe(clusterQueueSubject, recvCh); err != nil {
		errCh <- err
		return
	}

	if err := nc.Publish(clusterQueueSubject, []byte("heartbeat")); err != nil {
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

func (s *Server) handleClusterEvent(ctx context.Context, msg *nats.Msg) error {
	logrus.Debugf("event: %+v", msg)
	logrus.Debugf("data: %s", string(msg.Data))
	return nil
}

func (s *Server) getContainerIPFromCluster(ctx context.Context, containerID string) (*api.ContainerIP, error) {
	return nil, nil
}
