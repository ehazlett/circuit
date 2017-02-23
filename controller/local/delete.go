package local

import "github.com/sirupsen/logrus"

func (c *localController) DeleteNetwork(name string) error {
	peers, err := c.ds.GetNetworkPeers(name)
	if err != nil {
		return err
	}

	// remove peers
	for _, peer := range peers {
		if err := c.DisconnectNetwork(name, peer.ContainerPid); err != nil {
			logrus.Warnf("error disconnecting container from network: %s", err)
		}
	}

	if err := c.ds.DeleteNetwork(name); err != nil {
		return err
	}

	return nil
}

func (c *localController) ClearNetworkPeers(name string) error {
	return c.ds.ClearNetworkPeers(name)
}
