package behaviors

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////

// Task implements model.Task
type Task struct {
}

// Enter implements model.Task.Enter
func (tb *Task) Enter(ctx model.TaskContext) (enterResult model.EnterResult) {

	task := ctx.Task()
	log.Debugf("Enter Task '%s'", task.ID())

	ctx.SetStatus(model.TaskStatusEntered)

	//check if all predecessor links are done
	linkContexts := ctx.GetFromLinkInstances()

	ready := true
	skipped := false

	if len(linkContexts) == 0 {
		// has no predecessor links, so task is ready
		ready = true
	} else {
		skipped = true

		log.Debugf("Task '%s' has %d incoming Links", task.ID(), len(linkContexts))
		for _, linkContext := range linkContexts {

			log.Debugf("Task '%s' Link '%s' has status '%d'", task.ID(), linkContext.Link().ID(), linkContext.Status())
			if linkContext.Status() < model.LinkStatusFalse {
				ready = false
				break
			} else if linkContext.Status() == model.LinkStatusTrue {
				skipped = false
			}
		}
	}

	if ready {

		if skipped {
			log.Debugf("Task '%s' Skipped", ctx.Task().ID())
			ctx.SetStatus(model.TaskStatusSkipped)
			return model.ENTER_SKIP
		} else {
			log.Debugf("Task '%s' Ready", ctx.Task().ID())
			ctx.SetStatus(model.TaskStatusReady)
		}
		return model.ENTER_EVAL

	} else {
		log.Debugf("Task '%s' Not Ready", ctx.Task().ID())
	}

	return model.ENTER_NOTREADY
}

// Eval implements model.Task.Eval
func (tb *Task) Eval(ctx model.TaskContext) (evalResult model.EvalResult, err error) {

	if ctx.Status() == model.TaskStatusSkipped {
		return model.EVAL_SKIP, nil //todo introduce EVAL_SKIP?
	}

	task := ctx.Task()
	log.Debugf("Eval Task '%s'", task.ID())

	done, err := ctx.EvalActivity()

	if err != nil {
		log.Errorf("Error evaluating activity '%s'[%s] - %s", ctx.Task().ID(), ctx.Task().ActivityConfig().Ref(), err.Error())
		ctx.SetStatus(model.TaskStatusFailed)
		return model.EVAL_FAIL, err
	}

	if done {
		evalResult = model.EVAL_DONE
	} else {
		evalResult = model.EVAL_WAIT
	}

	return evalResult, nil
}

// PostEval implements model.Task.PostEval
func (tb *Task) PostEval(ctx model.TaskContext) (evalResult model.EvalResult, err error) {

	log.Debugf("PostEval Task '%s'", ctx.Task().ID())

	_, err = ctx.PostEvalActivity()

	//what to do if eval isn't "done"?
	if err != nil {
		log.Errorf("Error post evaluating activity '%s'[%s] - %s", ctx.Task().ID(), ctx.Task().ActivityConfig().Ref(), err.Error())
		ctx.SetStatus(model.TaskStatusFailed)
		return model.EVAL_FAIL, err
	}

	return model.EVAL_DONE, nil
}

// Done implements model.Task.Done
func (tb *Task) Done(ctx model.TaskContext) (notifyFlow bool, taskEntries []*model.TaskEntry, err error) {

	linkInsts := ctx.GetToLinkInstances()
	numLinks := len(linkInsts)

	ctx.SetStatus(model.TaskStatusDone)

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
				follow, err = ctx.EvalLink(linkInst.Link())

				if err != nil {
					return false, nil, err
				}
			}

			if follow {
				linkInst.SetStatus(model.LinkStatusTrue)

				taskEntry := &model.TaskEntry{Task: linkInst.Link().ToTask()}
				taskEntries = append(taskEntries, taskEntry)
			} else {
				linkInst.SetStatus(model.LinkStatusFalse)

				taskEntry := &model.TaskEntry{Task: linkInst.Link().ToTask()}
				taskEntries = append(taskEntries, taskEntry)
			}
		}

		//continue on to successor tasks
		return false, taskEntries, nil
	}

	log.Debugf("Notifying flow that task '%s' is done", ctx.Task().ID())

	// there are no outgoing links, so just notify parent that we are done
	return true, nil, nil
}

// Done implements model.Task.Skip
func (tb *Task) Skip(ctx model.TaskContext) (notifyFlow bool, taskEntries []*model.TaskEntry) {
	linkInsts := ctx.GetToLinkInstances()
	numLinks := len(linkInsts)

	ctx.SetStatus(model.TaskStatusSkipped)

	// process outgoing links
	if numLinks > 0 {

		taskEntries = make([]*model.TaskEntry, 0, numLinks)

		for _, linkInst := range linkInsts {
			linkInst.SetStatus(model.LinkStatusSkipped)
			taskEntry := &model.TaskEntry{Task: linkInst.Link().ToTask()}
			taskEntries = append(taskEntries, taskEntry)
		}

		return false, taskEntries
	}

	return true, nil
}

// Done implements model.Task.Error
func (tb *Task) Error(ctx model.TaskContext, err error) (handled bool, taskEntries []*model.TaskEntry) {

	linkInsts := ctx.GetToLinkInstances()
	numLinks := len(linkInsts)

	handled = false

	// process outgoing links
	if numLinks > 0 {

		for _, linkInst := range linkInsts {
			if linkInst.Link().Type() == definition.LtError {
				handled = true
			}
			break
		}

		if handled {
			taskEntries = make([]*model.TaskEntry, 0, numLinks)

			for _, linkInst := range linkInsts {

				if linkInst.Link().Type() == definition.LtError {
					linkInst.SetStatus(model.LinkStatusTrue)
				} else {
					linkInst.SetStatus(model.LinkStatusFalse)
				}

				taskEntry := &model.TaskEntry{Task: linkInst.Link().ToTask()}
				taskEntries = append(taskEntries, taskEntry)
			}

			return true, taskEntries
		}
	}

	return false, nil
}
