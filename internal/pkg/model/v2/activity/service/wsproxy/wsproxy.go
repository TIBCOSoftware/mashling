package wsproxy

import (
	"errors"

	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/gorilla/websocket"
)

var log = logger.GetLogger("service-wsproxy")

// WSProxy is websocket proxy service
type WSProxy struct {
	serviceName    string
	backendURL     string
	maxConnections int
	clientConn     *websocket.Conn
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

	// start proxy client as a goroutine
	go startProxyClient(wsp)

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
		case "maxConnections":
			i, ok := v.(float64)
			if !ok {
				return errors.New("invalid type for maxConnections")
			}
			wsp.maxConnections = int(i)
		default:
			// ignore and move on.
		}
	}

	return nil
}
