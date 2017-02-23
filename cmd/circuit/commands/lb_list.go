package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var lbListCmd = &cobra.Command{
	Use:   "list",
	Short: "List services",
	Long:  "List all services managed by Circuit",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := getController(cmd)
		if err != nil {
			logrus.Fatal(err)
		}

		services, err := c.ListServices()
		if err != nil {
			logrus.Fatal(err)
		}

		w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
		fmt.Fprintf(w, "NAME \tADDR \tPROTOCOL \tSCHEDULER \n")

		for _, s := range services {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s", s.Name, s.Addr, s.Protocol, s.Scheduler)
			if lbDetails && len(s.Targets) > 0 {
				fmt.Fprintf(w, "\n")
				for _, t := range s.Targets {
					fmt.Fprintf(w, "  -> %s\n", t)
				}
			} else {
				fmt.Fprintf(w, "\n")
			}
		}

		w.Flush()
	},
}
