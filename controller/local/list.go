package local

import (
	"github.com/containernetworking/cni/libcni"
	"github.com/ehazlett/circuit/config"
)

// ListNetworks returns all managed networks
func (c *localController) ListNetworks() ([]*libcni.NetworkConfig, error) {
	return c.ds.GetNetworks()
}

func (c *localController) ListNetworkPeers(name string) (map[string]*config.PeerInfo, error) {
	return c.ds.GetNetworkPeers(name)
}
