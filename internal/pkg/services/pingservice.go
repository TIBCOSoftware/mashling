package services

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
)

// DefaultPort is the default port for Ping service
const DefaultPort = "9090"

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
	pingResVal string
}

//PingResponse is to hold ping response
type PingResponse struct {
	Version        string
	Appversion     string
	Appdescription string
}

//Init intialises pingport if not configured
func (p *PingServiceConfig) Init(port string, pingRes PingResponse) error {
	if len(port) == 0 {
		port = DefaultPort
	}

	pingDataBytes, err := json.Marshal(pingRes)
	if err != nil {
		log.Println("[mashling-ping-service] Ping service data formation error")
	}

	p.pingResVal = string(pingDataBytes)

	p.Server = &http.Server{Addr: ":" + port}

	return nil
}

//Start starts ping  server on configured port
func (p *PingServiceConfig) Start() error {
	log.Println("[mashling-ping-service] Ping service starting...")
	http.HandleFunc("/ping", p.PingResponseHandlerShort)
	http.HandleFunc("/ping/details", p.PingResponseHandlerDetail)

	listener, err := net.Listen("tcp", p.Server.Addr)
	if err != nil {
		log.Println("[mashling-ping-service] Ping service failed to start due to error:", err)
		return err
	}

	err = listener.Close()
	if err != nil {
		log.Println("[mashling-ping-service] Ping service failed to start due to error:", err)
		return err
	}

	go func() {
		if err := p.ListenAndServe(); err != http.ErrServerClosed {
			log.Println("[mashling-ping-service] Ping service err:", err)
		}
	}()
	log.Println("[mashling-ping-service] Ping service started")

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
	if err := p.Shutdown(nil); err != nil {
		log.Println("[mashling-ping-service] Ping service error when stopping:", err)
		return err
	}
	log.Println("[mashling-ping-service] Ping service stopped")
	return nil
}
