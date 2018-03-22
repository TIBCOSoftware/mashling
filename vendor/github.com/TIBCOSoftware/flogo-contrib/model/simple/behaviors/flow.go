package behaviors

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("flowmodel-simple")

// Flow implements model.FlowBehavior
type Flow struct {
}

// Start implements model.Flow.Start
func (fb *Flow) Start(ctx model.FlowContext) (started bool, taskEntries []*model.TaskEntry) {
	return true, getFlowTaskEntries(ctx.FlowDefinition().Tasks(), true)
}

// StartErrorHandler implements model.Flow.StartErrorHandler
func (fb *Flow) StartErrorHandler(ctx model.FlowContext) (taskEntries []*model.TaskEntry) {
	return getFlowTaskEntries(ctx.FlowDefinition().GetErrorHandler().Tasks(), true)
}

// Resume implements model.Flow.Resume
func (fb *Flow) Resume(ctx model.FlowContext) (resumed bool) {
	return true
}

// TasksDone implements model.Flow.TasksDone
func (fb *Flow) TaskDone(ctx model.FlowContext) (flowDone bool) {
	tasks := ctx.TaskInstances()

	for _, taskInst := range tasks {

		if taskInst.Status() < model.TaskStatusDone { //ignore not started?

			log.Debugf("task %s not done or skipped", taskInst.Task().Name())
			return false
		}
	}

	log.Debug("all tasks done or skipped")

	// our tasks are done, so the flow is done
	return true
}

// Done implements model.Flow.Done
func (fb *Flow) Done(ctx model.FlowContext) {
	log.Debugf("Flow Done\n")
}

func getFlowTaskEntries(tasks []*definition.Task, leadingOnly bool) []*model.TaskEntry {

	var taskEntries []*model.TaskEntry

	for _, task := range tasks {

		if len(task.FromLinks()) == 0 || !leadingOnly {

			taskEntry := &model.TaskEntry{Task: task, EnterCode: 0}
			taskEntries = append(taskEntries, taskEntry)
		}
	}

	return taskEntries
}
