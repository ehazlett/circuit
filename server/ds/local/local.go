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
	"sort"

	"github.com/containernetworking/cni/libcni"
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
	nets, err := l.GetNetworks()
	if err != nil {
		return nil, err
	}

	for _, n := range nets {
		if n.Name == name {
			return n, nil
		}
	}

	return nil, errors.Errorf("network config for %s not found", name)
}

// GetNetworks returns all CNI network configs
func (l *Local) GetNetworks() ([]*libcni.NetworkConfigList, error) {
	files, err := libcni.ConfFiles(l.statePath, []string{".json"})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)

	var networks []*libcni.NetworkConfigList
	for _, confFile := range files {
		conf, err := libcni.ConfListFromFile(confFile)
		if err != nil {
			return nil, err
		}
		networks = append(networks, conf)
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
	files, err := libcni.ConfFiles(l.statePath, []string{".json"})
	if err != nil {
		return err
	}
	for _, confFile := range files {
		conf, err := libcni.ConfListFromFile(confFile)
		if err != nil {
			return err
		}
		if conf.Name == name {
			return os.Remove(confFile)
		}
	}
	return nil
}
