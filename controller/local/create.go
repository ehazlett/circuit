package local

import (
	"fmt"
	"net"

	"github.com/containernetworking/cni/libcni"
	"github.com/coreos/go-iptables/iptables"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

// CreateNetwork will create a network
// This only creates the network.
// To connect to a container, use `ConnectNetwork`.
func (c *localController) CreateNetwork(cfg *libcni.NetworkConfig) error {
	if err := c.ds.SaveNetwork(cfg); err != nil {
		return err
	}

	return nil
}

func getInterfaceAddr(name string) (*net.IPNet, error) {
	iface, err := netlink.LinkByName(name)
	if err != nil {
		return nil, err

	}
	addrs, err := netlink.AddrList(iface, netlink.FAMILY_V4)
	if err != nil {
		return nil, err

	}
	if len(addrs) == 0 {
		return nil, fmt.Errorf("unable to detect IP addresses for interface: %s", name)

	}
	if len(addrs) > 1 {
		logrus.Warnf("interface %s has more than 1 IP address; using %v", name, addrs[0])

	}

	return addrs[0].IPNet, nil

}

func (c *localController) addNat(ip string) error {
	logrus.Debugf("setting up nat for bridge: %s", ip)
	ipt, err := iptables.New()
	if err != nil {
		return err

	}
	spec := []string{
		"-s",
		ip,
		"-o",
		"eth0", // TODO: support custom nat interfaces
		"-j",
		"MASQUERADE",
	}
	if err := ipt.AppendUnique("nat", "POSTROUTING", spec...); err != nil {
		return err

	}

	return nil

}
