package hdmodule

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nodeset-org/hyperdrive-example/adapter/utils"
	"github.com/nodeset-org/hyperdrive-example/shared"
	hdconfig "github.com/nodeset-org/hyperdrive/config"
	"github.com/urfave/cli/v2"
)

// Request format for `start`
type startRequest struct {
	utils.KeyedRequest

	// The config settings being used
	Settings *hdconfig.HyperdriveSettings `json:"settings"`
}

// Handle the `start` command
func startServices(c *cli.Context) error {
	// Make sure the env vars are set
	if utils.ComposeDir == "" {
		return fmt.Errorf("%s not set", utils.ComposeDirEnvVarName)
	}
	if utils.ComposeProject == "" {
		return fmt.Errorf("%s not set", utils.ComposeProjectEnvVarName)
	}

	// Get the request
	_, err := utils.HandleKeyedRequest[*startRequest](c)
	if err != nil {
		return err
	}

	// Start the services - the example doesn't need to do anything interesting with the provided settings, so they are ignored
	args := []string{
		"compose",
		"-p",
		utils.ComposeProject,
		"-f",
		filepath.Join(utils.ComposeDir, shared.ServiceContainerName+".yml"),
		"up",
		"-d",
		"--quiet-pull",
	}
	cmd := exec.Command("docker", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error starting services: %w", err)
	}

	return nil
}
