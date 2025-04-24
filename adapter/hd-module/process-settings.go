package hdmodule

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/nodeset-org/hyperdrive-example/adapter/config"
	"github.com/nodeset-org/hyperdrive-example/adapter/config/ids"
	v0_1_0 "github.com/nodeset-org/hyperdrive-example/adapter/config/v0.1.0"
	"github.com/nodeset-org/hyperdrive-example/adapter/utils"
	hdconfig "github.com/nodeset-org/hyperdrive/config"
	modconfig "github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/urfave/cli/v2"
)

// Request format for `process-settings`
type processSettingsRequest struct {
	// The current config settings
	CurrentSettings *hdconfig.HyperdriveSettings `json:"currentSettings"`

	// The new (proposed) config settings
	NewSettings *hdconfig.HyperdriveSettings `json:"newSettings"`
}

// Response format for `process-settings`
type processSettingsResponse struct {
	// A list of errors that occurred during processing, if any
	Errors []string `json:"errors"`

	// A list of ports that will be exposed
	Ports map[string]uint16 `json:"ports"`

	// A list of services that need to be restarted as a result of the new settings
	ServicesToRestart []string `json:"servicesToRestart"`
}

// Handle the `process-settings` command
func processSettings(c *cli.Context) error {
	// Get the request
	request, err := utils.HandleRequest[*processSettingsRequest](c)
	if err != nil {
		return err
	}

	// Process the settings
	response, err := processSettingsImpl(request.CurrentSettings, request.NewSettings)
	if err != nil {
		return err
	}

	// Marshal it
	bytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("error marshalling process-settings response: %w", err)
	}

	// Print it
	fmt.Println(string(bytes))
	return nil
}

// Process the settings
func processSettingsImpl(oldHdSettings *hdconfig.HyperdriveSettings, newHdSettings *hdconfig.HyperdriveSettings) (*processSettingsResponse, error) {
	// Construct the old (current) module settings from the Hyperdrive config
	var oldSettings config.ExampleConfigSettings
	oldModInstance, exists := oldHdSettings.Modules[utils.FullyQualifiedModuleName]
	if !exists {
		// Create an instance with the default settings
		cfg := config.NewExampleConfig()
		oldModSettings := modconfig.CreateModuleSettings(cfg)
		err := oldModSettings.ConvertToKnownType(&oldSettings)
		if err != nil {
			return nil, fmt.Errorf("error creating default settings: %w", err)
		}
	} else {
		err := oldModInstance.DeserializeSettingsIntoKnownType(&oldSettings)
		if err != nil {
			return nil, fmt.Errorf("error loading old settings: %w", err)
		}
	}

	// Construct the new (proposed) module settings from the Hyperdrive config
	var newSettings config.ExampleConfigSettings
	newModInstance, exists := newHdSettings.Modules[utils.FullyQualifiedModuleName]
	if !exists {
		return nil, fmt.Errorf("could not find new settings for %s", utils.FullyQualifiedModuleName)
	}
	err := newModInstance.DeserializeSettingsIntoKnownType(&newSettings)
	if err != nil {
		return nil, fmt.Errorf("error loading new settings: %w", err)
	}

	// This is where any examples of validation will go when added
	errors := []string{}

	// Get the open ports
	ports := map[string]uint16{}
	if newSettings.ServerConfig.PortMode != v0_1_0.PortMode_Closed {
		ports[ids.ServerConfigID.String()+"/"+ids.PortID.String()] = uint16(newSettings.ServerConfig.Port)
	}

	// Get the list of services that need to be restarted
	servicesToRestart, err := newSettings.GetChangedServices(&oldSettings)
	if err != nil {
		return nil, fmt.Errorf("error getting changed services: %w", err)
	}
	if servicesToRestart == nil {
		servicesToRestart = []string{}
	}

	// Create the response
	response := &processSettingsResponse{
		Errors:            errors,
		Ports:             ports,
		ServicesToRestart: servicesToRestart,
	}
	return response, nil
}
