package wsproxy

import (
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
)

// WSProxy is websocket proxy service
type WSProxy struct {
	backendURL string
	clientConn *websocket.Conn
}

// InitializeWSProxy initializes an WSProxy service with provided settings.
func InitializeWSProxy(settings map[string]interface{}) (wspService *WSProxy, err error) {
	fmt.Println("InitializeWSProxy")
	wspService = &WSProxy{}
	err = wspService.setRequestValues(settings)
	return wspService, err
}

// Execute invokes this WSProxy service.
func (wsp *WSProxy) Execute() (err error) {
	fmt.Println("Execute()")
	fmt.Println("client address:", wsp.clientConn.RemoteAddr())
	fmt.Println("backendURL:", wsp.backendURL)

	//start proxy client
	go start(wsp.clientConn, wsp.backendURL)

	return nil
}

// UpdateRequest updates a request on an existing WSProxy service instance with new values.
func (wsp *WSProxy) UpdateRequest(values map[string]interface{}) (err error) {
	return wsp.setRequestValues(values)
}

func (wsp *WSProxy) setRequestValues(settings map[string]interface{}) (err error) {
	fmt.Println("setRequestValues")
	for k, v := range settings {
		switch k {
		case "wsconnection":
			conn, ok := v.(*websocket.Conn)
			if !ok {
				return errors.New("invalid type for web socket connection")
			}
			wsp.clientConn = conn
		case "url":
			url, ok := v.(string)
			if !ok {
				return errors.New("invalid type for url")
			}
			wsp.backendURL = url
		default:
			// ignore and move on.
		}
	}

	return nil
}
