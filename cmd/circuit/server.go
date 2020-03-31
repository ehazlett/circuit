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
	"github.com/ehazlett/circuit/server"
	cli "github.com/urfave/cli/v2"
)

var serverCommand = &cli.Command{
	Name:  "server",
	Usage: "start circuit server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "network-label",
			Aliases: []string{"l"},
			Usage:   "label used to autoconnect containers to networks",
			Value:   "io.circuit.network",
		},
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
			Name:  "cni-path",
			Usage: "cni path",
			Value: "/opt/containerd/bin",
		},
		&cli.StringFlag{
			Name:  "nats-addr",
			Usage: "join nats cluster",
			Value: "",
		},
	},
	Action: serverAction,
}

func serverAction(clix *cli.Context) error {
	cfg := &server.Config{
		GRPCAddress:           clix.String("address"),
		DsURI:                 clix.String("datastore"),
		ContainerdAddr:        clix.String("containerd"),
		ContainerdNamespace:   clix.String("containerd-namespace"),
		NetworkLabel:          clix.String("network-label"),
		CNIPath:               clix.String("cni-path"),
		TLSServerCertificate:  clix.String("tls-server-cert"),
		TLSServerKey:          clix.String("tls-server-key"),
		TLSInsecureSkipVerify: clix.Bool("tls-skip-verify"),
		NATSAddr:              clix.String("nats-addr"),
	}
	srv, err := server.NewServer(cfg)
	if err != nil {
		return err
	}

	return srv.Run()
}
