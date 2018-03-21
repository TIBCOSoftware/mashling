package instance

import (
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// ChgType denotes the type of change for an object in an instance
type ChgType int

const (
	// CtAdd denotes an addition
	CtAdd ChgType = 1
	// CtUpd denotes an update
	CtUpd ChgType = 2
	// CtDel denotes an deletion
	CtDel ChgType = 3
)

// WorkItemQueueChange represents a change in the WorkItem Queue
type WorkItemQueueChange struct {
	ChgType  ChgType
	ID       int
	WorkItem *WorkItem
}

// TaskDataChange represents a change to a TaskData
type TaskDataChange struct {
	ChgType  ChgType
	ID       string
	TaskData *TaskData
}

// LinkDataChange represents a change to a LinkData
type LinkDataChange struct {
	ChgType  ChgType
	ID       int
	LinkData *LinkData
}

// InstanceChange represents a change to the instance
type InstanceChange struct {
	State       int
	Status      Status
	Changes     int
	AttrChanges []*AttributeChange
}

// AttributeChange represents a change to an Attribute
type AttributeChange struct {
	ChgType   ChgType
	Attribute *data.Attribute
}

// InstanceChangeTracker is used to track all changes to an instance
type InstanceChangeTracker struct {
	wiqChanges map[int]*WorkItemQueueChange

	tdChanges map[string]*TaskDataChange
	ldChanges map[int]*LinkDataChange

	instChange *InstanceChange
}

// NewInstanceChangeTracker creates an InstanceChangeTracker
func NewInstanceChangeTracker() *InstanceChangeTracker {

	var changes InstanceChangeTracker
	changes.instChange = new(InstanceChange)
	return &changes
}

// SetState is called to track a state change on an instance
func (ict *InstanceChangeTracker) SetState(state int) {

	ict.instChange.State = state
	//ict.ctxChanges.Changes |= CHG_STATE
}

// SetStatus is called to track a status change on an instance
func (ict *InstanceChangeTracker) SetStatus(status Status) {

	ict.instChange.Status = status
	//ict.ctxChanges.Changes |= CHG_STATUS
}

// AttrChange is called to track a status change of an Attribute
func (ict *InstanceChangeTracker) AttrChange(chgType ChgType, attribute *data.Attribute) {

	var attrChange AttributeChange
	attrChange.ChgType = chgType

	//var attr data.Attribute
	//attr.Name = attribute.Name
	//
	//if chgType == CtAdd {
	//	attr.Type = attribute.Type
	//	attr.Value = attribute.Value
	//} else if chgType == CtUpd {
	//	attr.Value = attribute.Value
	//}

	//attrChange.Attribute = &attr

	attrChange.Attribute = attribute
	ict.instChange.AttrChanges = append(ict.instChange.AttrChanges, &attrChange)
}

// trackWorkItem records a WorkItem Queue change
func (ict *InstanceChangeTracker) trackWorkItem(wiChange *WorkItemQueueChange) {

	if ict.wiqChanges == nil {
		ict.wiqChanges = make(map[int]*WorkItemQueueChange)
	}
	ict.wiqChanges[wiChange.ID] = wiChange
}

// trackTaskData records a TaskData change
func (ict *InstanceChangeTracker) trackTaskData(tdChange *TaskDataChange) {

	if ict.tdChanges == nil {
		ict.tdChanges = make(map[string]*TaskDataChange)
	}

	ict.tdChanges[tdChange.ID] = tdChange
}

// trackLinkData records a LinkData change
func (ict *InstanceChangeTracker) trackLinkData(ldChange *LinkDataChange) {

	if ict.ldChanges == nil {
		ict.ldChanges = make(map[int]*LinkDataChange)
	}
	ict.ldChanges[ldChange.ID] = ldChange
}

// ResetChanges is used to reset any tracking data stored on instance objects
func (ict *InstanceChangeTracker) ResetChanges() {

	// reset TaskData objects
	if ict.tdChanges != nil {
		for _, v := range ict.tdChanges {
			if v.TaskData != nil {
				//v.TaskData.ResetChanges()
			}
		}
	}

	// reset LinkData objects
	if ict.ldChanges != nil {
		for _, v := range ict.ldChanges {
			if v.LinkData != nil {
				//v.LinkData.ResetChanges()
			}
		}
	}
}
