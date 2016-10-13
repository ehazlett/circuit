package commands

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var lbAddTargetsCmd = &cobra.Command{
	Use:   "add",
	Short: "Add one or more targets to a service",
	Long: `Add targets to a service
Example:
    circuit lb add <name> <target:port> [target:port]`,
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

		if err := c.AddTargetsToService(name, targets); err != nil {
			logrus.Fatalf("error adding targets to service: %s", err)
		}

		logrus.Infof("service %s updated", name)
	},
}
