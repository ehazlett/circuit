package commands

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var networkListCmd = &cobra.Command{
	Use:   "list",
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
		fmt.Fprintf(w, "NAME\tTYPE\tVERSION\tPEERS\n")

		for _, n := range networks {
			fmt.Fprintf(w, "%s\t%s\t%s\t", n.Network.Name, n.Network.Type, n.Network.CNIVersion)
			peers, err := c.ListNetworkPeers(n.Network.Name)
			if err != nil {
				logrus.Fatal(err)
			}

			netPeers := []string{}
			for ip, info := range peers {
				netPeers = append(netPeers, fmt.Sprintf("%s (%d)", ip, info.ContainerPid))
			}

			fmt.Fprintf(w, strings.Join(netPeers, ", "))
			fmt.Fprintf(w, "\n")
		}

		w.Flush()
	},
}
