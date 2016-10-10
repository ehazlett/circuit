package local

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

func getBridgeName(netName string) string {
	return bridgePrefix + "-" + netName
}

func getLocalPeerName(netName string) string {
	return fmt.Sprintf("veth-%s", netName)
}

func getContainerPeerName(netName string) string {
	return fmt.Sprintf("veth-%s-0", netName)
}

func createVethPair(netName, bridgeName string) (*netlink.Veth, error) {
	br, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return nil, err
	}

	attrs := netlink.NewLinkAttrs()
	attrs.Name = getLocalPeerName(netName)
	attrs.MasterIndex = br.Attrs().Index

	return &netlink.Veth{
		LinkAttrs: attrs,
		PeerName:  getContainerPeerName(netName),
	}, nil
}
