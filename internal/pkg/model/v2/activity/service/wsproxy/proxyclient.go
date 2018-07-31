package wsproxy

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

// ProxyClient is proxy between client websocket connection and server websocket connection
type ProxyClient struct {
	Name                           string
	startTime                      time.Time
	totalUpTime                    time.Duration
	upstreamBytes, downstreamBytes int64
	clientConn                     *websocket.Conn
	serverConn                     *websocket.Conn
}

// start creates new ProxyClient instance from client connection and server URL
func start(clientConn *websocket.Conn, serverURL string) (err error) {
	name := fmt.Sprintf("proxy-client-%v", clientConn.RemoteAddr())
	//create new proxy client
	pClient := &ProxyClient{
		Name:       name,
		startTime:  time.Now(),
		clientConn: clientConn,
	}
	//establish server connection
	log.Debugf("connecting to %s ", serverURL)
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Errorf("connection error: %s", err)
		return err
	}
	pClient.serverConn = conn

	defer pClient.closeConnections()

	//errors channels
	upstreamErr := make(chan error, 1)
	downstreamErr := make(chan error, 1)
	//run upstreamPump
	// go pClient.upstreamPump()
	go pumpStream(pClient.clientConn, pClient.serverConn, upstreamErr)
	//run downstreamPump
	// go pClient.downstreamPump()
	go pumpStream(pClient.serverConn, pClient.clientConn, downstreamErr)

	//register proxy client
	GetProxyHub().RegisterProxyClient(pClient)

	//wait until end of the streams
	var errMessage string
	select {
	case err = <-upstreamErr:
		errMessage = "error while copying from client to server: [%d] %v"
	case err = <-downstreamErr:
		errMessage = "error while copying from server to client: [%d] %v"
	}
	if e, ok := err.(*websocket.CloseError); ok {
		if websocket.IsUnexpectedCloseError(e) {
			log.Errorf(errMessage, e.Code, e.Text)
		}
	}

	return nil
}

// upstreamPump pumps message from client connection to server connection
func (pc *ProxyClient) upstreamPump() (err error) {
	defer pc.closeConnections()
	for {
		mt, message, err := pc.clientConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Errorf("[upstreamPump] read error: %s", err)
				return err
			}
			log.Infof("[upstreamPump] closing the server connection")
			return nil
		}
		log.Debugf("[upstreamPump] message: %s", string(message))
		err = pc.serverConn.WriteMessage(mt, []byte(message))
		if err != nil {
			log.Errorf("[upstreamPump] write error: %s", err)
			return err
		}
	}
}

// downstreamPump pumps messages from server connection to client connection
func (pc *ProxyClient) downstreamPump() (err error) {
	defer pc.closeConnections()
	for {
		mt, message, err := pc.serverConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Errorf("[downstreamPump] read error: %s", err)
				return err
			}
			log.Infof("[downstreamPump] closing the client connection")
			return nil
		}
		log.Debugf("[downstreamPump] message: %s", string(message))
		err = pc.clientConn.WriteMessage(mt, []byte(message))
		if err != nil {
			log.Errorf("[downstreamPump] write error: %s", err)
			return err
		}
	}
}

// closeConnections closes clinet & sever connections
func (pc *ProxyClient) closeConnections() {
	log.Info("closing websocket proxy...")
	pc.serverConn.Close()
	pc.clientConn.Close()
	pc.totalUpTime = time.Since(pc.startTime)

	pc.logProxyInstanceInfo()

	GetProxyHub().UnregisterProxyClient(pc)
}

func pumpStream(srcConn, dstConn *websocket.Conn, streamErr chan<- error) {
	for {
		messageType, message, err := srcConn.ReadMessage()
		if err != nil {
			errMessage := websocket.FormatCloseMessage(websocket.CloseMessage, fmt.Sprintf("%v", err))
			if e, ok := err.(*websocket.CloseError); ok {
				if e.Code != websocket.CloseNoStatusReceived {
					errMessage = websocket.FormatCloseMessage(e.Code, e.Text)
				}
			}
			streamErr <- err
			dstConn.WriteMessage(websocket.CloseMessage, errMessage)
			break
		}
		err = dstConn.WriteMessage(messageType, message)
		if err != nil {
			streamErr <- err
			break
		}
	}
}

func (pc *ProxyClient) logProxyInstanceInfo() {
	summaryMsg := `proxy instance closing summary:
	name: %s
	up time: %s
	upstream bytes: %d
	downstream bytes: %d`
	log.Infof(summaryMsg, pc.Name, pc.totalUpTime, pc.upstreamBytes, pc.downstreamBytes)
}
