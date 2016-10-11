package commands

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	debug     bool
	statePath string
)

func init() {
	//logrus.SetFormatter(&simplelog.SimpleFormatter{})
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "Enable debug logging")
	RootCmd.PersistentFlags().StringVarP(&statePath, "state", "s", "file:///var/lib/circuit", "Circuit configuration and database path")

	RootCmd.AddCommand(networksLsCmd)
	RootCmd.AddCommand(networksCreateCmd)
	RootCmd.AddCommand(networksConnectCmd)
	RootCmd.AddCommand(networksDisconnectCmd)
	RootCmd.AddCommand(networksQosCmd)
	RootCmd.AddCommand(networksRmCmd)
}

var RootCmd = &cobra.Command{
	Use:   "circuit",
	Short: "Container Network Management",
	Long:  "Circuit manages container networking",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
			logrus.Debug("debug enabled")
		}
	},
}
