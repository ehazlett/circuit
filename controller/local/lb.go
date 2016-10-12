package local

import "github.com/ehazlett/circuit/config"

func (l *localController) CreateService(s *config.Service) error {
	return l.lb.CreateService(s)
}
func (l *localController) RemoveService(s *config.Service) error {
	return l.lb.RemoveService(s)
}
func (l *localController) AddTargetsToService(serviceAddr string, protocol config.Protocol, targets []string) error {
	return l.lb.AddTargetsToService(serviceAddr, protocol, targets)
}
func (l *localController) RemoveTargetsFromService(serviceAddr string, protocol config.Protocol, targets []string) error {
	return l.lb.RemoveTargetsFromService(serviceAddr, protocol, targets)
}
func (l *localController) ClearServices() error {
	return l.lb.ClearServices()
}
