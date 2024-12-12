package config

const (
	// Field names
	ParametersKey string = "parameters"
	SectionsKey   string = "sections"
)

// SectionMetadata represents a section in a configuration metadata
type SectionMetadata struct {
	// Unique ID for referencing the section behind-the-scenes
	ID Identifier `json:"id" yaml:"id"`

	// Name of the section
	Name string `json:"name" yaml:"name"`

	// The description for the section
	Description DynamicProperty[string] `json:"description" yaml:"description"`

	// List of parameters in the section
	Parameters []IParameterMetadata `json:"parameters" yaml:"parameters"`

	// List of sections in the section
	Sections []SectionMetadata `json:"sections" yaml:"sections"`

	// Flag for disabling the section in the UI, graying it out
	Disabled DynamicProperty[bool] `json:"disabled,omitempty" yaml:"disabled,omitempty"`

	// Flag for hiding the section from the UI
	Hidden DynamicProperty[bool] `json:"hidden,omitempty" yaml:"hidden,omitempty"`
}

// Deserialize a section from a map
func DeserializeSectionMetadata(data map[string]any) (SectionMetadata, error) {
	section := SectionMetadata{}

	// Get the ID
	err := DeserializeIdentifier(data, IDKey, &section.ID, false)
	if err != nil {
		return section, err
	}

	// Get the name
	_, err = DeserializeProperty(data, NameKey, &section.Name, false)
	if err != nil {
		return section, err
	}

	// Get the description
	_, err = DeserializeDynamicProperty(data, DescriptionKey, &section.Description, false)
	if err != nil {
		return section, err
	}

	// Get the disabled flag
	_, err = DeserializeDynamicProperty(data, DisabledKey, &section.Disabled, true)
	if err != nil {
		return section, err
	}

	// Get the hidden flag
	_, err = DeserializeDynamicProperty(data, HiddenKey, &section.Hidden, true)
	if err != nil {
		return section, err
	}

	// Handle the parameters
	var parameters []map[string]any
	_, err = DeserializeProperty(data, ParametersKey, &parameters, false)
	if err != nil {
		return section, err
	}
	for _, parameterData := range parameters {
		parameter, err := DeserializeParameterMetadata(parameterData)
		if err != nil {
			return section, err
		}
		section.Parameters = append(section.Parameters, parameter)
	}

	// Handle subsections
	var subsections []map[string]any
	_, err = DeserializeProperty(data, SectionsKey, &subsections, false)
	if err != nil {
		return section, err
	}
	for _, subsectionData := range subsections {
		subsection, err := DeserializeSectionMetadata(subsectionData)
		if err != nil {
			return section, err
		}
		section.Sections = append(section.Sections, subsection)
	}

	return section, nil
}
