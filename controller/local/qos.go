package local

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/ehazlett/circuit/config"
	"github.com/vishvananda/netlink"
)

func (c *localController) SetNetworkQOS(networkName string, cfg *config.QOSConfig) error {
	logrus.Debugf("setting qos: net=%s cfg=%+v", networkName, cfg)

	var link netlink.Link
	if cfg.Interface == "" {
		bridgeName := getBridgeName(networkName)
		br, err := netlink.LinkByName(bridgeName)
		if err != nil {
			return fmt.Errorf("error getting bridge interface: %s", err)
		}
		link = br
	} else {
		iface, err := netlink.LinkByName(cfg.Interface)
		if err != nil {
			return fmt.Errorf("error getting interface: %s", err)
		}

		link = iface
	}

	attrs := netlink.QdiscAttrs{
		LinkIndex: link.Attrs().Index,
		Handle:    netlink.MakeHandle(1, 0),
		Parent:    netlink.HANDLE_ROOT,
	}
	qdisc := netlink.NewHtb(attrs)
	if err := netlink.QdiscAdd(qdisc); err != nil {
		return err
	}

	if cfg.Rate > 0 || cfg.Priority > 0 || cfg.Buffer > 0 || cfg.Cbuffer > 0 {
		logrus.Debugf("configuring htb class: %s", link.Attrs().Name)
		classattrs := netlink.ClassAttrs{
			LinkIndex: link.Attrs().Index,
			Handle:    netlink.MakeHandle(1, 1),
		}

		htbattrs := netlink.HtbClassAttrs{}
		if cfg.Rate > 0 {
			htbattrs.Rate = uint64(cfg.Rate * 10000)
		}
		if cfg.Ceiling > 0 {
			htbattrs.Ceil = uint64(cfg.Ceiling * 10000)
		}

		if cfg.Priority > 0 {
			htbattrs.Prio = uint32(cfg.Priority)
		}

		if cfg.Buffer > 0 {
			htbattrs.Buffer = uint32(cfg.Buffer)
		}
		if cfg.Cbuffer > 0 {
			htbattrs.Cbuffer = uint32(cfg.Cbuffer)
		}
		class := netlink.NewHtbClass(classattrs, htbattrs)
		_ = netlink.ClassDel(class)
		if err := netlink.ClassAdd(class); err != nil {
			return err
		}

		cclassattrs := netlink.ClassAttrs{
			LinkIndex: link.Attrs().Index,
			Handle:    netlink.MakeHandle(1, 0),
			Parent:    netlink.MakeHandle(1, 1),
		}
		cclass := netlink.NewHtbClass(cclassattrs, htbattrs)
		_ = netlink.ClassDel(cclass)
		if err := netlink.ClassAdd(cclass); err != nil {
			return err
		}
	}

	if cfg.Delay > 0 {
		latency := uint32(cfg.Delay.Nanoseconds() / 1000) // convert from duration to microseconds
		nattrs := netlink.NetemQdiscAttrs{}

		if cfg.Delay > 0 {
			nattrs.Latency = latency
		}

		netem := netlink.NewNetem(attrs, nattrs)
		_ = netlink.QdiscDel(netem)
		if err := netlink.QdiscAdd(netem); err != nil {
			return fmt.Errorf("error adding qdisc: %s", err)
		}
	}

	return nil
}

func (c *localController) ResetNetworkQOS(networkName string, iface string) error {
	var link netlink.Link
	if iface == "" {
		bridgeName := getBridgeName(networkName)
		br, err := netlink.LinkByName(bridgeName)
		if err != nil {
			return fmt.Errorf("error getting bridge interface: %s", err)
		}
		link = br
	} else {
		iface, err := netlink.LinkByName(iface)
		if err != nil {
			return fmt.Errorf("error getting interface: %s", err)
		}

		link = iface
	}

	attrs := netlink.QdiscAttrs{
		LinkIndex: link.Attrs().Index,
		Parent:    netlink.HANDLE_ROOT,
	}
	netem := netlink.NewNetem(attrs, netlink.NetemQdiscAttrs{})

	if err := netlink.QdiscDel(netem); err != nil {
		return err
	}

	return nil
}
