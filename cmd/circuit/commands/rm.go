package commands

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var networksRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a network",
	Long:  "Remove a network managed by Circuit",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Usage()
			return
		}

		c, err := getController(cmd)
		if err != nil {
			logrus.Fatal(err)
		}

		networkName := args[0]

		if networkName == "" {
			logrus.Fatalf("a network name must be specified")
		}

		if err := c.DeleteNetwork(networkName); err != nil {
			logrus.Fatal(err)
		}

		logrus.Infof("%s removed", networkName)
	},
}
