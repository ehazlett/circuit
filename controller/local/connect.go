package local

import (
	"fmt"
	"net"

	"github.com/ehazlett/circuit/config"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

// ConnectNetwork connects a container to a network.  Note, the network
// must be setup using `CreateNetwork`.  This creates a veth pair for use
// with the host and container.
func (c *localController) ConnectNetwork(name string, containerPid int) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	logrus.Debugf("connecting %s to container %d", name, containerPid)

	bridgeName := getBridgeName(name)
	vethPair, err := createVethPair(name, bridgeName, containerPid)
	if err != nil {
		return fmt.Errorf("error configuring veth pair %d -> %s: %s", containerPid, bridgeName, err)
	}

	logrus.Debugf("veth pair: %+v", vethPair)
	if err := netlink.LinkAdd(vethPair); err != nil {
		return fmt.Errorf("error creating veth pair: %s", err)
	}

	bridgeNet, err := getInterfaceAddr(bridgeName)
	if err != nil {
		return err
	}

	originalNS, err := netns.Get()
	if err != nil {
		return err
	}
	defer originalNS.Close()

	containerPeerName := getContainerPeerName(name)
	logrus.Debugf("using container peer: %s", containerPeerName)
	iface, err := netlink.LinkByName(containerPeerName)
	if err != nil {
		return fmt.Errorf("error getting container peer link: %s", err)
	}

	// set ns to pid
	logrus.Debugf("setting namespace: peer=%v pid=%d", iface, containerPid)
	if err := netlink.LinkSetNsPid(iface, containerPid); err != nil {
		return fmt.Errorf("error setting namespace to pid %d: %s", containerPid, err)
	}

	newns, err := netns.GetFromPid(containerPid)
	if err != nil {
		return fmt.Errorf("error getting namespace for container pid %d", containerPid)
	}
	defer newns.Close()

	logrus.Debugf("switching to new ns: %v", newns)
	if err := netns.Set(newns); err != nil {
		return err
	}

	// configure container interface
	if err := c.configureContainerInterface(iface, name, bridgeNet, containerPid); err != nil {
		return err
	}

	logrus.Debug("switching back to original netns")
	if err := netns.Set(originalNS); err != nil {
		return err
	}

	// configure local peer
	if err := c.configureLocalInterface(name, bridgeNet, containerPid); err != nil {
		return err
	}

	return nil
}

func (c *localController) configureContainerInterface(iface netlink.Link, networkName string, bridgeNet *net.IPNet, containerPid int) error {
	if err := netlink.LinkSetDown(iface); err != nil {
		return fmt.Errorf("error downing interface: %s", err)
	}

	if err := netlink.LinkSetName(iface, defaultContainerInterfaceName); err != nil {
		return fmt.Errorf("error setting container peer interface name: %s", err)
	}

	cIface, err := netlink.LinkByName(defaultContainerInterfaceName)
	if err != nil {
		return fmt.Errorf("error getting container peer interface: %s", err)
	}
	// allocate IP for peer
	ip, err := c.ipam.AllocateIP(bridgeNet, networkName, containerPid, config.ContainerPeer)
	if err != nil {
		return err
	}

	logrus.Debugf("allocated ip for container: %s", ip.String())

	_, n, err := net.ParseCIDR(ip.String() + "/16")
	if err != nil {
		return fmt.Errorf("error parsing allocated ip: %s", err)
	}

	n.IP = net.ParseIP(ip.String())

	ipAddr := &netlink.Addr{IPNet: n, Label: ""}
	logrus.Debugf("assigning ip to container peer: %s", ipAddr.IPNet.IP.String())
	if err := netlink.AddrAdd(cIface, ipAddr); err != nil {
		return fmt.Errorf("error assigning ip to container peer interface: s", err)
	}

	if err := netlink.LinkSetUp(cIface); err != nil {
		return fmt.Errorf("error bringing up container peer interface: %s", err)
	}

	a := ip.To4()
	gw := net.IPv4(a[0], a[1], a[2], byte(1))
	if err := netlink.RouteAdd(&netlink.Route{
		Scope:     netlink.SCOPE_UNIVERSE,
		LinkIndex: cIface.Attrs().Index,
		Gw:        gw,
	}); err != nil {
		return fmt.Errorf("error adding route to container peer: %s", err)
	}

	return nil
}

func (c *localController) configureLocalInterface(networkName string, bridgeNet *net.IPNet, containerPid int) error {
	localIP, err := c.ipam.AllocateIP(bridgeNet, networkName, containerPid, config.HostPeer)
	if err != nil {
		return err
	}

	logrus.Debugf("allocated ip for local peer: %s", localIP.String())

	// get interface from bridge
	peerName := getLocalPeerName(networkName, containerPid)
	peer, err := netlink.LinkByName(peerName)
	if err != nil {
		return fmt.Errorf("error getting local peer interface: %s", err)
	}

	// assign ip to local peer
	if err := netlink.LinkSetDown(peer); err != nil {
		return fmt.Errorf("error downing local interface: %s", err)
	}

	_, ln, err := net.ParseCIDR(localIP.String() + "/16")
	if err != nil {
		return err
	}

	ln.IP = net.ParseIP(localIP.String())

	peerAddr := &netlink.Addr{IPNet: ln, Label: ""}
	logrus.Debugf("assigning ip to container peer: %s", peerAddr.IPNet.IP.String())
	if err := netlink.AddrAdd(peer, peerAddr); err != nil {
		return fmt.Errorf("error assigning ip to local peer: %s", err)
	}

	if err := netlink.LinkSetUp(peer); err != nil {
		return fmt.Errorf("error bringing up local peer interface: %s", err)
	}

	return nil
}
