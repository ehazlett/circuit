package controller

import "github.com/ehazlett/circuit/config"

type Controller interface {
	CreateNetwork(c *config.Network) error
	ConnectNetwork(name string, containerPid int) error
	SetBandwidthLimit(device string, limitBytes, maxBytes int) error
	UpdateNetwork(name string, c *config.Network) error
	DeleteNetwork(name string) error
}
