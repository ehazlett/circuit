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
	"io/ioutil"
	"os"
	"text/tabwriter"

	api "github.com/ehazlett/circuit/api/circuit/v1"
	"github.com/ehazlett/circuit/client"
	"github.com/pkg/errors"
	cli "github.com/urfave/cli/v2"
)

var networkCommand = &cli.Command{
	Name:  "network",
	Usage: "network management",
	Flags: []cli.Flag{},
	Subcommands: []*cli.Command{
		networkCreateCommand,
		networkListCommand,
		networkInfoCommand,
		networkDeleteCommand,
		networkConnectCommand,
		networkDisconnectCommand,
	},
}

var networkCreateCommand = &cli.Command{
	Name:      "create",
	Usage:     "create network",
	ArgsUsage: "<name> <config-path>",
	Action: func(clix *cli.Context) error {
		c, err := client.NewClient(clix.String("address"))
		if err != nil {
			return err
		}
		defer c.Close()

		name := clix.Args().First()
		configPath := clix.Args().Get(1)
		if name == "" || configPath == "" {
			cli.ShowSubcommandHelp(clix)
			return fmt.Errorf("name and config must be specified")
		}

		data, err := ioutil.ReadFile(configPath)
		if err != nil {
			return errors.Wrap(err, "unable to read config")
		}

		ctx := context.Background()
		if _, err := c.CreateNetwork(ctx, &api.CreateNetworkRequest{
			Name: name,
			Data: data,
		}); err != nil {
			return err
		}

		return nil
	},
}

var networkListCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "list available networks",
	Action: func(clix *cli.Context) error {
		c, err := client.NewClient(clix.String("address"))
		if err != nil {
			return err
		}
		defer c.Close()

		ctx := context.Background()
		resp, err := c.ListNetworks(ctx, &api.ListNetworksRequest{})
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)
		const tfmt = "%s\t%s\n"
		fmt.Fprint(w, "NAME\tTYPE\n")
		for _, network := range resp.Networks {
			fmt.Fprintf(w, tfmt, network.Name, network.Type)
		}
		return w.Flush()
	},
}

var networkInfoCommand = &cli.Command{
	Name:  "info",
	Usage: "show network config",
	Action: func(clix *cli.Context) error {
		c, err := client.NewClient(clix.String("address"))
		if err != nil {
			return err
		}
		defer c.Close()

		name := clix.Args().First()
		if name == "" {
			cli.ShowSubcommandHelp(clix)
			return fmt.Errorf("name must be specified")
		}

		ctx := context.Background()
		resp, err := c.GetNetwork(ctx, &api.GetNetworkRequest{
			Name: name,
		})
		if err != nil {
			return err
		}

		fmt.Fprint(os.Stdout, string(resp.Network.Data))
		return nil
	},
}

var networkDeleteCommand = &cli.Command{
	Name:    "delete",
	Aliases: []string{"rm"},
	Usage:   "delete network",
	Action: func(clix *cli.Context) error {
		c, err := client.NewClient(clix.String("address"))
		if err != nil {
			return err
		}
		defer c.Close()

		name := clix.Args().First()
		if name == "" {
			cli.ShowSubcommandHelp(clix)
			return fmt.Errorf("name must be specified")
		}

		ctx := context.Background()
		if _, err := c.DeleteNetwork(ctx, &api.DeleteNetworkRequest{
			Name: name,
		}); err != nil {
			return err
		}
		fmt.Println(name)
		return nil
	},
}

var networkConnectCommand = &cli.Command{
	Name:      "connect",
	Usage:     "connect container to network",
	ArgsUsage: "<container> <network>",
	Action: func(clix *cli.Context) error {
		c, err := client.NewClient(clix.String("address"))
		if err != nil {
			return err
		}
		defer c.Close()

		container := clix.Args().First()
		network := clix.Args().Get(1)
		if container == "" || network == "" {
			cli.ShowSubcommandHelp(clix)
			return fmt.Errorf("container and network must be specified")
		}

		ctx := context.Background()
		resp, err := c.Connect(ctx, &api.ConnectRequest{
			Container: container,
			Network:   network,
		})
		if err != nil {
			return err
		}

		fmt.Printf("connected %s to %s with ip=%s\n", container, network, resp.IP)

		return nil
	},
}

var networkDisconnectCommand = &cli.Command{
	Name:      "disconnect",
	Usage:     "disconnect container from network",
	ArgsUsage: "<container> <network>",
	Action: func(clix *cli.Context) error {
		c, err := client.NewClient(clix.String("address"))
		if err != nil {
			return err
		}
		defer c.Close()

		container := clix.Args().First()
		network := clix.Args().Get(1)
		if container == "" || network == "" {
			cli.ShowSubcommandHelp(clix)
			return fmt.Errorf("container and network must be specified")
		}

		ctx := context.Background()
		if _, err := c.Disconnect(ctx, &api.DisconnectRequest{
			Network:   network,
			Container: container,
		}); err != nil {
			return err
		}

		fmt.Printf("disconnected %s from %s\n", container, network)

		return nil
	},
}
