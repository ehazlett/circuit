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
package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ehazlett/circuit/server"
	"github.com/pkg/errors"
	"github.com/pkg/profile"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

var serverCommand = &cli.Command{
	Name:  "server",
	Usage: "start circuit server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "containerd",
			Usage: "containerd address",
			Value: "/run/containerd/containerd.sock",
		},
		&cli.StringFlag{
			Name:  "containerd-namespace",
			Usage: "containerd namespace",
			Value: "default",
		},
		&cli.StringFlag{
			Name:  "datastore",
			Usage: "datastore uri",
			Value: "local:///var/lib/circuit",
		},
		&cli.StringFlag{
			Name:    "network-label",
			Aliases: []string{"l"},
			Usage:   "label used to autoconnect containers to networks",
			Value:   "io.circuit.network",
		},
		&cli.StringFlag{
			Name:  "cni-path",
			Usage: "cni path",
			Value: "/opt/containerd/bin",
		},
		&cli.StringFlag{
			Name:  "node-name",
			Usage: "cluster node name",
			Value: getNodeName(),
		},
		&cli.StringFlag{
			Name:  "redis-url",
			Usage: "redis url for clustering",
			Value: "",
		},
	},
	Action: serverAction,
}

func serverAction(clix *cli.Context) error {
	cfg := &server.Config{
		NodeName:              clix.String("node-name"),
		GRPCAddress:           clix.String("address"),
		DsURI:                 clix.String("datastore"),
		ContainerdAddr:        clix.String("containerd"),
		ContainerdNamespace:   clix.String("containerd-namespace"),
		NetworkLabel:          clix.String("network-label"),
		CNIPath:               clix.String("cni-path"),
		TLSServerCertificate:  clix.String("tls-server-cert"),
		TLSServerKey:          clix.String("tls-server-key"),
		TLSInsecureSkipVerify: clix.Bool("tls-skip-verify"),
		RedisURL:              clix.String("redis-url"),
	}
	if v := clix.String("profile"); v != "" {
		profileDir, err := ioutil.TempDir("", "circuit-profile-")
		if err != nil {
			return err
		}
		switch v {
		case "cpu":
			defer profile.Start(profile.Quiet, profile.ProfilePath(profileDir)).Stop()
		case "mem":
			defer profile.Start(profile.MemProfile, profile.Quiet, profile.ProfilePath(profileDir)).Stop()
		case "goroutine":
			defer profile.Start(profile.GoroutineProfile, profile.Quiet, profile.ProfilePath(profileDir)).Stop()
		case "block":
			defer profile.Start(profile.BlockProfile, profile.Quiet, profile.ProfilePath(profileDir)).Stop()
		case "mutex":
			defer profile.Start(profile.MutexProfile, profile.Quiet, profile.ProfilePath(profileDir)).Stop()
		default:
			return errors.Errorf("unknown profile %s", v)
		}
		logrus.Infof("generating profile to %s", filepath.Join(profileDir, v+".pprof"))
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		return err
	}

	return srv.Run()
}

func getNodeName() string {
	h, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return h
}
