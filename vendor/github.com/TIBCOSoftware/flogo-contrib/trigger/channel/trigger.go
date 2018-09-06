package channel

import (
	"context"
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/engine/channels"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("trigger-flogo-channel")

const ovData = "data"

// ChannelTrigger CHANNEL trigger struct
type ChannelTrigger struct {
	metadata *trigger.Metadata
	config   *trigger.Config
	handlers []*trigger.Handler
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &ChannelFactory{metadata: md}
}

// ChannelFactory CHANNEL Trigger factory
type ChannelFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *ChannelFactory) New(config *trigger.Config) trigger.Trigger {
	return &ChannelTrigger{metadata: t.metadata, config: config}
}

// Metadata implements trigger.Trigger.Metadata
func (t *ChannelTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

func (t *ChannelTrigger) Initialize(ctx trigger.InitContext) error {

	t.handlers = ctx.GetHandlers()

	// validate handlers
	for _, handler := range t.handlers {

		channel := strings.ToLower(handler.GetStringSetting("channel"))

		ch := channels.Get(channel)
		if ch == nil {
			return fmt.Errorf("unknown engine channel '%s'", channel)
		}

		l := &ChannelListener{handler: handler}
		ch.RegisterCallback(l.OnMessage)
	}

	return nil
}

func (t *ChannelTrigger) Start() error {
	return nil
}

// Stop implements util.Managed.Stop
func (t *ChannelTrigger) Stop() error {
	return nil
}

type ChannelListener struct {
	handler *trigger.Handler
}

func (l *ChannelListener) OnMessage(msg interface{}) {
	triggerData := make(map[string]interface{})

	if attrs, ok := msg.(map[string]*data.Attribute); ok {

		vals := make(map[string]interface{})
		for name, attr := range attrs {
			vals[name] = attr.Value()
		}
		triggerData[ovData] = vals
	} else {
		triggerData[ovData] = msg
	}

	//todo what should we do with the results?
	_, err := l.handler.Handle(context.TODO(), triggerData)

	if err != nil {
		log.Error(err)
	}
}
