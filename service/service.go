package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nodeset-org/hyperdrive-example/service/api"
	"github.com/nodeset-org/hyperdrive-example/shared"
	"github.com/nodeset-org/hyperdrive-example/shared/config"
	"github.com/urfave/cli/v2"
)

// Run
func main() {
	// Initialize application
	app := cli.NewApp()

	// Set application info
	app.Name = "Example Module Service"
	app.Usage = "Service for the Hyperdrive example module"
	app.Version = shared.Version
	app.Authors = []*cli.Author{
		{
			Name:  "Nodeset",
			Email: "info@nodeset.io",
		},
	}
	app.Copyright = "(C) 2024 NodeSet LLC"

	configFileFlag := &cli.StringFlag{
		Name:     "config-file",
		Aliases:  []string{"c"},
		Usage:    "The path of the configuration file to load",
		Required: true,
	}
	ipFlag := &cli.StringFlag{
		Name:    "ip",
		Aliases: []string{"i"},
		Usage:   "The IP address to bind the API server to",
		Value:   "127.0.0.1",
	}
	portFlag := &cli.UintFlag{
		Name:    "port",
		Aliases: []string{"p"},
		Usage:   "The port to bind the API server to",
		Value:   uint(shared.DefaultServerApiPort),
	}
	apiKeyFlag := &cli.StringFlag{
		Name:     "api-key",
		Aliases:  []string{"k"},
		Usage:    "Path of the key to use for authenticating incoming API requests",
		Required: true,
	}

	app.Flags = []cli.Flag{
		configFileFlag,
		ipFlag,
		portFlag,
		apiKeyFlag,
	}
	app.Action = func(c *cli.Context) error {
		// Get the config
		configFile := c.String(configFileFlag.Name)
		cfgMgr := config.NewConfigManager(configFile)
		cfg, err := cfgMgr.LoadConfigFromFile()
		if err != nil {
			return fmt.Errorf("error loading config file: %w", err)
		}
		if cfg == nil {
			fmt.Fprintf(os.Stderr, "Config file [%s] has not been created yet. Please configure the service first.\n", configFile)
			return nil
		}

		// Create the logger
		loggerImpl, err := shared.NewFileLogger(shared.ServiceLogFile)
		if err != nil {
			return fmt.Errorf("error creating logger: %w", err)
		}
		logger := slog.New(slog.NewTextHandler(loggerImpl, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

		// Wait group to handle graceful stopping
		stopWg := new(sync.WaitGroup)

		// Start the server after the task loop so it can log into NodeSet before this starts serving registration status checks
		ip := c.String(ipFlag.Name)
		port := c.Uint64(portFlag.Name)
		serverMgr, err := api.NewApiServer(ip, uint16(port), cfgMgr, logger, stopWg)
		if err != nil {
			return fmt.Errorf("error creating server manager: %w", err)
		}

		// Handle process closures
		termListener := make(chan os.Signal, 1)
		signal.Notify(termListener, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-termListener
			fmt.Println("Shutting down daemon...")
			serverMgr.Stop()
		}()

		// Run the daemon until closed
		fmt.Println("Example module online.")
		fmt.Printf("API calls are being logged to: %s\n", shared.ServiceLogFile)
		fmt.Printf("To view them, use `hyperdrive service daemon-logs %s\n", shared.ServiceContainerName)
		stopWg.Wait()

		_ = loggerImpl.Close()
		fmt.Println("Daemon stopped.")
		return nil
	}

	// Run application
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
