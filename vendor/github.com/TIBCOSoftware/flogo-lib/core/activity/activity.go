package activity

import "github.com/TIBCOSoftware/flogo-lib/core/data"


// Activity is an interface for defining a custom Activity Execution
type Activity interface {

	// Eval is called when an Activity is being evaluated.  Returning true indicates
	// that the task is done.
	Eval(ctx Context) (done bool, err error)

	// Metadata returns the metadata of the activity
	Metadata() *Metadata
}

// DynamicIO is an optional interface that can be implemented by an activity.  If implemented,
// IOMetadata() will be invoked to determine the inputs/outputs of the activity instead of
// relying on the static information from the Activity's Metadata
type DynamicIO interface {

	// IOMetadata get the input/output metadata
	IOMetadata(ctx Context) (*data.IOMetadata, error)
}
