package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var networksLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List networks",
	Long:  "List all networks managed by Circuit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: list networks")
	},
}
