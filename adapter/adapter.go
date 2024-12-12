package main

import (
	"fmt"
	"os"

	"github.com/nodeset-org/hyperdrive-example/adapter/config"
	hdmodule "github.com/nodeset-org/hyperdrive-example/adapter/hd-module"
	"github.com/nodeset-org/hyperdrive-example/adapter/utils"
	"github.com/nodeset-org/hyperdrive-example/shared"
	"github.com/urfave/cli/v2"
)

const ()

// Flags
var (
	keyFileFlag *cli.StringFlag = &cli.StringFlag{
		Name:     "secret",
		Aliases:  []string{"s"},
		Usage:    "The path to the secret key file for authentication",
		Required: true,
	}
)

func main() {
	// Initialize application
	app := cli.NewApp()

	// Set application info
	app.Name = "Example Adapter"
	app.Usage = "Adapter for the Hyperdrive example module"
	app.Version = shared.Version
	app.Authors = []*cli.Author{
		{
			Name:  "Nodeset",
			Email: "info@nodeset.io",
		},
	}
	app.Copyright = "(c) 2024 NodeSet LLC"

	// Enable Bash Completion
	app.EnableBashCompletion = true

	// Set application flags
	app.Flags = []cli.Flag{
		keyFileFlag,
	}

	// Register commands
	config.RegisterCommands(app)
	hdmodule.RegisterCommands(app)

	app.Before = func(c *cli.Context) error {
		// Make the authenticator
		keyFile := c.String(keyFileFlag.Name)
		if keyFile == "" {
			return fmt.Errorf("secret key file is required")
		}
		auth, err := utils.NewAuthenticator(keyFile)
		if err != nil {
			return err
		}
		c.App.Metadata[utils.AuthenticatorMetadataKey] = auth

		return nil
	}
	app.BashComplete = func(c *cli.Context) {
		// Load the context and flags prior to autocomplete
		err := app.Before(c)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}

		// Run the default autocomplete
		cli.DefaultAppComplete(c)
	}

	// Run application
	fmt.Println()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println()
}
