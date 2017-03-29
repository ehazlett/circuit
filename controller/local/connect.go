package local

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/ehazlett/circuit/ds"
	"github.com/sirupsen/logrus"
)

// ConnectNetwork connects a container to a network.  Note, the network
// must be setup using `CreateNetwork`.
func (c *localController) ConnectNetwork(name string, containerPid int) error {
	logrus.Debugf("connecting %s to container %d", name, containerPid)
	peer, err := c.ds.GetNetworkPeer(name, containerPid)
	if err != nil && err != ds.ErrNetworkPeerDoesNotExist {
		return err
	}
	if peer != nil {
		return fmt.Errorf("container %d is already connected to network %s", containerPid, name)
	}

	tmpConfDir, err := ioutil.TempDir("", "circuit-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpConfDir)

	ifaceName, err := c.generateIfaceName(containerPid)
	if err != nil {
		return err
	}

	cninet, nc, rt, err := c.getCniConfig(name, tmpConfDir, containerPid, ifaceName)
	if err != nil {
		return err
	}

	r, err := cninet.AddNetwork(nc, rt)
	if err != nil {
		return err
	}

	res, err := current.GetResult(r)
	if err != nil {
		return err
	}

	result, err := res.GetAsVersion("0.3.0")
	if err != nil {
		return err
	}

	cr := result.(*current.Result)
	if len(cr.IPs) == 0 {
		return fmt.Errorf("container did not receive an IP")
	}

	ip := cr.IPs[0]
	addr := ip.Address.IP.String()
	if err := c.ds.SaveNetworkPeer(name, containerPid, addr, ifaceName); err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"name":   name,
		"pid":    containerPid,
		"result": ip,
	}).Debug("container connected to network")

	return nil
}
