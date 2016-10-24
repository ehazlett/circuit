package local

import "github.com/ehazlett/circuit/config"

func (c *localController) CreateService(s *config.Service) error {
	return c.lb.CreateService(s)
}

func (c *localController) RemoveService(name string) error {
	return c.lb.RemoveService(name)
}

func (c *localController) AddTargetsToService(serviceName string, targets []string) error {
	return c.lb.AddTargetsToService(serviceName, targets)
}

func (c *localController) RemoveTargetsFromService(serviceName string, targets []string) error {
	return c.lb.RemoveTargetsFromService(serviceName, targets)
}

func (c *localController) ClearServices() error {
	return c.lb.ClearServices()
}

func (c *localController) ListServices() ([]*config.Service, error) {
	return c.lb.GetServices()
}
