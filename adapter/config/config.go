package config

import (
	"github.com/nodeset-org/hyperdrive-example/adapter/config/ids"
	v0_1_0 "github.com/nodeset-org/hyperdrive-example/adapter/config/v0.1.0"
	"github.com/nodeset-org/hyperdrive-example/shared"
	nativecfg "github.com/nodeset-org/hyperdrive-example/shared/config"
	hdconfig "github.com/nodeset-org/hyperdrive/modules/config"
)

const (
	DefaultUintValue uint64 = 42
)

type ExampleConfig struct {
	v0_1_0.ExampleConfig
	ExampleUint hdconfig.UintParameter
}

type ExampleConfigSettings struct {
	v0_1_0.ExampleConfigSettings
	ExampleUint uint64 `json:"exampleUint"`
}

func NewExampleConfig() *ExampleConfig {
	cfg := &ExampleConfig{
		ExampleConfig: *v0_1_0.NewExampleConfig(),
	}

	// ExampleUint
	cfg.ExampleUint.ID = ids.ExampleUintID
	cfg.ExampleUint.Name = "Example Unsigned Integer"
	cfg.ExampleUint.Description.Default = "This is an example of an unsigned integer parameter."
	cfg.ExampleUint.AffectedContainers = []string{shared.ServiceContainerName}
	cfg.ExampleUint.Default = DefaultUintValue

	return cfg
}

func (cfg ExampleConfig) GetParameters() []hdconfig.IParameter {
	return []hdconfig.IParameter{
		&cfg.ExampleBool,
		&cfg.ExampleInt,
		&cfg.ExampleUint,
		&cfg.ExampleFloat,
		&cfg.ExampleString,
		&cfg.ExampleChoice,
	}
}

func (cfg ExampleConfig) GetSections() []hdconfig.ISection {
	return cfg.ExampleConfig.GetSections()
}

func CreateInstanceFromNativeConfig(native *nativecfg.NativeExampleConfig) *ExampleConfigSettings {
	oldInstance := v0_1_0.CreateInstanceFromNativeConfig(native)
	instance := &ExampleConfigSettings{
		ExampleConfigSettings: *oldInstance,
		ExampleUint:           native.ExampleUint,
	}
	return instance
}

func ConvertInstanceToNativeConfig(instance *ExampleConfigSettings) *nativecfg.NativeExampleConfig {
	native := v0_1_0.ConvertInstanceToNativeConfig(&instance.ExampleConfigSettings)
	native.ExampleUint = instance.ExampleUint
	return native
}

func UpgradeSettings(old *v0_1_0.ExampleConfigSettings) *ExampleConfigSettings {
	return &ExampleConfigSettings{
		ExampleConfigSettings: *old,
		ExampleUint:           DefaultUintValue,
	}
}
