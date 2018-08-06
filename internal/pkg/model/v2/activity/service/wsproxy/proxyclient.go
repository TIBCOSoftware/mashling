package wsproxy

import (
	"errors"
	"fmt"
	"sync"
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

// ProxyService holds ongoing ProxyClient instances
type ProxyService struct {
	name           string
	proxyclients   map[string]*ProxyClient
	backendURL     string
	maxConnections int
	sync.RWMutex
}

// CreateProxyClient creates proxy client instance with the supplied name & client connection
func (p *ProxyService) CreateProxyClient(name string, conn *websocket.Conn) (*ProxyClient, error) {
	p.Lock()
	defer p.Unlock()

	pClient := p.proxyclients[name]
	if pClient != nil {
		errMessage := fmt.Sprintf("connection [%p] already being handled with the client[%s]", conn, name)
		return nil, errors.New(errMessage)
	}
	if len(p.proxyclients) >= p.maxConnections {
		errMessage := fmt.Sprintf("proxy service[%s] utilized maximum[%d] allowed concurrent connections, can't accept any more connections", p.name, p.maxConnections)
		closeMessage := websocket.FormatCloseMessage(websocket.CloseNormalClosure, errMessage)
		conn.WriteMessage(websocket.CloseMessage, closeMessage)
		conn.Close()
		return nil, errors.New(errMessage)
	}
	pClient = &ProxyClient{
		name:          name,
		startTime:     time.Now(),
		clientConn:    conn,
		upstreamErr:   make(chan error, 1),
		downstreamErr: make(chan error, 1),
	}
	p.proxyclients[name] = pClient

	return pClient, nil
}

// ReleaseProxyClient removes proxy client instance from proxy service
func (p *ProxyService) ReleaseProxyClient(pc *ProxyClient) {
	p.Lock()
	defer p.Unlock()
	delete(p.proxyclients, pc.name)
}

// ProxyServices holds multiple ProxyService instances
type ProxyServices struct {
	services map[string]*ProxyService
	sync.RWMutex
}

// proxyServices holds ProxyService instances created accross gateway
var proxyServices = ProxyServices{
	services: make(map[string]*ProxyService),
}

// GetService returns proxy service corresponding to supplied name
// it creates new service it doesn't exist already
func (p *ProxyServices) GetService(name string, backendURL string, maxConnections int) *ProxyService {
	p.Lock()
	defer p.Unlock()
	pService := p.services[name]
	if pService == nil {
		pService = &ProxyService{
			name:           name,
			proxyclients:   make(map[string]*ProxyClient),
			backendURL:     backendURL,
			maxConnections: maxConnections,
		}
		p.services[name] = pService
	}
	return pService
}

// ReleaseService delete the service it doesn't hold any proxy clients
func (p *ProxyServices) ReleaseService(name string) {
	p.Lock()
	defer p.Unlock()
	if pService, ok := p.services[name]; ok {
		if len(pService.proxyclients) <= 0 {
			delete(p.services, name)
		}
	}
}

// start creates new ProxyClient instance and handles upstream & downstream flow
func startProxyClient(wsp *WSProxy) error {
	log.Debugf("starting proxy between the connection:%p & backendURL:%s ...", wsp.clientConn, wsp.backendURL)
	// get proxy service
	pService := proxyServices.GetService(wsp.serviceName, wsp.backendURL, wsp.maxConnections)
	defer proxyServices.ReleaseService(wsp.serviceName)

	// create proxy client
	clientName := fmt.Sprintf("%s-%p-%s", wsp.serviceName, wsp.clientConn, wsp.clientConn.RemoteAddr())
	pClient, err := pService.CreateProxyClient(clientName, wsp.clientConn)
	if err != nil {
		log.Errorf("error while creating proxy - %s", err.Error())
		return nil
	}
	defer pService.ReleaseProxyClient(pClient)
	defer pClient.clientConn.Close()

	// establish backend connection
	log.Debugf("connecting to %s ", pService.backendURL)
	conn, _, err := websocket.DefaultDialer.Dial(pService.backendURL, nil)
	if err != nil {
		log.Errorf("connection error: %s", err)
		m := fmt.Sprintf("failed to connect backend url[%s]", pService.backendURL)
		closeMessage := websocket.FormatCloseMessage(websocket.CloseNormalClosure, m)
		pClient.clientConn.WriteMessage(websocket.CloseMessage, closeMessage)
		return nil
	}
	pClient.serverConn = conn
	defer pClient.serverConn.Close()
	log.Infof("proxy[%s] started", clientName)

	// handle upstream & downstream on saparate goroutines
	go pClient.upstreamPump()
	go pClient.downstreamPump()

	// wait until end of the streams
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
			log.Debugf(infoMessageTemplate, e.Code, e.Text)
		}
	}
	log.Debug(pClient.status())
	log.Infof("proxy[%s] stopped", pClient.name)

	return nil
}

// upstreamPump pumps message from client connection to server connection
func (pc *ProxyClient) upstreamPump() {
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

// status returns status of the proxy client
func (pc *ProxyClient) status() string {
	statusTemplate := `proxy instance status:
	name: %s
	up time: %s
	upstream bytes transferred: %d
	downstream bytes transferred: %d`
	status := fmt.Sprintf(statusTemplate, pc.name, time.Since(pc.startTime), pc.upstreamBytes, pc.downstreamBytes)
	return status
}
