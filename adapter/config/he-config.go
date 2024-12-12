package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strconv"

	"github.com/nodeset-org/hyperdrive-example/adapter/config/ids"
	"github.com/nodeset-org/hyperdrive-example/adapter/utils"
	hdconfig "github.com/nodeset-org/hyperdrive-example/hyperdrive/config"
	"github.com/nodeset-org/hyperdrive-example/shared"
	nativecfg "github.com/nodeset-org/hyperdrive-example/shared/config"
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
	hdconfig.ConfigurationMetadata `json:"-" yaml:"-"`

	ExampleBool hdconfig.BoolParameterMetadata `json:"-" yaml:"-"`

	ExampleInt hdconfig.IntParameterMetadata `json:"-" yaml:"-"`

	ExampleUint hdconfig.UintParameterMetadata `json:"-" yaml:"-"`

	ExampleFloat hdconfig.FloatParameterMetadata `json:"-" yaml:"-"`

	ExampleString hdconfig.StringParameterMetadata `json:"-" yaml:"-"`

	ExampleChoice hdconfig.ChoiceParameterMetadata[nativecfg.ExampleOption] `json:"-" yaml:"-"`

	SubConfig *SubConfig `json:"-" yaml:"-"`

	ServerConfig *ServerConfig `json:"-" yaml:"-"`
}

func NewExampleConfig() *ExampleConfig {
	cfg := &ExampleConfig{}

	// ExampleBool
	cfg.ExampleBool.ID = hdconfig.Identifier(ids.ExampleBoolID)
	cfg.ExampleBool.Name = "Example Boolean"
	cfg.ExampleBool.Description.Default = "This is an example of a boolean parameter. It doesn't directly affect the service, but it does control the behavior of some other config parameters."
	cfg.ExampleBool.AffectedContainers = []string{shared.ServiceContainerName}
	cfg.ExampleBool.Value = cfg.ExampleBool.Default

	// ExampleInt
	cfg.ExampleInt.ID = hdconfig.Identifier(ids.ExampleIntID)
	cfg.ExampleInt.Name = "Example Integer"
	cfg.ExampleInt.Description.Default = "This is an example of an integer parameter."
	cfg.ExampleInt.AffectedContainers = []string{shared.ServiceContainerName}
	cfg.ExampleInt.Value = cfg.ExampleInt.Default

	// ExampleUint
	cfg.ExampleUint.ID = hdconfig.Identifier(ids.ExampleUintID)
	cfg.ExampleUint.Name = "Example Unsigned Integer"
	cfg.ExampleUint.Description.Default = "This is an example of an unsigned integer parameter."
	cfg.ExampleInt.AffectedContainers = []string{shared.ServiceContainerName}
	cfg.ExampleUint.Value = cfg.ExampleUint.Default

	// ExampleFloat
	cfg.ExampleFloat.ID = hdconfig.Identifier(ids.ExampleFloatID)
	cfg.ExampleFloat.Name = "Example Float"
	cfg.ExampleFloat.Description.Default = "This is an example of a float parameter with a minimum and maximum set."
	cfg.ExampleFloat.Default = 50
	cfg.ExampleFloat.MinValue = 0.0
	cfg.ExampleFloat.MaxValue = 100.0
	cfg.ExampleFloat.Value = cfg.ExampleFloat.Default
	cfg.ExampleInt.AffectedContainers = []string{shared.ServiceContainerName}

	// ExampleString
	cfg.ExampleString.ID = hdconfig.Identifier(ids.ExampleStringID)
	cfg.ExampleString.Name = "Example String"
	cfg.ExampleString.Description.Default = "This is an example of a string parameter. It has a max length and regex pattern set."
	cfg.ExampleString.MaxLength = 10
	cfg.ExampleString.Regex = "^[a-zA-Z]*$"
	cfg.ExampleString.Value = cfg.ExampleString.Default
	cfg.ExampleInt.AffectedContainers = []string{shared.ServiceContainerName}

	// Options for ExampleChoice
	options := make([]hdconfig.ParameterMetadataOption[nativecfg.ExampleOption], 3)
	options[0].Name = "One"
	options[0].Description.Default = "This is the first option."
	options[0].Value = nativecfg.ExampleOption_One

	thresholdString := strconv.FormatFloat(FloatThreshold, 'f', -1, 64)
	options[1].Name = "Two"
	options[2].Description.Default = "This is the second option. It is hidden when ExampleFloat is less than " + thresholdString + "."
	options[2].Description.Template = fmt.Sprintf("{{if lt .GetValue %s %s}}This option is hidden because the float is less than %s.{{else}}This option is visible because the float is greater than or equal to %s.{{end}}", ids.ExampleFloatID, thresholdString, thresholdString, thresholdString)
	options[1].Value = nativecfg.ExampleOption_Two
	options[1].Disabled.Default = true
	options[1].Disabled.Template = "{{if eq .GetValue " + ids.ExampleBoolID + " true}}false{{else}}{{.UseDefault}}{{end}}"

	options[2].Name = "Three"
	options[0].Description.Default = "This is the third option."
	options[2].Value = nativecfg.ExampleOption_Three

	// ExampleChoice
	cfg.ExampleChoice.ID = hdconfig.Identifier(ids.ExampleChoiceID)
	cfg.ExampleChoice.Name = "Example Choice"
	cfg.ExampleChoice.Description.Default = "This is an example of a choice parameter between multiple options."
	cfg.ExampleChoice.Options = options
	cfg.ExampleChoice.Default = options[0].Value
	cfg.ExampleChoice.Value = cfg.ExampleChoice.Default

	cfg.Parameters = []hdconfig.IParameterMetadata{
		&cfg.ExampleBool,
		&cfg.ExampleInt,
		&cfg.ExampleUint,
		&cfg.ExampleFloat,
		&cfg.ExampleString,
		&cfg.ExampleChoice,
	}

	// Subconfigs
	cfg.SubConfig = NewSubConfig()
	cfg.ServerConfig = NewServerConfig()
	cfg.Sections = []hdconfig.SectionMetadata{
		cfg.SubConfig.SectionMetadata,
		cfg.ServerConfig.SectionMetadata,
	}

	return cfg
}

type SubConfig struct {
	hdconfig.SectionMetadata

	SubExampleBool hdconfig.BoolParameterMetadata `json:"-" yaml:"-"`

	SubExampleChoice hdconfig.ChoiceParameterMetadata[nativecfg.ExampleOption] `json:"-" yaml:"-"`
}

func NewSubConfig() *SubConfig {
	cfg := &SubConfig{}
	cfg.ID = hdconfig.Identifier(ids.SubConfigID)
	cfg.Name = "Sub Config"
	cfg.Description.Default = "This is a sub-section of the main configuration."
	cfg.Hidden.Default = true
	cfg.Hidden.Template = "{{if eq .GetValue " + ids.ExampleBoolID + " true}}false{{else}}true{{end}}"

	// SubExampleBool
	cfg.SubExampleBool.ID = hdconfig.Identifier(ids.SubExampleBoolID)
	cfg.SubExampleBool.Name = "Sub Example Boolean"
	cfg.SubExampleBool.Description.Default = "This is an example of a boolean parameter in a sub-section."
	cfg.SubExampleBool.Value = cfg.SubExampleBool.Default

	// Options for SubExampleChoice
	options := make([]hdconfig.ParameterMetadataOption[nativecfg.ExampleOption], 2)
	options[0].Name = "One"
	options[0].Description.Default = "This is the first option."
	options[0].Value = nativecfg.ExampleOption_One

	options[1].Name = "Two"
	options[1].Description.Default = "This is the second option."
	options[1].Value = nativecfg.ExampleOption_Two

	// SubExampleChoice
	cfg.SubExampleChoice.ID = hdconfig.Identifier(ids.SubExampleChoiceID)
	cfg.SubExampleChoice.Name = "Sub Example Choice"
	cfg.SubExampleChoice.Description.Default = "This is an example of a choice parameter between multiple options in a sub-section."
	cfg.SubExampleChoice.Options = options
	cfg.SubExampleChoice.Default = options[1].Value
	cfg.SubExampleChoice.Value = cfg.SubExampleChoice.Default

	cfg.Parameters = []hdconfig.IParameterMetadata{
		&cfg.SubExampleBool,
		&cfg.SubExampleChoice,
	}
	return cfg
}

type ServerConfig struct {
	hdconfig.SectionMetadata

	Port hdconfig.UintParameterMetadata `json:"-" yaml:"port"`

	PortMode hdconfig.ChoiceParameterMetadata[PortMode] `json:"-" yaml:"-"`
}

func NewServerConfig() *ServerConfig {
	cfg := &ServerConfig{}
	cfg.ID = hdconfig.Identifier(ids.ServerConfigID)
	cfg.Name = "Service Config"
	cfg.Description.Default = "This is the configuration for the module's service. This isn't used by the service directly, but it is used by Hyperdrive itself in the service's Docker Compose file template to configure the service during its starting process."
	cfg.Sections = []hdconfig.SectionMetadata{}
	cfg.Parameters = []hdconfig.IParameterMetadata{}

	// Port
	cfg.Port.ID = hdconfig.Identifier(ids.PortID)
	cfg.Port.Name = "API Port"
	cfg.Port.Description.Default = "This is the API port the server should run on."
	cfg.Port.Default = uint64(shared.DefaultServerApiPort)
	cfg.Port.MinValue = 0
	cfg.Port.MaxValue = 65535
	cfg.Port.Value = cfg.Port.Default
	cfg.Port.AffectedContainers = []string{shared.ServiceContainerName}

	// Options for PortMode
	options := make([]hdconfig.ParameterMetadataOption[PortMode], 3)
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
	cfg.PortMode.ID = hdconfig.Identifier(ids.PortModeID)
	cfg.PortMode.Name = "Expose API Port"
	cfg.PortMode.Description.Default = "Determine how the server's HTTP API restricts its access from various sources."
	cfg.PortMode.Options = options
	cfg.PortMode.Default = options[0].Value
	cfg.PortMode.Value = cfg.PortMode.Default
	cfg.PortMode.AffectedContainers = []string{shared.ServiceContainerName}

	cfg.Parameters = []hdconfig.IParameterMetadata{
		&cfg.Port,
		&cfg.PortMode,
	}
	return cfg
}

func ConvertToMetadata(native *nativecfg.NativeExampleConfig) *ExampleConfig {
	cfg := NewExampleConfig()

	cfg.ExampleBool.Value = native.ExampleBool
	cfg.ExampleInt.Value = native.ExampleInt
	cfg.ExampleUint.Value = native.ExampleUint
	cfg.ExampleFloat.Value = native.ExampleFloat
	cfg.ExampleString.Value = native.ExampleString
	cfg.ExampleChoice.Value = native.ExampleChoice

	cfg.SubConfig.SubExampleBool.Value = native.SubConfig.SubExampleBool
	cfg.SubConfig.SubExampleChoice.Value = native.SubConfig.SubExampleChoice

	return cfg
}

func (cfg *ExampleConfig) ConvertToNative() *nativecfg.NativeExampleConfig {
	native := &nativecfg.NativeExampleConfig{}

	native.ExampleBool = cfg.ExampleBool.Value
	native.ExampleInt = cfg.ExampleInt.Value
	native.ExampleUint = cfg.ExampleUint.Value
	native.ExampleFloat = cfg.ExampleFloat.Value
	native.ExampleString = cfg.ExampleString.Value
	native.ExampleChoice = cfg.ExampleChoice.Value
	native.SubConfig.SubExampleBool = cfg.SubConfig.SubExampleBool.Value
	native.SubConfig.SubExampleChoice = cfg.SubConfig.SubExampleChoice.Value
	return native
}

func (cfg *ExampleConfig) ConvertToInstance() map[string]any {
	instance := map[string]any{}
	instance[ids.ExampleBoolID] = cfg.ExampleBool.Value
	instance[ids.ExampleIntID] = cfg.ExampleInt.Value
	instance[ids.ExampleUintID] = cfg.ExampleUint.Value
	instance[ids.ExampleFloatID] = cfg.ExampleFloat.Value
	instance[ids.ExampleStringID] = cfg.ExampleString.Value
	instance[ids.ExampleChoiceID] = cfg.ExampleChoice.Value

	subInstance := map[string]any{}
	subInstance[ids.SubExampleBoolID] = cfg.SubConfig.SubExampleBool.Value
	subInstance[ids.SubExampleChoiceID] = cfg.SubConfig.SubExampleChoice.Value
	instance[ids.SubConfigID] = subInstance

	return instance
}

func ConvertFromInstance(instance map[string]any) (*ExampleConfig, error) {
	cfg := NewExampleConfig()

	// Top-level parameters
	var subConfig map[string]any
	var serviceConfig map[string]any
	errs := []error{
		procParam(instance, ids.ExampleBoolID, &cfg.ExampleBool.Value),
		procParam(instance, ids.ExampleIntID, &cfg.ExampleInt.Value),
		procParam(instance, ids.ExampleUintID, &cfg.ExampleUint.Value),
		procParam(instance, ids.ExampleFloatID, &cfg.ExampleFloat.Value),
		procParam(instance, ids.ExampleStringID, &cfg.ExampleString.Value),
		procParam(instance, ids.ExampleChoiceID, &cfg.ExampleChoice.Value),
		procParam(instance, ids.SubConfigID, &subConfig),
		procParam(instance, ids.ServerConfigID, &serviceConfig),
	}
	if err := errors.Join(errs...); err != nil {
		return nil, fmt.Errorf("error processing parameters: %w", err)
	}

	// Sub-config
	errs = []error{
		procParam(subConfig, ids.SubExampleBoolID, &cfg.SubConfig.SubExampleBool.Value),
		procParam(subConfig, ids.SubExampleChoiceID, &cfg.SubConfig.SubExampleChoice.Value),
	}
	if err := errors.Join(errs...); err != nil {
		return nil, fmt.Errorf("error processing sub-config parameters: %w", err)
	}

	// Service config
	errs = []error{
		procParam(serviceConfig, ids.PortID, &cfg.ServerConfig.Port.Value),
		procParam(serviceConfig, ids.PortModeID, &cfg.ServerConfig.PortMode.Value),
	}
	if err := errors.Join(errs...); err != nil {
		return nil, fmt.Errorf("error processing service config parameters: %w", err)
	}
	return cfg, nil
}

func procParam[ParamType any](instance map[string]any, paramID string, store *ParamType) error {
	paramAny, exists := instance[paramID]
	if !exists {
		return errors.New("missing required parameter: " + paramID)
	}
	paramTyped, ok := paramAny.(ParamType)
	if !ok {
		return fmt.Errorf("invalid type for parameter [%s]: %T", ids.ExampleBoolID, paramAny)
	}
	*store = paramTyped
	return nil
}

// Configuration manager
type AdapterConfigManager struct {
	// The adapter configuration
	AdapterConfig *ExampleConfig

	// The native configuration manager
	nativeConfigManager *nativecfg.ConfigManager

	// The path to the adapter configuration file
	adapterConfigPath string
}

// Create a new configuration manager for the adapter
func NewAdapterConfigManager() *AdapterConfigManager {
	return &AdapterConfigManager{
		nativeConfigManager: nativecfg.NewConfigManager(utils.ServiceConfigPath),
		adapterConfigPath:   utils.AdapterConfigPath,
	}
}

// Load the configuration from disk
func (m *AdapterConfigManager) LoadConfigFromDisk() (*ExampleConfig, error) {
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
	serverCfg := ServerConfig{}
	err = yaml.Unmarshal(bytes, &serverCfg)
	if err != nil {
		return nil, fmt.Errorf("error deserializing adapter config file [%s]: %w", m.adapterConfigPath, err)
	}

	// Merge the configs
	modCfg := ConvertToMetadata(nativeCfg)
	modCfg.ServerConfig = &serverCfg
	m.AdapterConfig = modCfg
	return modCfg, nil
}

// Save the configuration to a file. If the config hasn't been loaded yet, this doesn't do anything.
func (m *AdapterConfigManager) SaveConfigToDisk() error {
	if m.AdapterConfig == nil {
		return nil
	}

	// Save the native config
	nativeCfg := m.AdapterConfig.ConvertToNative()
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
