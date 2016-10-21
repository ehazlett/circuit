package commands

import "github.com/spf13/cobra"

var (
	networkDetails bool
)

func init() {
	networkLsCmd.Flags().BoolVar(&networkDetails, "details", false, "Show network details")

	networkCmd.AddCommand(networkCreateCmd)
	networkCmd.AddCommand(networkConnectCmd)
	networkCmd.AddCommand(networkDisconnectCmd)
	networkCmd.AddCommand(networkLsCmd)
	networkCmd.AddCommand(networkQosCmd)
	networkCmd.AddCommand(networkRmCmd)

}

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Manage networks",
	Long: `Manage networks
Details:
    circuit network help`,
}
