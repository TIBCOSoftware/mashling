package services

import (
	"io"
	"log"
	"net/http"

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
	return nil
}

//Start starts ping  server on configured port
func (p *PingServiceConfig) Start() error {
	http.HandleFunc("/ping", p.PingSimpleServer)
	http.HandleFunc("/ping/details", p.PingDetailedServer)
	go func() {
		err := http.ListenAndServe(":"+p.pingPort, nil)
		if err != nil {
			log.Println("error occured while running ping functionality :", err)
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
	return nil
}
