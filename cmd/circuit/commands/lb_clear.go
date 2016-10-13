package commands

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var lbClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Remove all services from load balancer",
	Long: `Remove all services from load balancer
Example:
    circuit lb clear`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := getController(cmd)
		if err != nil {
			logrus.Fatal(err)
		}

		if err := c.ClearServices(); err != nil {
			logrus.Fatalf("error clearing services: %s", err)
		}

		logrus.Info("services cleared")
	},
}
