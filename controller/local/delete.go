package local

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

func (c *localController) DeleteNetwork(name string) error {
	// stop and remove veth pair
	logrus.Debugf("removing veth pair")
	peerName := getLocalPeerName(name)
	peer, err := netlink.LinkByName(peerName)
	if err != nil {
		return fmt.Errorf("error getting local peer interface: %s", err)
	}

	if err := netlink.LinkSetDown(peer); err != nil {
		return err
	}

	if err := netlink.LinkDel(peer); err != nil {
		return err
	}

	// TODO: remove tc (tc qdisc del dev <veth> root)

	// TODO: "release" IPs back to pool

	bridgeName := getBridgeName(name)
	// TODO: remove only if there are no other networks
	//// remove nat
	//addr, err := getInterfaceAddr(bridgeName)
	//if err != nil {
	//	return err
	//}

	//logrus.Debugf("removing nat for network %s with IP %s", name, addr.String())
	//ipt, err := iptables.New()
	//if err != nil {
	//	return err
	//}
	//spec := []string{
	//	"-s",
	//	addr.String(),
	//	"-o",
	//	"eth0", // TODO: support custom nat interfaces
	//	"-j",
	//	"MASQUERADE",
	//}
	//if err := ipt.Delete("nat", "POSTROUTING", spec...); err != nil {
	//	return err
	//}

	// stop and remove bridge
	logrus.Debugf("removing bridge: %s", bridgeName)
	br, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return err
	}

	if err := netlink.LinkSetDown(br); err != nil {
		return err
	}

	if err := netlink.LinkDel(br); err != nil {
		return err
	}

	return nil
}
