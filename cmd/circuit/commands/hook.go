package commands

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/containernetworking/cni/libcni"
	"github.com/ehazlett/circuit/controller"
	"github.com/sirupsen/logrus"
)

type Conf struct {
	CNIVersion string
	Name       string
	Type       string
	IPAMType   string
	IPAMSubnet string
	BridgeName string
	IsGateway  bool
	IPMasq     bool
}

var netTmpl = `
{
    "cniVersion": "{{.CNIVersion}}",
    "name": "{{.Name}}",
    "type": "{{.Type}}",
    "bridge": "{{.BridgeName}}",
    "ipMasq": {{.IPMasq}},
    "isGateway": {{.IsGateway}},
    "ipam": {
        "type": "{{.IPAMType}}",
        "subnet": "{{.IPAMSubnet}}",
        "routes": [
            {
                "dst": "0.0.0.0/0"
            }
        ]
    }
}
`

func handleHook(c controller.Controller, hook *runcHook) error {
	// check for env vars to override
	networkName := os.Getenv("CIRCUIT_NETWORK")
	if networkName == "" {
		networkName = hook.ID
	}

	data := fmt.Sprintf("hook: %s name: %s", hook.ID, networkName)
	ioutil.WriteFile("/tmp/circuit", []byte(data), 0644)

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

		cniVersion := os.Getenv("CNI_VERSION")
		cniType := os.Getenv("CNI_TYPE")
		ipamType := os.Getenv("CIRCUIT_IPAM_TYPE")
		ipamSubnet := os.Getenv("CIRCUIT_IPAM_SUBNET")
		bridgeName := os.Getenv("CIRCUIT_BRIDGE_NAME")
		isGateway := false
		if os.Getenv("CIRCUIT_IS_GATEWAY") != "" {
			isGateway = true
		}
		ipMasq := false
		if os.Getenv("CIRCUIT_IP_MASQ") != "" {
			ipMasq = true
		}

		conf := &Conf{
			CNIVersion: cniVersion,
			Name:       networkName,
			Type:       cniType,
			IPAMType:   ipamType,
			IPAMSubnet: ipamSubnet,
			BridgeName: bridgeName,
			IsGateway:  isGateway,
			IPMasq:     ipMasq,
		}

		t, err := template.New("conf").Parse(netTmpl)
		if err != nil {
			logrus.Fatal(err)
		}

		var buf bytes.Buffer
		if err := t.Execute(&buf, conf); err != nil {
			logrus.Fatal(err)
		}

		cfg, err := libcni.ConfFromBytes(buf.Bytes())
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
