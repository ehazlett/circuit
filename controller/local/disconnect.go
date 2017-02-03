package local

import (
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// DisconnectNetwork disconnects a container from a network
func (c *localController) DisconnectNetwork(name string, containerPid int) error {
	logrus.Debugf("disconnecting %d from networks %s", containerPid, name)

	tmpConfDir, err := ioutil.TempDir("", "circuit-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpConfDir)

	// TODO: get iface for network
	peer, err := c.ds.GetNetworkPeer(name, containerPid)
	if err != nil {
		return err
	}
	cninet, nc, rt, err := c.getCniConfig(name, tmpConfDir, containerPid, peer.IfaceName)
	if err != nil {
		return err
	}

	if err := cninet.DelNetwork(nc, rt); err != nil {
		return err
	}

	if err := c.ds.DeleteNetworkPeer(name, containerPid); err != nil {
		return err
	}

	return nil
}
