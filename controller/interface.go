package controller

import "github.com/ehazlett/circuit/config"

type Controller interface {
	CreateNetwork(c *config.Network) error
	ConnectNetwork(name string, containerPid int) error
	DeleteNetwork(name string) error
	SetNetworkQOS(name string, cfg *config.QOSConfig) error
	ResetNetworkQOS(name, iface string) error
}
