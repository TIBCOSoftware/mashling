package runner

import (
	"context"
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Based off: http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html

// RequestType is value that indicates the type of Request
type RequestType int

const (
	// RtRun denotes a run action request
	RtRun RequestType = 10
)

// ActionWorkRequest describes a Request that Worker should handle
type ActionWorkRequest struct {
	ReqType    RequestType
	ID         string
	actionData *ActionData
}

// ActionData action related data to pass along in a ActionWorkRequest
type ActionData struct {
	context        context.Context
	action         action.Action
	options        map[string]interface{}
	arc            chan (*ActionResult)
}

// ActionResult is a simple struct to hold the results for an Action
type ActionResult struct {
	results map[string]*data.Attribute
	err     error
}

// A ActionWorker handles WorkRequest, work requests consist of start, restart
// and resume of FlowInstances
type ActionWorker struct {
	ID          int
	runner      *DirectRunner
	Work        chan ActionWorkRequest
	WorkerQueue chan chan ActionWorkRequest
	QuitChan    chan bool
}

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, runner *DirectRunner, workerQueue chan chan ActionWorkRequest) ActionWorker {
	// Create, and return the worker.
	worker := ActionWorker{
		ID:          id,
		runner:      runner,
		Work:        make(chan ActionWorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool)}

	return worker
}

// Start function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.  This is where all the request are handled
func (w ActionWorker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue.
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				// Receive a work request.
				logger.Debugf("worker-%d: Received Request\n", w.ID)

				switch work.ReqType {
				default:

					err := fmt.Errorf("Unsupported work request type: '%d'", work.ReqType)
					actionData := work.actionData
					actionData.arc <- &ActionResult{err: err}

				case RtRun:

					actionData := work.actionData

					handler := &AsyncResultHandler{result: make(chan *ActionResult), done: make(chan bool, 1)}

					act := actionData.action

					var ctxData *trigger.ContextData

					if actionData.context != nil {
						var exists bool
						ctxData, exists = trigger.ExtractContextData(actionData.context)

						if !exists {
							logger.Warn("Trigger data not applied to context")
						}
					}

					inputs := generateInputs(act, ctxData)

					err := act.Run(actionData.context, inputs, actionData.options, handler)

					if err != nil {
						logger.Debugf("worker-%d: Action Run error: %s\n", w.ID, err.Error())
						// error so just return
						actionData.arc <- &ActionResult{err: err}
					} else {
						done := false
						//wait for reply
						for !done {
							select {
							case result := <-handler.result:
								logger.Debugf("*** Worker received result: %v\n", result)
								result.results = generateOutputs(act, ctxData, result.results)
								actionData.arc <- result
							case <-handler.done:
								if !handler.replied {
									actionData.arc <- &ActionResult{}
								}
								done = true
							}
						}
					}

					logger.Debugf("worker-%d: Completed Request\n", w.ID)
				}

			case <-w.QuitChan:
				// We have been asked to stop.
				logger.Debugf("worker-%d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (w ActionWorker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

// AsyncResultHandler simple ResultHandler to use in the asynchronous case
type AsyncResultHandler struct {
	done    chan (bool)
	result  chan (*ActionResult)
	replied bool
}

// HandleResult implements action.ResultHandler.HandleResult
func (rh *AsyncResultHandler) HandleResult(results map[string]*data.Attribute, err error) {
	rh.replied = true
	rh.result <- &ActionResult{results: results, err: err}
}

// Done implements action.ResultHandler.Done
func (rh *AsyncResultHandler) Done() {
	rh.done <- true
}
