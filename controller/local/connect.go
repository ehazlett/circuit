package local

import (
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// ConnectNetwork connects a container to a network.  Note, the network
// must be setup using `CreateNetwork`.  This creates a veth pair for use
// with the host and container.
func (c *localController) ConnectNetwork(name string, containerPid int) error {
	logrus.Debugf("connecting %s to container %d", name, containerPid)

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

	res, err := cninet.AddNetwork(nc, rt)
	if err != nil {
		return err
	}

	ip := res.IP4.IP.IP.String()
	if err := c.ds.SaveNetworkPeer(name, containerPid, ip, ifaceName); err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"name":   name,
		"pid":    containerPid,
		"result": ip,
	}).Debug("container connected to network")

	return nil
}
