package hdmodule

import (
	"bytes"
	"os"
	"testing"

	"github.com/nodeset-org/hyperdrive-example/adapter/app"
	hdtemplate "github.com/nodeset-org/hyperdrive/shared/templates"
	"github.com/urfave/cli/v2"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/require"
)

const ExampleSettingsJson = `{
	"modules": {
		"NodeSet/example-module": {
			"enabled": true,
			"version": "0.2.0",
			"settings": {
				"exampleUint": 42,
				"server": {
					"port": 8080,
					"portMode": "closed"
				},
				"subConfig": {
					"subConfigBool": true,
					"subConfigChoice": "one"
				}
			}
		}
	}
}`

func TestCallConfigFunction(t *testing.T) {
	// Build the request
	request := CallConfigFunctionRequest{
		FuncName: "GetDerivedValued",
	}
	err := json.Unmarshal([]byte(ExampleSettingsJson), &request.Settings)
	require.NoError(t, err)

	// Marshal the request to JSON
	inputBytes, err := json.Marshal(request)
	require.NoError(t, err)

	// Add newline to simulate ENTER key
	inputBytes = append(inputBytes, '\n')

	// Redirect stdout to a new pipe
	oldStdout := os.Stdout
	rOut, wOut, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = wOut
	// Change it back at the end
	defer func() {
		os.Stdout = oldStdout
	}()

	// Run the function
	app := app.CreateApp()
	app.Reader = bytes.NewReader(inputBytes)

	ctx := cli.NewContext(app, nil, nil)
	err = callConfigFunction(ctx)
	require.NoError(t, err)

	// Read from stdout
	wOut.Close()
	outputBuf := new(bytes.Buffer)
	_, err = outputBuf.ReadFrom(rOut)
	require.NoError(t, err)

	var result hdtemplate.CallConfigFunctionResponse
	err = json.Unmarshal(outputBuf.Bytes(), &result)
	require.NoError(t, err)

	require.Equal(t, "1337", result.Result)
	t.Log("call-config-function ran successfully")
}
