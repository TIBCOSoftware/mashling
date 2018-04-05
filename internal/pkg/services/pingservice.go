package services

import (
	"errors"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/TIBCOSoftware/mashling/lib/util"
)

//PingService interface for ping services
type PingService interface {
	Init(string, string) error
	Start() error
	Stop() error
}

//GetPingService returns pingserviceconfig struct
func GetPingService() PingService {
	psc := &PingServiceConfig{}
	return psc
}

//PingServiceConfig holds ping related variables
type PingServiceConfig struct {
	*http.Server
	listener     net.Listener
	pingPort     string
	pingResponse string
}

//Init intialises pingport if not configured
func (p *PingServiceConfig) Init(pingPort, pingResponse string) error {
	if len(pingPort) != 0 {
		p.pingPort = pingPort
	} else {
		p.pingPort = util.Mashling_Default_Ping_Port_Val
	}
	p.pingResponse = pingResponse

	p.Server = &http.Server{Addr: ":" + p.pingPort}

	return nil
}

//Start starts ping  server on configured port
func (p *PingServiceConfig) Start() error {
	http.HandleFunc("/ping", p.PingSimpleServer)
	http.HandleFunc("/ping/details", p.PingDetailedServer)

	if p.listener != nil {
		return errors.New("Server already started")
	}

	addr := p.Addr
	if addr == "" {
		addr = ":http"
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	p.listener = listener

	go func() {
		err := p.Serve(listener)
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}
		}
	}()

	return nil
}

//PingSimpleServer handles simple response
func (p *PingServiceConfig) PingSimpleServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "{\"response\":\"Ping successful\"}\n")
}

//PingDetailedServer handles simple response
func (p *PingServiceConfig) PingDetailedServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, p.pingResponse+"\n")
}

//Stop handles nullifying configured port
func (p *PingServiceConfig) Stop() error {
	p.pingPort = ""
	if p.listener == nil {
		return errors.New("Server not started")
	}

	err := p.listener.Close()
	if err != nil {
		return err
	}

	return nil
}
