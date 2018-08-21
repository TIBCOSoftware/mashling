package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"reflect"

	"github.com/TIBCOSoftware/flogo-lib/config"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("trigger-flogo-cli")

var singleton *CliTrigger

// CliTrigger CLI trigger struct
type CliTrigger struct {
	metadata     *trigger.Metadata
	config       *trigger.Config
	handlerInfos []*handlerInfo
	defHandler   *trigger.Handler
}

type handlerInfo struct {
	Invoke  bool
	handler *trigger.Handler
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &CliTriggerFactory{metadata: md}
}

// CliTriggerFactory CLI Trigger factory
type CliTriggerFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *CliTriggerFactory) New(config *trigger.Config) trigger.Trigger {
	singleton = &CliTrigger{metadata: t.metadata, config: config}

	return singleton
}

// Metadata implements trigger.Trigger.Metadata
func (t *CliTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

func (t *CliTrigger) Initialize(ctx trigger.InitContext) error {

	level, err := logger.GetLevelForName(config.GetLogLevel())

	if err == nil {
		log.SetLogLevel(level)
	}

	//if len(t.config.Settings) == 0 {
	//	return fmt.Errorf("no Settings found for trigger '%s'", t.config.Id)
	//}

	if len(ctx.GetHandlers()) == 0 {
		return fmt.Errorf("no Handlers found for trigger '%s'", t.config.Id)
	}

	hasDefault := false

	// Init handlers
	for _, handler := range ctx.GetHandlers() {

		cmdString := "default"

		aInfo := &handlerInfo{Invoke: false, handler: handler}

		if cmd, ok := handler.GetSetting("command"); ok && cmd != nil {
			cmdString = cmd.(string)
		}

		if cmd, set := handler.GetSetting("default"); set {
			if def, ok := cmd.(bool); ok && def {
				t.defHandler = handler
				hasDefault = true
			}
		}

		t.handlerInfos = append(t.handlerInfos, aInfo)

		xv := reflect.ValueOf(aInfo).Elem()
		addr := xv.FieldByName("Invoke").Addr().Interface()

		switch ptr := addr.(type) {
		case *bool:
			flag.BoolVar(ptr, cmdString, false, "")
		}
	}

	if !hasDefault && len(t.handlerInfos) > 0 {
		t.defHandler = t.handlerInfos[0].handler
	}

	return nil
}

func (t *CliTrigger) Start() error {
	return nil
}

// Stop implements util.Managed.Stop
func (t *CliTrigger) Stop() error {
	return nil
}

func Invoke() (string, error) {

	var args []string
	flag.Parse()

	// if we have additional args (after the cmd name and the flow cmd switch)
	// stuff those into args and pass to Invoke(). The action will only receive the
	// aditional args that were intending for the action logic.
	if arg := flag.Args(); len(arg) >= 2 {
		args = flag.Args()[2:]
	}

	for _, info := range singleton.handlerInfos {

		if info.Invoke {
			return singleton.Invoke(info.handler, args)
		}
	}

	return singleton.Invoke(singleton.defHandler, args)
}

func (t *CliTrigger) Invoke(handler *trigger.Handler, args []string) (string, error) {

	log.Infof("invoking handler '%s'", handler)

	data := map[string]interface{}{
		"args": args,
	}

	results, err := handler.Handle(context.Background(), data)

	if err != nil {
		log.Debugf("error: %s", err.Error())
		return "", err
	}

	var replyData interface{}

	if len(results) != 0 {
		dataAttr, ok := results["data"]
		if ok {
			replyData = dataAttr.Value()
		}
	}

	if replyData != nil {
		data, err := json.Marshal(replyData)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	return "", nil
}
