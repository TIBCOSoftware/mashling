package logger

import (
	"errors"
	"sync"

	"github.com/Sirupsen/logrus"
)

var logrusHooks sync.Map

type LogHook interface {
	GetHook() (logrus.Hook, error)
	Id() (string, error)
}

// Initialize sets up log hooks based off of type.
func Initialize(hookType string, settings map[string]interface{}) (hook LogHook, err error) {
	switch hookType {
	case "kafka":
		return InitializeKafkaHook(settings)
	default:
		return nil, errors.New("unknown log hook type")
	}
}
