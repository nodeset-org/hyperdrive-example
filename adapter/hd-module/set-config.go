package hdmodule

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive-example/adapter/config"
	"github.com/nodeset-org/hyperdrive-example/adapter/utils"
	"github.com/urfave/cli/v2"
)

// Request format for `set-config`
type setConfigRequest struct {
	utils.KeyedRequest

	// The config instance to process
	Config map[string]any `json:"config"`
}

// Handle the `set-config` command
func setConfig(c *cli.Context) error {
	// Get the request
	request, err := utils.HandleKeyedRequest[*setConfigRequest](c)
	if err != nil {
		return err
	}

	// Get the config
	cfg, err := config.ConvertFromInstance(request.Config)
	if err != nil {
		return err
	}

	// Make a config manager
	cfgMgr := config.NewAdapterConfigManager()
	cfgMgr.AdapterConfig = cfg

	// Save it
	err = cfgMgr.SaveConfigToDisk()
	if err != nil {
		return fmt.Errorf("error saving config: %w", err)
	}
	return nil
}
