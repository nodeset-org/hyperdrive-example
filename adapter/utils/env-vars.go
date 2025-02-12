package utils

import "os"

// The mode to run the adapter in
type AdapterMode string

const (
	// Unknown (blank) mode; should be treated as an error since Hyperdrive always sets this
	AdapterMode_Unknown AdapterMode = ""

	// Global mode: stateless, no user-specific configuration.
	AdapterMode_Global AdapterMode = "global"

	// Project mode: stateful, user-specific configuration.
	AdapterMode_Project AdapterMode = "project"
)

// Adapter environment variables
const (
	// The mode to run the adapter in - global or project
	AdapterModeEnvVarName string = "HD_ADAPTER_MODE"

	// In project mode, the path of the user's config dir
	ConfigDirEnvVarName string = "HD_CONFIG_DIR"

	// In project mode, the path of the user's log dir
	LogDirEnvVarName string = "HD_LOG_DIR"

	// The path to the secret key file for authentication in project mode
	KeyFileEnvVarName string = "HD_KEY_FILE"

	// In project mode, the path to the instantiated compose file dir
	ComposeDirEnvVarName string = "HD_COMPOSE_DIR"

	// The name of the Docker Compose project for this module
	ComposeProjectEnvVarName string = "HD_COMPOSE_PROJECT"
)

var (
	// The mode the adapter has been started in
	Mode AdapterMode = AdapterMode(os.Getenv(AdapterModeEnvVarName))

	// If in project mode, the path to the user's config dir
	ConfigDir string = os.Getenv(ConfigDirEnvVarName)

	// If in project mode, the path to the user's log dir
	LogDir string = os.Getenv(LogDirEnvVarName)

	// The path to the secret key file for authentication in project mode
	KeyFile string = os.Getenv(KeyFileEnvVarName)

	// If in project mode, the path to the instantiated compose file dir
	ComposeDir string = os.Getenv(ComposeDirEnvVarName)

	// The name of the Docker Compose project for this module in project mode
	ComposeProject string = os.Getenv(ComposeProjectEnvVarName)
)
