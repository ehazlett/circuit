package ds

import "github.com/ehazlett/circuit/config"

type Backend interface {
	// networks
	GetNetwork(name string) (*config.Network, error)
	GetNetworks() ([]*config.Network, error)
	SaveNetwork(network *config.Network) error
	SaveIPAddr(ip, network string, containerPid int, peerType config.PeerType) error
	DeleteIPAddr(ip, network string) error
	GetNetworkPeers(name string) (map[string]*config.IPPeer, error)
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
