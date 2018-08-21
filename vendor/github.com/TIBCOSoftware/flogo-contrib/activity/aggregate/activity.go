package aggregate

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/TIBCOSoftware/flogo-contrib/activity/aggregate/window"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/flogo-oss/stream/pipeline/support"
)

// activityLogger is the default logger for the Aggregate Activity
var activityLogger = logger.GetLogger("activity-aggregate")

const (
	sFunction           = "function"
	sWindowType         = "windowType"
	sWindowSize         = "windowSize"
	sResolution         = "resolution"
	sProceedOnlyOnEmit  = "proceedOnlyOnEmit"
	sAdditionalSettings = "additionalSettings"

	ivValue = "value"

	ovResult = "result"
	ovReport = "report"
)

//we can generate json from this! - we could also create a "validate-able" object from this
type Settings struct {
	Function           string `md:"required,allowed(avg,sum,min,max,count)"`
	WindowType         string `md:"required,allowed(tumbling,sliding,timeTumbling,timeSliding)"`
	WindowSize         int    `md:"required"`
	ProceedOnlyOnEmit  bool
	Resolution         int
	AdditionalSettings map[string]string
}


func init() {
	activityLogger.SetLogLevel(logger.InfoLevel)
}

var metadata *activity.Metadata

func New(config *activity.Config) (activity.Activity, error) {
	act := &AggregateActivity{}

	//todo implement
	//config.Settings

	return act, nil
}

// AggregateActivity is an Activity that is used to Aggregate a message to the console
type AggregateActivity struct {
	settings *Settings
	mutex    *sync.RWMutex
}

// NewActivity creates a new AppActivity
func NewActivity(md *activity.Metadata) activity.Activity {
	metadata = md
	activity.RegisterFactory(md.ID, New)
	return &AggregateActivity{mutex: &sync.RWMutex{}}
}

// Metadata returns the activity's metadata
func (a *AggregateActivity) Metadata() *activity.Metadata {
	return metadata
}

// Eval implements api.Activity.Eval - Aggregates the Message
func (a *AggregateActivity) Eval(ctx activity.Context) (done bool, err error) {

	//todo move to Activity instance creation
	settings, err := getSettings(ctx)
	if err != nil {
		return false, err
	}

	ss, ok := activity.GetSharedTempDataSupport(ctx)
	if !ok {
		return false, fmt.Errorf("AggregateActivity not supported by this activity host")
	}

	sharedData := ss.GetSharedTempData()
	wv, defined := sharedData["window"]

	timerSupport, timerSupported := support.GetTimerSupport(ctx)

	var w window.Window

	if !defined {
		//create the window & associated timer if necessary

		windowSettings := &window.Settings{Size: settings.WindowSize, ExternalTimer: timerSupported, Resolution: settings.Resolution}
		windowSettings.SetAdditionalSettings(settings.AdditionalSettings)

		timerSupport, timerSupported := support.GetTimerSupport(ctx)
		wType := strings.ToLower(settings.WindowType)

		switch wType {
		case "tumbling":
			w, err = NewTumblingWindow(settings.Function, windowSettings)
		case "sliding":
			w, err = NewSlidingWindow(settings.Function, windowSettings)
		case "timetumbling":
			w, err = NewTumblingTimeWindow(settings.Function, windowSettings)
			if timerSupported {
				timerSupport.CreateTimer(time.Duration(settings.WindowSize)*time.Millisecond, moveWindow, true)
			}
		case "timesliding":
			w, err = NewSlidingTimeWindow(settings.Function, windowSettings)
			if timerSupported {
				timerSupport.CreateTimer(time.Duration(settings.Resolution)*time.Millisecond, moveWindow, true)
			}
		default:
			return false, fmt.Errorf("unsupported window type: '%s'", settings.WindowType)
		}

		sharedData["window"] = w
	} else {
		w = wv.(window.Window)
	}

	in := ctx.GetInput(ivValue)

	emit, result := w.AddSample(in)

	if timerSupported {
		timerSupport.UpdateTimer(true)
	}

	ctx.SetOutput(ovResult, result)
	ctx.SetOutput(ovReport, emit)

	done = !(settings.ProceedOnlyOnEmit && !emit)

	return done, nil
}

func (a *AggregateActivity) PostEval(ctx activity.Context, userData interface{}) (done bool, err error) {
	return true, nil
}

func moveWindow(ctx activity.Context) bool {

	ss, _ := activity.GetSharedTempDataSupport(ctx)
	sharedData := ss.GetSharedTempData()

	wv, _ := sharedData["window"]

	w, _ := wv.(window.TimeWindow)

	emit, result := w.NextBlock()

	ctx.SetOutput(ovResult, result)
	ctx.SetOutput(ovReport, emit)

	poe := true // by default only proceed on emit
	poeSetting, exists := ctx.GetSetting(sProceedOnlyOnEmit)
	if exists {
		poe, _ = data.CoerceToBoolean(poeSetting)
	}

	return !(poe && !emit)
}

func getSettings(ctx activity.Context) (*Settings, error) {

	settings := &Settings{}

	settings.Function = "avg" // default function
	setting, exists := ctx.GetSetting(sFunction)
	if exists {
		val, err := data.CoerceToString(setting)
		if err == nil {
			settings.Function = val
		}
	}

	settings.WindowType = "tumbling" // default window type
	setting, exists = ctx.GetSetting(sWindowType)
	if exists {
		val, err := data.CoerceToString(setting)
		if err == nil {
			settings.WindowType = val
		}
	}

	settings.WindowSize = 5 // default window resolution
	setting, exists = ctx.GetSetting(sWindowSize)
	if exists {
		val, err := data.CoerceToInteger(setting)
		if err == nil {
			settings.WindowSize = val
		}
	}

	settings.Resolution = 1 // default window resolution
	setting, exists = ctx.GetSetting(sResolution)
	if exists {
		val, err := data.CoerceToInteger(setting)
		if err == nil {
			settings.Resolution = val
		}
	}

	settings.ProceedOnlyOnEmit = true // by default only proceed on emit
	setting, exists = ctx.GetSetting(sProceedOnlyOnEmit)
	if exists {
		val, err := data.CoerceToBoolean(setting)
		if err == nil {
			settings.ProceedOnlyOnEmit = val
		}
	}

	setting, exists = ctx.GetSetting(sAdditionalSettings)
	if exists {
		val, err := data.CoerceToString(setting)
		if err == nil {

			settings.AdditionalSettings, err = toParams(val)
			if err != nil {
				return nil, err
			}
		}
	}

	// settings validation can be done here once activities are created on configuration instead of
	// setting up during runtime

	return settings, nil
}

func toParams(values string) (map[string]string, error) {

	var params map[string]string

	result := strings.Split(values, ",")
	params = make(map[string]string)
	for _, pair := range result {
		nv := strings.Split(pair, "=")
		if len(nv) != 2 {
			return nil, fmt.Errorf("invalid settings")
		}
		params[nv[0]] = nv[1]
	}

	return params, nil
}
