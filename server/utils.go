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
	"strings"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/typeurl"
	"github.com/containernetworking/cni/libcni"
	"github.com/ehazlett/circuit"
	api "github.com/ehazlett/circuit/api/circuit/v1"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

const (
	maxInterfaceCount = 10
)

var (
	// ErrNetworkConfigExtensionNotFound is returned when the network config containerd extension is not found
	ErrNetworkConfigExtensionNotFound = errors.New("network config extension not found")
)

func (s *Server) getCniConfig(networkName string, containerPid int, ifaceName string) (*libcni.CNIConfig, *libcni.NetworkConfigList, *libcni.RuntimeConf, error) {
	cfg, err := s.ds.GetNetwork(networkName)
	if err != nil {
		return nil, nil, nil, err
	}

	netConf, err := libcni.ConfListFromBytes(cfg.Bytes)
	if err != nil {
		return nil, nil, nil, err
	}

	cninet := &libcni.CNIConfig{
		Path: []string{s.config.CNIPath},
	}

	rt := &libcni.RuntimeConf{
		ContainerID: fmt.Sprintf("%d", containerPid),
		NetNS:       fmt.Sprintf("/proc/%d/ns/net", containerPid),
		IfName:      ifaceName,
	}

	return cninet, netConf, rt, nil
}

func (s *Server) generateIfaceName(containerPid int) (string, error) {
	originalNs, err := netns.Get()
	if err != nil {
		return "", err

	}
	defer originalNs.Close()

	cntNs, err := netns.GetFromPid(containerPid)
	if err != nil {
		return "", err
	}
	defer cntNs.Close()

	ifaceName := ""
	netns.Set(cntNs)
	for i := 0; i < maxInterfaceCount; i++ {
		n := fmt.Sprintf("eth%d", i)
		if _, err := netlink.LinkByName(n); err != nil {
			if !strings.Contains(err.Error(), "no such network interface") {
				ifaceName = n
				break
			}
		}
	}
	netns.Set(originalNs)

	if ifaceName == "" {
		return "", fmt.Errorf("unable to generate device name; maximum number of devices reached (%d)", maxInterfaceCount)
	}

	return ifaceName, nil
}

func (s *Server) getContainerIfaceNames(containerPid int) ([]string, error) {
	originalNs, err := netns.Get()
	if err != nil {
		return nil, err
	}
	defer originalNs.Close()

	cntNs, err := netns.GetFromPid(containerPid)
	if err != nil {
		return nil, err
	}
	defer cntNs.Close()

	ifaces := []string{}
	netns.Set(cntNs)
	for i := 0; i < maxInterfaceCount; i++ {
		n := fmt.Sprintf("eth%d", i)
		if _, err := netlink.LinkByName(n); err == nil {
			ifaces = append(ifaces, n)
		}
	}
	netns.Set(originalNs)

	return ifaces, nil
}

func (s *Server) loadNetworkConfig(ctx context.Context, c containerd.Container) (*api.NetworkConfig, error) {
	extensions, err := c.Extensions(ctx)
	if err != nil {
		return nil, err
	}

	ext, ok := extensions[circuit.NetworkConfigExtension]
	if !ok {
		return nil, ErrNetworkConfigExtensionNotFound
	}

	v, err := typeurl.UnmarshalAny(&ext)
	if err != nil {
		return nil, err
	}

	e, ok := v.(*api.NetworkConfig)
	if !ok {
		return nil, errors.Errorf("expected type 'v1.NetworkConfig'; received %T", v)
	}

	return e, nil
}

func withUpdateExtension(name string, extension interface{}) containerd.UpdateContainerOpts {
	return func(ctx context.Context, _ *containerd.Client, c *containers.Container) error {
		if name == "" {
			return errors.Wrapf(errdefs.ErrInvalidArgument, "extension name must not be zero-length")
		}
		any, err := typeurl.MarshalAny(extension)
		if err != nil {
			if errors.Cause(err) == typeurl.ErrNotFound {
				return errors.Wrapf(err, "extension %q is not registered with the typeurl package, see `typeurl.Register`", name)
			}
			return errors.Wrap(err, "error marshalling extension")
		}

		if c.Extensions == nil {
			c.Extensions = make(map[string]types.Any)
		}
		c.Extensions[name] = *any
		return nil
	}
}
