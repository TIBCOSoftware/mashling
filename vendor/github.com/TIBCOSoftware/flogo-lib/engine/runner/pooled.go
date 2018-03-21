package runner

import (
	"context"
	"errors"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// PooledRunner is a action runner that queues and runs a action in a worker pool
type PooledRunner struct {
	workerQueue chan chan ActionWorkRequest
	workQueue   chan ActionWorkRequest
	numWorkers  int
	workers     []*ActionWorker
	active      bool

	directRunner *DirectRunner
}

// PooledConfig is the configuration object for a PooledRunner
type PooledConfig struct {
	NumWorkers    int `json:"numWorkers"`
	WorkQueueSize int `json:"workQueueSize"`
}

// NewPooledRunner create a new pooled
func NewPooled(config *PooledConfig) *PooledRunner {

	var pooledRunner PooledRunner
	pooledRunner.directRunner = NewDirect()

	// config via engine config
	pooledRunner.numWorkers = config.NumWorkers
	pooledRunner.workQueue = make(chan ActionWorkRequest, config.WorkQueueSize)

	return &pooledRunner
}

// Start will start the engine, by starting all of its workers
func (runner *PooledRunner) Start() error {

	if !runner.active {

		runner.workerQueue = make(chan chan ActionWorkRequest, runner.numWorkers)

		runner.workers = make([]*ActionWorker, runner.numWorkers)

		for i := 0; i < runner.numWorkers; i++ {
			id := i + 1
			logger.Debugf("Starting worker with id '%d'", id)
			worker := NewWorker(id, runner.directRunner, runner.workerQueue)
			runner.workers[i] = &worker
			worker.Start()
		}

		go func() {
			for {
				select {
				case work := <-runner.workQueue:
					logger.Debug("Received work request")

					//todo fix, this creates unbounded go routines waiting to be serviced by worker queue
					go func() {
						worker := <-runner.workerQueue

						logger.Debug("Dispatching work request")
						worker <- work
					}()
				}
			}
		}()

		runner.active = true
	}

	return nil
}

// Stop will stop the engine, by stopping all of its workers
func (runner *PooledRunner) Stop() error {

	if runner.active {

		runner.active = false

		for _, worker := range runner.workers {
			logger.Debug("Stopping worker", worker.ID)
			worker.Stop()
		}
	}

	return nil
}

//Deprecated
func (runner *PooledRunner) Run(ctx context.Context, act action.Action, uri string, options interface{}) (code int, data interface{}, err error) {

	if act == nil {
		return 0, nil, errors.New("Action not specified")
	}

	newOptions := make(map[string]interface{})
	newOptions["deprecated_options"] = options

	if runner.active {

		actionData := &ActionData{context: ctx, action: act, options: newOptions, arc: make(chan *ActionResult, 1)}
		work := ActionWorkRequest{ReqType: RtRun, actionData: actionData}

		runner.workQueue <- work
		reply := <-actionData.arc

		ndata := reply.results
		err := reply.err
		//return reply.results, reply.err

		if len(ndata) != 0 {
			defData, ok := ndata["data"]
			if ok {
				data = defData.Value()
			}
			defCode, ok := ndata["code"]
			if ok && defCode.Value() != nil {
				code = defCode.Value().(int)
			}
		}

		return code, data, err
	}

	return 0, nil, errors.New("Runner not active")
}

// Run implements action.Runner.Run
func (runner *PooledRunner) RunAction(ctx context.Context, act action.Action, options map[string]interface{}) (results map[string]*data.Attribute, err error) {

	if act == nil {
		return nil, errors.New("Action not specified")
	}

	if runner.active {

		data := &ActionData{context: ctx, action: act, options: options, arc: make(chan *ActionResult, 1)}
		work := ActionWorkRequest{ReqType: RtRun, actionData: data}

		runner.workQueue <- work
		logger.Debugf("Run Action '%s' queued", act.Config().Id)

		reply := <-data.arc
		logger.Debugf("Run Action '%s' complete", act.Config().Id)

		return reply.results, reply.err
	}

	//Run rejected
	return nil, errors.New("Runner not active")
}
