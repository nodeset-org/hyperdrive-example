package utils

import "github.com/urfave/cli/v2"

// Flags
var (
	ConfigDirFlag *cli.StringFlag = &cli.StringFlag{
		Name:        "config-dir",
		Aliases:     []string{"c"},
		Usage:       "The path to the directory for module configuration files",
		DefaultText: "/hd/config",
	}
	LogDirFlag *cli.StringFlag = &cli.StringFlag{
		Name:        "log-dir",
		Aliases:     []string{"l"},
		Usage:       "The path to the directory for module log files",
		DefaultText: "/hd/logs",
	}
	KeyFileFlag *cli.StringFlag = &cli.StringFlag{
		Name:        "secret",
		Aliases:     []string{"s"},
		Usage:       "The path to the secret key file for authentication",
		DefaultText: "/hd/secret",
	}
)
