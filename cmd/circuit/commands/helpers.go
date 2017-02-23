package commands

import (
	"github.com/containernetworking/cni/libcni"
	"github.com/ehazlett/circuit/controller"
	"github.com/ehazlett/circuit/controller/local"
	"github.com/spf13/cobra"
)

func getControllerConfigFromCmd(c *cobra.Command) *controller.ControllerConfig {
	return &controller.ControllerConfig{
		DsURI:   statePath,
		CNIPath: cniPath,
	}
}

func getController(c *cobra.Command) (controller.Controller, error) {
	cfg := getControllerConfigFromCmd(c)
	// TODO: support different controller backends
	return local.NewLocalController(cfg)
}

func loadConfig(path string) (*libcni.NetworkConfig, error) {
	return libcni.ConfFromFile(path)
}
