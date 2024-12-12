package hdmodule

import (
	"github.com/nodeset-org/hyperdrive-example/adapter/utils"
	"github.com/urfave/cli/v2"
)

// Handles `hd-module` commands
func RegisterCommands(app *cli.App) {
	app.Commands = append(app.Commands, &cli.Command{
		Name:    "hd-module",
		Aliases: []string{"hd"},
		Usage:   "Handle Hyperdrive module commands",
		Subcommands: []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Flags:   []cli.Flag{},
				Usage:   "Print the module version.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return version()
				},
			},
			{
				Name:    "get-log-file",
				Aliases: []string{"l"},
				Flags:   []cli.Flag{},
				Usage:   "Get the path to a log file.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return getLogFile(c)
				},
			},
			{
				Name:    "get-config-metadata",
				Aliases: []string{"c"},
				Flags:   []cli.Flag{},
				Usage:   "Get the metadata for the module's configuration, representing how to present the parameters to the user.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return getConfigMetadata(c)
				},
			},
			{
				Name:    "get-config-instance",
				Aliases: []string{"i"},
				Flags:   []cli.Flag{},
				Usage:   "Get the instance of the module's configuration, with the values all set.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return getConfigInstance(c)
				},
			},
			{
				Name:    "process-config",
				Aliases: []string{"p"},
				Flags:   []cli.Flag{},
				Usage:   "Process the module's configuration, validating it without saving.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return processConfig(c)
				},
			},
			{
				Name:    "set-config",
				Aliases: []string{"s"},
				Flags:   []cli.Flag{},
				Usage:   "Sets the module's configuration, saving it to disk.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return setConfig(c)
				},
			},
			{
				Name:    "get-containers",
				Aliases: []string{"t"},
				Flags:   []cli.Flag{},
				Usage:   "Get the list of containers owned by this module.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return getContainers(c)
				},
			},
			{
				Name:    "run",
				Aliases: []string{"r"},
				Flags:   []cli.Flag{},
				Usage:   "Run a command.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return run(c)
				},
			},
		},
	})
}