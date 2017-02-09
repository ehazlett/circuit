package commands

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/containernetworking/cni/libcni"
	"github.com/ehazlett/circuit/controller"
	"github.com/sirupsen/logrus"
)

func handleHook(c controller.Controller, hook *runcHook) error {
	// check for env vars to override
	networkName := os.Getenv("CIRCUIT_NETWORK")
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
		cniConf := os.Getenv("CNI_CONF")

		// deserialize to override values
		f, err := ioutil.ReadFile(cniConf)
		if err != nil {
			logrus.Fatal(err)
		}

		var conf map[string]interface{}
		if err := json.Unmarshal(f, &conf); err != nil {
			logrus.Fatal(err)
		}

		// override the network name
		conf["name"] = networkName

		data, err := json.Marshal(conf)
		if err != nil {
			return err
		}

		cfg, err := libcni.ConfFromBytes(data)
		if err != nil {
			logrus.Fatal(err)
		}

		if err := c.CreateNetwork(cfg); err != nil {
			logrus.Fatal(err)
		}

		if err := c.ConnectNetwork(cfg.Network.Name, hook.Pid); err != nil {
			logrus.Fatal(err)
		}
	}

	return nil
}
