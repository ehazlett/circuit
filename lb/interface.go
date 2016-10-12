package lb

import "github.com/ehazlett/circuit/config"

type LoadBalancer interface {
	CreateService(s *config.Service) error
	RemoveService(s *config.Service) error
	AddTargetsToService(serviceAddr string, protocol config.Protocol, targets []string) error
	RemoveTargetsFromService(serviceAddr string, protocol config.Protocol, targets []string) error
	ClearServices() error
}
