package config

import (
	"github.com/nodeset-org/hyperdrive-example/adapter/config/ids"
	"github.com/nodeset-org/hyperdrive-example/shared"
	hdconfig "github.com/nodeset-org/hyperdrive/modules/config"
)

type ServerConfig struct {
	hdconfig.SectionHeader

	Port hdconfig.UintParameter

	PortMode hdconfig.ChoiceParameter[PortMode]
}

type ServerConfigSettings struct {
	Port     uint64   `json:"port" yaml:"port"`
	PortMode PortMode `json:"portMode" yaml:"portMode"`
}

func NewServerConfig() *ServerConfig {
	cfg := &ServerConfig{}
	cfg.ID = ids.ServerConfigID
	cfg.Name = "Service Config"
	cfg.Description.Default = "This is the configuration for the module's service. This isn't used by the service directly, but it is used by Hyperdrive itself in the service's Docker Compose file template to configure the service during its starting process."

	// Port
	cfg.Port.ID = ids.PortID
	cfg.Port.Name = "API Port"
	cfg.Port.Description.Default = "This is the API port the server should run on."
	cfg.Port.Default = uint64(shared.DefaultServerApiPort)
	cfg.Port.MinValue = 0
	cfg.Port.MaxValue = 65535
	cfg.Port.AffectedContainers = []string{shared.ServiceContainerName}

	// Options for PortMode
	options := make([]hdconfig.ParameterOption[PortMode], 3)
	options[0].Name = "Closed"
	options[0].Description.Default = "The API is only accessible to internal Docker container traffic."
	options[0].Value = PortMode_Closed

	options[1].Name = "Localhost Only"
	options[1].Description.Default = "The API is accessible from internal Docker containers and your own local machine, but no other external machines."
	options[1].Value = PortMode_Localhost

	options[2].Name = "All External Traffic"
	options[2].Description.Default = "The port is accessible to everything, including external machines.\n\n[orange]Use with caution!"
	options[2].Value = PortMode_External

	// PortMode
	cfg.PortMode.ID = ids.PortModeID
	cfg.PortMode.Name = "Expose API Port"
	cfg.PortMode.Description.Default = "Determine how the server's HTTP API restricts its access from various sources."
	cfg.PortMode.Options = options
	cfg.PortMode.Default = options[0].Value
	cfg.PortMode.AffectedContainers = []string{shared.ServiceContainerName}

	return cfg
}

func (cfg ServerConfig) GetParameters() []hdconfig.IParameter {
	return []hdconfig.IParameter{
		&cfg.Port,
		&cfg.PortMode,
	}
}

func (cfg ServerConfig) GetSections() []hdconfig.ISection {
	return []hdconfig.ISection{}
}
