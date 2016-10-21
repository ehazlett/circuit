package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Sirupsen/logrus"
	"github.com/ehazlett/circuit/config"
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

		w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
		fmt.Fprintf(w, "NAME \tSUBNET")
		if networkDetails {
			fmt.Fprintf(w, "\tCONTAINER PEERS")
		}
		fmt.Fprintf(w, "\n")

		for _, n := range networks {
			fmt.Fprintf(w, "%s\t%s", n.Name, n.Subnet)
			if networkDetails {
				fmt.Fprintf(w, "\t")
				for _, p := range n.IPs {
					if p.Type == config.ContainerPeer {
						fmt.Fprintf(w, "%s ", p.IP)
					}
				}
				fmt.Fprintf(w, "\n")
			} else {
				fmt.Fprintf(w, "\n")
			}
		}

		w.Flush()
	},
}
