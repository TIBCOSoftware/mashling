package activity

import "github.com/TIBCOSoftware/flogo-lib/core/action"

// Context describes the execution context for an Activity.
// It provides access to attributes, task and Flow information.
type Context interface {

	// FlowDetails gets the action.Context under with the activity is executing
	ActionContext() action.Context

	// TaskName returns the name of the Task the Activity is currently executing
	TaskName() string

	// GetInput gets the value of the specified input attribute
	GetInput(name string) interface{}

	// GetOutput gets the value of the specified output attribute
	GetOutput(name string) interface{}

	// SetOutput sets the value of the specified output attribute
	SetOutput(name string, value interface{})

	//Deprecated
	// FlowDetails returns the details fo the Flow Instance
	FlowDetails() FlowDetails
}

// Deprecated
// FlowDetails details of the flow that is being executed
type FlowDetails interface {

	// ID returns the ID of the Flow Instance
	ID() string

	// FlowName returns the name of the Flow
	Name() string

	// ReplyHandler returns the reply handler for the flow Instance
	ReplyHandler() ReplyHandler
}
