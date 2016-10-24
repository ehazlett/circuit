package commands

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
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

	RootCmd.AddCommand(networkCmd)
	RootCmd.AddCommand(lbCmd)
	RootCmd.AddCommand(restoreCmd)
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

			if err := handleHook(c, hook); err != nil {
				logrus.Fatal(err)
			}
		} else {
			cmd.Help()
		}
	},
}
