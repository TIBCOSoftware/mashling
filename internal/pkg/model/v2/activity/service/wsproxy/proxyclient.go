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

// RegisterProxyClient registers proxy client instance with proxy service
func (p *ProxyService) RegisterProxyClient(pc *ProxyClient) (err error) {
	p.Lock()
	if len(p.proxyclients) < p.maxConnections {
		p.proxyclients[pc.name] = pc
	} else {
		errMessage := fmt.Sprintf("proxy service[%s] utilized maximum[%d] concurrent connections, can't accept any more connections", p.name, p.maxConnections)
		err = errors.New(errMessage)
	}
	p.Unlock()
	return err
}

// DeregisterProxyClient deregisters proxy client instance from proxy service
func (p *ProxyService) DeregisterProxyClient(pc *ProxyClient) {
	p.Lock()
	delete(p.proxyclients, pc.name)
	p.Unlock()
}

// GetProxyClient returns proxy client instance corresponding to supplied name
func (p *ProxyService) GetProxyClient(name string) *ProxyClient {
	p.RLock()
	defer p.RUnlock()

	if pc, ok := p.proxyclients[name]; ok {
		return pc
	}
	return nil
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
func (p *ProxyServices) GetService(name string) *ProxyService {
	p.RLock()
	defer p.RUnlock()
	if ps, ok := p.services[name]; ok {
		return ps
	}
	return nil
}

// AddService adds proxy service to proxyServices
func (p *ProxyServices) AddService(ps *ProxyService) {
	p.Lock()
	p.services[ps.name] = ps
	p.Unlock()
}

// start creates new ProxyClient instance and handles upstream & downstream flow
func startProxyClient(wsp *WSProxy) error {
	// get proxy service
	pService := proxyServices.GetService(wsp.serviceName)
	if pService == nil {
		pService = &ProxyService{
			name:           wsp.serviceName,
			proxyclients:   make(map[string]*ProxyClient),
			backendURL:     wsp.backendURL,
			maxConnections: wsp.maxConnections,
		}
		proxyServices.AddService(pService)
	}

	// get proxy client
	clientName := fmt.Sprintf("%s-%p-%s", wsp.serviceName, wsp.clientConn, wsp.clientConn.RemoteAddr())
	pClient := pService.GetProxyClient(clientName)
	if pClient != nil {
		log.Warnf("connection [%p] already utilized", wsp.clientConn)
		return nil
	}

	// create new proxy client
	pClient = &ProxyClient{
		name:          clientName,
		startTime:     time.Now(),
		clientConn:    wsp.clientConn,
		upstreamErr:   make(chan error, 1),
		downstreamErr: make(chan error, 1),
	}
	defer pClient.clientConn.Close()
	log.Infof("handling websocket connection with the proxy client[%s]", clientName)

	// register proxy client
	err := pService.RegisterProxyClient(pClient)
	if err != nil {
		log.Error(err)
		closeMessage := websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error())
		pClient.clientConn.WriteMessage(websocket.CloseMessage, closeMessage)
		return nil
	}
	defer pService.DeregisterProxyClient(pClient)

	// establish backend connection
	log.Debugf("connecting to %s ", pService.backendURL)
	conn, _, err := websocket.DefaultDialer.Dial(pService.backendURL, nil)
	if err != nil {
		log.Errorf("connection error: %s", err)
		m := fmt.Sprintf("Mashling failed to connect backend url[%s]", pService.backendURL)
		closeMessage := websocket.FormatCloseMessage(websocket.CloseNormalClosure, m)
		pClient.clientConn.WriteMessage(websocket.CloseMessage, closeMessage)
		return nil
	}
	pClient.serverConn = conn
	defer pClient.serverConn.Close()

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
			log.Infof(infoMessageTemplate, e.Code, e.Text)
		}
	}
	log.Info(pClient.status())
	log.Infof("proxy[%s->%s] closed", pService.name, pClient.name)

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
	upstream bytes copyied: %d
	downstream bytes copyied: %d`
	status := fmt.Sprintf(statusTemplate, pc.name, time.Since(pc.startTime), pc.upstreamBytes, pc.downstreamBytes)
	return status
}
