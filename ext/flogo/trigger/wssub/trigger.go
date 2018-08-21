/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package wssub

import (
	"context"
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/mashling/lib/util"
	"github.com/gorilla/websocket"
)

// log is the default package logger
var log = logger.GetLogger("trigger-tibco-wssub")

// WssubTrigger is websocket event trigger
type WssubTrigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	config   *trigger.Config
	wsconn   *websocket.Conn
}

// NewFactory creates a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &WssubFactory{metadata: md}
}

// WssubFactory websocket event subscriber Trigger factory
type WssubFactory struct {
	metadata *trigger.Metadata
}

// New Creates a new trigger instance for a given id
func (t *WssubFactory) New(config *trigger.Config) trigger.Trigger {
	return &WssubTrigger{metadata: t.metadata, config: config}
}

// Metadata implements trigger.Trigger.Metadata
func (t *WssubTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements ext.Trigger.Init
func (t *WssubTrigger) Init(runner action.Runner) {
	t.runner = runner
}

// Start implements ext.Trigger.Start
func (t *WssubTrigger) Start() error {

	// extract websocket endpoint url
	urlSetting := t.config.Settings["url"]
	if urlSetting == nil || urlSetting.(string) == "" {
		return fmt.Errorf("server url not provided")
	}

	url := urlSetting.(string)
	log.Infof("dialing websocket endpoint[%s]...", url)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("error while connecting to websocket endpoint[%s] - %s", url, err)
	}
	t.wsconn = conn

	err = run(t)

	return err
}

// Stop implements ext.Trigger.Stop
func (t *WssubTrigger) Stop() error {
	t.wsconn.Close()
	return nil
}

func run(t *WssubTrigger) error {

	// listen to websocket connection in a go routine
	go func() {
		for {
			_, message, err := t.wsconn.ReadMessage()
			if err != nil {
				log.Errorf("error while reading websocket message: %s", err)
				break
			}

			var content interface{}
			err = util.Unmarshal(util.MIMEApplicationJSON, message, &content)
			if err != nil {
				log.Errorf("received message is not a valid JSON, skipping this message")
				continue
			}

			for _, handler := range t.config.Handlers {

				actionID := handler.ActionId
				action := action.Get(actionID)
				log.Infof("Found action: ' %+x' for the actionID[%d]", action, actionID)

				data := make(map[string]interface{})
				data["content"] = content
				startAttrs, errorAttrs := t.metadata.OutputsToAttrs(data, true)
				if errorAttrs != nil || startAttrs == nil {
					log.Errorf("Failed to create output attributes for websocket message for ActionID [%s] for reason [%s] message lost", handler.ActionId, errorAttrs)
				}

				ctx := trigger.NewContext(context.Background(), startAttrs)
				_, _, err := t.runner.Run(ctx, action, actionID, nil)

				if err != nil {
					log.Errorf("Run action for ActionID [%s] failed for reason [%s] message lost", err, handler.ActionId)
				}
			}
		}
		log.Infof("stopped listening to websocket endpoint")
	}()

	return nil
}
