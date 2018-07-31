package wsproxy

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

// ProxyClient is proxy between client websocket connection and server websocket connection
type ProxyClient struct {
	name                           string
	startTime                      time.Time
	upstreamBytes, downstreamBytes int64
	clientConn                     *websocket.Conn
	serverConn                     *websocket.Conn
	upstreamErr, downstreamErr     chan error
}

// start creates new ProxyClient instance from client connection and server URL
func start(name string, clientConn *websocket.Conn, serverURL string) (err error) {

	//do nothing if already proxy client exists with the given name
	if GetProxyHub().GetProxyClient(name) != nil {
		return nil
	}

	defer clientConn.Close()

	//create new proxy client
	pClient := &ProxyClient{
		name:          name,
		startTime:     time.Now(),
		clientConn:    clientConn,
		upstreamErr:   make(chan error, 1),
		downstreamErr: make(chan error, 1),
	}

	//establish server connection
	log.Debugf("connecting to %s ", serverURL)
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Errorf("connection error: %s", err)
		return err
	}
	pClient.serverConn = conn
	defer conn.Close()

	//register proxy client
	GetProxyHub().RegisterProxyClient(pClient)
	defer GetProxyHub().DeregisterProxyClient(pClient)

	//run upstreamPump
	go pClient.upstreamPump()
	//run downstreamPump
	go pClient.downstreamPump()

	//wait until end of the streams
	var errMessageTemplate string
	var infoMessageTemplate string
	select {
	case err = <-pClient.upstreamErr:
		errMessageTemplate = "error while copying from client to server: [%d] %v"
		infoMessageTemplate = "close initiated from client: [%d] %v"
	case err = <-pClient.downstreamErr:
		errMessageTemplate = "error while copying from server to client: [%d] %v"
		infoMessageTemplate = "close initiated from backend: [%d] %v"
	}
	if e, ok := err.(*websocket.CloseError); ok {
		if websocket.IsUnexpectedCloseError(e, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure) {
			log.Errorf(errMessageTemplate, e.Code, e.Text)
		} else {
			log.Infof(infoMessageTemplate, e.Code, e.Text)
		}
	}
	log.Info(pClient.status())
	log.Infof("proxy[%s] closed", pClient.name)

	return nil
}

// upstreamPump pumps message from client connection to server connection
func (pc *ProxyClient) upstreamPump() {
	// defer pc.closeConnections()
	for {
		mt, message, err := pc.clientConn.ReadMessage()
		if err != nil {
			errMessage := websocket.FormatCloseMessage(websocket.CloseMessage, fmt.Sprintf("%v", err))
			if e, ok := err.(*websocket.CloseError); ok {
				if e.Code != websocket.CloseNoStatusReceived {
					errMessage = websocket.FormatCloseMessage(e.Code, e.Text)
				}
			}
			pc.upstreamErr <- err
			pc.serverConn.WriteMessage(websocket.CloseMessage, errMessage)
			break
		}
		err = pc.serverConn.WriteMessage(mt, []byte(message))
		if err != nil {
			pc.upstreamErr <- err
			break
		}
		copiedBytes := len(message)
		pc.upstreamBytes += int64(copiedBytes)
	}
}

// downstreamPump pumps messages from server connection to client connection
func (pc *ProxyClient) downstreamPump() {
	// defer pc.closeConnections()
	for {
		mt, message, err := pc.serverConn.ReadMessage()
		if err != nil {
			errMessage := websocket.FormatCloseMessage(websocket.CloseMessage, fmt.Sprintf("%v", err))
			if e, ok := err.(*websocket.CloseError); ok {
				if e.Code != websocket.CloseNoStatusReceived {
					errMessage = websocket.FormatCloseMessage(e.Code, e.Text)
				}
			}
			pc.downstreamErr <- err
			pc.clientConn.WriteMessage(websocket.CloseMessage, errMessage)
			break
		}
		err = pc.clientConn.WriteMessage(mt, []byte(message))
		if err != nil {
			pc.downstreamErr <- err
			break
		}
		copiedBytes := len(message)
		pc.downstreamBytes += int64(copiedBytes)
	}
}

//Status returns status of the proxy client
func (pc *ProxyClient) status() string {
	statusTemplate := `proxy instance status:
	name: %s
	up time: %s
	upstream bytes copyied: %d
	downstream bytes copyied: %d`
	status := fmt.Sprintf(statusTemplate, pc.name, time.Since(pc.startTime), pc.upstreamBytes, pc.downstreamBytes)
	return status
}
