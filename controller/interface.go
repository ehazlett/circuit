package controller

import "github.com/ehazlett/circuit/config"

type Controller interface {
	CreateNetwork(c *config.Network) error
	ListNetworks() ([]*config.Network, error)
	ConnectNetwork(name string, containerPid int) error
	DisconnectNetwork(name string, containerPid int) error
	DeleteNetwork(name string) error
	// qos
	SetNetworkQOS(name string, cfg *config.QOSConfig) error
	ResetNetworkQOS(name, iface string) error
	// lb
	CreateService(s *config.Service) error
	RemoveService(s *config.Service) error
	AddTargetsToService(serviceAddr string, protocol config.Protocol, targets []string) error
	RemoveTargetsFromService(serviceAddr string, protocol config.Protocol, targets []string) error
	ClearServices() error
}
