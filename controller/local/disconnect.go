package local

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ehazlett/circuit/ds"
	"github.com/sirupsen/logrus"
)

// DisconnectNetwork disconnects a container from a network
func (c *localController) DisconnectNetwork(name string, containerPid int) error {
	logrus.Debugf("disconnecting %d from networks %s", containerPid, name)

	peer, err := c.ds.GetNetworkPeer(name, containerPid)
	if err != nil {
		if err == ds.ErrNetworkPeerDoesNotExist {
			return fmt.Errorf("container %d is not connected to network %s", containerPid, name)
		}

		return err
	}
	tmpConfDir, err := ioutil.TempDir("", "circuit-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpConfDir)

	cninet, nc, rt, err := c.getCniConfig(name, tmpConfDir, containerPid, peer.IfaceName)
	if err != nil {
		logrus.Warnf("unable to detect peer: %s", err)
	}

	if err := cninet.DelNetwork(nc, rt); err != nil {
		logrus.Warnf("unable to disconnect: %s", err)
	}

	if err := c.ds.DeleteNetworkPeer(name, containerPid); err != nil {
		logrus.Fatal(err)
	}

	return nil
}
