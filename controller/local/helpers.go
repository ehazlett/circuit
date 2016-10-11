package local

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/ehazlett/circuit/controller"
	"github.com/vishvananda/netlink"
)

func getBridgeName(netName string) string {
	return controller.InterfacePrefix + "-" + netName
}

func getLocalPeerName(netName string, containerPid int) string {
	return fmt.Sprintf("veth-%d", containerPid)
}

func getContainerPeerName(netName string) string {
	return fmt.Sprintf("veth-%s-0", netName)
}

func createVethPair(netName, bridgeName string, containerPid int) (*netlink.Veth, error) {
	logrus.Debugf("creating veth pair: parent=%s pid=%d", bridgeName, containerPid)
	br, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return nil, err
	}

	linkName := getLocalPeerName(netName, containerPid)

	logrus.Debugf("creating local peer: name=%s parent=%d", linkName, br.Attrs().Index)
	attrs := netlink.NewLinkAttrs()
	attrs.Name = linkName
	attrs.MasterIndex = br.Attrs().Index

	return &netlink.Veth{
		LinkAttrs: attrs,
		PeerName:  getContainerPeerName(netName),
	}, nil
}
