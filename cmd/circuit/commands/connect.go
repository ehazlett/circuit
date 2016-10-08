package commands

import (
	"log"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var networksConnectCmd = &cobra.Command{
	Use:   "connect <pid> <network>",
	Short: "Connect a container to a network",
	Long: `Connect a container to a network
Example:
    circuit connect 12345 foo`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			cmd.Usage()
			return
		}

		c, err := getController(cmd)
		if err != nil {
			logrus.Fatal(err)
		}

		pidNum := args[0]
		networkName := args[1]

		if pidNum == "" {
			log.Fatal("ERR: you must specify a container pid")
		}

		if networkName == "" {
			log.Fatal("ERR: you must specify a network name")
		}

		pid, err := strconv.Atoi(pidNum)
		if err != nil {
			log.Fatalf("ERR: unable to detect pid: %s", err)
		}

		if err := c.ConnectNetwork(networkName, pid); err != nil {
			logrus.Fatal(err)
		}

		logrus.Infof("connected container %d to network %s", pid, networkName)
	},
}
