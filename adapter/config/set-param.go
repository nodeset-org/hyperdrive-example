package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/nodeset-org/hyperdrive-example/adapter/utils"
	"github.com/nodeset-org/hyperdrive-example/shared"
	"github.com/nodeset-org/hyperdrive-example/shared/api"
	"github.com/urfave/cli/v2"
)

// Set one of the config parameters
func setParam(c *cli.Context, param string, value string) error {
	// Create the logger
	logDir := c.String(utils.LogDirFlag.Name)
	if logDir == "" {
		return fmt.Errorf("log directory flag is required")
	}
	logHandler, err := shared.NewFileLogger(filepath.Join(logDir, utils.AdapterLogFile))
	if err != nil {
		return fmt.Errorf("error creating logger: %w", err)
	}
	logger := slog.New(slog.NewTextHandler(logHandler, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	defer logHandler.Close()

	// Create the configuration manager
	cfgMgr, err := NewAdapterConfigManager(c)
	if err != nil {
		return fmt.Errorf("error creating config manager: %w", err)
	}
	cfg, err := cfgMgr.LoadConfigFromDisk()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}
	if cfg == nil {
		return fmt.Errorf("config has not been created yet")
	}

	// Create an API client
	serviceName := os.Getenv(TestServerEndpointEnvVarName)
	if serviceName == "" {
		projectName := os.Getenv("HD_PROJECT_NAME")
		serviceName = projectName + "_" + shared.ServiceContainerName
	}
	apiClient, err := api.NewApiClient(logger, serviceName, uint(cfg.ServerConfig.Port))
	if err != nil {
		return fmt.Errorf("error creating API client: %w", err)
	}

	// Run the set call
	_, err = apiClient.SetParam(param, value)
	if err != nil {
		return fmt.Errorf("error getting parameter: %w", err)
	}
	fmt.Println("Parameter set successfully")
	return nil
}
