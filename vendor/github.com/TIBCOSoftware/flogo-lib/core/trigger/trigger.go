package trigger

import (
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/util/managed"
)

// Factory is used to create new instances for a trigger
type Factory interface {
	New(config *Config) Trigger
}

// Trigger is object that triggers/starts flow instances and
// is managed by an engine
type Trigger interface {
	managed.Managed

	// Metadata returns the metadata of the trigger
	Metadata() *Metadata
}

// Initializable interface should be implemented by all Triggers, the Initialize method
// will eventually move up to Trigger to replace the the old "Init" method
type Initializable interface {

	// Initialize is called to initialize the Trigger
	Initialize(ctx InitContext) error
}

// InitContext is the initialization context for the trigger instance
type InitContext interface {

	// GetHandlers gets the handlers associated with the trigger
	GetHandlers() []*Handler
}

// Deprecated: No longer used
type InitOld interface {

	// Deprecated: Triggers should implement trigger.Initializable interface
	Init(actionRunner action.Runner)
}