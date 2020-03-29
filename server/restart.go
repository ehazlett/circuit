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

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/errdefs"
	"github.com/ehazlett/circuit"
	"github.com/sirupsen/logrus"
)

func (s *Server) restartWatcher() {
	t := time.NewTicker(5 * time.Second)
	for range t.C {
		c, err := s.containerd()
		if err != nil {
			logrus.Error(err)
			continue
		}
		ctx := context.Background()
		containers, err := c.Containers(ctx)
		if err != nil {
			logrus.Error(err)
			continue
		}

		for _, container := range containers {
			labels, err := container.Labels(ctx)
			if err != nil {
				logrus.Error(err)
				continue
			}

			if _, ok := labels[circuit.RestartLabel]; !ok {
				continue
			}

			// check to restart
			t, err := container.Task(ctx, nil)
			if err != nil {
				if !errdefs.IsNotFound(err) {
					logrus.Error(err)
					continue
				}
			}

			// check for existing task
			if t != nil {
				// check status; start if necessary
				st, err := t.Status(ctx)
				if err != nil {
					logrus.Error(err)
					continue
				}

				switch st.Status {
				case containerd.Running:
					continue
				case containerd.Stopped:
					if _, err := t.Delete(ctx); err != nil {
						logrus.Error(err)
						continue
					}
				}
			}

			// create and start
			task, err := container.NewTask(ctx, cio.NullIO)
			if err != nil {
				logrus.Error(err)
				continue
			}
			if err := task.Start(ctx); err != nil {
				logrus.Error(err)
				continue
			}

		}
	}
}
