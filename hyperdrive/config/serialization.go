package config

import "fmt"

// Deserialize a named property from a map. Assumes the value deserialized by the underlying JSON unmarshaller
// converts the property to the right type.
func DeserializeProperty[Type any](data map[string]any, propertyName string, property *Type, optional bool) (bool, error) {
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
