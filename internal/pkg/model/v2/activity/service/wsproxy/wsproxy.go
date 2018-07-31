package wsproxy

import (
	"errors"
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/gorilla/websocket"
)

var log = logger.GetLogger("service-wsproxy")

// WSProxy is websocket proxy service
type WSProxy struct {
	serviceName string
	backendURL  string
	clientConn  *websocket.Conn
}

// InitializeWSProxy initializes an WSProxy service with provided settings.
func InitializeWSProxy(name string, settings map[string]interface{}) (wspService *WSProxy, err error) {
	wspService = &WSProxy{
		serviceName: name,
	}
	err = wspService.setRequestValues(settings)
	return wspService, err
}

// Execute invokes this WSProxy service.
func (wsp *WSProxy) Execute() (err error) {

	//proxy client name is combination of connection + service name
	proxyName := fmt.Sprintf("%p-%s", wsp.clientConn, wsp.serviceName)
	log.Infof("starting proxy (name: %s, client address: %s & server url: %s)...", proxyName, wsp.clientConn.RemoteAddr(), wsp.backendURL)
	//start proxy client
	go start(proxyName, wsp.clientConn, wsp.backendURL)

	return nil
}

// UpdateRequest updates a request on an existing WSProxy service instance with new values.
func (wsp *WSProxy) UpdateRequest(values map[string]interface{}) (err error) {
	return wsp.setRequestValues(values)
}

func (wsp *WSProxy) setRequestValues(settings map[string]interface{}) (err error) {
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
