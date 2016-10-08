package commands

import (
	"log"

	"github.com/Sirupsen/logrus"
	"github.com/ehazlett/circuit/config"
	"github.com/spf13/cobra"
)

func init() {
	networksCreateCmd.Flags().IntVarP(&networkBandwidth, "bandwidth", "b", 0, "Network bandwidth (default: unlimited)")
}

var networksCreateCmd = &cobra.Command{
	Use:   "create <name> <subnet>",
	Short: "Create a network",
	Long: `Create a container network
Example:
    circuit create sandbox 10.254.0.0/16`,
	ValidArgs: []string{"name", "subnet"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			cmd.Usage()
			return
		}

		c, err := getController(cmd)
		if err != nil {
			logrus.Fatal(err)
		}

		networkName := args[0]
		networkSubnet := args[1]

		if networkName == "" {
			log.Fatal("ERR: you must specify a network name")
		}

		if networkSubnet == "" {
			log.Fatal("ERR: you must specify a network subnet (i.e. 10.254.0.0/16)")
		}

		logrus.Debugf("name: %s subnet: %s", networkName, networkSubnet)
		n := &config.Network{
			Name:           networkName,
			Subnet:         networkSubnet,
			BandwidthBytes: -1,
		}

		if err := c.CreateNetwork(n); err != nil {
			logrus.Fatal(err)
		}

		logrus.Infof("%s created", networkName)
	},
}
