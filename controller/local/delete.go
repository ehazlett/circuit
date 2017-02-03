package local

import "fmt"

func (c *localController) DeleteNetwork(name string) error {
	peers, err := c.ds.GetNetworkPeers(name)
	if err != nil {
		return err
	}

	if len(peers) > 0 {
		return fmt.Errorf("cannot delete network; network has peers")
	}

	if err := c.ds.DeleteNetwork(name); err != nil {
		return err
	}

	return nil
}

func (c *localController) ClearNetworkPeers(name string) error {
	return c.ds.ClearNetworkPeers(name)
}
