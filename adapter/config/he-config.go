package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"

	"github.com/nodeset-org/hyperdrive-example/adapter/config/ids"
	"github.com/nodeset-org/hyperdrive-example/adapter/utils"
	"github.com/nodeset-org/hyperdrive-example/shared"
	nativecfg "github.com/nodeset-org/hyperdrive-example/shared/config"
	hdconfig "github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

const (
	FloatThreshold float64 = 75.0
)

type PortMode string

const (
	PortMode_Closed    PortMode = "closed"
	PortMode_Localhost PortMode = "localhost"
	PortMode_External  PortMode = "external"
)

type ExampleConfig struct {
	ExampleBool hdconfig.BoolParameter

	ExampleInt hdconfig.IntParameter

	ExampleUint hdconfig.UintParameter

	ExampleFloat hdconfig.FloatParameter

	ExampleString hdconfig.StringParameter

	ExampleChoice hdconfig.ChoiceParameter[nativecfg.ExampleOption]

	SubConfig *SubConfig

	ServerConfig *ServerConfig
}

type ExampleConfigInstance struct {
	Version       string                  `json:"version"`
	ExampleBool   bool                    `json:"exampleBool"`
	ExampleInt    int64                   `json:"exampleInt"`
	ExampleUint   uint64                  `json:"exampleUint"`
	ExampleFloat  float64                 `json:"exampleFloat"`
	ExampleString string                  `json:"exampleString"`
	ExampleChoice nativecfg.ExampleOption `json:"exampleChoice"`

	SubConfig struct {
		SubExampleBool   bool                    `json:"subConfigBool"`
		SubExampleChoice nativecfg.ExampleOption `json:"subConfigChoice"`
	} `json:"subConfig"`

	ServerConfig struct {
		Port     uint64   `json:"port" yaml:"port"`
		PortMode PortMode `json:"portMode" yaml:"portMode"`
	} `json:"server" yaml:"server"`
}

func NewExampleConfig() *ExampleConfig {
	cfg := &ExampleConfig{}

	// ExampleBool
	cfg.ExampleBool.ID = ids.ExampleBoolID
	cfg.ExampleBool.Name = "Example Boolean"
	cfg.ExampleBool.Description.Default = "This is an example of a boolean parameter. It doesn't directly affect the service, but it does control the behavior of some other config parameters."
	cfg.ExampleBool.AffectedContainers = []string{shared.ServiceContainerName}

	// ExampleInt
	cfg.ExampleInt.ID = ids.ExampleIntID
	cfg.ExampleInt.Name = "Example Integer"
	cfg.ExampleInt.Description.Default = "This is an example of an integer parameter."
	cfg.ExampleInt.AffectedContainers = []string{shared.ServiceContainerName}

	// ExampleUint
	cfg.ExampleUint.ID = ids.ExampleUintID
	cfg.ExampleUint.Name = "Example Unsigned Integer"
	cfg.ExampleUint.Description.Default = "This is an example of an unsigned integer parameter."
	cfg.ExampleUint.AffectedContainers = []string{shared.ServiceContainerName}

	// ExampleFloat
	cfg.ExampleFloat.ID = ids.ExampleFloatID
	cfg.ExampleFloat.Name = "Example Float"
	cfg.ExampleFloat.Description.Default = "This is an example of a float parameter with a minimum and maximum set."
	cfg.ExampleFloat.Default = 50
	cfg.ExampleFloat.MinValue = 0.0
	cfg.ExampleFloat.MaxValue = 100.0
	cfg.ExampleFloat.AffectedContainers = []string{shared.ServiceContainerName}

	// ExampleString
	cfg.ExampleString.ID = ids.ExampleStringID
	cfg.ExampleString.Name = "Example String"
	cfg.ExampleString.Description.Default = "This is an example of a string parameter. It has a max length and regex pattern set."
	cfg.ExampleString.MaxLength = 10
	cfg.ExampleString.Regex = "^[a-zA-Z]*$"
	cfg.ExampleString.AffectedContainers = []string{shared.ServiceContainerName}

	// Options for ExampleChoice
	options := make([]hdconfig.ParameterOption[nativecfg.ExampleOption], 3)
	options[0].Name = "One"
	options[0].Description.Default = "This is the first option."
	options[0].Value = nativecfg.ExampleOption_One

	thresholdString := strconv.FormatFloat(FloatThreshold, 'f', -1, 64)
	options[1].Name = "Two"
	options[1].Description.Default = "This is the second option. It is hidden when ExampleFloat is less than " + thresholdString + "."
	options[1].Description.Template = fmt.Sprintf("{{if lt .GetValue %s %s}}This option is hidden because the float is less than %s.{{else}}This option is visible because the float is greater than or equal to %s.{{end}}", ids.ExampleFloatID, thresholdString, thresholdString, thresholdString)
	options[1].Value = nativecfg.ExampleOption_Two
	options[1].Disabled.Default = true
	options[1].Disabled.Template = "{{if eq .GetValue " + ids.ExampleBoolID.String() + " true}}false{{else}}{{.UseDefault}}{{end}}"

	options[2].Name = "Three"
	options[2].Description.Default = "This is the third option."
	options[2].Value = nativecfg.ExampleOption_Three

	// ExampleChoice
	cfg.ExampleChoice.ID = ids.ExampleChoiceID
	cfg.ExampleChoice.Name = "Example Choice"
	cfg.ExampleChoice.Description.Default = "This is an example of a choice parameter between multiple options."
	cfg.ExampleChoice.Options = options
	cfg.ExampleChoice.Default = options[0].Value
	cfg.ExampleChoice.AffectedContainers = []string{}

	// Subconfigs
	cfg.SubConfig = NewSubConfig()
	cfg.ServerConfig = NewServerConfig()

	return cfg
}

type SubConfig struct {
	hdconfig.SectionHeader

	SubExampleBool hdconfig.BoolParameter

	SubExampleChoice hdconfig.ChoiceParameter[nativecfg.ExampleOption]
}

func NewSubConfig() *SubConfig {
	cfg := &SubConfig{}
	cfg.ID = ids.SubConfigID
	cfg.Name = "Sub Config"
	cfg.Description.Default = "This is a sub-section of the main configuration."
	cfg.Hidden.Default = true
	cfg.Hidden.Template = "{{if eq .GetValue " + ids.ExampleBoolID.String() + " true}}false{{else}}true{{end}}"

	// SubExampleBool
	cfg.SubExampleBool.ID = ids.SubExampleBoolID
	cfg.SubExampleBool.Name = "Sub Example Boolean"
	cfg.SubExampleBool.Description.Default = "This is an example of a boolean parameter in a sub-section."

	// Options for SubExampleChoice
	options := make([]hdconfig.ParameterOption[nativecfg.ExampleOption], 2)
	options[0].Name = "One"
	options[0].Description.Default = "This is the first option."
	options[0].Value = nativecfg.ExampleOption_One

	options[1].Name = "Two"
	options[1].Description.Default = "This is the second option."
	options[1].Value = nativecfg.ExampleOption_Two

	// SubExampleChoice
	cfg.SubExampleChoice.ID = ids.SubExampleChoiceID
	cfg.SubExampleChoice.Name = "Sub Example Choice"
	cfg.SubExampleChoice.Description.Default = "This is an example of a choice parameter between multiple options in a sub-section."
	cfg.SubExampleChoice.Options = options
	cfg.SubExampleChoice.Default = options[1].Value

	return cfg
}

type ServerConfig struct {
	hdconfig.SectionHeader

	Port hdconfig.UintParameter

	PortMode hdconfig.ChoiceParameter[PortMode]
}

func NewServerConfig() *ServerConfig {
	cfg := &ServerConfig{}
	cfg.ID = ids.ServerConfigID
	cfg.Name = "Service Config"
	cfg.Description.Default = "This is the configuration for the module's service. This isn't used by the service directly, but it is used by Hyperdrive itself in the service's Docker Compose file template to configure the service during its starting process."

	// Port
	cfg.Port.ID = ids.PortID
	cfg.Port.Name = "API Port"
	cfg.Port.Description.Default = "This is the API port the server should run on."
	cfg.Port.Default = uint64(shared.DefaultServerApiPort)
	cfg.Port.MinValue = 0
	cfg.Port.MaxValue = 65535
	cfg.Port.AffectedContainers = []string{shared.ServiceContainerName}

	// Options for PortMode
	options := make([]hdconfig.ParameterOption[PortMode], 3)
	options[0].Name = "Closed"
	options[0].Description.Default = "The API is only accessible to internal Docker container traffic."
	options[0].Value = PortMode_Closed

	options[1].Name = "Localhost Only"
	options[1].Description.Default = "The API is accessible from internal Docker containers and your own local machine, but no other external machines."
	options[1].Value = PortMode_Localhost

	options[2].Name = "All External Traffic"
	options[2].Description.Default = "The port is accessible to everything, including external machines.\n\n[orange]Use with caution!"
	options[2].Value = PortMode_External

	// PortMode
	cfg.PortMode.ID = ids.PortModeID
	cfg.PortMode.Name = "Expose API Port"
	cfg.PortMode.Description.Default = "Determine how the server's HTTP API restricts its access from various sources."
	cfg.PortMode.Options = options
	cfg.PortMode.Default = options[0].Value
	cfg.PortMode.AffectedContainers = []string{shared.ServiceContainerName}

	return cfg
}

func (cfg ExampleConfig) GetVersion() string {
	return shared.Version
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
	return []hdconfig.ISection{
		cfg.SubConfig,
		cfg.ServerConfig,
	}
}

func (cfg SubConfig) GetParameters() []hdconfig.IParameter {
	return []hdconfig.IParameter{
		&cfg.SubExampleBool,
		&cfg.SubExampleChoice,
	}
}

func (cfg SubConfig) GetSections() []hdconfig.ISection {
	return []hdconfig.ISection{}
}

func (cfg ServerConfig) GetParameters() []hdconfig.IParameter {
	return []hdconfig.IParameter{
		&cfg.Port,
		&cfg.PortMode,
	}
}

func (cfg ServerConfig) GetSections() []hdconfig.ISection {
	return []hdconfig.ISection{}
}

func CreateInstanceFromNativeConfig(native *nativecfg.NativeExampleConfig) *ExampleConfigInstance {
	instance := &ExampleConfigInstance{
		ExampleBool:   native.ExampleBool,
		ExampleInt:    native.ExampleInt,
		ExampleUint:   native.ExampleUint,
		ExampleFloat:  native.ExampleFloat,
		ExampleString: native.ExampleString,
		ExampleChoice: native.ExampleChoice,
	}
	instance.SubConfig.SubExampleBool = native.SubConfig.SubExampleBool
	instance.SubConfig.SubExampleChoice = native.SubConfig.SubExampleChoice

	/*
		cfg := NewExampleConfig()
		instance := hdconfig.CreateModuleConfigurationInstance(cfg)
		mustSetParameter(instance, cfg.ExampleBool.ID, native.ExampleBool)
		mustSetParameter(instance, cfg.ExampleInt.ID, native.ExampleInt)
		mustSetParameter(instance, cfg.ExampleUint.ID, native.ExampleUint)
		mustSetParameter(instance, cfg.ExampleFloat.ID, native.ExampleFloat)
		mustSetParameter(instance, cfg.ExampleString.ID, native.ExampleString)
		mustSetParameter(instance, cfg.ExampleChoice.ID, native.ExampleChoice)
	*/
	return instance
}

// Sets a parameter value on an instance. Panics if it fails.
func mustSetParameter(instance hdconfig.IInstanceContainer, id hdconfig.Identifier, value any) {
	param, err := instance.GetParameter(id)
	if err != nil {
		panic(fmt.Errorf("can't retrieve parameter instance [%s]: %w", id, err))
	}
	err = param.SetValue(value)
	if err != nil {
		panic(fmt.Errorf("can't set parameter instance [%s]: %w", id, err))
	}
}

func ConvertInstanceToNativeConfig(instance *ExampleConfigInstance) *nativecfg.NativeExampleConfig {
	native := &nativecfg.NativeExampleConfig{
		ExampleBool:   instance.ExampleBool,
		ExampleInt:    instance.ExampleInt,
		ExampleUint:   instance.ExampleUint,
		ExampleFloat:  instance.ExampleFloat,
		ExampleString: instance.ExampleString,
		ExampleChoice: instance.ExampleChoice,
	}
	native.SubConfig.SubExampleBool = instance.SubConfig.SubExampleBool
	native.SubConfig.SubExampleChoice = instance.SubConfig.SubExampleChoice
	/*
		native.ExampleBool = mustGetNativeValue[bool](instance, ids.ExampleBoolID)
		native.ExampleInt = mustGetNativeValue[int64](instance, ids.ExampleIntID)
		native.ExampleUint = mustGetNativeValue[uint64](instance, ids.ExampleUintID)
		native.ExampleFloat = mustGetNativeValue[float64](instance, ids.ExampleFloatID)
		native.ExampleString = mustGetNativeValue[string](instance, ids.ExampleStringID)
		native.ExampleChoice = nativecfg.ExampleOption(mustGetNativeValue[string](instance, ids.ExampleChoiceID))

		subCfg, err := instance.GetSection(ids.SubConfigID)
		if err != nil {
			panic(fmt.Errorf("can't retrieve sub-config section: %w", err))
		}
		native.SubConfig.SubExampleBool = mustGetNativeValue[bool](subCfg, ids.SubExampleBoolID)
		native.SubConfig.SubExampleChoice = nativecfg.ExampleOption(mustGetNativeValue[string](subCfg, ids.SubExampleChoiceID))
	*/
	return native
}

func mustGetNativeValue[Type any](instance hdconfig.IInstanceContainer, id hdconfig.Identifier) Type {
	param, err := instance.GetParameter(id)
	if err != nil {
		panic(fmt.Errorf("can't retrieve parameter instance [%s]: %w", id, err))
	}
	return param.GetValue().(Type)
}

// Configuration manager
type AdapterConfigManager struct {
	// The adapter configuration instance
	AdapterConfig *ExampleConfigInstance

	// The native configuration manager
	nativeConfigManager *nativecfg.ConfigManager

	// The path to the adapter configuration file
	adapterConfigPath string
}

// Create a new configuration manager for the adapter
func NewAdapterConfigManager(c *cli.Context) (*AdapterConfigManager, error) {
	configDir := c.String(utils.ConfigDirFlag.Name)
	if configDir == "" {
		return nil, fmt.Errorf("config directory is required")
	}
	return &AdapterConfigManager{
		nativeConfigManager: nativecfg.NewConfigManager(filepath.Join(configDir, utils.ServiceConfigFile)),
		adapterConfigPath:   filepath.Join(configDir, utils.AdapterConfigFile),
	}, nil
}

// Load the configuration from disk
func (m *AdapterConfigManager) LoadConfigFromDisk() (*ExampleConfigInstance, error) {
	// Load the native config
	nativeCfg, err := m.nativeConfigManager.LoadConfigFromFile()
	if err != nil {
		return nil, fmt.Errorf("error loading service config: %w", err)
	}

	// Check if the adapter config exists
	_, err = os.Stat(m.adapterConfigPath)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}

	// Load it
	bytes, err := os.ReadFile(m.adapterConfigPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file [%s]: %w", m.adapterConfigPath, err)
	}

	// Deserialize it
	modCfg := CreateInstanceFromNativeConfig(nativeCfg)
	err = yaml.Unmarshal(bytes, &modCfg.ServerConfig)
	if err != nil {
		return nil, fmt.Errorf("error deserializing adapter config file [%s]: %w", m.adapterConfigPath, err)
	}
	m.AdapterConfig = modCfg
	return modCfg, nil
}

// Save the configuration to a file. If the config hasn't been loaded yet, this doesn't do anything.
func (m *AdapterConfigManager) SaveConfigToDisk() error {
	if m.AdapterConfig == nil {
		return nil
	}

	// Save the native config
	nativeCfg := ConvertInstanceToNativeConfig(m.AdapterConfig)
	m.nativeConfigManager.Config = nativeCfg
	err := m.nativeConfigManager.SaveConfigToFile()
	if err != nil {
		return fmt.Errorf("error saving service config: %w", err)
	}

	// Serialize the adapter config
	bytes, err := yaml.Marshal(m.AdapterConfig.ServerConfig)
	if err != nil {
		return fmt.Errorf("error serializing adapter config: %w", err)
	}

	// Write it
	err = os.WriteFile(m.adapterConfigPath, bytes, nativecfg.ConfigFileMode)
	if err != nil {
		return fmt.Errorf("error writing adapter config file [%s]: %w", m.adapterConfigPath, err)
	}
	return nil
}
