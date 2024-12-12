package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/nodeset-org/hyperdrive-example/adapter/utils"
	"github.com/nodeset-org/hyperdrive-example/shared"
	"github.com/nodeset-org/hyperdrive-example/shared/api"
)

// Get one of the config parameters
func getParam(param string) error {
	// Create the logger
	logHandler, err := shared.NewFileLogger(utils.AdapterLogPath)
	if err != nil {
		return fmt.Errorf("error creating logger: %w", err)
	}
	logger := slog.New(slog.NewTextHandler(logHandler, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	defer logHandler.Close()

	// Create the configuration manager
	cfgMgr := NewAdapterConfigManager()
	cfg, err := cfgMgr.LoadConfigFromDisk()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// Create an API client
	projectName := os.Getenv("HD_PROJECT_NAME")
	serviceName := projectName + "_" + shared.ServiceContainerName
	apiClient, err := api.NewApiClient(logger, serviceName, uint(cfg.ServerConfig.Port.Value))
	if err != nil {
		return fmt.Errorf("error creating API client: %w", err)
	}

	// Run the get call
	resp, err := apiClient.GetParam(param)
	if err != nil {
		return fmt.Errorf("error getting parameter: %w", err)
	}

	// Print the response
	fmt.Println(resp.Data.Value)
	return nil
}
