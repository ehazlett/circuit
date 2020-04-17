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
package ds

import (
	"github.com/containernetworking/cni/libcni"
	"github.com/pkg/errors"
)

var (
	// ErrNetworkDoesNotExist is returned when a network does not exist
	ErrNetworkDoesNotExist = errors.New("network does not exist")
)

// Datastore is the interface to implement for the circuit datastore
type Datastore interface {
	// GetNetwork returns the CNI network config for the specified network
	GetNetwork(name string) (*libcni.NetworkConfigList, error)
	// GetNetworks returns all CNI network configs
	GetNetworks() ([]*libcni.NetworkConfigList, error)
	// SaveNetwork saves a CNI network config as byte array to the datastore
	SaveNetwork(name string, data []byte) error
	// DeleteNetwork removes the network from the datastore
	DeleteNetwork(name string) error
}
