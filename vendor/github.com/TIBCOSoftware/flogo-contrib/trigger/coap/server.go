package coap

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-coap"
)

// Graceful shutdown CoapServer derived from: https://github.com/corneldamian/httpway/blob/master/server.go

// NewServer create a new server instance
func NewServer(n, addr string, handler coap.Handler) *Server {
	srv := &Server{n: n, Addr: addr, Handler: handler}
	return srv
}

//Server the server  structure
type Server struct {
	Handler coap.Handler
	Addr    string
	n       string

	listener     *net.UDPConn
	lastError    error
	serverGroup  *sync.WaitGroup
	clientsGroup chan bool
}

// Start will start server
// command isn't blocking, will exit after run
func (s *Server) Start() error {
	if s.Handler == nil {
		return errors.New("No server handler set")
	}

	if s.listener != nil {
		return errors.New("Server already started")
	}

	addr := s.Addr
	if addr == "" {
		addr = ":5683"
	}

	uaddr, err := net.ResolveUDPAddr(s.n, addr)
	if err != nil {
		return err
	}

	s.listener, err = net.ListenUDP(s.n, uaddr)
	if err != nil {
		return err
	}

	s.serverGroup = &sync.WaitGroup{}
	s.clientsGroup = make(chan bool, 50000)

	s.Handler = &serverHandler{s.Handler, s.clientsGroup}

	s.serverGroup.Add(1)
	go func() {

		err := coap.Serve(s.listener, s.Handler)
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}

			s.lastError = err
		}
	}()

	return nil
}

// Stop sends stop command to the server
func (s *Server) Stop() error {
	if s.listener == nil {
		return errors.New("Server not started")
	}

	if err := s.listener.Close(); err != nil {
		return err
	}

	return s.lastError
}

// IsStarted checks if the server is started
// will return true even if the server is stopped but there are still some requests to finish
func (s *Server) IsStarted() bool {
	if s.listener != nil {
		return true
	}

	if len(s.clientsGroup) > 0 {
		return true
	}

	return false
}

// WaitStop waits until server is stopped and all requests are finish
// timeout - is the time to wait for the requests to finish after the server is stopped
// will return error if there are still some requests not finished
func (s *Server) WaitStop(timeout time.Duration) error {
	if s.listener == nil {
		return errors.New("Server not started")
	}

	s.serverGroup.Wait()

	checkClients := time.Tick(100 * time.Millisecond)
	timeoutTime := time.NewTimer(timeout)

	for {
		select {
		case <-checkClients:
			if len(s.clientsGroup) == 0 {
				return s.lastError
			}
		case <-timeoutTime.C:
			return fmt.Errorf("WaitStop error, timeout after %s waiting for %d client(s) to finish", timeout, len(s.clientsGroup))
		}
	}
}

type serverHandler struct {
	handler      coap.Handler
	clientsGroup chan bool
}

func (sh *serverHandler) ServeCOAP(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
	sh.clientsGroup <- true
	defer func() {
		<-sh.clientsGroup
	}()

	return sh.handler.ServeCOAP(l, a, m)
}
