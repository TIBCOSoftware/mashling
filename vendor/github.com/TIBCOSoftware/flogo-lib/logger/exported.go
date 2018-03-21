package logger

import (
	"fmt"
)

func Debug(args ...interface{}) {
	GetDefaultLogger().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	GetDefaultLogger().Debugf(format, args...)
}

func Info(args ...interface{}) {
	GetDefaultLogger().Info(args...)
}

func Infof(format string, args ...interface{}) {
	GetDefaultLogger().Infof(format, args...)
}

func Warn(args ...interface{}) {
	GetDefaultLogger().Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	GetDefaultLogger().Warnf(format, args...)
}

func Error(args ...interface{}) {
	GetDefaultLogger().Error(args...)
}

func Errorf(format string, args ...interface{}) {
	GetDefaultLogger().Errorf(format, args...)
}

func SetLogLevel(level Level) {
	GetDefaultLogger().SetLogLevel(level)
}

func GetDefaultLogger() Logger {
	defLogger := GetLogger("engine")
	if defLogger == nil {
		errorMsg := fmt.Sprintf("Engine: Error Getting engine Logger null")
		panic(errorMsg)
	}
	return defLogger
}
