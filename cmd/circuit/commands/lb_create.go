package commands

import (
	"github.com/sirupsen/logrus"
	"github.com/ehazlett/circuit/config"
	"github.com/spf13/cobra"
)

var lbCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new service",
	Long: `Create a new service
Example:
    circuit lb create <name> <ip:port>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			cmd.Help()
			return
		}

		name := args[0]
		addr := args[1]

		if name == "" {
			cmd.Help()
			logrus.Fatal("you must specify a service name")
		}

		if addr == "" {
			cmd.Help()
			logrus.Fatal("you must specify a service address")
		}

		c, err := getController(cmd)
		if err != nil {
			logrus.Fatal(err)
		}

		var protocol config.Protocol
		switch lbProtocol {
		case "tcp":
			protocol = config.ProtocolTCP
		case "udp":
			protocol = config.ProtocolUDP
		default:
			logrus.Fatalf("unknown service protocol: %s", lbProtocol)
		}

		var scheduler config.Scheduler
		switch lbScheduler {
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
			logrus.Fatalf("unknown service scheduler: %s", lbScheduler)
		}

		svc := &config.Service{
			Name:      name,
			Addr:      addr,
			Protocol:  protocol,
			Scheduler: scheduler,
		}

		if err := c.CreateService(svc); err != nil {
			logrus.Fatalf("error creating service: %s", err)
		}

		logrus.Infof("service %s created", name)
	},
}
