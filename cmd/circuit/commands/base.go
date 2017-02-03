package commands

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ehazlett/simplelog"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	debug     bool
	statePath string
	cniPath   []string
)

type runcHook struct {
	ID         string `json:"id"`
	Pid        int    `json:"pid"`
	OciVersion string `json:"ociVersion"`
	Root       string `json:"root"`
	BundlePath string `json:"bundlePath"`
}

func init() {
	logrus.SetFormatter(&simplelog.SimpleFormatter{})
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "Enable debug logging")
	RootCmd.PersistentFlags().StringVarP(&statePath, "state", "s", "file:///var/lib/circuit", "Circuit configuration and database path")
	cniPaths := strings.Split(os.Getenv("CNI_PATH"), ":")
	cniPaths = append(cniPaths, "/var/lib/circuit/cni-plugins")

	RootCmd.PersistentFlags().StringSliceVarP(&cniPath, "cni-path", "c", cniPaths, "CNI plugin path")

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
			logrus.SetFormatter(&logrus.TextFormatter{})
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
