package wsmessage

import (
	"net/url"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/gorilla/websocket"
)

// log is the default package logger
var log = logger.GetLogger("activity-sendWSMessage")

// WsMsgActivity is a stub for your Activity implementation
type WsMsgActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new WsMsgActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &WsMsgActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *WsMsgActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval - Sends a message to a WebSocket enabled server like TIBCO eFTL
func (a *WsMsgActivity) Eval(context activity.Context) (done bool, err error) {

	// Get the activity data from the context
	wsHost := context.GetInput("Server").(string)
	wsChannel := context.GetInput("Channel").(string)
	wsDestination := context.GetInput("Destination").(string)
	wsMessage := context.GetInput("Message").(string)
	wsUser := context.GetInput("Username").(string)
	wsPassword := context.GetInput("Password").(string)

	wsURL := url.URL{Scheme: "ws", Host: wsHost, Path: wsChannel}
	log.Debugf("connecting to %s", wsURL.String())

	wsConn, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		log.Debugf("Error while dialing to wsHost: ", err)
	}

	loginMessage := `{"op": 1, "client_type": "js", "client_version": "3.0.0   V9", "user":"` + wsUser + `", "password":"` + wsPassword + `", "login_options": {"_qos": "true"}}`

	log.Debugf("Preparing to send login message: [%s]", loginMessage)

	err = wsConn.WriteMessage(websocket.TextMessage, []byte(loginMessage))
	if err != nil {
		log.Debugf("Error while sending login message to wsHost: [%s]", err)
		return
	}

	textMessage := `{"op": 8, "body": {"_dest":"` + wsDestination + `", "text":"` + wsMessage + `"}, "seq": 1}`

	log.Debugf("Preparing to send message: [%s]", textMessage)

	err = wsConn.WriteMessage(websocket.TextMessage, []byte(textMessage))
	if err != nil {
		log.Debugf("Error while sending message to wsHost: [%s]", err)
		return
	}
	wsConn.Close()
	context.SetOutput("result", "OK")

	return true, nil
}
