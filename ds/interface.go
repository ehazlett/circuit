package ds

import (
	"net"

	"github.com/ehazlett/circuit/config"
)

type Backend interface {
	GetNetwork(name string) (*config.Network, error)
	SaveNetwork(network *config.Network) error
	SaveIPAddr(ip, network string) error
	GetNetworkIPs(name string) ([]net.IP, error)
}
