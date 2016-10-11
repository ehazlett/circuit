package local

import "github.com/ehazlett/circuit/config"

// ListNetworks returns all managed networks
func (c *localController) ListNetworks() ([]*config.Network, error) {
	return c.ds.GetNetworks()
}
