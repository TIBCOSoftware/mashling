package logger

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	flogger "github.com/TIBCOSoftware/flogo-lib/logger"
)

var logFactory = &MashlingLoggerFactory{}
var loggerMap sync.Map
var mashlingLoggerName = "mashling"
var logLevel = flogger.InfoLevel

type MashlingLoggerFactory struct {
}

type MashlingLogger struct {
	loggerName string
	loggerImpl *logrus.Logger
}

type LogFormatter struct {
	loggerName string
}

func init() {
	flogger.RegisterLoggerFactory(logFactory)
}

func Configure(levelName string, hooks []LogHook) error {
	level, err := flogger.GetLevelForName(strings.ToUpper(levelName))
	if err != nil {
		return nil
	}
	SetLogLevel(level)
	for _, hook := range hooks {
		logrusHook, err := hook.GetHook()
		if err != nil {
			return err
		}
		logrusHookID, err := hook.Id()
		if err != nil {
			return err
		}
		logrusHooks.Store(logrusHookID, logrusHook)
		loggerMap.Range(func(key, value interface{}) bool {
			logr, ok := value.(*MashlingLogger)
			if !ok {
				return false
			}
			logr.loggerImpl.AddHook(logrusHook)
			return true
		})
	}
	return nil
}

func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	logEntry := fmt.Sprintf("%s %-6s [%s] - %s\n", entry.Time.Format("2006-01-02 15:04:05.000"), getLevel(entry.Level), f.loggerName, entry.Message)
	return []byte(logEntry), nil
}

func Debug(args ...interface{}) {
	GetMashlingLogger().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	GetMashlingLogger().Debugf(format, args...)
}

func Info(args ...interface{}) {
	GetMashlingLogger().Info(args...)
}

func Infof(format string, args ...interface{}) {
	GetMashlingLogger().Infof(format, args...)
}

func Warn(args ...interface{}) {
	GetMashlingLogger().Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	GetMashlingLogger().Warnf(format, args...)
}

func Error(args ...interface{}) {
	GetMashlingLogger().Error(args...)
}

func Errorf(format string, args ...interface{}) {
	GetMashlingLogger().Errorf(format, args...)
}

func SetLogLevel(level flogger.Level) {
	logLevel = level
	loggerMap.Range(func(key, value interface{}) bool {
		logr, ok := value.(*MashlingLogger)
		if !ok {
			// skip logger
			return true
		}
		logr.SetLogLevel(logLevel)
		return true
	})
}

func GetMashlingLogger() flogger.Logger {
	defLogger := GetLogger(mashlingLoggerName)
	if defLogger == nil {
		errorMsg := fmt.Sprintf("error getting Mashling logger '%s'", mashlingLoggerName)
		panic(errorMsg)
	}
	return defLogger
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
	lStored, exists := loggerMap.Load(name)
	if !exists {
		logImpl := logrus.New()
		logImpl.Formatter = &LogFormatter{
			loggerName: name,
		}
		l := &MashlingLogger{
			loggerName: name,
			loggerImpl: logImpl,
		}
		l.SetLogLevel(logLevel)
		// Add hooks
		logrusHooks.Range(func(key, value interface{}) bool {
			hook, ok := value.(logrus.Hook)
			if !ok {
				return false
			}
			l.loggerImpl.AddHook(hook)
			return true
		})
		loggerMap.Store(name, l)
		return l
	}
	return lStored.(flogger.Logger)
}

func GetLogger(name string) flogger.Logger {
	return logFactory.GetLogger(name)
}
