package commands

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/ehazlett/circuit/config"
	"github.com/ehazlett/circuit/controller"
)

func handleHook(c controller.Controller, hook *runcHook) error {
	// check for env vars to override
	networkName := os.Getenv("NETWORK")
	if networkName == "" {
		networkName = hook.ID
	}
	subnet := os.Getenv("SUBNET")

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
		n := &config.Network{
			Name: networkName,
		}

		if subnet != "" {
			n.Subnet = subnet
		}

		if err := c.CreateNetwork(n); err != nil {
			logrus.Fatal(err)
		}

		if err := c.ConnectNetwork(n.Name, hook.Pid); err != nil {
			logrus.Fatal(err)
		}
	}

	return nil
}
