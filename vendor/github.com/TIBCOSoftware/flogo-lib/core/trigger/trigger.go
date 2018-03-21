package trigger

import (
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/util"
)

// Factory is used to create new instances for a trigger
type Factory interface {
	New(config *Config) Trigger
}

// Trigger is object that triggers/starts flow instances and
// is managed by an engine
type Trigger interface {
	util.Managed

	// Metadata returns the metadata of the trigger
	Metadata() *Metadata

	// Init sets up the trigger, it is called before Start()
	Init(actionRunner action.Runner)
}

// Trigger is object that triggers/starts flow instances and
// is managed by an engine
type TriggerDeprecated interface {
	util.Managed

	// TriggerMetadata returns the metadata of the trigger
	Metadata() *Metadata

	// Init sets up the trigger, it is called before Start()
	Init(config *Config, actionRunner action.Runner)
}

type Status string

const (
	Started Status = "Started"
	Stopped        = "Stopped"
	Failed         = "Failed"
)

//TriggerInstance contains all the information for a Trigger Instance, configuration and interface
type TriggerInstance struct {
	Config *Config
	Interf Trigger
	Status Status
	Error  error
}

type TriggerInstanceInfo struct {
	Name   string
	Status Status
	Error  error
}
