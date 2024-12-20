package config

import (
	"github.com/nodeset-org/hyperdrive-example/adapter/utils"
	"github.com/urfave/cli/v2"
)

const (
	TestServerEndpointEnvVarName string = "HD_TEST_SERVER_ENDPOINT"
)

// Handles `hd-module` commands
func RegisterCommands(app *cli.App) {
	app.Commands = append(app.Commands, &cli.Command{
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "Commands for interacting with the module's configuration",
		Subcommands: []*cli.Command{
			{
				Name:    "get-param",
				Aliases: []string{"g"},
				Flags:   []cli.Flag{},
				Usage:   "Get the value of a parameter.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 1)
					param := c.Args().Get(0)

					// Run
					return getParam(c, param)
				},
			},
			{
				Name:    "set-param",
				Aliases: []string{"s"},
				Flags:   []cli.Flag{},
				Usage:   "Sets the value for a parameter.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 2)
					param := c.Args().Get(0)
					value := c.Args().Get(1)

					// Run
					return setParam(c, param, value)
				},
			},
		},
	})
}
