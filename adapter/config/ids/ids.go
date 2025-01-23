package ids

import (
	hdconfig "github.com/nodeset-org/hyperdrive/modules/config"
)

const (
	ExampleBoolID      hdconfig.Identifier = "exampleBool"
	ExampleIntID       hdconfig.Identifier = "exampleInt"
	ExampleUintID      hdconfig.Identifier = "exampleUint"
	ExampleFloatID     hdconfig.Identifier = "exampleFloat"
	ExampleStringID    hdconfig.Identifier = "exampleString"
	ExampleChoiceID    hdconfig.Identifier = "exampleChoice"
	SubConfigID        hdconfig.Identifier = "subConfig"
	SubExampleBoolID   hdconfig.Identifier = "subConfigBool"
	SubExampleChoiceID hdconfig.Identifier = "subConfigChoice"
	ServerConfigID     hdconfig.Identifier = "server"
	PortModeID         hdconfig.Identifier = "portMode"
	PortID             hdconfig.Identifier = "port"
)
