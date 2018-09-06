package channel

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/engine/channels"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// activityLogger is the default logger for the Channel Activity
var activityLogger = logger.GetLogger("activity-channel")

const (
	sChannel = "channel"

	ivChannel = "channel"
	ivValue   = "value"
)

//we can generate json from this! - we could also create a "validate-able" object from this
type Settings struct {
	Channel string `md:"required"`
}

func init() {
	activityLogger.SetLogLevel(logger.InfoLevel)
}

// ChannelActivity is an Activity that is used to Channel a message to the console
type ChannelActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &ChannelActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *ChannelActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Channels the Message
func (a *ChannelActivity) Eval(ctx activity.Context) (done bool, err error) {

	//todo move to Activity instance creation
	settings, err := getSettings(ctx)
	if err != nil {
		return false, err
	}

	name := settings.Channel

	iName := ctx.GetInput(ivChannel)

	if iName != nil && iName != "" {
		name = iName.(string)
	}

	ch := channels.Get(name)

	if ch == nil {
		return false, fmt.Errorf("channel '%s' not registered with engine", name)
	}

	in := ctx.GetInput(ivValue)

	blocking := true

	//should we allow
	if blocking {
		ch.Publish(in)
	} else {
		ch.PublishNoWait(in)
	}

	return true, nil
}

func getSettings(ctx activity.Context) (*Settings, error) {

	settings := &Settings{}

	setting, exists := ctx.GetSetting(sChannel)
	if exists {
		val, err := data.CoerceToString(setting)
		if err == nil {
			settings.Channel = val
		}
	} else {
		//throw error if //ivChannel is not defined
	}

	// settings validation can be done here once activities are created on configuration instead of
	// setting up during runtime

	return settings, nil
}
