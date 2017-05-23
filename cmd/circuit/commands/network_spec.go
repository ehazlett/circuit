package commands

import (
	"encoding/json"
	"fmt"

	"github.com/containernetworking/cni/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var networkSpecCmd = &cobra.Command{
	Use:   "spec",
	Short: "Generate a skeleton spec for a network",
	Long: `Spec a container network
Example:
    circuit network spec`,
	Run: func(cmd *cobra.Command, args []string) {
		c := &types.NetConf{
			CNIVersion: "0.3.0",
			Name:       "test-network",
			Type:       "bridge",
			IPAM: struct {
				Type string `json:"type,omitempty"`
			}{
				Type: "host-local",
			},
		}

		data, err := json.MarshalIndent(c, "", "    ")
		if err != nil {
			logrus.Fatal(err)
		}

		fmt.Println(string(data))
	},
}
