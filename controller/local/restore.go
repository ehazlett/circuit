package local

// Restore restores all networks, lb, qos
func (c *localController) Restore() error {
	// networks
	networks, err := c.ListNetworks()
	if err != nil {
		return err
	}

	for _, cfg := range networks {
		// reset network ips
		if err := c.CreateNetwork(cfg); err != nil {
			return err
		}

		// TODO: restore container connectivity
		if err := c.ClearNetworkPeers(cfg.Network.Name); err != nil {
			return err
		}
	}

	// load balancers
	services, err := c.ListServices()
	if err != nil {
		return err
	}

	for _, svc := range services {
		if err := c.CreateService(svc); err != nil {
			return err
		}
	}

	return nil
}
