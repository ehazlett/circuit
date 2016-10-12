package ds

import (
	"net"

	"github.com/ehazlett/circuit/config"
)

type Backend interface {
	// networks
	GetNetwork(name string) (*config.Network, error)
	GetNetworks() ([]*config.Network, error)
	SaveNetwork(network *config.Network) error
	SaveIPAddr(ip, network string) error
	DeleteIPAddr(ip, network string) error
	GetNetworkIPs(name string) ([]net.IP, error)
	DeleteNetwork(name string) error
	// lb
	SaveService(s *config.Service) error
	DeleteService(name string) error
	GetService(name string) (*config.Service, error)
	GetServices() ([]*config.Service, error)
	GetServiceTargets(serviceName string) ([]string, error)
	AddTargetToService(serviceName, target string) error
	RemoveTargetFromService(serviceName, target string) error
}
