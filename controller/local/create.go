package local

import (
	"fmt"
	"net"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/coreos/go-iptables/iptables"
	"github.com/ehazlett/circuit/config"
	"github.com/vishvananda/netlink"
)

// CreateNetwork will create and setup a network bridge
// This only creates the network.
// To connect to a container, use `ConnectNetwork`.
func (c *localController) CreateNetwork(cfg *config.Network) error {
	// create bridge
	if err := c.createBridge(cfg); err != nil {
		return err
	}

	// save
	if err := c.ds.SaveNetwork(cfg); err != nil {
		return err
	}

	return nil
}

func (c *localController) createBridge(cfg *config.Network) error {
	logrus.Debugf("creating network: name=%s subnet=%s", cfg.Name, cfg.Subnet)
	bridgeName := getBridgeName(cfg.Name)
	_, brErr := net.InterfaceByName(bridgeName)
	if brErr == nil {
		logrus.Infof("network appears to be configured")
		return nil
	}
	if !strings.Contains(brErr.Error(), "no such network interface") {
		return brErr
	}

	attrs := netlink.NewLinkAttrs()
	attrs.Name = bridgeName
	attrs.MTU = defaultMTU
	br := &netlink.Bridge{LinkAttrs: attrs}
	if err := netlink.LinkAdd(br); err != nil {
		return fmt.Errorf("error creating bridge: %s", err)
	}

	// assign ip
	ip, _, err := net.ParseCIDR(cfg.Subnet)
	if err != nil {
		return fmt.Errorf("error parsing network subnet: %s", err)
	}
	a := ip.To4()
	ipAddr := net.IPv4(a[0], a[1], a[2], byte(1))
	addr, err := netlink.ParseAddr(ipAddr.String() + "/16")
	if err != nil {
		return fmt.Errorf("error parsing address %s: %s", ipAddr, err)
	}
	if err := netlink.AddrAdd(br, addr); err != nil {
		return fmt.Errorf("error assigning ip address %v to bridge %s: %s", addr, bridgeName, err)
	}
	if _, err := getInterfaceAddr(bridgeName); err != nil {
		return fmt.Errorf("error detecting ip address for bridge %s: %s", bridgeName, err)
	}

	// add rule to masquerade
	if err := c.addNat(addr.String()); err != nil {
		return fmt.Errorf("error configuring nat: %s", err)
	}

	// bring up interface
	if err := netlink.LinkSetUp(br); err != nil {
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
