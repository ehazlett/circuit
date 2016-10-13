package commands

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var lbRemoveTargetsCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete one or more targets from a service",
	Long: `Delete targets from a service
Example:
    circuit lb rm <name> <target:port> [target:port]`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Help()
			return
		}

		name := args[0]
		targets := args[1:]

		if name == "" {
			cmd.Help()
			logrus.Fatal("you must specify a service name")
		}

		if len(targets) == 0 {
			cmd.Help()
			logrus.Fatal("you must specify at least one target")
		}

		c, err := getController(cmd)
		if err != nil {
			logrus.Fatal(err)
		}

		if err := c.RemoveTargetsFromService(name, targets); err != nil {
			logrus.Fatalf("error removing targets from service: %s", err)
		}

		logrus.Infof("service %s updated", name)
	},
}
