package engine

import (
	"os"
	"strconv"

	"github.com/TIBCOSoftware/flogo-lib/engine/runner"
)

const (
	ENV_RUNNER_TYPE_KEY       = "FLOGO_RUNNER_TYPE"
	RUNNER_TYPE_DEFAULT       = "POOLED"
	ENV_RUNNER_WORKERS_KEY    = "FLOGO_RUNNER_WORKERS"
	RUNNER_WORKERS_DEFAULT    = 5
	ENV_RUNNER_QUEUE_SIZE_KEY = "FLOGO_RUNNER_QUEUE"
	RUNNER_QUEUE_SIZE_DEFAULT = 50
)

//GetRunnerType returns the runner type
func GetRunnerType() string {
	runnerTypeEnv := os.Getenv(ENV_RUNNER_TYPE_KEY)
	if len(runnerTypeEnv) > 0 {
		return runnerTypeEnv
	}
	return RUNNER_TYPE_DEFAULT
}

//GetRunnerWorkers returns the number of workers to use
func GetRunnerWorkers() int {
	numWorkers := RUNNER_WORKERS_DEFAULT
	workersEnv := os.Getenv(ENV_RUNNER_WORKERS_KEY)
	if len(workersEnv) > 0 {
		i, err := strconv.Atoi(workersEnv)
		if err == nil {
			numWorkers = i
		}
	}
	return numWorkers
}

//GetRunnerQueueSize returns the runner queue size
func GetRunnerQueueSize() int {
	queueSize := RUNNER_QUEUE_SIZE_DEFAULT
	queueSizeEnv := os.Getenv(ENV_RUNNER_QUEUE_SIZE_KEY)
	if len(queueSizeEnv) > 0 {
		i, err := strconv.Atoi(queueSizeEnv)
		if err == nil {
			queueSize = i
		}
	}
	return queueSize
}

//NewPooledRunnerConfig creates a new Pooled config, looks for environment variables to override default values
func NewPooledRunnerConfig() *runner.PooledConfig {
	return &runner.PooledConfig{NumWorkers: GetRunnerWorkers(), WorkQueueSize: GetRunnerQueueSize()}
}
