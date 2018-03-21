package definition

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Definition is the object that describes the definition of
// a flow.  It contains its data (attributes) and
// structure (tasks & links).
type Definition struct {
	name          string
	modelID       string
	explicitReply bool
	rootTask      *Task
	ehTask        *Task

	attrs map[string]*data.Attribute

	links       map[int]*Link
	tasks       map[string]*Task

	linkExprMgr LinkExprManager
}

// Name returns the name of the definition
func (pd *Definition) Name() string {
	return pd.name
}

// ModelID returns the ID of the model the definition uses
func (pd *Definition) ModelID() string {
	return pd.modelID
}

// RootTask returns the root task of the definition
func (pd *Definition) RootTask() *Task {
	return pd.rootTask
}

func (pd *Definition) ExplicitReply() bool {
	return pd.explicitReply
}

// ErrorHandler returns the error handler task of the definition
func (pd *Definition) ErrorHandlerTask() *Task {
	return pd.ehTask
}

// GetAttr gets the specified attribute
func (pd *Definition) GetAttr(attrName string) (attr *data.Attribute, exists bool) {

	if pd.attrs != nil {
		attr, found := pd.attrs[attrName]
		if found {
			return attr, true
		}
	}

	return nil, false
}

// GetTask returns the task with the specified ID
func (pd *Definition) GetTask(taskID string) *Task {
	task := pd.tasks[taskID]
	return task
}

// GetLink returns the link with the specified ID
func (pd *Definition) GetLink(linkID int) *Link {
	task := pd.links[linkID]
	return task
}

// SetLinkExprManager sets the Link Expression Manager for the definition
func (pd *Definition) SetLinkExprManager(mgr LinkExprManager) {
	// todo revisit
	pd.linkExprMgr = mgr
}

// GetLinkExprManager gets the Link Expression Manager for the definition
func (pd *Definition) GetLinkExprManager() LinkExprManager {
	return pd.linkExprMgr
}

////////////////////////////////////////////////////////////////////////////
// Task

// Task is the object that describes the definition of
// a task.  It contains its data (attributes) and its
// nested structure (child tasks & child links).
type Task struct {
	id           string
	typeID       int
	activityType string
	activityRef  string
	name         string
	tasks        []*Task
	links        []*Link
	isScope      bool

	definition *Definition
	parent     *Task

	settings    map[string]interface{}
	inputAttrs  map[string]*data.Attribute
	outputAttrs map[string]*data.Attribute

	inputMapper  data.Mapper
	outputMapper data.Mapper

	toLinks   []*Link
	fromLinks []*Link
}

// ID gets the id of the task
func (task *Task) ID() string {
	return task.id
}

// Name gets the name of the task
func (task *Task) Name() string {
	return task.name
}

// TypeID gets the id of the task type
func (task *Task) TypeID() int {
	return task.typeID
}

// ActivityType gets the activity type
func (task *Task) ActivityType() string {
	return task.activityType
}

// ActivityRef gets the activity ref
func (task *Task) ActivityRef() string {
	return task.activityRef
}

// Parent gets the parent task of the task
func (task *Task) Parent() *Task {
	return task.parent
}

// ChildTasks gets the child tasks of the task
func (task *Task) ChildTasks() []*Task {
	return task.tasks
}

// ChildLinks gets the child tasks of the task
func (task *Task) ChildLinks() []*Link {
	return task.links
}

func (task *Task) GetSetting(attrName string) (value interface{}, exists bool) {
	value, exists = task.settings[attrName]
	return value,exists
}

// GetAttr gets the specified attribute
// DEPRECATED
func (task *Task) GetAttr(attrName string) (attr *data.Attribute, exists bool) {

	if task.inputAttrs != nil {
		attr, found := task.inputAttrs[attrName]
		if found {
			return attr, true
		}
	}

	return nil, false
}

// GetAttr gets the specified input attribute
func (task *Task) GetInputAttr(attrName string) (attr *data.Attribute, exists bool) {

	if task.inputAttrs != nil {
		attr, found := task.inputAttrs[attrName]
		if found {
			return attr, true
		}
	}

	return nil, false
}

// GetOutputAttr gets the specified output attribute
func (task *Task) GetOutputAttr(attrName string) (attr *data.Attribute, exists bool) {

	if task.outputAttrs != nil {
		attr, found := task.outputAttrs[attrName]
		if found {
			return attr, true
		}
	}

	return nil, false
}

// ToLinks returns the predecessor links of the task
func (task *Task) ToLinks() []*Link {
	return task.toLinks
}

// FromLinks returns the successor links of the task
func (task *Task) FromLinks() []*Link {
	return task.fromLinks
}

// InputMapper returns the InputMapper of the task
func (task *Task) InputMapper() data.Mapper {
	return task.inputMapper
}

// OutputMapper returns the OutputMapper of the task
func (task *Task) OutputMapper() data.Mapper {
	return task.outputMapper
}

func (task *Task) String() string {
	return fmt.Sprintf("Task[%d]:'%s'", task.id, task.name)
}

// IsScope returns flag indicating if the Task is a scope task (a container of attributes)
func (task *Task) IsScope() bool {
	return task.isScope
}

////////////////////////////////////////////////////////////////////////////
// Link

// LinkType is an enum for possible Link Types
type LinkType int

const (
	// LtDependency denotes an normal dependency link
	LtDependency LinkType = 0

	// LtExpression denotes a link with an expression
	LtExpression LinkType = 1 //expr language on the model or def?

	// LtLabel denotes 'labelled' link
	LtLabel LinkType = 2

	// LtError denotes an error link
	LtError LinkType = 3
)

// Link is the object that describes the definition of
// a link.
type Link struct {
	id       int
	name     string
	fromTask *Task
	toTask   *Task
	linkType LinkType
	value    string //expression or label

	definition *Definition
	parent     *Task
}

// ID gets the id of the link
func (link *Link) ID() int {
	return link.id
}

// Type gets the link type
func (link *Link) Type() LinkType {
	return link.linkType
}

// Value gets the "value" of the link
func (link *Link) Value() string {
	return link.value
}

// FromTask returns the task the link is coming from
func (link *Link) FromTask() *Task {
	return link.fromTask
}

// ToTask returns the task the link is going to
func (link *Link) ToTask() *Task {
	return link.toTask
}

func (link *Link) String() string {
	return fmt.Sprintf("Link[%d]:'%s' - [from:%d, to:%d]", link.id, link.name, link.fromTask.id, link.toTask.id)
}
