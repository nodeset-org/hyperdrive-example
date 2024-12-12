package config

// Top-level object of a module configuration
type ConfigurationMetadata struct {
	// List of parameters in the configuration
	Parameters []IParameterMetadata `json:"parameters" yaml:"parameters"`

	// List of sections in the configuration
	Sections []SectionMetadata `json:"sections" yaml:"sections"`
}

// Deserialize a configuration from a map
func DeserializeConfigurationMetadata(data map[string]any) (ConfigurationMetadata, error) {
	configuration := ConfigurationMetadata{}

	// Handle the parameters
	var parameters []map[string]any
	_, err := DeserializeProperty(data, ParametersKey, &parameters, false)
	if err != nil {
		return configuration, err
	}
	for _, parameterData := range parameters {
		parameter, err := DeserializeParameterMetadata(parameterData)
		if err != nil {
			return configuration, err
		}
		configuration.Parameters = append(configuration.Parameters, parameter)
	}

	// Handle subsections
	var subsections []map[string]any
	_, err = DeserializeProperty(data, SectionsKey, &subsections, false)
	if err != nil {
		return configuration, err
	}
	for _, subsectionData := range subsections {
		subsection, err := DeserializeSectionMetadata(subsectionData)
		if err != nil {
			return configuration, err
		}
		configuration.Sections = append(configuration.Sections, subsection)
	}

	return configuration, nil
}
