package hdmodule

import (
	"fmt"

	"github.com/goccy/go-json"

	"github.com/nodeset-org/hyperdrive-example/shared"
)

// Response format for `version`
type versionResponse struct {
	// Version of the module
	Version string `json:"version"`
}

// Handle the `version` command
func version() error {
	// Create the response
	version := versionResponse{
		Version: shared.Version,
	}

	// Marshal it
	bytes, err := json.Marshal(version)
	if err != nil {
		return fmt.Errorf("error marshalling version response: %w", err)
	}

	// Print it
	fmt.Println(string(bytes))
	return nil
}
