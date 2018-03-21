// This package contains all the default values for the configuration
package config

import (
	"os"
	"strconv"
)

const (
	LOG_DATE_FORMAT_KEY         = "FLOGO_LOG_DTFORMAT"
	LOG_DATE_FORMAT_DEFAULT     = "2006-01-02 15:04:05.000"
	LOG_LEVEL_KEY               = "FLOGO_LOG_LEVEL"
	LOG_LEVEL_DEFAULT           = "INFO"
	RUNNER_TYPE_KEY             = "FLOGO_RUNNER_TYPE"
	RUNNER_TYPE_DEFAULT         = "POOLED"
	RUNNER_WORKERS_KEY          = "FLOGO_RUNNER_WORKERS"
	RUNNER_WORKERS_DEFAULT      = 5
	RUNNER_QUEUE_SIZE_KEY       = "FLOGO_RUNNER_QUEUE"
	RUNNER_QUEUE_SIZE_DEFAULT   = 50
	APP_CONFIG_LOCATION_KEY     = "FLOGO_CONFIG_PATH"
	APP_CONFIG_LOCATION_DEFAULT = "flogo.json"
	STOP_ENGINE_ON_ERROR_KEY    = "STOP_ENGINE_ON_ERROR"
)

var defaultLogLevel = LOG_LEVEL_DEFAULT

//GetFlogoConfigPath returns the flogo config path
func GetFlogoConfigPath() string {
	flogoConfigPathEnv := os.Getenv(APP_CONFIG_LOCATION_KEY)
	if len(flogoConfigPathEnv) > 0 {
		return flogoConfigPathEnv
	}
	return APP_CONFIG_LOCATION_DEFAULT
}

//GetRunnerType returns the runner type
func GetRunnerType() string {
	runnerTypeEnv := os.Getenv(RUNNER_TYPE_KEY)
	if len(runnerTypeEnv) > 0 {
		return runnerTypeEnv
	}
	return RUNNER_TYPE_DEFAULT
}

//GetRunnerWorkers returns the number of workers to use
func GetRunnerWorkers() int {
	numWorkers := RUNNER_WORKERS_DEFAULT
	workersEnv := os.Getenv(RUNNER_WORKERS_KEY)
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
	queueSizeEnv := os.Getenv(RUNNER_QUEUE_SIZE_KEY)
	if len(queueSizeEnv) > 0 {
		i, err := strconv.Atoi(queueSizeEnv)
		if err == nil {
			queueSize = i
		}
	}
	return queueSize
}

func SetDefaultLogLevel(logLevel string) {
	defaultLogLevel = logLevel
}

//GetLogLevel returns the log level
func GetLogLevel() string {
	logLevelEnv := os.Getenv(LOG_LEVEL_KEY)
	if len(logLevelEnv) > 0 {
		return logLevelEnv
	}
	return defaultLogLevel
}

func GetLogDateTimeFormat() string {
	logLevelEnv := os.Getenv(LOG_DATE_FORMAT_KEY)
	if len(logLevelEnv) > 0 {
		return logLevelEnv
	}
	return LOG_DATE_FORMAT_DEFAULT
}

func StopEngineOnError() bool {
	stopEngineOnError := os.Getenv(STOP_ENGINE_ON_ERROR_KEY)
	if len(stopEngineOnError) == 0 {
		return true
	}
	b, _ := strconv.ParseBool(stopEngineOnError)
	return b
}
