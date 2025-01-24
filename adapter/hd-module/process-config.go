package hdmodule

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/nodeset-org/hyperdrive-example/adapter/config"
	"github.com/nodeset-org/hyperdrive-example/adapter/config/ids"
	"github.com/nodeset-org/hyperdrive-example/adapter/utils"
	hdconfig "github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/urfave/cli/v2"
)

// Request format for `process-config`
type processConfigRequest struct {
	utils.KeyedRequest

	// The config instance to process
	Config *hdconfig.HyperdriveConfigInstance `json:"config"`
	//Config *config.ExampleConfigInstance `json:"config"`
}

// Response format for `process-config`
type processConfigResponse struct {
	// A list of errors that occurred during processing, if any
	Errors []string `json:"errors"`

	// A list of ports that will be exposed
	Ports map[string]uint16 `json:"ports"`
}

// Handle the `process-config` command
func processConfig(c *cli.Context) error {
	// Get the request
	request, err := utils.HandleKeyedRequest[*processConfigRequest](c)
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

	// This is where any examples of validation will go when added
	errors := []string{}

	// Get the open ports
	ports := map[string]uint16{}
	if settings.ServerConfig.PortMode != config.PortMode_Closed {
		ports[ids.ServerConfigID.String()+"/"+ids.PortID.String()] = uint16(settings.ServerConfig.Port)
	}

	// Create the response
	response := processConfigResponse{
		Errors: errors,
		Ports:  ports,
	}

	// Marshal it
	bytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("error marshalling process-config response: %w", err)
	}

	// Print it
	fmt.Println(string(bytes))
	return nil
}
