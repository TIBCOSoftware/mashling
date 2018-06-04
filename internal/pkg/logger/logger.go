package logger

import (
	"fmt"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/TIBCOSoftware/flogo-lib/config"
	flogger "github.com/TIBCOSoftware/flogo-lib/logger"
)

var loggerMap = make(map[string]flogger.Logger)
var mutex = &sync.RWMutex{}
var mashlingLoggerName = "mashling"

type MashlingLoggerFactory struct {
}

type MashlingLogger struct {
	loggerName string
	loggerImpl *logrus.Logger
}

type LogFormatter struct {
	loggerName string
}

func Register() {
	flogger.RegisterLoggerFactory(&MashlingLoggerFactory{})
}

func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	logEntry := fmt.Sprintf("LULZ %s %-6s [%s] - %s\n", entry.Time.Format("2006-01-02 15:04:05.000"), getLevel(entry.Level), f.loggerName, entry.Message)
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
func (logger *MashlingLogger) Debug(args ...interface{}) {
	logger.loggerImpl.Debug(args...)
}

// DebugEnabled checks if Debug level is enabled.
func (logger *MashlingLogger) DebugEnabled() bool {
	return logger.loggerImpl.Level >= logrus.DebugLevel
}

// Info logs message at Info level.
func (logger *MashlingLogger) Info(args ...interface{}) {
	logger.loggerImpl.Info(args...)
}

// InfoEnabled checks if Info level is enabled.
func (logger *MashlingLogger) InfoEnabled() bool {
	return logger.loggerImpl.Level >= logrus.InfoLevel
}

// Warn logs message at Warning level.
func (logger *MashlingLogger) Warn(args ...interface{}) {
	logger.loggerImpl.Warn(args...)
}

// WarnEnabled checks if Warning level is enabled.
func (logger *MashlingLogger) WarnEnabled() bool {
	return logger.loggerImpl.Level >= logrus.WarnLevel
}

// Error logs message at Error level.
func (logger *MashlingLogger) Error(args ...interface{}) {
	logger.loggerImpl.Error(args...)
}

// ErrorEnabled checks if Error level is enabled.
func (logger *MashlingLogger) ErrorEnabled() bool {
	return logger.loggerImpl.Level >= logrus.ErrorLevel
}

// Debugf logs message at Debug level.
func (logger *MashlingLogger) Debugf(format string, args ...interface{}) {
	logger.loggerImpl.Debugf(format, args...)
}

// Infof logs message at Info level.
func (logger *MashlingLogger) Infof(format string, args ...interface{}) {
	logger.loggerImpl.Infof(format, args...)
}

// Warnf logs message at Warning level.
func (logger *MashlingLogger) Warnf(format string, args ...interface{}) {
	logger.loggerImpl.Warnf(format, args...)
}

// Errorf logs message at Error level.
func (logger *MashlingLogger) Errorf(format string, args ...interface{}) {
	logger.loggerImpl.Errorf(format, args...)
}

//SetLogLevel sets the log level to be compliant with Flogo
func (logger *MashlingLogger) SetLogLevel(logLevel flogger.Level) {
	switch logLevel {
	case flogger.DebugLevel:
		logger.loggerImpl.Level = logrus.DebugLevel
	case flogger.InfoLevel:
		logger.loggerImpl.Level = logrus.InfoLevel
	case flogger.ErrorLevel:
		logger.loggerImpl.Level = logrus.ErrorLevel
	case flogger.WarnLevel:
		logger.loggerImpl.Level = logrus.WarnLevel
	default:
		logger.loggerImpl.Level = logrus.ErrorLevel
	}
}

func (logfactory *MashlingLoggerFactory) GetLogger(name string) flogger.Logger {
	mutex.RLock()
	l := loggerMap[name]
	mutex.RUnlock()
	if l == nil {
		logImpl := logrus.New()
		logImpl.Formatter = &LogFormatter{
			loggerName: name,
		}
		l = &MashlingLogger{
			loggerName: name,
			loggerImpl: logImpl,
		}
		// Get log level from config
		logLevelName := config.GetLogLevel()
		// Get log level for name
		level, err := flogger.GetLevelForName(logLevelName)
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
