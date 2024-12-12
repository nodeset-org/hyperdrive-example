package hdmodule

import (
	"fmt"
	"strings"

	"github.com/kballard/go-shellquote"
	"github.com/nodeset-org/hyperdrive-example/adapter/utils"
	"github.com/urfave/cli/v2"
)

// Request format for `run`
type runRequest struct {
	utils.KeyedRequest

	// The command to run
	Command string `json:"command"`
}

// Handle the `run` command
func run(c *cli.Context) error {
	// Get the request
	request, err := utils.HandleKeyedRequest[*runRequest](c)
	if err != nil {
		return err
	}

	// Prevent recursive calls
	if strings.HasPrefix(request.Command, "hd-module") || strings.HasPrefix(request.Command, "hd") {
		return fmt.Errorf("recursive calls to `run` are not allowed")
	}

	// Run the app with the new command
	args, err := shellquote.Split(request.Command)
	if err != nil {
		return fmt.Errorf("error parsing command: %w", err)
	}
	return c.App.Run(args)
}
