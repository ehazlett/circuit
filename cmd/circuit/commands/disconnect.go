package commands

import (
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var networksDisconnectCmd = &cobra.Command{
	Use:   "disconnect <pid> <network>",
	Short: "Disconnect a container from a network",
	Long: `Disconnect a container from a network
Example:
    circuit disconnect 12345 foo`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			cmd.Help()
			return
		}

		c, err := getController(cmd)
		if err != nil {
			logrus.Fatal(err)
		}

		pidNum := args[0]
		networkName := args[1]

		if pidNum == "" {
			logrus.Fatal("ERR: you must specify a container pid")
		}

		if networkName == "" {
			logrus.Fatal("ERR: you must specify a network name")
		}

		pid, err := strconv.Atoi(pidNum)
		if err != nil {
			logrus.Fatalf("ERR: unable to detect pid: %s", err)
		}

		if err := c.DisconnectNetwork(networkName, pid); err != nil {
			logrus.Fatal(err)
		}

		logrus.Infof("disconnected container %d from network %s", pid, networkName)
	},
}
