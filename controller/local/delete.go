package local

import (
	"github.com/Sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

func (c *localController) DeleteNetwork(name string) error {
	// stop and remove veth pair
	logrus.Debugf("removing veth pair")
	peerName := getLocalPeerName(name)
	// ignore errors when trying to find peer; peer is removed
	// upon container stop so it might not exist
	peer, _ := netlink.LinkByName(peerName)
	if peer != nil {
		if err := netlink.LinkSetDown(peer); err != nil {
			return err
		}

		if err := netlink.LinkDel(peer); err != nil {
			return err
		}
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
	// warn only on missing bridge as it might have been removed manually
	if err != nil {
		logrus.Warn(err)
	}

	if br != nil {
		if err := netlink.LinkSetDown(br); err != nil {
			return err
		}

		if err := netlink.LinkDel(br); err != nil {
			return err
		}
	}

	if err := c.ds.DeleteNetwork(name); err != nil {
		return err
	}

	return nil
}
