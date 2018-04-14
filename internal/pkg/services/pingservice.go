package services

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/TIBCOSoftware/mashling/lib/util"
)

//PingService interface for ping services
type PingService interface {
	Init(string, PingResponse) error
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
	listener   net.Listener
	pingPort   string
	pingResVal string
}

//PingResponse is to hold ping response
type PingResponse struct {
	Version        string
	Appversion     string
	Appdescription string
}

//Init intialises pingport if not configured
func (p *PingServiceConfig) Init(pingPort string, pingRes PingResponse) error {
	if len(pingPort) != 0 {
		p.pingPort = pingPort
	} else {
		p.pingPort = util.Mashling_Default_Ping_Port_Val
	}

	pingDataBytes, err := json.Marshal(pingRes)
	if err != nil {
		log.Println("[mashling] ping data formation error")
	}

	p.pingResVal = string(pingDataBytes)

	p.Server = &http.Server{Addr: ":" + p.pingPort}

	return nil
}

//Start starts ping  server on configured port
func (p *PingServiceConfig) Start() error {
	log.Println("[mashling-ping-service] Ping service Starting...")
	http.HandleFunc("/ping", p.PingResponseHandlerShort)
	http.HandleFunc("/ping/details", p.PingResponseHandlerDetail)

	if p.listener != nil {
		return errors.New("Server already started")
	}

	addr := p.Addr
	if addr == "" {
		addr = ":http"
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("[mashling-ping-service] failed to start Ping service due to error [%v]", err)
		os.Exit(1)
	}

	p.listener = listener

	go func() {
		err := p.Serve(listener)
		if err != nil {
			log.Printf("[mashling-ping-service] failed to start Ping service due to error [%v]", err)
			os.Exit(1)
		}
	}()
	log.Println("[mashling-ping-service] Ping service Started")

	return nil
}

//PingResponseHandlerShort handles simple response
func (p *PingServiceConfig) PingResponseHandlerShort(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "{\"response\":\"Ping successful\"}\n")
}

//PingResponseHandlerDetail handles simple response
func (p *PingServiceConfig) PingResponseHandlerDetail(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, p.pingResVal+"\n")
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
	log.Printf("[mashling-ping-service] Ping service Stoped")

	return nil
}
