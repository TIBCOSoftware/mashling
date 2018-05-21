package action

import (
	"context"
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Action is an action to perform as a result of a trigger
type Action interface {
	// Metadata get the Action's metadata
	Metadata() *Metadata

	// IOMetadata get the Action's IO metadata
	IOMetadata() *data.IOMetadata
}

// SyncAction is a synchronous action to perform as a result of a trigger
type SyncAction interface {
	Action

	// Run this Action
	Run(context context.Context, inputs map[string]*data.Attribute) (map[string]*data.Attribute, error)
}

// AsyncAction is an asynchronous action to perform as a result of a trigger, the action can asynchronously
// return results as it runs.  It returns immediately, but will continue to run.
type AsyncAction interface {
	Action

	// Run this Action
	Run(context context.Context, inputs map[string]*data.Attribute, handler ResultHandler) error
}

// Factory is used to create new instances for an action
type Factory interface {

	// New create a new Action
	New(config *Config) (Action, error)
}

// GetMetadata method to ensure we have metadata, remove in future
func GetMetadata(act Action) *Metadata {
	if act.Metadata() == nil {
		_, async := act.(AsyncAction)
		return &Metadata{ID: fmt.Sprintf("%T", act), Async: async}
	} else {
		return act.Metadata()
	}
}

// Runner runs actions
type Runner interface {
	// Deprecated: Use Execute() instead
	Run(context context.Context, act Action, uri string, options interface{}) (code int, data interface{}, err error)

	// Deprecated: Use Execute() instead
	RunAction(ctx context.Context, act Action, options map[string]interface{}) (results map[string]*data.Attribute, err error)

	// Execute the specified Action
	Execute(ctx context.Context, act Action, inputs map[string]*data.Attribute) (results map[string]*data.Attribute, err error)
}

// ResultHandler used to handle results from the Action
type ResultHandler interface {

	// HandleResult is invoked when there are results available
	HandleResult(results map[string]*data.Attribute, err error)

	// Done indicates that the action has completed
	Done()
}
