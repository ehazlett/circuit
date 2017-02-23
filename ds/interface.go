package ds

import (
	"github.com/containernetworking/cni/libcni"
	"github.com/ehazlett/circuit/config"
)

type Backend interface {
	// networks
	GetNetwork(name string) (*libcni.NetworkConfig, error)
	GetNetworks() ([]*libcni.NetworkConfig, error)
	SaveNetwork(network *libcni.NetworkConfig) error
	DeleteNetwork(name string) error
	SaveNetworkPeer(name string, containerPid int, ip string, ifaceName string) error
	DeleteNetworkPeer(name string, containerPid int) error
	ClearNetworkPeers(name string) error
	GetNetworkPeer(name string, containerPid int) (*config.PeerInfo, error)
	// GetNetworkPeers returns a IP to Pid map for the network peers
	GetNetworkPeers(name string) (map[string]*config.PeerInfo, error)
	// lb
	SaveService(s *config.Service) error
	DeleteService(name string) error
	GetService(name string) (*config.Service, error)
	GetServices() ([]*config.Service, error)
	GetServiceTargets(serviceName string) ([]string, error)
	AddTargetToService(serviceName, target string) error
	RemoveTargetFromService(serviceName, target string) error
}
