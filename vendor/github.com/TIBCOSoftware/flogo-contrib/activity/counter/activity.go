package counter

import (
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("activity-tibco-counter")

const (
	ivCounterName = "counterName"
	ivIncrement   = "increment"
	ivReset       = "reset"

	ovValue = "value"
)

// CounterActivity is a Counter Activity implementation
type CounterActivity struct {
	sync.Mutex
	metadata *activity.Metadata
	counters map[string]int
}

// NewActivity creates a new CounterActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &CounterActivity{metadata: metadata, counters: make(map[string]int)}
}

// Metadata implements activity.Activity.Metadata
func (a *CounterActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *CounterActivity) Eval(context activity.Context) (done bool, err error) {

	counterName := context.GetInput(ivCounterName).(string)

	var increment, reset bool

	if context.GetInput(ivIncrement) != nil {
		increment = context.GetInput(ivIncrement).(bool)
	}
	if context.GetInput(ivReset) != nil {
		reset = context.GetInput(ivReset).(bool)
	}

	var count int

	if reset {
		count = a.resetCounter(counterName)

		log.Debugf("Counter [%s] reset", counterName)
	} else if increment {
		count = a.incrementCounter(counterName)

		log.Debugf("Counter [%s] incremented: %d", counterName, count)
	} else {
		count = a.getCounter(counterName)

		log.Debugf("Counter [%s] = %d", counterName, count)
	}

	context.SetOutput(ovValue, count)

	return true, nil
}

func (a *CounterActivity) incrementCounter(counterName string) int {
	a.Lock()
	defer a.Unlock()

	count := 1

	if counter, exists := a.counters[counterName]; exists {
		count = counter + 1
	}

	a.counters[counterName] = count

	return count
}

func (a *CounterActivity) resetCounter(counterName string) int {
	a.Lock()
	defer a.Unlock()

	if _, exists := a.counters[counterName]; exists {
		a.counters[counterName] = 0
	}

	return 0
}

func (a *CounterActivity) getCounter(counterName string) int {
	a.Lock()
	defer a.Unlock()

	return a.counters[counterName]
}
