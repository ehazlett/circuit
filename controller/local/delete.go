package local

import (
	"github.com/Sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

func (c *localController) DeleteNetwork(name string) error {
	// stop and remove veth pair
	logrus.Debugf("removing veth pair")

	// TODO: remove nat only if there are no other networks

	bridgeName := getBridgeName(name)

	// stop and remove bridge
	br, err := netlink.LinkByName(bridgeName)
	// warn only on missing bridge as it might have been removed manually
	if err != nil {
		logrus.Warn(err)
	}

	if br != nil {
		// remove existing links from bridge
		links, err := netlink.LinkList()
		if err != nil {
			return err
		}
		for _, link := range links {
			if link.Attrs().MasterIndex == br.Attrs().Index {
				logrus.Debugf("removing link from bridge: %s", link.Attrs().Name)
				if err := netlink.LinkSetDown(link); err != nil {
					return err
				}

				if err := netlink.LinkDel(link); err != nil {
					return err
				}
			}
		}

		logrus.Debugf("removing bridge: %s", bridgeName)
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
