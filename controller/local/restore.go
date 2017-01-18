package local

// Restore restores all networks, lb, qos
func (c *localController) Restore() error {
	// networks
	networks, err := c.ListNetworks()
	if err != nil {
		return err
	}

	for _, network := range networks {
		// reset network ips
		network.Peers = nil
		if err := c.CreateNetwork(network); err != nil {
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
