package config

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/nodeset-org/hyperdrive-example/adapter/utils"
	"github.com/nodeset-org/hyperdrive-example/shared"
	"github.com/nodeset-org/hyperdrive-example/shared/api"
	"github.com/urfave/cli/v2"
)

var (
	getParamFlag *cli.StringFlag = &cli.StringFlag{
		Name:  "param",
		Usage: "The parameter to get",
	}
)

// Get one of the config parameters
func getParam(c *cli.Context) error {
	// Create the logger
	logDir := utils.LogDir
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
		serviceName = utils.ComposeProject + "_" + shared.ServiceContainerName
	}
	apiClient, err := api.NewApiClient(logger, serviceName, uint(cfg.ServerConfig.Port))
	if err != nil {
		return fmt.Errorf("error creating API client: %w", err)
	}

	// Get the param
	param := c.String(getParamFlag.Name)
	if param == "" {
		fmt.Println("Please provide a parameter to get:")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		param = scanner.Text()
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
