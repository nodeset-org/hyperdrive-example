package hdmodule

import (
	"encoding/json"
	"fmt"

	//"github.com/goccy/go-json"
	"github.com/nodeset-org/hyperdrive-example/adapter/config"
	"github.com/nodeset-org/hyperdrive-example/adapter/utils"
	"github.com/urfave/cli/v2"
)

func getConfigMetadata(c *cli.Context) error {
	// Get the request
	_, err := utils.HandleKeyedRequest[*utils.KeyedRequest](c)
	if err != nil {
		return err
	}

	// Get the config
	cfgMgr := config.NewAdapterConfigManager()
	cfg, err := cfgMgr.LoadConfigFromDisk()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// Handle no config file by using the default
	if cfg == nil {
		cfg = config.NewExampleConfig()
	}

	// Test
	bytes, err := json.Marshal(cfg.ServerConfig)
	if err != nil {
		return fmt.Errorf("error marshalling config: %w", err)
	}
	bytes, err = json.Marshal(cfg.SubConfig)
	if err != nil {
		return fmt.Errorf("error marshalling config: %w", err)
	}
	// EndTest

	// Create the response
	bytes, err = json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshalling config: %w", err)
	}

	// Print it
	fmt.Println(string(bytes))
	return nil
}
