package hdmodule

import (
	"testing"

	"github.com/goccy/go-json"
	"github.com/nodeset-org/hyperdrive-example/shared"
	hdconfig "github.com/nodeset-org/hyperdrive/config"
	"github.com/stretchr/testify/require"
)

const (
	oldSettingsJson string = `{"version":"","projectName":"hde-test","apiPort":8080,"enableIPv6":false,"userDataPath":"/tmp/hde-adapter-test/data","additionalDockerNetworks":"","clientTimeout":30,"containerTag":"nodeset/hyperdrive:v2.0.0-dev","logging":{"level":"info","format":"logfmt","addSource":false,"maxSize":20,"maxBackups":3,"maxAge":90,"localTime":false,"compress":true},"modules":{"NodeSet/example-module":{"enabled":true,"version":"0.2.0","settings":{"exampleBool":false,"exampleChoice":"one","exampleDerivedValue":1337,"exampleFloat":50,"exampleInt":0,"exampleString":"","exampleUint":42,"server":{"port":8080,"portMode":"closed"},"subConfig":{"subConfigBool":false,"subConfigChoice":"two"}}}}}`

	newSettingsJson string = `{"version":"","projectName":"hde-test","apiPort":8080,"enableIPv6":false,"userDataPath":"/tmp/hde-adapter-test/data","additionalDockerNetworks":"","clientTimeout":10,"containerTag":"nodeset/hyperdrive:v2.0.0-dev","logging":{"level":"info","format":"logfmt","addSource":false,"maxSize":20,"maxBackups":3,"maxAge":90,"localTime":false,"compress":true},"modules":{"NodeSet/example-module":{"enabled":true,"version":"0.2.0","settings":{"exampleBool":false,"exampleChoice":"one","exampleDerivedValue":1337,"exampleFloat":80,"exampleInt":0,"exampleString":"","exampleUint":42,"server":{"port":8085,"portMode":"open"},"subConfig":{"subConfigBool":false,"subConfigChoice":"two"}}}}}`
)

func TestProcessSettings(t *testing.T) {
	oldSettings := new(hdconfig.HyperdriveSettings)
	err := json.Unmarshal([]byte(oldSettingsJson), oldSettings)
	require.NoError(t, err)

	newSettings := new(hdconfig.HyperdriveSettings)
	err = json.Unmarshal([]byte(newSettingsJson), newSettings)
	require.NoError(t, err)

	// Process the settings
	response, err := processSettingsImpl(oldSettings, newSettings)
	require.NoError(t, err)

	// Check the response
	require.Len(t, response.Errors, 0)
	require.Len(t, response.Ports, 1)
	require.Equal(t, uint16(8085), response.Ports["server/port"])
	require.Len(t, response.ServicesToRestart, 1)
	require.Equal(t, shared.ServiceContainerName, response.ServicesToRestart[0])
	t.Log("Settings processed properly")
}
