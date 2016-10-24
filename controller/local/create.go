package local

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/coreos/go-iptables/iptables"
	"github.com/ehazlett/circuit/config"
	"github.com/vishvananda/netlink"
)

// CreateNetwork will create and setup a network bridge
// This only creates the network.
// To connect to a container, use `ConnectNetwork`.
func (c *localController) CreateNetwork(cfg *config.Network) error {
	// get a subnet if blank
	if cfg.Subnet == "" {
		s, err := c.getSubnet()
		if err != nil {
			return err
		}
		cfg.Subnet = s
	}

	// create bridge
	if err := c.createBridge(cfg); err != nil {
		if err == ErrNetworkExists {
			logrus.Debugf("network appears to be configured")
			return nil
		}

		return err
	}

	// if a new network reset IPs if specified
	// this can happen from a restore
	cfg.IPs = nil

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
		return ErrNetworkExists
	}
	if !strings.Contains(brErr.Error(), "no such network interface") {
		return brErr
	}

	attrs := netlink.NewLinkAttrs()
	attrs.Name = bridgeName
	attrs.MTU = defaultMTU
	br := &netlink.Bridge{LinkAttrs: attrs}
	if err := netlink.LinkAdd(br); err != nil {
		return fmt.Errorf("error creating bridge: attrs=%+v err=%s", attrs, err)
	}

	// assign ip
	ip, _, err := net.ParseCIDR(cfg.Subnet)
	if err != nil {
		return fmt.Errorf("error parsing network subnet: %s", err)
	}
	// split and create the router (x.x.x.1) IP to assign to bridge
	a := ip.To4()
	ipAddr := net.IPv4(a[0], a[1], a[2], byte(1))
	// hack: manually split the subnet to join as /x
	parts := strings.Split(cfg.Subnet, "/")
	sub := parts[1]

	addr, err := netlink.ParseAddr(fmt.Sprintf("%s/%s", ipAddr.String(), sub))
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

func (c *localController) getSubnet() (string, error) {
	for {

		taken := false
		s := rand.NewSource(time.Now().UnixNano())
		r := rand.New(s)
		d := r.Intn(254)
		sub := fmt.Sprintf("10.254.%d.0/24", d)

		nets, err := c.ds.GetNetworks()
		if err != nil {
			return "", err
		}

		for _, n := range nets {
			if n.Subnet == sub {
				taken = true
				break
			}
		}

		if !taken {
			return sub, nil
		}
	}
}
