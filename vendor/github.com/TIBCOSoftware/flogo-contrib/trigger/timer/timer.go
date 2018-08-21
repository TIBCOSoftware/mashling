package timer

import (
	"context"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/carlescere/scheduler"
)

var log = logger.GetLogger("trigger-flogo-timer")

type TimerTrigger struct {
	metadata *trigger.Metadata
	config   *trigger.Config
	timers   []*scheduler.Job

	handlers []*trigger.Handler
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &TimerFactory{metadata: md}
}

// TimerFactory Timer Trigger factory
type TimerFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *TimerFactory) New(config *trigger.Config) trigger.Trigger {
	return &TimerTrigger{metadata: t.metadata, config: config}
}

// Metadata implements trigger.Trigger.Metadata
func (t *TimerTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements trigger.Init
func (t *TimerTrigger) Initialize(ctx trigger.InitContext) error {

	t.handlers = ctx.GetHandlers()
	return nil
}

// Start implements ext.Trigger.Start
func (t *TimerTrigger) Start() error {

	log.Debug("Start")
	handlers := t.handlers

	log.Debug("Processing handlers")
	for _, handler := range handlers {

		repeating := handler.GetStringSetting("repeating")
		log.Debug("Repeating: ", repeating)
		if repeating == "false" {
			t.scheduleOnce(handler)
		} else if repeating == "true" {
			t.scheduleRepeating(handler)
		} else {
			log.Error("No match for repeating: ", repeating)
		}
		log.Debug("Settings repeating: ", handler.GetStringSetting("repeating"))
		//log.Debugf("Processing Handler: %s", handler.ActionId)
	}

	return nil
}

// Stop implements ext.Trigger.Stop
func (t *TimerTrigger) Stop() error {

	log.Debug("Stopping endpoints")
	for _, timer := range t.timers {

		if timer.IsRunning() {
			//log.Debug("Stopping timer for : ", k)
			timer.Quit <- true
		}
	}

	t.timers = nil

	return nil
}

func (t *TimerTrigger) scheduleOnce(endpoint *trigger.Handler) {
	log.Info("Scheduling a run one time job")

	seconds := getInitialStartInSeconds(endpoint)
	log.Debug("Seconds until trigger fires: ", seconds)

	var timerJob *scheduler.Job

	if seconds != 0 {
		timerJob = scheduler.Every(int(seconds))
	} else {
		log.Debug("Start Date not specified, executing action immediately")
	}

	fn := func() {
		log.Debug("Executing \"Once\" timer trigger")

		_, err := endpoint.Handle(context.Background(), nil)
		if err != nil {
			log.Error("Error running handler: ", err.Error())
		}
		//act := action.Get(endpoint.ActionId)
		//
		//if act != nil {
		//	log.Debugf("Running action: %s", endpoint.ActionId)
		//
		//	_, err := t.runner.RunHandler(context.Background(), act, nil )
		//
		//	if err != nil {
		//		log.Error("Error running action: ", err.Error())
		//	}
		//} else {
		//	log.Errorf("Action '%s' not found", endpoint.ActionId)
		//}

		if timerJob != nil {
			timerJob.Quit <- true
		}
	}

	if seconds == 0 {
		fn()
	} else {
		timerJob, err := timerJob.Seconds().NotImmediately().Run(fn)
		if err != nil {
			log.Error("Error performing scheduleOnce: ", err.Error())
		}

		t.timers = append(t.timers, timerJob)
	}
}

func (t *TimerTrigger) scheduleRepeating(endpoint *trigger.Handler) {
	log.Info("Scheduling a repeating job")

	seconds := getInitialStartInSeconds(endpoint)

	fn2 := func() {
		log.Debug("-- Starting \"Repeating\" (repeat) timer action")

		_, err := endpoint.Handle(context.Background(), nil)
		if err != nil {
			log.Error("Error running handler: ", err.Error())
		}

		//act := action.Get(endpoint.ActionId)
		//log.Debugf("Found action: '%+x'", act)
		//log.Debugf("ActionID: '%s'", endpoint.ActionId)
		//_, _, err := t.runner.Run(context.Background(), act, endpoint.ActionId, nil)
		//
		//if err != nil {
		//	log.Error("Error starting flow: ", err.Error())
		//}
	}

	if endpoint.GetStringSetting("notImmediate") == "false" {
		t.scheduleJobEverySecond(endpoint, fn2)
	} else {

		log.Debug("Seconds till trigger fires: ", seconds)
		timerJob := scheduler.Every(seconds)
		if timerJob == nil {
			log.Error("timerJob is nil")
		}

		t.scheduleJobEverySecond(endpoint, fn2)

		timerJob, err := timerJob.Seconds().NotImmediately().Run(fn2)
		if err != nil {
			log.Error("Error scheduleRepeating (first) flo err: ", err.Error())
		}
		if timerJob == nil {
			log.Error("timerJob is nil")
		}

		t.timers = append(t.timers, timerJob)
	}
}

func getInitialStartInSeconds(endpoint *trigger.Handler) int {

	layout := time.RFC3339
	startDate := endpoint.GetStringSetting("startDate")

	if startDate == "" {
		return 0
	}

	idx := strings.LastIndex(startDate, "Z")
	timeZone := startDate[idx+1:]
	log.Debug("Time Zone: ", timeZone)
	startDate = strings.TrimSuffix(startDate, timeZone)
	log.Debug("startDate: ", startDate)

	// is timezone negative
	var isNegative bool
	isNegative = strings.HasPrefix(timeZone, "-")
	// remove sign
	timeZone = strings.TrimPrefix(timeZone, "-")

	triggerDate, err := time.Parse(layout, startDate)
	if err != nil {
		log.Error("Error parsing time err: ", err.Error())
	}
	log.Debug("Time parsed from settings: ", triggerDate)

	var hour int
	var minutes int

	sliceArray := strings.Split(timeZone, ":")
	if len(sliceArray) != 2 {
		log.Error("Time zone has wrong format: ", timeZone)
	} else {
		hour, _ = strconv.Atoi(sliceArray[0])
		minutes, _ = strconv.Atoi(sliceArray[1])

		log.Debug("Duration hour: ", time.Duration(hour)*time.Hour)
		log.Debug("Duration minutes: ", time.Duration(minutes)*time.Minute)
	}

	hours, _ := strconv.Atoi(timeZone)
	log.Debug("hours: ", hours)
	if isNegative {
		log.Debug("Adding to triggerDate")
		triggerDate = triggerDate.Add(time.Duration(hour) * time.Hour)
		triggerDate = triggerDate.Add(time.Duration(minutes) * time.Minute)
	} else {
		log.Debug("Subtracting to triggerDate")
		triggerDate = triggerDate.Add(time.Duration(hour * -1))
		triggerDate = triggerDate.Add(time.Duration(minutes))
	}

	currentTime := time.Now().UTC()
	log.Debug("Current time: ", currentTime)
	log.Debug("Setting start time: ", triggerDate)
	duration := time.Since(triggerDate)
    durSeconds := duration.Seconds()
    if durSeconds < 0 {
    	//Future date
	    return int(math.Abs(durSeconds))
    } else {
    	// Past date
	   return 0
	}
}

type PrintJob struct {
	Msg string
}

func (j *PrintJob) Run() error {
	log.Debug(j.Msg)
	return nil
}

func (t *TimerTrigger) scheduleJobEverySecond(tgrHandler *trigger.Handler, fn func()) {

	var interval int = 0
	if seconds := tgrHandler.GetStringSetting("seconds"); seconds != "" {
		seconds, _ := strconv.Atoi(seconds)
		interval = interval + seconds
	}
	if minutes := tgrHandler.GetStringSetting("minutes"); minutes != "" {
		minutes, _ := strconv.Atoi(minutes)
		interval = interval + minutes*60
	}
	if hours := tgrHandler.GetStringSetting("hours"); hours != "" {
		hours, _ := strconv.Atoi(hours)
		interval = interval + hours*3600
	}

	log.Debug("Repeating seconds: ", interval)
	// schedule repeating
	timerJob, err := scheduler.Every(interval).Seconds().Run(fn)
	if err != nil {
		log.Error("Error scheduleRepeating (repeat seconds) flo err: ", err.Error())
	}
	if timerJob == nil {
		log.Error("timerJob is nil")
	}

	t.timers = append(t.timers, timerJob)
}
