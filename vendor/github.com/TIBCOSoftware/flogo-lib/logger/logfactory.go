package logger

import (
	"fmt"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/TIBCOSoftware/flogo-lib/config"
)

var loggerMap = make(map[string]Logger)
var mutex = &sync.RWMutex{}

type DefaultLoggerFactory struct {
}

func init() {
	RegisterLoggerFactory(&DefaultLoggerFactory{})
}

type DefaultLogger struct {
	loggerName string
	loggerImpl *logrus.Logger
}

type LogFormatter struct {
	loggerName string
}

func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	logEntry := fmt.Sprintf("%s %-6s [%s] - %s\n", entry.Time.Format(config.GetLogDateTimeFormat()), getLevel(entry.Level), f.loggerName, entry.Message)
	return []byte(logEntry), nil
}

func getLevel(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return "DEBUG"
	case logrus.InfoLevel:
		return "INFO"
	case logrus.ErrorLevel:
		return "ERROR"
	case logrus.WarnLevel:
		return "WARN"
	case logrus.PanicLevel:
		return "PANIC"
	case logrus.FatalLevel:
		return "FATAL"
	}

	return "UNKNOWN"
}

// Debug logs message at Debug level.
func (logger *DefaultLogger) Debug(args ...interface{}) {
	logger.loggerImpl.Debug(args...)
}

// DebugEnabled checks if Debug level is enabled.
func (logger *DefaultLogger) DebugEnabled() bool {
	return logger.loggerImpl.Level >= logrus.DebugLevel
}

// Info logs message at Info level.
func (logger *DefaultLogger) Info(args ...interface{}) {
	logger.loggerImpl.Info(args...)
}

// InfoEnabled checks if Info level is enabled.
func (logger *DefaultLogger) InfoEnabled() bool {
	return logger.loggerImpl.Level >= logrus.InfoLevel
}

// Warn logs message at Warning level.
func (logger *DefaultLogger) Warn(args ...interface{}) {
	logger.loggerImpl.Warn(args...)
}

// WarnEnabled checks if Warning level is enabled.
func (logger *DefaultLogger) WarnEnabled() bool {
	return logger.loggerImpl.Level >= logrus.WarnLevel
}

// Error logs message at Error level.
func (logger *DefaultLogger) Error(args ...interface{}) {
	logger.loggerImpl.Error(args...)
}

// ErrorEnabled checks if Error level is enabled.
func (logger *DefaultLogger) ErrorEnabled() bool {
	return logger.loggerImpl.Level >= logrus.ErrorLevel
}

// Debug logs message at Debug level.
func (logger *DefaultLogger) Debugf(format string, args ...interface{}) {
	logger.loggerImpl.Debugf(format, args...)
}

// Info logs message at Info level.
func (logger *DefaultLogger) Infof(format string, args ...interface{}) {
	logger.loggerImpl.Infof(format, args...)
}

// Warn logs message at Warning level.
func (logger *DefaultLogger) Warnf(format string, args ...interface{}) {
	logger.loggerImpl.Warnf(format, args...)
}

// Error logs message at Error level.
func (logger *DefaultLogger) Errorf(format string, args ...interface{}) {
	logger.loggerImpl.Errorf(format, args...)
}

//SetLog Level
func (logger *DefaultLogger) SetLogLevel(logLevel Level) {
	switch logLevel {
	case DebugLevel:
		logger.loggerImpl.Level = logrus.DebugLevel
	case InfoLevel:
		logger.loggerImpl.Level = logrus.InfoLevel
	case ErrorLevel:
		logger.loggerImpl.Level = logrus.ErrorLevel
	case WarnLevel:
		logger.loggerImpl.Level = logrus.WarnLevel
	default:
		logger.loggerImpl.Level = logrus.ErrorLevel
	}
}

func (logfactory *DefaultLoggerFactory) GetLogger(name string) Logger {
	mutex.RLock()
	l := loggerMap[name]
	mutex.RUnlock()
	if l == nil {
		logImpl := logrus.New()
		logImpl.Formatter = &LogFormatter{
			loggerName: name,
		}
		l = &DefaultLogger{
			loggerName: name,
			loggerImpl: logImpl,
		}
		// Get log level from config
		logLevelName := config.GetLogLevel()
		// Get log level for name
		level, err := GetLevelForName(logLevelName)
		if err != nil {
			return nil
		}
		l.SetLogLevel(level)
		mutex.Lock()
		loggerMap[name] = l
		mutex.Unlock()
	}
	return l
}
