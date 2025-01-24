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
	Config *hdconfig.HyperdriveConfigInstance `json:"config"`
}

// Handle the `set-config` command
func setConfig(c *cli.Context) error {
	// Get the request
	request, err := utils.HandleKeyedRequest[*setConfigRequest](c)
	if err != nil {
		return err
	}

	// Construct the module instance from the Hyperdrive config
	var settings *config.ExampleConfigInstance
	for _, module := range request.Config.Modules {
		if module.Name != utils.FullyQualifiedModuleName {
			continue
		}

		modCfg := config.NewExampleConfig()
		modSettings, err := module.Settings.CreateSettingsFromMetadata(modCfg)
		if err != nil {
			return fmt.Errorf("error creating settings from metadata: %w", err)
		}
		settings = new(config.ExampleConfigInstance)
		err = modSettings.ConvertToKnownType(settings)
		if err != nil {
			return fmt.Errorf("error converting settings to known type: %w", err)
		}
	}

	// Make sure the config was found
	if settings == nil {
		return fmt.Errorf("could not find config for %s", utils.FullyQualifiedModuleName)
	}

	// Make a config manager
	cfgMgr, err := config.NewAdapterConfigManager(c)
	if err != nil {
		return fmt.Errorf("error creating config manager: %w", err)
	}
	cfgMgr.AdapterConfig = settings

	// Save it
	err = cfgMgr.SaveConfigToDisk()
	if err != nil {
		return fmt.Errorf("error saving config: %w", err)
	}
	return nil
}
