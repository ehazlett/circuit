package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var networksUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a network",
	Long:  "Update a container network",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: update network")
	},
}
