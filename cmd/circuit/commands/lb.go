package commands

import (
	"github.com/Sirupsen/logrus"
	"github.com/ehazlett/circuit/config"
	"github.com/spf13/cobra"
)

var (
	networkLBScheduler string
	networkLBProtocol  string
)

func init() {
	networksLBCreateCmd.Flags().StringVar(&networkLBProtocol, "protocol", "tcp", "Load balancer service protocol (tcp, udp)")
	networksLBCreateCmd.Flags().StringVar(&networkLBScheduler, "scheduler", "rr", "Load balancer service scheduler type (rr, wrr, lc, wlc)")

	networksLBCmd.AddCommand(networksLBCreateCmd)
	networksLBCmd.AddCommand(networksLBRemoveCmd)
	networksLBCmd.AddCommand(networksLBAddTargetsCmd)
	networksLBCmd.AddCommand(networksLBRemoveTargetsCmd)
	networksLBCmd.AddCommand(networksLBClearCmd)
}

var networksLBCmd = &cobra.Command{
	Use:   "lb",
	Short: "Manage Load Balancing",
	Long: `Manage load balancing
Details:
    circuit lb help`,
}

var networksLBCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new service",
	Long: `Create a new service
Example:
    circuit lb create <ip:port>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
			return
		}

		addr := args[0]

		if addr == "" {
			cmd.Help()
			logrus.Fatal("you must specify a service address")
		}

		c, err := getController(cmd)
		if err != nil {
			logrus.Fatal(err)
		}

		var protocol config.Protocol
		switch networkLBProtocol {
		case "tcp":
			protocol = config.ProtocolTCP
		case "udp":
			protocol = config.ProtocolUDP
		default:
			logrus.Fatalf("unknown service protocol: %s", networkLBProtocol)
		}

		var scheduler config.Scheduler
		switch networkLBScheduler {
		case "rr":
			scheduler = config.SchedulerRR
		case "wrr":
			scheduler = config.SchedulerWRR
		case "lc":
			scheduler = config.SchedulerLC
		case "wlc":
			scheduler = config.SchedulerWLC
		case "lblc":
			scheduler = config.SchedulerLBLC
		case "lblcr":
			scheduler = config.SchedulerLBLCR
		case "dh":
			scheduler = config.SchedulerDH
		case "sh":
			scheduler = config.SchedulerSH
		case "sed":
			scheduler = config.SchedulerSED
		case "nq":
			scheduler = config.SchedulerNQ
		default:
			logrus.Fatalf("unknown service scheduler: %s", networkLBScheduler)
		}

		// lblc|lblcr|dh|sh|sed|nq

		svc := &config.Service{
			Addr:      addr,
			Protocol:  protocol,
			Scheduler: scheduler,
		}

		if err := c.CreateService(svc); err != nil {
			logrus.Fatalf("error creating service: %s", err)
		}

		logrus.Infof("service %s created", addr)
	},
}

var networksLBRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a service",
	Long: `Remove a service
Example:
    circuit lb remove <ip:port>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
			return
		}

		addr := args[0]

		if addr == "" {
			cmd.Help()
			logrus.Fatal("you must specify a service address")
		}

		c, err := getController(cmd)
		if err != nil {
			logrus.Fatal(err)
		}

		var protocol config.Protocol
		switch networkLBProtocol {
		case "tcp":
			protocol = config.ProtocolTCP
		case "udp":
			protocol = config.ProtocolUDP
		default:
			logrus.Fatalf("unknown service protocol: %s", networkLBProtocol)
		}

		svc := &config.Service{
			Addr:     addr,
			Protocol: protocol,
		}
		if err := c.RemoveService(svc); err != nil {
			logrus.Fatalf("error removing service: %s", err)
		}

		logrus.Infof("service %s removed", addr)
	},
}

var networksLBAddTargetsCmd = &cobra.Command{
	Use:   "add",
	Short: "Add one or more targets to a service",
	Long: `Add targets to a service
Example:
    circuit lb add <ip:port> <target:port> [target:port]`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Help()
			return
		}

		addr := args[0]
		targets := args[1:]

		if addr == "" {
			cmd.Help()
			logrus.Fatal("you must specify a service address")
		}

		if len(targets) == 0 {
			cmd.Help()
			logrus.Fatal("you must specify at least one target")
		}

		c, err := getController(cmd)
		if err != nil {
			logrus.Fatal(err)
		}

		var protocol config.Protocol
		switch networkLBProtocol {
		case "tcp":
			protocol = config.ProtocolTCP
		case "udp":
			protocol = config.ProtocolUDP
		default:
			logrus.Fatalf("unknown service protocol: %s", networkLBProtocol)
		}

		if err := c.AddTargetsToService(addr, protocol, targets); err != nil {
			logrus.Fatalf("error adding targets to service: %s", err)
		}

		logrus.Infof("service %s updated", addr)
	},
}

var networksLBRemoveTargetsCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete one or more targets from a service",
	Long: `Delete targets from a service
Example:
    circuit lb rm <ip:port> <target:port> [target:port]`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Help()
			return
		}

		addr := args[0]
		targets := args[1:]

		if addr == "" {
			cmd.Help()
			logrus.Fatal("you must specify a service address")
		}

		if len(targets) == 0 {
			cmd.Help()
			logrus.Fatal("you must specify at least one target")
		}

		c, err := getController(cmd)
		if err != nil {
			logrus.Fatal(err)
		}

		var protocol config.Protocol
		switch networkLBProtocol {
		case "tcp":
			protocol = config.ProtocolTCP
		case "udp":
			protocol = config.ProtocolUDP
		default:
			logrus.Fatalf("unknown service protocol: %s", networkLBProtocol)
		}

		if err := c.RemoveTargetsFromService(addr, protocol, targets); err != nil {
			logrus.Fatalf("error removing targets from service: %s", err)
		}

		logrus.Infof("service %s updated", addr)
	},
}

var networksLBClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Remove all services from load balancer",
	Long: `Remove all services from load balancer
Example:
    circuit lb clear`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := getController(cmd)
		if err != nil {
			logrus.Fatal(err)
		}

		if err := c.ClearServices(); err != nil {
			logrus.Fatalf("error clearing services: %s", err)
		}

		logrus.Info("services cleared")
	},
}
