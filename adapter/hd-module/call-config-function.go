package hdmodule

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/nodeset-org/hyperdrive-example/adapter/config"
	"github.com/nodeset-org/hyperdrive-example/adapter/utils"
	hdconfig "github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/urfave/cli/v2"
)

type CallConfigFunctionRequest struct {
	FuncName string                       `json:"funcName"`
	Settings *hdconfig.HyperdriveSettings `json:"settings"`
}

func callConfigFunction(c *cli.Context) error {
	request, err := utils.HandleRequest[*CallConfigFunctionRequest](c)
	if err != nil {
		return err
	}

	modInstance, exists := request.Settings.Modules[utils.FullyQualifiedModuleName]
	if !exists {
		return fmt.Errorf("could not find config for %s", utils.FullyQualifiedModuleName)
	}

	var settings config.ExampleConfigSettings
	err = modInstance.DeserializeSettingsIntoKnownType(&settings)
	if err != nil {
		return fmt.Errorf("error loading settings: %w", err)
	}

	switch request.FuncName {
	case "GetDerivedValued":
		bytes, err := json.Marshal(settings.GetDerivedValue())
		if err != nil {
			return fmt.Errorf("error marshalling derived value: %w", err)
		}
		fmt.Println(string(bytes))
		return nil
	default:
		return fmt.Errorf("unknown function: %s", request.FuncName)
	}
}
