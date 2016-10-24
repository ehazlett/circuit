package lb

import "github.com/ehazlett/circuit/config"

type LoadBalancer interface {
	CreateService(s *config.Service) error
	RemoveService(serviceName string) error
	AddTargetsToService(serviceName string, targets []string) error
	RemoveTargetsFromService(serviceName string, targets []string) error
	ClearServices() error
	GetServices() ([]*config.Service, error)
	GetService(name string) (*config.Service, error)
}
