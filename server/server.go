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
	"crypto/tls"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/containerd/containerd"
	circuitapi "github.com/ehazlett/circuit/api/circuit/v1"
	"github.com/ehazlett/circuit/server/ds"
	"github.com/ehazlett/circuit/server/ds/local"
	"github.com/ehazlett/circuit/version"
	"github.com/ehazlett/ttlcache"
	ptypes "github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	defaultRuntime = "io.containerd.runc.v2"
	dsLocal        = "local"
)

var (
	// ErrUnsupportedDatastore is returned when an unsupported datastore is specified
	ErrUnsupportedDatastore = errors.New("unsupported datastore")

	heartbeatInterval = time.Second * 5
	empty             = &ptypes.Empty{}
)

// Config is the metrics server configuration
type Config struct {
	// NodeName is the name of the node when using clustering
	NodeName string
	// ContainerdAddr is the containerd address
	ContainerdAddr string
	// ContainerdNamespace is the containerd namespace to manage
	ContainerdNamespace string
	// CNIPath is the configured CNI path
	CNIPath string
	// NetworkLabel is the label used for autoconnecting containers to networks
	NetworkLabel string
	// GRPCAddress is the address of the GRPC server
	GRPCAddress string
	// DsUri is the datastore URI
	DsURI string
	// TLSCertificate is the certificate used for grpc communication
	TLSServerCertificate string
	// TLSKey is the key used for grpc communication
	TLSServerKey string
	// TLSInsecureSkipVerify disables certificate verification
	TLSInsecureSkipVerify bool
	// NATSAddr is the NATS address for clustering
	NATSAddr string
}

// Server is the circuit server
type Server struct {
	config     *Config
	ds         ds.Datastore
	grpcServer *grpc.Server
	cache      *ttlcache.TTLCache
}

// NewServer returns a new metrics server
func NewServer(cfg *Config) (*Server, error) {
	logrus.Infof("%s %s", version.Name, version.BuildVersion())
	logrus.Debugf("using config containerd=%s namespace=%s cniPath=%s", cfg.ContainerdAddr, cfg.ContainerdNamespace, cfg.CNIPath)
	d, err := getDatastore(cfg.DsURI)
	if err != nil {
		return nil, err
	}

	grpcOpts := []grpc.ServerOption{}
	if cfg.TLSServerCertificate != "" && cfg.TLSServerKey != "" {
		logrus.WithFields(logrus.Fields{
			"cert": cfg.TLSServerCertificate,
			"key":  cfg.TLSServerKey,
		}).Debug("configuring TLS for GRPC")
		cert, err := tls.LoadX509KeyPair(cfg.TLSServerCertificate, cfg.TLSServerKey)
		if err != nil {
			return nil, err
		}
		creds := credentials.NewTLS(&tls.Config{
			Certificates:       []tls.Certificate{cert},
			ClientAuth:         tls.RequestClientCert,
			InsecureSkipVerify: cfg.TLSInsecureSkipVerify,
		})
		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(grpcOpts...)

	c, err := ttlcache.NewTTLCache(heartbeatInterval * 2)
	if err != nil {
		return nil, err
	}

	srv := &Server{
		config:     cfg,
		ds:         d,
		grpcServer: grpcServer,
		cache:      c,
	}

	// register
	circuitapi.RegisterCircuitServer(grpcServer, srv)
	// register cluster service if configured
	if cfg.NATSAddr != "" {
		logrus.Debug("enabling cluster service")
		circuitapi.RegisterClusterServer(grpcServer, srv)
	}

	return srv, nil
}

// Run starts the circuit server
func (s *Server) Run() error {
	l, err := net.Listen("tcp", s.config.GRPCAddress)
	if err != nil {
		return err
	}

	ctx := context.Background()

	errCh := make(chan error)
	logrus.Debug("starting event handler")
	go s.eventListener(ctx, errCh)
	go s.clusterListener(ctx, errCh)
	go s.restartWatcher()

	logrus.Infof("starting server on %s", s.config.GRPCAddress)
	go func() {
		if err := s.grpcServer.Serve(l); err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) containerd() (*containerd.Client, error) {
	return containerd.New(s.config.ContainerdAddr,
		containerd.WithDefaultNamespace(s.config.ContainerdNamespace),
		containerd.WithDefaultRuntime(defaultRuntime),
	)
}

func (s *Server) clusterEnabled() bool {
	return s.config.NATSAddr != ""
}

func getDatastore(uri string) (ds.Datastore, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	var d ds.Datastore

	switch u.Scheme {
	case dsLocal:
		datastore, err := local.New(u.Path)
		if err != nil {
			return nil, err
		}
		d = datastore
	default:
		return nil, errors.Wrapf(ErrUnsupportedDatastore, "available: %s", strings.Join(supportedDatastores(), ","))
	}

	logrus.Debugf("datastore: using %s", u.Scheme)
	return d, nil
}

func supportedDatastores() []string {
	return []string{
		dsLocal,
	}
}
