package controller

import (
	"github.com/containernetworking/cni/libcni"
	"github.com/ehazlett/circuit/config"
)

type Controller interface {
	// network
	CreateNetwork(c *libcni.NetworkConfig) error
	ListNetworks() ([]*libcni.NetworkConfig, error)
	ConnectNetwork(name string, containerPid int) error
	DisconnectNetwork(name string, containerPid int) error
	DeleteNetwork(name string) error
	ListNetworkPeers(name string) (map[string]*config.PeerInfo, error)
	GetNetworkPeer(name string, pid int) (*config.PeerInfo, error)
	ClearNetworkPeers(name string) error
	// qos
	SetNetworkQOS(name string, cfg *config.QOSConfig) error
	ResetNetworkQOS(name, iface string) error
	// lb
	CreateService(s *config.Service) error
	RemoveService(name string) error
	AddTargetsToService(serviceName string, targets []string) error
	RemoveTargetsFromService(serviceName string, targets []string) error
	ClearServices() error
	GetService(serviceName string) (*config.Service, error)
	ListServices() ([]*config.Service, error)
	// util
	Restore() error
}
