package controller

import "github.com/ehazlett/circuit/config"

type Controller interface {
	// network
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
	RemoveService(name string) error
	AddTargetsToService(serviceName string, targets []string) error
	RemoveTargetsFromService(serviceName string, targets []string) error
	ClearServices() error
	ListServices() ([]*config.Service, error)
	// util
	Restore() error
}
