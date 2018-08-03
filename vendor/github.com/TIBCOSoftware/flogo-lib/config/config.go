// This package contains all the default values for the configuration
package config

import (
	"os"
	"strconv"
)

const (
	ENV_LOG_DATE_FORMAT_KEY      = "FLOGO_LOG_DTFORMAT"
	LOG_DATE_FORMAT_DEFAULT      = "2006-01-02 15:04:05.000"
	ENV_LOG_LEVEL_KEY            = "FLOGO_LOG_LEVEL"
	LOG_LEVEL_DEFAULT            = "INFO"
	ENV_APP_CONFIG_LOCATION_KEY  = "FLOGO_CONFIG_PATH"
	APP_CONFIG_LOCATION_DEFAULT  = "flogo.json"
	ENV_STOP_ENGINE_ON_ERROR_KEY = "FLOGO_ENGINE_STOP_ON_ERROR"
)

var defaultLogLevel = LOG_LEVEL_DEFAULT

//GetFlogoConfigPath returns the flogo config path
func GetFlogoConfigPath() string {
	flogoConfigPathEnv := os.Getenv(ENV_APP_CONFIG_LOCATION_KEY)
	if len(flogoConfigPathEnv) > 0 {
		return flogoConfigPathEnv
	}
	return APP_CONFIG_LOCATION_DEFAULT
}

func SetDefaultLogLevel(logLevel string) {
	defaultLogLevel = logLevel
}

//GetLogLevel returns the log level
func GetLogLevel() string {
	logLevelEnv := os.Getenv(ENV_LOG_LEVEL_KEY)
	if len(logLevelEnv) > 0 {
		return logLevelEnv
	}
	return defaultLogLevel
}

func GetLogDateTimeFormat() string {
	logLevelEnv := os.Getenv(ENV_LOG_DATE_FORMAT_KEY)
	if len(logLevelEnv) > 0 {
		return logLevelEnv
	}
	return LOG_DATE_FORMAT_DEFAULT
}

func StopEngineOnError() bool {
	stopEngineOnError := os.Getenv(ENV_STOP_ENGINE_ON_ERROR_KEY)
	if len(stopEngineOnError) == 0 {
		return true
	}
	b, _ := strconv.ParseBool(stopEngineOnError)
	return b
}
