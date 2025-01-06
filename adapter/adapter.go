package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nodeset-org/hyperdrive-example/adapter/app"
	hdmodule "github.com/nodeset-org/hyperdrive-example/adapter/hd-module"
	"github.com/urfave/cli/v2"
)

// Run the adapter
func main() {
	app := app.CreateApp()

	// Register extra commands
	hdmodule.RegisterCommands(app)

	// Include the idle command to run forever
	app.Commands = append(app.Commands, &cli.Command{
		Name:    "idle",
		Aliases: []string{"i"},
		Usage:   "Run the adapter in idle mode, waiting forever until stopped.",
		Action: func(c *cli.Context) error {
			// Wait group to handle graceful stopping
			stopWg := new(sync.WaitGroup)
			stopWg.Add(1)

			// Handle process closures
			termListener := make(chan os.Signal, 1)
			signal.Notify(termListener, os.Interrupt, syscall.SIGTERM)
			go func() {
				<-termListener
				fmt.Println("Shutting down adapter...")
				stopWg.Done()
			}()

			// Run the adapter until closed
			fmt.Println("Adapter online.")
			stopWg.Wait()
			fmt.Println("Adapter stopped.")
			return nil
		},
	})

	// Run application
	fmt.Println()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println()
}
