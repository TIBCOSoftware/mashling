package model

// FlowModel defines the execution Model for a Flow.  It contains the
// execution behaviors for Flows and Tasks.
type FlowModel struct {
	name          string
	flowBehavior  FlowBehavior
	taskBehaviors map[int]TaskBehavior
}

// New creates a new FlowModel from the specified Behaviors
func New(name string) *FlowModel {

	var flowModel FlowModel
	flowModel.name = name
	flowModel.taskBehaviors = make(map[int]TaskBehavior)

	return &flowModel
}

// Name returns the name of the FlowModel
func (pm *FlowModel) Name() string {
	return pm.name
}

// RegisterFlowBehavior registers the specified FlowBehavior with the Model
func (pm *FlowModel) RegisterFlowBehavior(flowBehavior FlowBehavior) {

	pm.flowBehavior = flowBehavior
}

// GetFlowBehavior returns FlowBehavior of the FlowModel
func (pm *FlowModel) GetFlowBehavior() FlowBehavior {
	return pm.flowBehavior
}

// RegisterTaskBehavior registers the specified TaskBehavior with the Model
func (pm *FlowModel) RegisterTaskBehavior(id int, taskBehavior TaskBehavior) {
	pm.taskBehaviors[id] = taskBehavior
}

// GetTaskBehavior returns TaskBehavior with the specified ID in he FlowModel
func (pm *FlowModel) GetTaskBehavior(id int) TaskBehavior {
	return pm.taskBehaviors[id]
}
