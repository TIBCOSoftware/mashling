package simple

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// SimpleIteratorTaskBehavior implements model.TaskBehavior
type SimpleIteratorTaskBehavior struct {
}

// Enter implements model.TaskBehavior.Enter
func (tb *SimpleIteratorTaskBehavior) Enter(context model.TaskContext, enterCode int) (eval bool, evalCode int) {

	//todo inherit this code from base task

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

type Iteration struct {
	Key interface{}
	Value interface{}
}

// Eval implements model.TaskBehavior.Eval
func (tb *SimpleIteratorTaskBehavior) Eval(context model.TaskContext, evalCode int) (evalResult model.EvalResult, doneCode int, err error) {

	if context.State() == STATE_SKIPPED {
		return model.EVAL_DONE, EC_SKIP, nil
	}

	task := context.Task()
	log.Debugf("Task Eval: %v\n", task)

	if context.HasActivity() {

		var itx Iterator

		itxAttr, ok := context.GetWorkingData("_iterator")
		iterationAttr, _ := context.GetWorkingData("iteration")

		if ok {
			itx = itxAttr.Value().(Iterator)
		} else {

			iterateOn, ok := context.GetSetting("iterate")

			if !ok {
				//todo if iterateOn is not defined, what should we do?
				//just skip for now
				return model.EVAL_DONE, 0, nil
			}

			switch t := iterateOn.(type) {
			case string:
				count, err := data.CoerceToInteger(iterateOn)
				if err != nil {
					return model.EVAL_FAIL, 0, err
				}
				itx = NewIntIterator(count)
			case int:
				count := iterateOn.(int)
				itx = NewIntIterator(count)
			case float64:
				count := int(iterateOn.(float64))
				itx = NewIntIterator(count)
			case map[string]interface{}:
				itx = NewObjectIterator(t)
			case []interface{}:
				itx = NewArrayIterator(t)
			default:
				return model.EVAL_FAIL, 0, fmt.Errorf("unsupported type '%s' for iterateOn", t)
			}

			itxAttr, _ = data.NewAttribute("_iterator", data.ANY, itx)
			context.AddWorkingData(itxAttr)

			iteration := map[string]interface{}{
				"key": nil,
				"value":   nil,
			}

			iterationAttr, _ = data.NewAttribute("iteration", data.OBJECT, iteration)
			context.AddWorkingData(iterationAttr)
		}

		repeat := itx.next()

		if repeat {
			log.Debugf("Repeat:%s, Key:%s, Value:%v", repeat, itx.Key(), itx.Value())

			iteration,_ := iterationAttr.Value().(map[string]interface{})
			iteration["key"] = itx.Key()
			iteration["value"] = itx.Value()

			_, err := context.EvalActivity()

			//what to do if eval isn't "done"?
			if err != nil {
				log.Errorf("Error evaluating activity '%s'[%s] - %s", context.Task().Name(), context.Task().ActivityType(), err.Error())
				context.SetState(STATE_FAILED)
				return model.EVAL_FAIL, 0, err
			}

			evalResult = model.EVAL_REPEAT
		} else {
			evalResult = model.EVAL_DONE
		}

		return evalResult, 0, nil
	}

	//no-op
	return model.EVAL_DONE, 0, nil
}

// PostEval implements model.TaskBehavior.PostEval
func (tb *SimpleIteratorTaskBehavior) PostEval(context model.TaskContext, evalCode int, data interface{}) (done bool, doneCode int, err error) {

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
func (tb *SimpleIteratorTaskBehavior) Done(context model.TaskContext, doneCode int) (notifyParent bool, childDoneCode int, taskEntries []*model.TaskEntry, err error) {

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
func (tb *SimpleIteratorTaskBehavior) Error(context model.TaskContext) (handled bool, taskEntry *model.TaskEntry) {

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
func (tb *SimpleIteratorTaskBehavior) ChildDone(context model.TaskContext, childTask *definition.Task, childDoneCode int) (done bool, doneCode int) {

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

///////////////////////////////////
// Iterators

type Iterator interface {
	Key() interface{}
	Value() interface{}
	next() bool
}

type ArrayIterator struct {
	current int
	data    []interface{}
}

func (itx *ArrayIterator) Key() interface{} {
	return itx.current
}

func (itx *ArrayIterator) Value() interface{} {
	return itx.data[itx.current]
}
func (itx *ArrayIterator) next() bool {
	itx.current++
	if itx.current >= len(itx.data) {
		return false
	}
	return true
}

func NewArrayIterator(data []interface{}) *ArrayIterator {
	return &ArrayIterator{data: data, current: -1}
}

type IntIterator struct {
	current int
	count   int
}

func (itx *IntIterator) Key() interface{} {
	return itx.current
}

func (itx *IntIterator) Value() interface{} {
	return itx.current
}

func (itx *IntIterator) next() bool {
	itx.current++
	if itx.current >= itx.count {
		return false
	}
	return true
}

func NewIntIterator(count int) *IntIterator {
	return &IntIterator{count: count, current: -1}
}

type ObjectIterator struct {
	current int
	keyMap  map[int]string
	data    map[string]interface{}
}

func (itx *ObjectIterator) Key() interface{} {
	return itx.keyMap[itx.current]
}

func (itx *ObjectIterator) Value() interface{} {
	key := itx.keyMap[itx.current]
	return itx.data[key]
}

func (itx *ObjectIterator) next() bool {
	itx.current++
	if itx.current >= len(itx.data) {
		return false
	}
	return true
}

func NewObjectIterator(data map[string]interface{}) *ObjectIterator {
	keyMap := make(map[int]string, len(data))
	i := 0
	for key := range data {
		keyMap[i] = key
		i++
	}

	return &ObjectIterator{keyMap: keyMap, data: data, current: -1}
}
