package hdmodule

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/nodeset-org/hyperdrive-example/adapter/config"
	"github.com/nodeset-org/hyperdrive-example/adapter/config/ids"
	v0_1_0 "github.com/nodeset-org/hyperdrive-example/adapter/config/v0.1.0"
	"github.com/nodeset-org/hyperdrive-example/adapter/utils"
	hdconfig "github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/urfave/cli/v2"
)

// Request format for `process-config`
type processConfigRequest struct {
	utils.KeyedRequest

	// The config instance to process
	Config *hdconfig.HyperdriveSettings `json:"config"`
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

	// This is where any examples of validation will go when added
	errors := []string{}

	// Get the open ports
	ports := map[string]uint16{}
	if settings.ServerConfig.PortMode != v0_1_0.PortMode_Closed {
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
