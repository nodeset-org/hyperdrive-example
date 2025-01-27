package hdmodule

import (
	"fmt"

	"github.com/nodeset-org/hyperdrive-example/adapter/config"
	"github.com/nodeset-org/hyperdrive-example/adapter/utils"
	hdconfig "github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/urfave/cli/v2"
)

// Request format for `set-settings`
type setSettingsRequest struct {
	utils.KeyedRequest

	// The config settings to save
	Settings *hdconfig.HyperdriveSettings `json:"settings"`
}

// Handle the `set-settings` command
func setSettings(c *cli.Context) error {
	// Get the request
	request, err := utils.HandleKeyedRequest[*setSettingsRequest](c)
	if err != nil {
		return err
	}

	// Construct the module settings from the Hyperdrive config
	modInstance, exists := request.Settings.Modules[utils.FullyQualifiedModuleName]
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
		return fmt.Errorf("error saving settings: %w", err)
	}
	return nil
}
