package commands

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/ehazlett/circuit/config"
	"github.com/spf13/cobra"
)

var (
	debug     bool
	statePath string
)

type runcHook struct {
	ID         string `json:"id"`
	Pid        int    `json:"pid"`
	OciVersion string `json:"ociVersion"`
	Root       string `json:"root"`
	BundlePath string `json:"bundlePath"`
}

func init() {
	//logrus.SetFormatter(&simplelog.SimpleFormatter{})
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "Enable debug logging")
	RootCmd.PersistentFlags().StringVarP(&statePath, "state", "s", "file:///var/lib/circuit", "Circuit configuration and database path")

	RootCmd.AddCommand(networksLsCmd)
	RootCmd.AddCommand(networksCreateCmd)
	RootCmd.AddCommand(networksConnectCmd)
	RootCmd.AddCommand(networksDisconnectCmd)
	RootCmd.AddCommand(networksLBCmd)
	RootCmd.AddCommand(networksQosCmd)
	RootCmd.AddCommand(networksRmCmd)
}

var RootCmd = &cobra.Command{
	Use:   "circuit",
	Short: "Container Network Management",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
			logrus.Debug("debug enabled")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// if data is being piped in, use "hook" mode
			data, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				logrus.Fatal(err)
			}

			var hook *runcHook
			if err := json.Unmarshal(data, &hook); err != nil {
				logrus.Fatal(err)
			}

			c, err := getController(cmd)
			if err != nil {
				logrus.Fatal(err)
			}

			switch hook.Pid {
			case 0:
				// if hook is passed and pid == 0, container
				// is stopped.  we remove the network.
				if err := c.DeleteNetwork(hook.ID); err != nil {
					logrus.Fatal(err)
				}
			default:
				// if hook is passed and pid != 0, we do the following:
				// 1. create a network with the container name
				// 2. connect the container to the network
				n := &config.Network{
					Name: hook.ID,
				}

				if err := c.CreateNetwork(n); err != nil {
					logrus.Fatal(err)
				}

				if err := c.ConnectNetwork(hook.ID, hook.Pid); err != nil {
					logrus.Fatal(err)
				}
			}
		} else {
			cmd.Help()
		}
	},
}
