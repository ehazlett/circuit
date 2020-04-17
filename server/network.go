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
	"net"
	"strings"

	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/ehazlett/circuit"
	api "github.com/ehazlett/circuit/api/circuit/v1"
	ptypes "github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (s *Server) CreateNetwork(ctx context.Context, req *api.CreateNetworkRequest) (*ptypes.Empty, error) {
	if err := s.ds.SaveNetwork(req.Name, req.Data); err != nil {
		return nil, err
	}
	logrus.WithFields(logrus.Fields{"name": req.Name}).Info("created network")
	return empty, nil
}

func (s *Server) DeleteNetwork(ctx context.Context, req *api.DeleteNetworkRequest) (*ptypes.Empty, error) {
	if err := s.ds.DeleteNetwork(req.Name); err != nil {
		return nil, err
	}
	logrus.WithFields(logrus.Fields{"name": req.Name}).Info("deleted network")
	return empty, nil
}

func (s *Server) Connect(ctx context.Context, req *api.ConnectRequest) (*api.ConnectResponse, error) {
	ip, err := s.connect(ctx, req.Container, req.Network)
	if err != nil {
		return nil, err
	}
	return &api.ConnectResponse{
		IP: ip.String(),
	}, nil
}

func (s *Server) Disconnect(ctx context.Context, req *api.DisconnectRequest) (*ptypes.Empty, error) {
	c, err := s.containerd()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	container, err := c.LoadContainer(ctx, req.Container)
	if err != nil {
		return nil, errors.Wrapf(err, "error loading container %s", req.Container)
	}

	task, err := container.Task(ctx, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "no task found for container %s", req.Container)
	}

	pids, err := task.Pids(ctx)
	if err != nil {
		return nil, err
	}

	if len(pids) == 0 {
		return nil, errors.Errorf("no pids found for task in container %s", req.Container)
	}

	containerPid := pids[0].Pid

	if err := s.disconnect(ctx, req.Container, req.Network, containerPid); err != nil {
		return nil, err
	}

	return empty, nil
}

func (s *Server) ListNetworks(ct context.Context, req *api.ListNetworksRequest) (*api.ListNetworksResponse, error) {
	nets, err := s.ds.GetNetworks()
	if err != nil {
		return nil, err
	}

	networks := []*api.Network{}
	for _, n := range nets {
		networks = append(networks, &api.Network{
			Name: n.Name,
			Data: n.Bytes,
		})
	}
	return &api.ListNetworksResponse{
		Networks: networks,
	}, nil
}

func (s *Server) GetNetwork(ctx context.Context, req *api.GetNetworkRequest) (*api.GetNetworkResponse, error) {
	network, err := s.ds.GetNetwork(req.Name)
	if err != nil {
		return nil, err
	}

	return &api.GetNetworkResponse{
		Network: &api.Network{
			Name: req.Name,
			Data: network.Bytes,
		},
	}, nil
}

func (s *Server) GetContainerIPs(ctx context.Context, req *api.GetContainerIPsRequest) (*api.GetContainerIPsResponse, error) {
	cIPs, err := s.getContainerIPs(ctx, req.Container)
	if err != nil {
		return nil, err
	}

	return &api.GetContainerIPsResponse{
		IPs: cIPs,
	}, nil
}

func (s *Server) getContainerIPs(ctx context.Context, containerID string) ([]*api.ContainerIP, error) {
	// resolve via cluster if enabled; otherwise lookup locally
	if !s.clusterEnabled() {
		return s.getLocalContainerIPs(ctx, containerID)
	}
	// cluster enabled
	cIPs, err := s.getClusterContainerIPs(ctx, containerID)
	if err != nil {
		return nil, err
	}

	return cIPs, nil
}

func (s *Server) getLocalContainerIPs(ctx context.Context, containerID string) ([]*api.ContainerIP, error) {
	c, err := s.containerd()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	container, err := c.LoadContainer(ctx, containerID)
	if err != nil {
		return nil, err
	}

	networkConfig, err := s.loadNetworkConfig(ctx, container)
	if err != nil {
		if err == ErrNetworkConfigExtensionNotFound {
			return nil, nil
		}
		return nil, err
	}

	cIPs := []*api.ContainerIP{}
	for network, cfg := range networkConfig.Networks {
		cIPs = append(cIPs, &api.ContainerIP{
			Network:   network,
			IP:        cfg.IP,
			Interface: cfg.Interface,
		})
	}

	return cIPs, nil
}

func (s *Server) connect(ctx context.Context, containerID, networkName string) (net.IP, error) {
	c, err := s.containerd()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	container, err := c.LoadContainer(ctx, containerID)
	if err != nil {
		return nil, errors.Wrapf(err, "error loading container %s", containerID)
	}

	task, err := container.Task(ctx, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "no task found for container %s", container.ID())
	}

	pids, err := task.Pids(ctx)
	if err != nil {
		return nil, err
	}

	if len(pids) == 0 {
		return nil, errors.Errorf("no pids found for task in container %s", container.ID())
	}

	containerPid := int(pids[0].Pid)
	ifaceName, err := s.generateIfaceName(containerPid)
	if err != nil {
		return nil, err
	}

	cninet, nc, rt, err := s.getCniConfig(networkName, containerPid, ifaceName)
	if err != nil {
		return nil, errors.Wrap(err, "error getting cni config")
	}

	r, err := cninet.AddNetworkList(ctx, nc, rt)
	if err != nil {
		return nil, errors.Wrap(err, "error adding cni network")
	}

	res, err := current.GetResult(r)
	if err != nil {
		return nil, errors.Wrap(err, "error getting result from cninet")
	}

	result, err := res.GetAsVersion("0.3.0")
	if err != nil {
		return nil, errors.Wrap(err, "error getting result as version")
	}

	cr := result.(*current.Result)
	if len(cr.IPs) == 0 {
		return nil, fmt.Errorf("container did not receive an IP")
	}

	ipConfig := cr.IPs[0]
	ip := ipConfig.Address.IP

	networkConfig, err := s.loadNetworkConfig(ctx, container)
	if err != nil {
		if err != ErrNetworkConfigExtensionNotFound {
			return nil, errors.Wrap(err, "error loading network config")
		}
		networkConfig = &api.NetworkConfig{
			Networks: map[string]*api.ContainerNetworkConfig{},
		}
	}
	if networkConfig.Networks == nil {
		networkConfig.Networks = map[string]*api.ContainerNetworkConfig{}
	}

	networkConfig.Networks[networkName] = &api.ContainerNetworkConfig{
		Interface: ifaceName,
		IP:        ip.String(),
	}

	if err := container.Update(ctx, withUpdateExtension(circuit.NetworkConfigExtension, networkConfig)); err != nil {
		return nil, errors.Wrap(err, "error updating container extension")
	}

	return ip, nil
}

func (s *Server) disconnect(ctx context.Context, containerID, networkName string, pid uint32) error {
	c, err := s.containerd()
	if err != nil {
		return err
	}
	defer c.Close()

	container, err := c.LoadContainer(ctx, containerID)
	if err != nil {
		return errors.Wrap(err, "error loading container")
	}

	networkConfig, err := s.loadNetworkConfig(ctx, container)
	if err != nil {
		return errors.Wrap(err, "error loading network config")
	}

	network, ok := networkConfig.Networks[networkName]
	if !ok {
		return errors.Errorf("%s is not connected to %s", containerID, networkName)
	}

	ifaceName := network.Interface

	cninet, nc, rt, err := s.getCniConfig(networkName, int(pid), ifaceName)
	if err != nil {
		return errors.Wrapf(err, "error getting cni config for container %s", containerID)
	}

	if err := cninet.DelNetworkList(ctx, nc, rt); err != nil {
		// check for "no such file" to see if the netns path exists.  if not the container is removed
		// cni does not have a known code for not exists
		// https://github.com/containernetworking/cni/blob/master/SPEC.md#well-known-error-codes
		if strings.Contains(err.Error(), "no such file") {
			return nil
		}
		return errors.Wrapf(err, "error disconnecting %s from %s", containerID, network)
	}

	// remove network from network config
	delete(networkConfig.Networks, networkName)
	if err := container.Update(ctx, withUpdateExtension(circuit.NetworkConfigExtension, networkConfig)); err != nil {
		return errors.Wrap(err, "error updating container extension")
	}

	return nil
}
