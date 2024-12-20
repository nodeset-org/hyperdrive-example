package config

import (
	"fmt"
)

// Deserialize a named property from a map. Assumes the value deserialized by the underlying JSON unmarshaller
// converts the property to the right type.
func deserializeProperty[Type any](data map[string]any, propertyName string, property *Type, optional bool) (bool, error) {
	// Get the property by its name
	value, exists := data[propertyName]
	if !exists {
		if optional {
			return false, nil
		}
		return false, fmt.Errorf("missing property %s", propertyName)
	}

	// Convert it to the right type
	propertyTyped, ok := value.(Type)
	if !ok {
		return true, fmt.Errorf("invalid property %s [%v]", propertyName, value)
	}

	// Set the property
	*property = propertyTyped
	return true, nil
}

//TODO: Prob move this to the main repo and need a way to do this for individual subconfigs (like the server cfg so it can be serialized with just the server stuff)
/*
func ConvertFromInstance(instance map[string]any) (*ExampleConfig, error) {
	cfg := NewExampleConfig()

	// Top-level parameters
	var subConfig map[string]any
	var serviceConfig map[string]any
	errs := []error{
		procParam(instance, ids.ExampleBoolID, &cfg.ExampleBool.Value),
		parseNumberParam(instance, ids.ExampleIntID, &cfg.ExampleInt.NumberParameterMetadata),
		parseNumberParam(instance, ids.ExampleUintID, &cfg.ExampleUint.NumberParameterMetadata),
		parseNumberParam(instance, ids.ExampleFloatID, &cfg.ExampleFloat.NumberParameterMetadata),
		procParam(instance, ids.ExampleStringID, &cfg.ExampleString.Value),
		parseChoiceParam(instance, ids.ExampleChoiceID, &cfg.ExampleChoice),
		procParam(instance, ids.SubConfigID, &subConfig),
		procParam(instance, ids.ServerConfigID, &serviceConfig),
	}
	if err := errors.Join(errs...); err != nil {
		return nil, fmt.Errorf("error processing parameters: %w", err)
	}

	// Sub-config
	errs = []error{
		procParam(subConfig, ids.SubExampleBoolID, &cfg.SubConfig.SubExampleBool.Value),
		parseChoiceParam(subConfig, ids.SubExampleChoiceID, &cfg.SubConfig.SubExampleChoice),
	}
	if err := errors.Join(errs...); err != nil {
		return nil, fmt.Errorf("error processing sub-config parameters: %w", err)
	}

	// Service config
	errs = []error{
		parseNumberParam(serviceConfig, ids.PortID, &cfg.ServerConfig.Port.NumberParameterMetadata),
		parseChoiceParam(serviceConfig, ids.PortModeID, &cfg.ServerConfig.PortMode),
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
		return fmt.Errorf("invalid type for parameter [%s]: %T", paramID, paramAny)
	}
	*store = paramTyped
	return nil
}

func parseNumberParam[ParamType NumberParameterType](instance map[string]any, paramID string, param *NumberParameterMetadata[ParamType]) error {
	paramAny, exists := instance[paramID]
	if !exists {
		return errors.New("missing required parameter: " + paramID)
	}
	paramFloat, ok := paramAny.(float64)
	if !ok {
		return fmt.Errorf("invalid type for parameter [%s]: %T", paramID, paramAny)
	}
	param.Value = ParamType(paramFloat)
	return nil
}

func parseChoiceParam[ChoiceType ~string](instance map[string]any, paramID string, param *ChoiceParameterMetadata[ChoiceType]) error {
	paramAny, exists := instance[paramID]
	if !exists {
		return errors.New("missing required parameter: " + paramID)
	}
	paramString, ok := paramAny.(string)
	if !ok {
		return fmt.Errorf("invalid type for parameter [%s]: %T", paramID, paramAny)
	}
	param.Value = ChoiceType(paramString)
	return nil
}
*/
