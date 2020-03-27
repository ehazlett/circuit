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

	eventsapi "github.com/containerd/containerd/api/events"
	"github.com/containerd/containerd/events"
	"github.com/containerd/typeurl"
	"github.com/sirupsen/logrus"
)

type EventAction string

const (
	EventNone       EventAction = ""
	EventConnect    EventAction = "connect"
	EventDisconnect EventAction = "disconnect"
)

func (s *Server) eventListener(ctx context.Context, errCh chan error) {
	c, err := s.containerd()
	if err != nil {
		errCh <- err
		return
	}
	defer c.Close()

	// handle events
	ch, eventErrCh := c.Subscribe(ctx)
	for {
		select {
		case evt := <-ch:
			if err := s.eventHandler(evt); err != nil {
				logrus.Error(err)
				continue
			}
		case err := <-eventErrCh:
			errCh <- err
			return
		}
	}
}

func (s *Server) eventHandler(evt *events.Envelope) error {
	t, err := typeurl.UnmarshalAny(evt.Event)
	if err != nil {
		return err
	}
	var (
		action      EventAction
		containerID string
		pid         uint32
	)
	switch v := t.(type) {
	case *eventsapi.TaskStart:
		// connect
		logrus.Debugf("task start: container=%s pid=%d", v.ContainerID, v.Pid)
		containerID = v.ContainerID
		pid = v.Pid
		action = EventConnect
	case *eventsapi.TaskExit:
		// disconnect
		logrus.Debugf("task exit: container=%s pid=%d", v.ContainerID, v.Pid)
		containerID = v.ContainerID
		pid = v.Pid
		action = EventDisconnect
	default:
		action = EventNone
	}
	if action == EventNone {
		return nil
	}

	c, err := s.containerd()
	if err != nil {
		return err
	}
	defer c.Close()

	ctx := context.Background()

	container, err := c.LoadContainer(ctx, containerID)
	if err != nil {
		return err
	}

	labels, err := container.Labels(ctx)
	if err != nil {
		return err
	}

	network, ok := labels[s.config.NetworkLabel]
	if !ok {
		logrus.Debugf("container %s does not have network label; ignoring event", containerID)
		return nil
	}

	switch action {
	case EventConnect:
		ip, err := s.connect(ctx, containerID, network)
		if err != nil {
			return err
		}
		logrus.Infof("connected %s to %s with ip %s", containerID, network, ip.String())
	case EventDisconnect:
		if err := s.disconnect(ctx, containerID, network, pid); err != nil {
			return err
		}
	}

	return nil
}
