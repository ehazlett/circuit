package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var networkLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List networks",
	Long:  "List all networks managed by Circuit",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := getController(cmd)
		if err != nil {
			logrus.Fatal(err)
		}

		networks, err := c.ListNetworks()
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Debug(networks)
		w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
		fmt.Fprintf(w, "NAME \tSUBNET \n")

		for _, n := range networks {
			fmt.Fprintf(w, "%s\t%s\n", n.Name, n.Subnet)
		}

		w.Flush()
	},
}
