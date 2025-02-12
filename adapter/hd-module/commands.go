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
				Name:    "upgrade-instance",
				Aliases: []string{"u"},
				Flags:   []cli.Flag{},
				Usage:   "Upgrade an instance of the module's configuration to the latest version - used when the configuration was generated with an older version of this module.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return upgradeInstance(c)
				},
			},
			{
				Name:    "process-settings",
				Aliases: []string{"ps"},
				Flags:   []cli.Flag{},
				Usage:   "Process the current settings for the module's configuration, validating it without saving.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return processSettings(c)
				},
			},
			{
				Name:    "set-settings",
				Aliases: []string{"ss"},
				Flags:   []cli.Flag{},
				Usage:   "Sets the settings for the module's configuration, saving it to disk.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return setSettings(c)
				},
			},
			{
				Name:    "start",
				Aliases: []string{"s"},
				Flags:   []cli.Flag{},
				Usage:   "Start the module's services (stopping and restarting any that changed due to the new configuration).",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return startServices(c)
				},
			},
			{
				Name:    "stop",
				Aliases: []string{"t"},
				Flags:   []cli.Flag{},
				Usage:   "Stop the module's services.",
				Action: func(c *cli.Context) error {
					// Validate args
					utils.ValidateArgCount(c, 0)

					// Run
					return stopServices(c)
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
