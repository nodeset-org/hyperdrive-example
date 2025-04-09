package config

import (
	"fmt"
	"strconv"

	"github.com/nodeset-org/hyperdrive-example/adapter/config/ids"
	"github.com/nodeset-org/hyperdrive-example/shared"
	nativecfg "github.com/nodeset-org/hyperdrive-example/shared/config"
	hdconfig "github.com/nodeset-org/hyperdrive/modules/config"
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

	ExampleFloat hdconfig.FloatParameter

	ExampleString hdconfig.StringParameter

	ExampleChoice hdconfig.ChoiceParameter[nativecfg.ExampleOption]

	ExampleDerivedValue hdconfig.StringParameter

	SubConfig *SubConfig

	ServerConfig *ServerConfig
}

type ExampleConfigSettings struct {
	ExampleBool   bool                    `json:"exampleBool"`
	ExampleInt    int64                   `json:"exampleInt"`
	ExampleFloat  float64                 `json:"exampleFloat"`
	ExampleString string                  `json:"exampleString"`
	ExampleChoice nativecfg.ExampleOption `json:"exampleChoice"`

	ExampleDerivedValue string                `json:"exampleDerivedValue"`
	SubConfig           *SubConfigSettings    `json:"subConfig"`
	ServerConfig        *ServerConfigSettings `json:"server" yaml:"server"`
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

	var thresholdString string
	if FloatThreshold == float64(int(FloatThreshold)) {
		// Needed to preserve trailing zero so the template parser knows this is a float
		thresholdString = strconv.FormatFloat(FloatThreshold, 'f', 1, 64)
	} else {
		thresholdString = strconv.FormatFloat(FloatThreshold, 'f', -1, 64)
	}
	options[1].Name = "Two"
	options[1].Description.Default = fmt.Sprintf("This is the second option. It is hidden when %s is less than %s.", cfg.ExampleFloat.Name, thresholdString)
	options[1].Description.Template = fmt.Sprintf(`{{if lt (.GetValue "%s") %s}}This option is hidden because %s is less than %s.{{else}}This option is visible because the float is greater than or equal to %s.{{end}}`, ids.ExampleFloatID, thresholdString, cfg.ExampleFloat.Name, thresholdString, thresholdString)
	options[1].Value = nativecfg.ExampleOption_Two
	options[1].Hidden.Default = false
	options[1].Hidden.Template = fmt.Sprintf(`{{if lt (.GetValue "%s") %s}}true{{else}}{{.UseDefault}}{{end}}`, ids.ExampleFloatID, thresholdString)

	options[2].Name = "Three"
	options[2].Description.Default = "This is the third option."
	options[2].Value = nativecfg.ExampleOption_Three

	// ExampleChoice
	cfg.ExampleChoice.ID = ids.ExampleChoiceID
	cfg.ExampleChoice.Name = "Example Choice"
	cfg.ExampleChoice.Description.Default = fmt.Sprintf("This is an example of a choice parameter between multiple options. The second option is hidden when %s is less than "+thresholdString+".", cfg.ExampleFloat.Name)
	cfg.ExampleChoice.Options = options
	cfg.ExampleChoice.Default = options[0].Value
	cfg.ExampleChoice.AffectedContainers = []string{}

	// ExampleDerivedValue
	cfg.ExampleDerivedValue.ID = ids.ExampleDerivedValueID
	cfg.ExampleDerivedValue.Hidden.Default = true
	cfg.ExampleDerivedValue.Name = "Sub Derived Value"
	cfg.ExampleDerivedValue.Description.Default = "This is a derived value based on the sub-section's parameters."
	cfg.ExampleDerivedValue.Default = ""

	// Subconfigs
	cfg.SubConfig = NewSubConfig()
	cfg.ServerConfig = NewServerConfig()

	return cfg
}

func (cfg ExampleConfig) GetParameters() []hdconfig.IParameter {
	return []hdconfig.IParameter{
		&cfg.ExampleBool,
		&cfg.ExampleInt,
		&cfg.ExampleFloat,
		&cfg.ExampleString,
		&cfg.ExampleChoice,
		&cfg.ExampleDerivedValue,
	}
}

func (cfg ExampleConfig) GetSections() []hdconfig.ISection {
	return []hdconfig.ISection{
		cfg.SubConfig,
		cfg.ServerConfig,
	}
}

// GetChangedServices returns a list of services that would be affected by the new settings
func (s *ExampleConfigSettings) GetChangedServices(oldSettings *ExampleConfigSettings) ([]string, error) {
	cfg := NewExampleConfig()
	newModSettings := hdconfig.CreateModuleSettings(cfg)
	err := newModSettings.CopySettingsFromKnownType(s)
	if err != nil {
		return nil, fmt.Errorf("error copying new settings: %w", err)
	}

	oldModSettings := hdconfig.CreateModuleSettings(cfg)
	err = oldModSettings.CopySettingsFromKnownType(oldSettings)
	if err != nil {
		return nil, fmt.Errorf("error copying old settings: %w", err)
	}

	// Compare the settings - if there are no differences, return nil
	diff := hdconfig.CompareSettings(cfg, oldModSettings, newModSettings)
	if len(diff.ParameterDifferences) == 0 && len(diff.SectionDifferences) == 0 {
		return nil, nil
	}

	// Any parameter changes will affect the service, so just return it
	changedServices := []string{
		shared.ServiceContainerName,
	}
	return changedServices, nil
}

func CreateInstanceFromNativeConfig(native *nativecfg.NativeExampleConfig) *ExampleConfigSettings {
	instance := &ExampleConfigSettings{
		ExampleBool:         native.ExampleBool,
		ExampleInt:          native.ExampleInt,
		ExampleFloat:        native.ExampleFloat,
		ExampleString:       native.ExampleString,
		ExampleChoice:       native.ExampleChoice,
		ExampleDerivedValue: native.ExampleDerivedValue,
		SubConfig: &SubConfigSettings{
			SubExampleBool:   native.SubConfig.SubExampleBool,
			SubExampleChoice: native.SubConfig.SubExampleChoice,
		},
		ServerConfig: &ServerConfigSettings{},
	}
	return instance
}

func ConvertInstanceToNativeConfig(instance *ExampleConfigSettings) *nativecfg.NativeExampleConfig {
	native := &nativecfg.NativeExampleConfig{
		ExampleBool:         instance.ExampleBool,
		ExampleInt:          instance.ExampleInt,
		ExampleFloat:        instance.ExampleFloat,
		ExampleString:       instance.ExampleString,
		ExampleChoice:       instance.ExampleChoice,
		ExampleDerivedValue: instance.ExampleDerivedValue,
		SubConfig: nativecfg.NativeSubConfig{
			SubExampleBool:   instance.SubConfig.SubExampleBool,
			SubExampleChoice: instance.SubConfig.SubExampleChoice,
		},
	}
	return native
}
