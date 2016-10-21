package commands

import "github.com/spf13/cobra"

func init() {
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
