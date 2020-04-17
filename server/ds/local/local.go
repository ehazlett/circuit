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
package local

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/containernetworking/cni/libcni"
	"github.com/ehazlett/circuit/server/ds"
	"github.com/pkg/errors"
)

// Local is a local datastore
type Local struct {
	statePath string
}

func New(statePath string) (*Local, error) {
	if err := os.MkdirAll(statePath, 0700); err != nil {
		return nil, err
	}
	return &Local{
		statePath: statePath,
	}, nil
}

// GetNetwork returns the CNI network config for the specified network
func (l *Local) GetNetwork(name string) (*libcni.NetworkConfigList, error) {
	configPath := filepath.Join(l.statePath, name+".json")
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			return nil, ds.ErrNetworkDoesNotExist
		} else {
			return nil, err
		}
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	networkConfig, err := libcni.ConfListFromBytes(data)
	if err != nil {
		return nil, err
	}

	return networkConfig, nil
}

// GetNetworks returns all CNI network configs
func (l *Local) GetNetworks() ([]*libcni.NetworkConfigList, error) {
	nets, err := ioutil.ReadDir(l.statePath)
	if err != nil {
		return nil, err
	}

	var networks []*libcni.NetworkConfigList
	for _, p := range nets {
		name := strings.TrimSuffix(p.Name(), filepath.Ext(p.Name()))
		n, err := l.GetNetwork(name)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to get info for network %s", name)
		}

		networks = append(networks, n)
	}

	return networks, nil
}

// SaveNetwork saves a CNI network config as byte array to the datastore
func (l *Local) SaveNetwork(name string, data []byte) error {
	configPath := filepath.Join(l.statePath, name+".json")
	f, err := ioutil.TempFile("", "circuit-ds-local-")
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	if err := os.Rename(f.Name(), configPath); err != nil {
		return err
	}

	return nil
}

// DeleteNetwork removes the network from the datastore
func (l *Local) DeleteNetwork(name string) error {
	if err := os.Remove(filepath.Join(l.statePath, name+".json")); err != nil {
		if os.IsNotExist(err) {
			return errors.Errorf("network %s does not exist", name)
		}
		return err
	}
	return nil
}
