package aggregate

import (
	"errors"
	"sync"

	"github.com/TIBCOSoftware/flogo-contrib/activity/aggregate/aggregator"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// activityLogger is the default logger for the Aggregate Activity
var activityLogger = logger.GetLogger("activity-tibco-aggregate")

const (
	ivFunction   = "function"
	ivWindowSize = "windowSize"
	ivValue      = "value"

	ovResult = "result"
	ovReport = "report"
)

func init() {
	activityLogger.SetLogLevel(logger.InfoLevel)
}

// AggregateActivity is an Activity that is used to Aggregate a message to the console
// inputs : {function, windowSize, autoRest, value}
// outputs: {result, report}
type AggregateActivity struct {
	metadata *activity.Metadata
	mutex    *sync.RWMutex

	// aggregators stateful map of aggregators
	aggregators map[string]aggregator.Aggregator
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &AggregateActivity{metadata: metadata, aggregators: make(map[string]aggregator.Aggregator), mutex: &sync.RWMutex{}}
}

// Metadata returns the activity's metadata
func (a *AggregateActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Aggregates the Message
func (a *AggregateActivity) Eval(context activity.Context) (done bool, err error) {

	aggregatorKey := context.ActivityHost().Name() + ":" + context.Name()

	a.mutex.RLock()
	//get aggregator for activity, assumes flow & task names are unique
	aggr, ok := a.aggregators[aggregatorKey]

	a.mutex.RUnlock()

	//if window not create for this flow, create it

	if !ok {

		//go doesn't have lock upgrades or try, so do same check again

		a.mutex.Lock()
		aggr, ok = a.aggregators[aggregatorKey]

		if !ok {
			windowSize, _ := context.GetInput(ivWindowSize).(int)
			aggrName, _ := context.GetInput(ivFunction).(string)

			factory := aggregator.GetFactory(aggrName)

			if factory == nil {
				return false, errors.New("Unknown aggregator: " + aggrName)
			}

			aggr = factory(windowSize)
			a.aggregators[aggregatorKey] = aggr

			activityLogger.Debug("Aggregator created for ", aggregatorKey)
		}

		a.mutex.Unlock()
	}

	value, ok := context.GetInput(ivValue).(float64)

	if !ok {
		value, _ = data.CoerceToNumber(context.GetInput(ivValue))
	}

	report, result := aggr.Add(value)

	context.SetOutput(ovReport, report)
	context.SetOutput(ovResult, result)

	return true, nil
}
