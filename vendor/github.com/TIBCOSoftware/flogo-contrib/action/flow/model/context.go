package model

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// FlowContext is the execution context of the Flow when executing
// a Flow Behavior function
type FlowContext interface {
	// FlowDefinition returns the Flow definition associated with this context
	FlowDefinition() *definition.Definition

	//State gets the state of the Flow instance
	State() int

	//SetState sets the state of the Flow instance
	SetState(state int)
}

// TaskContext is the execution context of the Task when executing
// a Task Behavior function
type TaskContext interface {

	// State gets the state of the Task instance
	State() int

	// SetState sets the state of the Task instance
	SetState(state int)

	// Task returns the Task associated with this context
	Task() *definition.Task

	// FromInstLinks returns the instances of predecessor Links of the current
	// task.
	FromInstLinks() []LinkInst

	// ToInstLinks returns the instances of successor Links of the current
	// task.
	ToInstLinks() []LinkInst

	// EnterLeadingChildren enters the set of child Tasks that
	// do not have any incoming links.
	// todo: should we allow cross-boundary links?
	EnterLeadingChildren(enterCode int)

	// EnterChildren enters the set of child Tasks specified,
	// If single TaskEntry with nil Task is supplied,
	// all the child tasks are entered with the specified code.
	EnterChildren(taskEntries []*TaskEntry)

	// ChildTaskInsts gets all the instances of child tasks of the
	// current task
	ChildTaskInsts() (taskInsts []TaskInst, hasChildTasks bool)

	// EvalLink evaluates the specified link
	EvalLink(link *definition.Link) (bool, error)

	// HasActivity flag indicating if the task has an Activity
	HasActivity() bool

	// EvalActivity evaluates the Activity associated with the Task
	EvalActivity() (done bool, err error)

	// Failed marks the Activity as failed
	Failed(err error)

	GetSetting(setting string) (value interface{}, exists bool)

	AddWorkingData(attr *data.Attribute)

	UpdateWorkingData(key string, value interface{}) error

	GetWorkingData(key string) (*data.Attribute, bool)
}

// LinkInst is the instance of a link
type LinkInst interface {

	// Link returns the Link associated with this Link Instance
	Link() *definition.Link

	// State gets the state of the Link instance
	State() int

	// SetState sets the state of the Link instance
	SetState(state int)
}

type TaskInst interface {

	// Task returns the Task associated with this Task Instance
	Task() *definition.Task

	// State gets the state of the Task instance
	State() int
}
