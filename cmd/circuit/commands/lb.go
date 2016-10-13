package commands

import "github.com/spf13/cobra"

var (
	lbScheduler string
	lbProtocol  string
	lbDetails   bool
)

func init() {
	lbCreateCmd.Flags().StringVar(&lbProtocol, "protocol", "tcp", "Load balancer service protocol (tcp, udp)")
	lbCreateCmd.Flags().StringVar(&lbScheduler, "scheduler", "rr", "Load balancer service scheduler type (rr, wrr, lc, wlc)")

	lbLsCmd.Flags().BoolVar(&lbDetails, "details", false, "Show load balancer details")

	lbCmd.AddCommand(lbLsCmd)
	lbCmd.AddCommand(lbCreateCmd)
	lbCmd.AddCommand(lbDeleteCmd)
	lbCmd.AddCommand(lbAddTargetsCmd)
	lbCmd.AddCommand(lbRemoveTargetsCmd)
	lbCmd.AddCommand(lbClearCmd)
}

var lbCmd = &cobra.Command{
	Use:   "lb",
	Short: "Manage Load Balancing",
	Long: `Manage load balancing
Details:
    circuit lb help`,
}
