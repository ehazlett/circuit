package commands

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore Circuit networks and load balancers",
	Long: `Restore networks and load balancers
Details:
    circuit restore help`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := getController(cmd)
		if err != nil {
			logrus.Fatal(err)
		}

		if err := c.Restore(); err != nil {
			logrus.Fatal(err)
		}
	},
}
