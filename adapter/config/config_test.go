package config

import (
	"testing"

	"github.com/goccy/go-json"
	"github.com/nodeset-org/hyperdrive-example/shared"
	"github.com/stretchr/testify/require"
)

const (
	oldSettingsJson string = `{"exampleBool":false,"exampleChoice":"one","exampleFloat":50,"exampleInt":0,"exampleString":"","exampleUint":42,"server":{"port":8080,"portMode":"closed"},"subConfig":{"subConfigBool":false,"subConfigChoice":"two"}}`

	newSettingsJson string = `{"exampleBool":false,"exampleChoice":"one","exampleFloat":80,"exampleInt":0,"exampleString":"","exampleUint":42,"server":{"port":8085,"portMode":"open"},"subConfig":{"subConfigBool":false,"subConfigChoice":"two"}}`
)

func TestChangedServices(t *testing.T) {
	oldSettings := new(ExampleConfigSettings)
	err := json.Unmarshal([]byte(oldSettingsJson), oldSettings)
	require.NoError(t, err)

	newSettings := new(ExampleConfigSettings)
	err = json.Unmarshal([]byte(newSettingsJson), newSettings)
	require.NoError(t, err)

	// Compare the settings
	services, err := newSettings.GetChangedServices(oldSettings)
	require.NoError(t, err)

	// Check the services
	require.Len(t, services, 1)
	require.Equal(t, shared.ServiceContainerName, services[0])
}
