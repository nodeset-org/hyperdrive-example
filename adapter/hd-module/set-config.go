package hdmodule

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive-example/adapter/config"
	"github.com/nodeset-org/hyperdrive-example/adapter/utils"
	hdconfig "github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/urfave/cli/v2"
)

// Request format for `set-config`
type setConfigRequest struct {
	utils.KeyedRequest

	// The config instance to process
	Config *hdconfig.HyperdriveSettings `json:"config"`
}

// Handle the `set-config` command
func setConfig(c *cli.Context) error {
	// Get the request
	request, err := utils.HandleKeyedRequest[*setConfigRequest](c)
	if err != nil {
		return err
	}

	// Construct the module settings from the Hyperdrive config
	modInstance, exists := request.Config.Modules[utils.FullyQualifiedModuleName]
	if !exists {
		return fmt.Errorf("could not find config for %s", utils.FullyQualifiedModuleName)
	}
	var settings config.ExampleConfigSettings
	err = modInstance.DeserializeSettingsIntoKnownType(&settings)
	if err != nil {
		return fmt.Errorf("error loading settings: %w", err)
	}

	// Make a config manager
	cfgMgr, err := config.NewAdapterConfigManager(c)
	if err != nil {
		return fmt.Errorf("error creating config manager: %w", err)
	}
	cfgMgr.AdapterConfig = &settings

	// Save it
	err = cfgMgr.SaveConfigToDisk()
	if err != nil {
		return fmt.Errorf("error saving config: %w", err)
	}
	return nil
}
