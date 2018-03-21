package model

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
)

// TaskEntry is a struct used to specify what Task to
// enter and its corresponding enter code
type TaskEntry struct {
	Task      *definition.Task
	EnterCode int
}

// FlowBehavior is the execution behavior of the Flow.
type FlowBehavior interface {

	// Start the flow instance.  Returning true indicates that the
	// flow can start and eval will be scheduled on the Root Task.
	// Return false indicates that the flow could not be started
	// at this time.
	Start(context FlowContext) (start bool, evalCode int)

	// Resume the flow instance.  Returning true indicates that the
	// flow can resume.  Return false indicates that the flow
	// could not be resumed at this time.
	Resume(context FlowContext) bool //<---

	//do we need the following two

	// TasksDone is called when the RootTask is Done.
	TasksDone(context FlowContext, doneCode int)

	// Done is called when the flow is done.
	Done(context FlowContext) //maybe return something to the state server?
}

type EvalResult int

const (
	EVAL_FAIL EvalResult = iota
	EVAL_DONE
	EVAL_REPEAT
	EVAL_WAIT
)

// TaskBehavior is the execution behavior of a Task.
type TaskBehavior interface {

	// Enter determines if a Task is ready to be evaluated, returning true
	// indicates that the task is ready to be evaluated.
	Enter(context TaskContext, enterCode int) (eval bool, evalCode int)

	// Eval is called when a Task is being evaluated.  Returning true indicates
	// that the task is done.  If err is set, it indicates that the
	// behavior intends for the flow ErrorHandler to handle the error
	Eval(context TaskContext, evalCode int) (evalResult EvalResult, doneCode int, err error)

	// PostEval is called when a task that didn't complete during the Eval
	// needs to be notified.  Returning true indicates that the task is done.
	// If err is set, it indicates that the  behavior intends for the
	// flow ErrorHandler to handle the error
	PostEval(context TaskContext, evalCode int, data interface{}) (done bool, doneCode int, err error)

	// Done is called when Eval, PostEval or ChildDone return true, indicating
	// that the task is done.  This step is used to finalize the task and
	// determine the next set of tasks to be entered.  Returning true indicates
	// that the parent task should be notified.  Also returns the set of Tasks
	// that should be entered next.
	Done(context TaskContext, doneCode int) (notifyParent bool, childDoneCode int, taskEntries []*TaskEntry, err error)

	// Error is called when there is an issue executing Eval, it returns a boolean indicating
	// if it handled the error, otherwise the error is handled by the global error handler
	Error(context TaskContext) (handled bool, taskEntry *TaskEntry)

	// ChildDone is called when child task is Done and has indicated that its
	// parent should be notified.  Returning true indicates that the task
	// is done.
	ChildDone(context TaskContext, childTask *definition.Task, childDoneCode int) (done bool, doneCode int)
}
