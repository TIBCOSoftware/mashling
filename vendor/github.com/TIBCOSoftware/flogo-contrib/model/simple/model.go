package simple

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("model-tibco-simple")

const (
	MODEL_NAME = "tibco-simple"
)

func init() {
	model.Register(New())
}

////////////////////////////////////////////////////////////////////////////////////////////////////////

func New() *model.FlowModel {
	m := model.New(MODEL_NAME)
	m.RegisterFlowBehavior(&SimpleFlowBehavior{})
	m.RegisterTaskBehavior(1, &SimpleTaskBehavior{})
	m.RegisterTaskBehavior(2, &SimpleIteratorTaskBehavior{})
	return m
}

// SimpleFlowBehavior implements model.FlowBehavior
type SimpleFlowBehavior struct {
}

// Start implements model.FlowBehavior.Start
func (pb *SimpleFlowBehavior) Start(context model.FlowContext) (start bool, evalCode int) {
	// just schedule the root task
	return true, 0
}

// Resume implements model.FlowBehavior.Resume
func (pb *SimpleFlowBehavior) Resume(context model.FlowContext) bool {
	return true
}

// TasksDone implements model.FlowBehavior.TasksDone
func (pb *SimpleFlowBehavior) TasksDone(context model.FlowContext, doneCode int) {
	// all tasks are done
}

// Done implements model.FlowBehavior.Done
func (pb *SimpleFlowBehavior) Done(context model.FlowContext) {
	log.Debugf("Flow Done\n")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////

// SimpleTaskBehavior implements model.TaskBehavior
type SimpleTaskBehavior struct {
}

// Enter implements model.TaskBehavior.Enter
func (tb *SimpleTaskBehavior) Enter(context model.TaskContext, enterCode int) (eval bool, evalCode int) {

	task := context.Task()
	log.Debugf("Task Enter: %s\n", task.Name())

	context.SetState(STATE_ENTERED)

	//check if all predecessor links are done

	linkContexts := context.FromInstLinks()

	ready := true
	skipped := false

	if len(linkContexts) == 0 {
		// has no predecessor links, so task is ready
		ready = true
	} else {
		skipped = true

		log.Debugf("Num Links: %d\n", len(linkContexts))
		for _, linkContext := range linkContexts {

			log.Debugf("Task: %s, linkData: %v\n", task.Name(), linkContext)
			if linkContext.State() < STATE_LINK_FALSE {
				ready = false
				break
			} else if linkContext.State() == STATE_LINK_TRUE {
				skipped = false
			}
		}
	}

	if ready {

		if skipped {
			log.Debugf("Task Skipped\n")
			context.SetState(STATE_SKIPPED)
			//todo hack, wait for explicit skip support from engine
			return ready, -666
		} else {
			log.Debugf("Task Ready\n")
			context.SetState(STATE_READY)
		}

	} else {
		log.Debugf("Task Not Ready\n")
	}

	return ready, 0
}

// Eval implements model.TaskBehavior.Eval
func (tb *SimpleTaskBehavior) Eval(context model.TaskContext, evalCode int) (evalResult model.EvalResult, doneCode int, err error) {

	if context.State() == STATE_SKIPPED {
		return model.EVAL_DONE, EC_SKIP, nil
	}

	task := context.Task()
	log.Debugf("Task Eval: %v\n", task)

	if len(task.ChildTasks()) > 0 {
		log.Debugf("Has Children\n")

		//has children, so set to waiting
		context.SetState(STATE_WAITING)

		context.EnterLeadingChildren(0)

		return model.EVAL_WAIT, 0, nil

	} else {

		if context.HasActivity() {

			done, err := context.EvalActivity()

			if err != nil {
				log.Errorf("Error evaluating activity '%s'[%s] - %s", context.Task().Name(), context.Task().ActivityType(), err.Error())
				context.SetState(STATE_FAILED)
				return model.EVAL_FAIL, 0, err
			}

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
}

// PostEval implements model.TaskBehavior.PostEval
func (tb *SimpleTaskBehavior) PostEval(context model.TaskContext, evalCode int, data interface{}) (done bool, doneCode int, err error) {

	log.Debugf("Task PostEval\n")

	if context.HasActivity() { //if activity is async

		//done := activity.PostEval(activityContext, data)
		done := true
		return done, 0, nil
	}

	//no-op
	return true, 0, nil
}

// Done implements model.TaskBehavior.Done
func (tb *SimpleTaskBehavior) Done(context model.TaskContext, doneCode int) (notifyParent bool, childDoneCode int, taskEntries []*model.TaskEntry, err error) {

	task := context.Task()

	linkInsts := context.ToInstLinks()
	numLinks := len(linkInsts)

	if context.State() == STATE_SKIPPED {
		log.Debugf("skipped task: %s\n", task.Name())

		// skip outgoing links
		if numLinks > 0 {

			taskEntries = make([]*model.TaskEntry, 0, numLinks)
			for _, linkInst := range linkInsts {

				linkInst.SetState(STATE_LINK_SKIPPED)

				//todo: engine should not eval mappings for skipped tasks, skip
				//todo: needs to be a state/op understood by the engine
				taskEntry := &model.TaskEntry{Task: linkInst.Link().ToTask(), EnterCode: EC_SKIP}
				taskEntries = append(taskEntries, taskEntry)
			}

			//continue on to successor tasks
			return false, 0, taskEntries, nil
		}
	} else {
		log.Debugf("done task: %s", task.Name())

		context.SetState(STATE_DONE)
		//context.SetTaskDone() for task garbage collection

		// process outgoing links
		if numLinks > 0 {

			taskEntries = make([]*model.TaskEntry, 0, numLinks)

			for _, linkInst := range linkInsts {

				follow := true

				if linkInst.Link().Type() == definition.LtError {
					//todo should we skip or ignore?
					continue
				}

				if linkInst.Link().Type() == definition.LtExpression {
					//todo handle error
					follow, err = context.EvalLink(linkInst.Link())

					if err != nil {
						return false, 0, nil, err
					}
				}

				if follow {
					linkInst.SetState(STATE_LINK_TRUE)

					taskEntry := &model.TaskEntry{Task: linkInst.Link().ToTask(), EnterCode: 0}
					taskEntries = append(taskEntries, taskEntry)
				} else {
					linkInst.SetState(STATE_LINK_FALSE)

					taskEntry := &model.TaskEntry{Task: linkInst.Link().ToTask(), EnterCode: EC_SKIP}
					taskEntries = append(taskEntries, taskEntry)
				}
			}

			//continue on to successor tasks
			return false, 0, taskEntries, nil
		}
	}

	log.Debug("notifying parent that task is done")

	// there are no outgoing links, so just notify parent that we are done
	return true, 0, nil, nil
}

// Done implements model.TaskBehavior.Error
func (tb *SimpleTaskBehavior) Error(context model.TaskContext) (handled bool, taskEntry *model.TaskEntry) {

	linkInsts := context.ToInstLinks()
	numLinks := len(linkInsts)

	// process outgoing links
	if numLinks > 0 {

		for _, linkInst := range linkInsts {

			if linkInst.Link().Type() == definition.LtError {
				linkInst.SetState(STATE_LINK_TRUE)
				taskEntry := &model.TaskEntry{Task: linkInst.Link().ToTask(), EnterCode: 0}
				return true, taskEntry
			}
		}
	}

	// there are no outgoing error links, so just return false
	return false, nil
}

// ChildDone implements model.TaskBehavior.ChildDone
func (tb *SimpleTaskBehavior) ChildDone(context model.TaskContext, childTask *definition.Task, childDoneCode int) (done bool, doneCode int) {

	childTasks, hasChildren := context.ChildTaskInsts()

	if !hasChildren {
		log.Debug("Task ChildDone - No Children")
		return true, 0
	}

	for _, taskInst := range childTasks {

		if taskInst.State() < STATE_DONE {

			log.Debugf("task %s not done or skipped", taskInst.Task().Name())
			return false, 0
		}
	}

	log.Debug("all child tasks done or skipped")

	// our children are done, so just transition ourselves to done
	return true, 0
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// State
const (
	EC_SKIP = 1

	STATE_NOT_STARTED int = 0

	STATE_LINK_FALSE   int = 1
	STATE_LINK_TRUE    int = 2
	STATE_LINK_SKIPPED int = 3

	STATE_ENTERED int = 10
	STATE_READY   int = 20
	STATE_WAITING int = 30
	STATE_DONE    int = 40
	STATE_SKIPPED int = 50
	STATE_FAILED  int = 100
)
