package commands

import (
	"os"

	"github.com/containernetworking/cni/libcni"
	"github.com/ehazlett/circuit/controller"
	"github.com/sirupsen/logrus"
)

func handleHook(c controller.Controller, hook *runcHook) error {
	// check for env vars to override
	networkName := os.Getenv("NETWORK")
	if networkName == "" {
		networkName = hook.ID
	}

	switch hook.Pid {
	case 0:
		// if hook is passed and pid == 0, container
		// is stopped.  we remove the network.
		if err := c.DeleteNetwork(networkName); err != nil {
			logrus.Fatal(err)
		}
	default:
		// if hook is passed and pid != 0, we do the following:
		// 1. create a network with the container name
		// 2. connect the container to the network

		// TODO: generate CNI config
		cfg := &libcni.NetworkConfig{}

		if err := c.CreateNetwork(cfg); err != nil {
			logrus.Fatal(err)
		}

		if err := c.ConnectNetwork(cfg.Network.Name, hook.Pid); err != nil {
			logrus.Fatal(err)
		}
	}

	return nil
}
