package commands

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var networkCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a network",
	Long: `Create a container network
Example:
    circuit network create <config>`,
	ValidArgs: []string{"name", "subnet"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
			return
		}

		c, err := getController(cmd)
		if err != nil {
			logrus.Fatal(err)
		}

		conf := args[0]

		if _, err := os.Stat(conf); err != nil {
			logrus.Fatal("ERR: you must specify a valid config path")
		}

		cfg, err := loadConfig(conf)
		if err != nil {
			logrus.Fatalf("ERR: unable to parse config: %s", err)
		}

		if err := c.CreateNetwork(cfg); err != nil {
			logrus.Fatal(err)
		}

		logrus.Infof("%s created", cfg.Network.Name)
	},
}
