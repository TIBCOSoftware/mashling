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

// Initializable is an optional interface that can be implemented by an activity.  If implemented,
// it will be invoked for each corresponding activity configuration that has settings.
//type Initializable interface {
//
//	// Initialize is called to initialize the Activity for a particular configuration
//	Initialize(ctx InitContext) error
//}

// DynamicIO is an optional interface that can be implemented by an activity.  If implemented,
// IOMetadata() will be invoked to determine the inputs/outputs of the activity instead of
// relying on the static information from the Activity's Metadata
type DynamicIO interface {

	// IOMetadata get the input/output metadata
	IOMetadata(ctx Context) (*data.IOMetadata, error)
}
