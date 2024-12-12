package utils

import "github.com/nodeset-org/hyperdrive-example/shared"

const (
	// Directory for module logs
	AdapterLogDir string = "/hd/logs"

	// Name of the log file for the adapter
	AdapterLogFile string = "adapter.log"

	// Full path to the service log file
	ServiceLogPath string = AdapterLogDir + "/" + shared.ServiceLogFile

	// Full path to the adapter log file
	AdapterLogPath string = AdapterLogDir + "/" + AdapterLogFile

	// Directory for configuration files
	AdapterConfigDir string = "/hd/config"

	// Service configuration file
	ServiceConfigFile string = "native.cfg"

	// Adapter configuration file
	AdapterConfigFile string = "adapter.cfg"

	// Full path to the service configuration file
	ServiceConfigPath string = AdapterConfigDir + "/" + ServiceConfigFile

	// Full path to the adapter configuration file
	AdapterConfigPath string = AdapterConfigDir + "/" + AdapterConfigFile
)
