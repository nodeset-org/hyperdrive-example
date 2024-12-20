package main

import (
	"fmt"
	"os"

	"github.com/nodeset-org/hyperdrive-example/adapter/app"
	hdmodule "github.com/nodeset-org/hyperdrive-example/adapter/hd-module"
)

// Run the adapter
func main() {
	app := app.CreateApp()

	// Register commands
	hdmodule.RegisterCommands(app)

	// Run application
	fmt.Println()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println()
}
