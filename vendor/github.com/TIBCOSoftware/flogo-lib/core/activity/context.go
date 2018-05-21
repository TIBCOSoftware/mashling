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

	// GetSetting gets the value of the specified setting
	GetSetting(setting string) (value interface{}, exists bool)

	// GetInitValue gets the specified initialization value
	GetInitValue(key string) (value interface{}, exists bool)

	// GetInput gets the value of the specified input attribute
	GetInput(name string) interface{}

	// GetOutput gets the value of the specified output attribute
	GetOutput(name string) interface{}

	// SetOutput sets the value of the specified output attribute
	SetOutput(name string, value interface{})

	// Deprecated: Use ActivityHost().Name() instead.
	TaskName() string

	// Deprecated: Use ActivityHost() instead.
	FlowDetails() FlowDetails
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

//type InitContext interface {
//
//	// GetSetting gets the value of the specified setting
//	GetSetting(setting string) (value interface{}, exists bool)
//
//	// GetResolver gets the resolver associated with the activity host
//	GetResolver() data.Resolver
//
//	// SetInitValue sets the value associated with this initialization
//	SetInitValue(key string, value interface{})
//}

// Deprecated: Use ActivityHost() instead.
type FlowDetails interface {

	// ID returns the ID of the Flow Instance
	ID() string

	// FlowName returns the name of the Flow
	Name() string

	// ReplyHandler returns the reply handler for the flow Instance
	ReplyHandler() ReplyHandler
}
