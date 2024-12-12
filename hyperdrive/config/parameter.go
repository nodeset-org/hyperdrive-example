package config

import "fmt"

// Parameter types
type ParameterType string

const (
	ParameterType_Bool   ParameterType = "bool"
	ParameterType_Int    ParameterType = "int"
	ParameterType_Uint   ParameterType = "uint"
	ParameterType_Float  ParameterType = "float"
	ParameterType_String ParameterType = "string"
	ParameterType_Choice ParameterType = "choice"
)

// Common interface for all Parameter metadata structs
type IParameterMetadata interface {
	GetID() Identifier
	GetName() string
	GetDescription() DynamicProperty[string]
	GetType() ParameterType
	GetDefaultAsAny() any
	GetValueAsAny() any
	GetAdvanced() bool
	GetHidden() DynamicProperty[bool]
	GetOverwriteOnUpgrade() bool
	GetAffectedContainers() []string
	Deserialize(data map[string]any) error
}

// ===========================
/// === Parameter Metadata ===
// ===========================

const (
	// Field names
	IDKey                 string = "id"
	NameKey               string = "name"
	DescriptionKey        string = "description"
	TypeKey               string = "type"
	DefaultKey            string = "default"
	ValueKey              string = "value"
	AdvancedKey           string = "advanced"
	DisabledKey           string = "disabled"
	HiddenKey             string = "hidden"
	OverwriteOnUpgradeKey string = "overwriteOnUpgrade"
	AffectedContainersKey string = "affectedContainers"
)

// Parameter metadata implementation according to the spec
type ParameterMetadata[Type any] struct {
	// Unique ID for referencing the parameter behind-the-scenes
	ID Identifier `json:"id" yaml:"id"`

	// Human-readable name for the parameter
	Name string `json:"name" yaml:"name"`

	// Description of the parameter
	Description DynamicProperty[string] `json:"description" yaml:"description"`

	// Default value for the parameter
	Default Type `json:"default" yaml:"default"`

	// Current value assigned to the parameter, if configured
	Value Type `json:"value" yaml:"value"`

	// Flag for hiding the parameter behind the "advanced mode" toggle
	Advanced bool `json:"advanced,omitempty" yaml:"advanced,omitempty"`

	// Flag for disabling the parameter in the UI, graying it out
	Disabled DynamicProperty[bool] `json:"disabled,omitempty" yaml:"disabled,omitempty"`

	// Dynamic flag for hiding the parameter from the UI
	Hidden DynamicProperty[bool] `json:"hidden,omitempty" yaml:"hidden,omitempty"`

	// Flag for overwriting the value with the default on an upgrade
	OverwriteOnUpgrade bool `json:"overwriteOnUpgrade" yaml:"overwriteOnUpgrade"`

	// List of containers affected if this parameter is changed
	AffectedContainers []string `json:"affectedContainers" yaml:"affectedContainers"`
}

func (p *ParameterMetadata[Type]) GetID() Identifier {
	return p.ID
}

func (p *ParameterMetadata[Type]) GetName() string {
	return p.Name
}

func (p *ParameterMetadata[Type]) GetDescription() DynamicProperty[string] {
	return p.Description
}

func (p *ParameterMetadata[Type]) GetDefaultAsAny() any {
	return p.Default
}

func (p *ParameterMetadata[Type]) GetValueAsAny() any {
	return p.Value
}

func (p *ParameterMetadata[Type]) GetAdvanced() bool {
	return p.Advanced
}

func (p *ParameterMetadata[Type]) GetHidden() DynamicProperty[bool] {
	return p.Hidden
}

func (p *ParameterMetadata[Type]) GetOverwriteOnUpgrade() bool {
	return p.OverwriteOnUpgrade
}

func (p *ParameterMetadata[Type]) GetAffectedContainers() []string {
	return p.AffectedContainers
}

// DeserializeImpl the parameter metadata from a map
func (p *ParameterMetadata[Type]) DeserializeImpl(data map[string]any) error {
	// Get the ID
	err := DeserializeIdentifier(data, IDKey, &p.ID, false)
	if err != nil {
		return err
	}

	// Get the name
	_, err = DeserializeProperty(data, NameKey, &p.Name, false)
	if err != nil {
		return err
	}

	// Get the description
	_, err = DeserializeDynamicProperty(data, DescriptionKey, &p.Description, false)
	if err != nil {
		return err
	}

	// Get the default value
	_, err = DeserializeProperty(data, DefaultKey, &p.Default, false)
	if err != nil {
		return err
	}

	// Get the current value
	_, err = DeserializeProperty(data, ValueKey, &p.Value, false)
	if err != nil {
		return err
	}

	// Get the advanced flag
	_, err = DeserializeProperty(data, AdvancedKey, &p.Advanced, true)
	if err != nil {
		return err
	}

	// Get the disabled flag
	_, err = DeserializeDynamicProperty(data, DisabledKey, &p.Disabled, true)
	if err != nil {
		return err
	}

	// Get the hidden flag
	_, err = DeserializeDynamicProperty(data, HiddenKey, &p.Hidden, true)
	if err != nil {
		return err
	}

	// Get the overwriteOnUpgrade flag
	_, err = DeserializeProperty(data, OverwriteOnUpgradeKey, &p.OverwriteOnUpgrade, false)
	if err != nil {
		return err
	}

	// Get the affectedContainers list
	_, err = DeserializeProperty(data, AffectedContainersKey, &p.AffectedContainers, false)
	if err != nil {
		return err
	}
	return nil
}

/// =======================
/// === Bool Parameters ===
/// =======================

// A boolean parameter's metadata
type BoolParameterMetadata struct {
	ParameterMetadata[bool]
}

func (p *BoolParameterMetadata) GetType() ParameterType {
	return ParameterType_Bool
}

func (p *BoolParameterMetadata) Deserialize(data map[string]any) error {
	return p.ParameterMetadata.DeserializeImpl(data)
}

/// =======================================
/// === Prototype for Number Parameters ===
/// =======================================

const (
	// Field names
	MinValueKey string = "minValue"
	MaxValueKey string = "maxValue"
)

type NumberParameterType interface {
	int64 | uint64 | float64
}

// An integer parameter's metadata
type NumberParameterMetadata[Type NumberParameterType] struct {
	ParameterMetadata[Type]

	// Minimum value for the parameter
	MinValue Type `json:"minValue,omitempty" yaml:"minValue,omitempty"`

	// Maximum value for the parameter
	MaxValue Type `json:"maxValue,omitempty" yaml:"maxValue,omitempty"`
}

func (p *NumberParameterMetadata[Type]) Deserialize(data map[string]any) error {
	err := p.ParameterMetadata.DeserializeImpl(data)
	if err != nil {
		return err
	}

	// Get the min value
	_, err = DeserializeProperty(data, MinValueKey, &p.MinValue, true)
	if err != nil {
		return err
	}

	// Get the max value
	_, err = DeserializeProperty(data, MaxValueKey, &p.MaxValue, true)
	if err != nil {
		return err
	}
	return nil
}

/// ======================
/// === Int Parameters ===
/// ======================

// An integer parameter's metadata
type IntParameterMetadata struct {
	NumberParameterMetadata[int64]
}

func (p *IntParameterMetadata) GetType() ParameterType {
	return ParameterType_Int
}

/// =======================
/// === Uint Parameters ===
/// =======================

// An unsigned integer parameter's metadata
type UintParameterMetadata struct {
	NumberParameterMetadata[uint64]
}

func (p *UintParameterMetadata) GetType() ParameterType {
	return ParameterType_Uint
}

/// ========================
/// === Float Parameters ===
/// ========================

// A float parameter's metadata
type FloatParameterMetadata struct {
	NumberParameterMetadata[float64]
}

func (p *FloatParameterMetadata) GetType() ParameterType {
	return ParameterType_Float
}

/// =========================
/// === String Parameters ===
/// =========================

const (
	// Field names
	MaxLengthKey string = "maxLength"
	RegexKey     string = "regex"
)

// A string parameter's metadata
type StringParameterMetadata struct {
	ParameterMetadata[string]

	// The max length of the string
	MaxLength uint `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`

	// The pattern for the regular expression the string must match
	Regex string `json:"regex,omitempty" yaml:"regex,omitempty"`
}

func (p *StringParameterMetadata) GetType() ParameterType {
	return ParameterType_String
}

func (p *StringParameterMetadata) Deserialize(data map[string]any) error {
	err := p.ParameterMetadata.DeserializeImpl(data)
	if err != nil {
		return err
	}

	// Get the max length
	_, err = DeserializeProperty(data, MaxLengthKey, &p.MaxLength, true)
	if err != nil {
		return err
	}

	// Get the regex pattern
	_, err = DeserializeProperty(data, RegexKey, &p.Regex, true)
	if err != nil {
		return err
	}
	return nil
}

/// =========================
/// === Choice Parameters ===
/// =========================

type ChoiceParameterMetadata[ChoiceType any] struct {
	ParameterMetadata[ChoiceType]

	// The choices available for the parameter
	Options []ParameterMetadataOption[ChoiceType] `json:"options" yaml:"options"`
}

func (p *ChoiceParameterMetadata[ChoiceType]) GetType() ParameterType {
	return ParameterType_Choice
}

// Unmarshal the choice parameter from a map
func (p *ChoiceParameterMetadata[ChoiceType]) Deserialize(data map[string]any) error {
	err := p.ParameterMetadata.DeserializeImpl(data)
	if err != nil {
		return err
	}

	// Get the options
	var options []map[string]any
	_, err = DeserializeProperty(data, "options", &options, false)
	if err != nil {
		return err
	}

	// Unmarshal the options
	for _, optionData := range options {
		option := ParameterMetadataOption[ChoiceType]{}
		err = option.Unmarshal(optionData)
		if err != nil {
			return err
		}
		p.Options = append(p.Options, option)
	}
	return nil
}

/// =========================
/// === Parameter Options ===
/// =========================

// A single option for a choice parameter
type ParameterMetadataOption[ChoiceType any] struct {
	// The option's name
	Name string `json:"name" yaml:"name"`

	// The description for the option
	Description DynamicProperty[string] `json:"description" yaml:"description"`

	// The value for the option
	Value ChoiceType `json:"value" yaml:"value"`

	// Flag for disabling the option in the UI, graying it out
	Disabled DynamicProperty[bool] `json:"disabled,omitempty" yaml:"disabled,omitempty"`

	// Flag for hiding the option from the UI
	Hidden DynamicProperty[bool] `json:"hidden,omitempty" yaml:"hidden,omitempty"`
}

// Unmarshal the option from a map
func (o *ParameterMetadataOption[ChoiceType]) Unmarshal(data map[string]any) error {
	// Get the name
	_, err := DeserializeProperty(data, NameKey, &o.Name, false)
	if err != nil {
		return err
	}

	// Get the description
	_, err = DeserializeDynamicProperty(data, DescriptionKey, &o.Description, false)
	if err != nil {
		return err
	}

	// Get the value
	_, err = DeserializeProperty(data, ValueKey, &o.Value, false)
	if err != nil {
		return err
	}

	// Get the disabled flag
	_, err = DeserializeDynamicProperty(data, DisabledKey, &o.Disabled, true)
	if err != nil {
		return err
	}

	// Get the hidden flag
	_, err = DeserializeDynamicProperty(data, HiddenKey, &o.Hidden, true)
	if err != nil {
		return err
	}
	return nil
}

/// =========================
/// === Utility Functions ===
/// =========================

func DeserializeParameterMetadata(serializedParam map[string]any) (IParameterMetadata, error) {
	// Get the type
	var paramType string
	_, err := DeserializeProperty(serializedParam, TypeKey, &paramType, false)
	if err != nil {
		return nil, err
	}

	// Create the parameter based on the type
	var param IParameterMetadata
	switch ParameterType(paramType) {
	case ParameterType_Bool:
		param = &BoolParameterMetadata{}
	case ParameterType_Int:
		param = &IntParameterMetadata{}
	case ParameterType_Uint:
		param = &UintParameterMetadata{}
	case ParameterType_Float:
		param = &FloatParameterMetadata{}
	case ParameterType_String:
		param = &StringParameterMetadata{}
	case ParameterType_Choice:
		param = &ChoiceParameterMetadata[string]{}
	default:
		return nil, fmt.Errorf("invalid parameter type: %s", paramType)
	}

	// Deserialize the parameter
	err = param.Deserialize(serializedParam)
	if err != nil {
		return nil, err
	}
	return param, nil
}