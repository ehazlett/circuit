package local

func (c *localController) DeleteNetwork(name string) error {
	// stop and remove bridge
	//cfg, err := c.ds.GetNetwork(name)
	//if err != nil {
	//	return err
	//}

	// TODO: get interface and check for peers; if none remove
	if err := c.ds.DeleteNetwork(name); err != nil {
		return err
	}

	return nil
}

func (c *localController) ClearNetworkPeers(name string) error {
	return c.ds.ClearNetworkPeers(name)
}
