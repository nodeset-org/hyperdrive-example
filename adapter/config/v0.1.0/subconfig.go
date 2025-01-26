package config

import (
	"github.com/nodeset-org/hyperdrive-example/adapter/config/ids"
	nativecfg "github.com/nodeset-org/hyperdrive-example/shared/config"
	hdconfig "github.com/nodeset-org/hyperdrive/modules/config"
)

type SubConfig struct {
	hdconfig.SectionHeader

	SubExampleBool hdconfig.BoolParameter

	SubExampleChoice hdconfig.ChoiceParameter[nativecfg.ExampleOption]
}

type SubConfigSettings struct {
	SubExampleBool   bool                    `json:"subConfigBool"`
	SubExampleChoice nativecfg.ExampleOption `json:"subConfigChoice"`
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

func (cfg SubConfig) GetParameters() []hdconfig.IParameter {
	return []hdconfig.IParameter{
		&cfg.SubExampleBool,
		&cfg.SubExampleChoice,
	}
}

func (cfg SubConfig) GetSections() []hdconfig.ISection {
	return []hdconfig.ISection{}
}
