package hdmodule

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nodeset-org/hyperdrive-example/adapter/utils"
	"github.com/nodeset-org/hyperdrive-example/shared"
	"github.com/urfave/cli/v2"
)

// Request format for `stop`
type stopRequest struct {
	utils.KeyedRequest

	// The services to stop. If empty, all services will be stopped.
	Services []string `json:"services"`
}

// Handle the `stop` command
func stopServices(c *cli.Context) error {
	// Make sure the env vars are set
	if utils.ComposeDir == "" {
		return fmt.Errorf("%s not set", utils.ComposeDirEnvVarName)
	}
	if utils.ComposeProject == "" {
		return fmt.Errorf("%s not set", utils.ComposeProjectEnvVarName)
	}

	// Get the request
	request, err := utils.HandleKeyedRequest[*stopRequest](c)
	if err != nil {
		return err
	}

	// Get the list of files to stop
	filesTopStop := []string{}
	if len(request.Services) == 0 {
		filesTopStop = []string{shared.ServiceContainerName + ".yml"}
	} else {
		for _, service := range request.Services {
			switch service {
			case shared.ServiceContainerName:
				filesTopStop = append(filesTopStop, shared.ServiceContainerName+".yml")
			}
		}
	}

	// Create the compose command args
	args := []string{
		"compose",
		"-p",
		utils.ComposeProject,
	}
	for _, file := range filesTopStop {
		args = append(args, "-f", filepath.Join(utils.ComposeDir, file))
	}
	args = append(args, "stop")

	// Stop the services
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
