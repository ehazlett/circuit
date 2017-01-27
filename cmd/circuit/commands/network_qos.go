package commands

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ehazlett/circuit/config"
	"github.com/spf13/cobra"
)

var (
	networkQosDelay     time.Duration
	networkQosRate      int
	networkQosCeiling   int
	networkQosBuffer    int
	networkQosCbuffer   int
	networkQosPriority  int
	networkQosInterface string
)

func init() {
	networkQosSetCmd.Flags().DurationVarP(&networkQosDelay, "delay", "d", time.Second*0, "Network delay (default: 0ms)")
	networkQosSetCmd.Flags().IntVar(&networkQosRate, "rate", 0, "Network class rate in kbit (default: unlimited)")
	networkQosSetCmd.Flags().IntVar(&networkQosCeiling, "ceiling", 0, "Network class ceiling in kbit (default: unlimited)")
	networkQosSetCmd.Flags().IntVar(&networkQosBuffer, "buffer", 0, "Network class buffer")
	networkQosSetCmd.Flags().IntVar(&networkQosCbuffer, "cbuffer", 0, "Network class cbuffer")
	networkQosSetCmd.Flags().IntVar(&networkQosPriority, "priority", 0, "Network class priority (default: 0)")
	networkQosSetCmd.Flags().StringVarP(&networkQosInterface, "interface", "i", "", "Specify network interface to use instead of entire bridge")

	networkQosResetCmd.Flags().StringVarP(&networkQosInterface, "interface", "i", "", "Specify network interface to use instead of entire bridge")

	networkQosCmd.AddCommand(networkQosSetCmd)
	networkQosCmd.AddCommand(networkQosResetCmd)
}

var networkQosCmd = &cobra.Command{
	Use:   "qos",
	Short: "Manage QOS for a network",
	Long: `Manage quality of service for networks
Details:
    circuit qos help`,
}

var networkQosSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set QOS for a network",
	Long: `Setup quality of service for a network
Example:
    circuit qos set --delay 100ms foo
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
			return
		}

		c, err := getController(cmd)
		if err != nil {
			logrus.Fatal(err)
		}

		networkName := args[0]

		if networkName == "" {
			logrus.Fatal("ERR: you must specify a network name")
		}

		cfg := &config.QOSConfig{
			Delay:     networkQosDelay,
			Rate:      networkQosRate,
			Ceiling:   networkQosCeiling,
			Buffer:    networkQosBuffer,
			Cbuffer:   networkQosCbuffer,
			Priority:  networkQosPriority,
			Interface: networkQosInterface,
		}

		if err := c.SetNetworkQOS(networkName, cfg); err != nil {
			logrus.Fatal(err)
		}

		iface := networkName
		if cfg.Interface != "" {
			iface = cfg.Interface
		}

		logrus.Infof("qos configured for %s", iface)
	},
}

var networkQosResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset QOS for a network",
	Long: `Reset quality of service for a network
Example:
    circuit qos reset foo
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
			return
		}

		c, err := getController(cmd)
		if err != nil {
			logrus.Fatal(err)
		}

		networkName := args[0]

		if networkName == "" {
			logrus.Fatal("ERR: you must specify a network name")
		}

		if err := c.ResetNetworkQOS(networkName, networkQosInterface); err != nil {
			logrus.Fatal(err)
		}

		iface := networkName
		if networkQosInterface != "" {
			iface = networkQosInterface
		}

		logrus.Infof("qos reset for %s", iface)
	},
}
