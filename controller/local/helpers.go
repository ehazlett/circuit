package local

import "fmt"

func getBridgeName(netName string) string {
	return bridgePrefix + "-" + netName
}

func getLocalPeerName(netName string) string {
	return fmt.Sprintf("veth-%s", netName)
}

func getContainerPeerName(netName string) string {
	return fmt.Sprintf("veth-%s-0", netName)
}
