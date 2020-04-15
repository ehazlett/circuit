/*
  Copyright (c) Evan Hazlett

  Permission is hereby granted, free of charge, to any person
  obtaining a copy of this software and associated documentation
  files (the "Software"), to deal in the Software without
  restriction, including without limitation the rights to use, copy,
  modify, merge, publish, distribute, sublicense, and/or sell copies
  of the Software, and to permit persons to whom the Software is
  furnished to do so, subject to the following conditions:
  The above copyright notice and this permission notice shall be
  included in all copies or substantial portions of the Software.

  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
  EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
  OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
  IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE
  OR OTHER DEALINGS IN THE SOFTWARE.
*/
package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	api "github.com/ehazlett/circuit/api/circuit/v1"
	"github.com/ehazlett/circuit/client"
	cli "github.com/urfave/cli/v2"
)

var clusterCommand = &cli.Command{
	Name:  "cluster",
	Usage: "cluster management",
	Flags: []cli.Flag{},
	Subcommands: []*cli.Command{
		clusterNodesCommand,
	},
}

var clusterNodesCommand = &cli.Command{
	Name:  "nodes",
	Usage: "list cluster nodes",
	Action: func(clix *cli.Context) error {
		c, err := client.NewClient(clix.String("address"))
		if err != nil {
			return err
		}
		defer c.Close()

		ctx := context.Background()
		resp, err := c.Nodes(ctx, &api.NodesRequest{})
		if err != nil {
			return err
		}

		sort.Slice(resp.Nodes, func(i, j int) bool { return resp.Nodes[i].Name < resp.Nodes[j].Name })

		w := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)
		const tfmt = "%s\t%s\n"
		fmt.Fprint(w, "NAME\tVERSION\n")
		for _, node := range resp.Nodes {
			fmt.Fprintf(w, tfmt, node.Name, node.Version)
		}
		return w.Flush()

		return nil
	},
}
