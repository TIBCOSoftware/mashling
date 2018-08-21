package activity

import (
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	)

// Context describes the execution context for an Activity.
// It provides access to attributes, task and Flow information.
type Context interface {

	// ActivityHost gets the "host" under with the activity is executing
	ActivityHost() Host

	//Name the name of the activity that is currently executing
	Name() string

	// GetInput gets the value of the specified input attribute
	GetInput(name string) interface{}

	// GetOutput gets the value of the specified output attribute
	GetOutput(name string) interface{}

	// SetOutput sets the value of the specified output attribute
	SetOutput(name string, value interface{})

	//// GetSharedTempData get shared temporary data for activity, lifespan
	//// of the data dependent on the activity host implementation
	//GetSharedTempData() map[string]interface{}

	/////////////////
	// Deprecated

	// GetSetting gets the value of the specified setting
	// Deprecated
	GetSetting(setting string) (value interface{}, exists bool)

	// GetInitValue gets the specified initialization value
	// Deprecated
	GetInitValue(key string) (value interface{}, exists bool)

	// Deprecated: Use ActivityHost().Name() instead.
	TaskName() string

	// Deprecated: Use ActivityHost() instead.
	FlowDetails() FlowDetails

	/////////////////
}

type Host interface {

	// ID returns the ID of the Activity Host
	ID() string

	// Name the name of the Activity Host
	Name() string

	// IOMetadata get the input/output metadata of the activity host
	IOMetadata() *data.IOMetadata

	// Reply is used to reply to the activity Host with the results of the execution
	Reply(replyData map[string]*data.Attribute, err error)

	// Return is used to indicate to the activity Host that it should complete and return the results of the execution
	Return(returnData map[string]*data.Attribute, err error)

	//todo rename, essentially the flow's attrs for now
	WorkingData() data.Scope

	// GetResolver gets the resolver associated with the activity host
	GetResolver() data.Resolver

	//Map with action specific details/properties, flowId, etc.
	//GetDetails() map[string]string
}


// Deprecated: Use ActivityHost() instead.
type FlowDetails interface {

	// ID returns the ID of the Flow Instance
	ID() string

	// FlowName returns the name of the Flow
	Name() string

	// ReplyHandler returns the reply handler for the flow Instance
	ReplyHandler() ReplyHandler
}

//SharedTempDataSupport - temporary interface until we transition this to activity.Context
//Deprecated
type SharedTempDataSupport interface {

	// GetSharedTempData get shared temporary data for activity, lifespan
	// of the data dependent on the activity host implementation
	GetSharedTempData() map[string]interface{}
}

// GetSharedTempDataSupport for the activity
func GetSharedTempDataSupport(ctx Context) (SharedTempDataSupport, bool) {

	ts, ok :=  ctx.(SharedTempDataSupport)
	return ts, ok
}
