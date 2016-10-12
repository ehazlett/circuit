package local

import "github.com/ehazlett/circuit/config"

func (l *localController) CreateService(s *config.Service) error {
	return l.lb.CreateService(s)
}

func (l *localController) RemoveService(name string) error {
	return l.lb.RemoveService(name)
}

func (l *localController) AddTargetsToService(serviceName string, targets []string) error {
	return l.lb.AddTargetsToService(serviceName, targets)
}

func (l *localController) RemoveTargetsFromService(serviceName string, targets []string) error {
	return l.lb.RemoveTargetsFromService(serviceName, targets)
}

func (l *localController) ClearServices() error {
	return l.lb.ClearServices()
}

func (l *localController) ListServices() ([]*config.Service, error) {
	return l.lb.GetServices()
}
