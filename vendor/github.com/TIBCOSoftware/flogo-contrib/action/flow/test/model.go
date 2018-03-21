package test

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

func init() {
	model.Register(NewTestModel())
}

func NewTestModel() *model.FlowModel {
	m := model.New("test")
	m.RegisterFlowBehavior(&SimpleFlowBehavior{})
	m.RegisterTaskBehavior(1, &SimpleTaskBehavior{})

	return m
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// SimpleFlow

type SimpleFlowBehavior struct {
}

func (b *SimpleFlowBehavior) Start(context model.FlowContext) (start bool, evalCode int) {

	//just schedule the root task
	return true, 0
}

func (b *SimpleFlowBehavior) Resume(context model.FlowContext) bool {

	return true
}

func (b *SimpleFlowBehavior) TasksDone(context model.FlowContext, doneCode int) {
	logger.Debugf("Flow TasksDone\n")

}

func (b *SimpleFlowBehavior) Done(context model.FlowContext) {
	logger.Debugf("Flow Done\n")

}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// SimpleTask

type SimpleTaskBehavior struct {
}

func (b *SimpleTaskBehavior) Enter(context model.TaskContext, enterCode int) (eval bool, evalCode int) {

	task := context.Task()
	//check if all predecessor links are done
	logger.Debugf("Task Enter: %s\n", task.Name())

	context.SetState(STATE_ENTERED)

	linkContexts := context.FromInstLinks()

	ready := true

	if len(linkContexts) == 0 {
		ready = true
	} else {

		logger.Debugf("Num Links: %d\n", len(linkContexts))
		for _, linkContext := range linkContexts {

			logger.Debugf("Task: %s, linkData: %v\n", task.Name(), linkContext)
			if linkContext.State() != STATE_LINK_TRUE {
				ready = false
				break
			}
		}
	}

	if ready {
		logger.Debugf("Task Ready\n")
		context.SetState(STATE_READY)
	} else {
		logger.Debugf("Task Not Ready\n")
	}

	return ready, 0
}

func (b *SimpleTaskBehavior) Eval(context model.TaskContext, evalCode int) (evalResult model.EvalResult, doneCode int, err error) {

	task := context.Task()
	logger.Debugf("Task Eval: %v\n", task)

	if len(task.ChildTasks()) > 0 {
		logger.Debugf("Has Children\n")

		context.SetState(STATE_WAITING)

		//for now enter all children (bpel style) - costly
		context.EnterChildren(nil)

		return model.EVAL_WAIT, 0, nil
	}

	if context.HasActivity() {

		//log.Debug("Evaluating Activity: ", activity.GetType())
		done, _ := context.EvalActivity()
		if done {
			evalResult = model.EVAL_DONE
		} else {
			evalResult = model.EVAL_WAIT
		}

		return evalResult, 0, nil
	}

	//no-op
	return model.EVAL_DONE, 0, nil
}

func (b *SimpleTaskBehavior) PostEval(context model.TaskContext, evalCode int, data interface{}) (done bool, doneCode int, err error) {
	logger.Debugf("Task PostEval\n")

	if context.HasActivity() { //and is async

		//done := activity.PostEval(activityContext, data)
		done := true

		return done, 0, nil
	}
	//no-op
	return true, 0, nil
}

func (b *SimpleTaskBehavior) Done(context model.TaskContext, doneCode int) (notifyParent bool, childDoneCode int, taskEntries []*model.TaskEntry, err error) {

	context.SetState(STATE_DONE)
	//context.SetTaskDone() for task garbage collection

	task := context.Task()

	logger.Debugf("done task:%s\n", task.Name())

	links := task.ToLinks()

	numLinks := len(links)

	if numLinks > 0 {

		taskEntries := make([]*model.TaskEntry, 0, numLinks)

		for _, link := range links {

			follow, _ := context.EvalLink(link)
			if follow {

				taskEntry := &model.TaskEntry{Task: link.ToTask(), EnterCode: 0}
				taskEntries = append(taskEntries, taskEntry)
			}
		}

		//continue on to successor links
		return false, 0, taskEntries, nil
	}

	//notify parent that we are done
	return true, 0, nil, nil
}

func (b *SimpleTaskBehavior) ChildDone(context model.TaskContext, childTask *definition.Task, childDoneCode int) (done bool, doneCode int) {
	logger.Debugf("Task ChildDone\n")

	return true, 0
}

func (b *SimpleTaskBehavior) Error(context model.TaskContext) (handled bool, taskEntry *model.TaskEntry) {
	return false, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// SimpleLink

type SimpleLinkBehavior struct {
}

func (b *SimpleLinkBehavior) Eval(context model.LinkInst, evalCode int) {

	logger.Debugf("Link Eval\n")

	context.SetState(STATE_LINK_TRUE)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// State
const (
	STATE_NOT_STARTED int = 0

	STATE_LINK_FALSE int = 1
	STATE_LINK_TRUE  int = 2

	STATE_ENTERED int = 10
	STATE_READY   int = 20
	STATE_WAITING int = 30
	STATE_DONE    int = 40
)
