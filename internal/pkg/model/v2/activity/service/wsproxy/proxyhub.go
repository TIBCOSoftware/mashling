package wsproxy

import "fmt"

//ProxyHub holds all ongoing websocket proxy clients instances information
type ProxyHub struct {
	proxys map[string]*ProxyClient
}

//Hub operations allowed on the ProxyHub
type Hub interface {
	RegisterProxyClient(*ProxyClient)
	DeregisterProxyClient(*ProxyClient)
	GetProxyClient(string) *ProxyClient
	Status() string
}

//holds proxy hub instance
var pHub *ProxyHub

//initialize proxy hub once
func init() {
	log.Infof("initialize proxy hub")
	pHub = &ProxyHub{
		proxys: make(map[string]*ProxyClient),
	}
}

//GetProxyHub returns ProxyHub instance
func GetProxyHub() *ProxyHub {
	return pHub
}

//RegisterProxyClient registers proxy client with the hub
func (hub *ProxyHub) RegisterProxyClient(pClient *ProxyClient) {
	log.Debugf("register proxy client[%s]", pClient.name)
	hub.proxys[pClient.name] = pClient
	log.Debugf(hub.Status())
}

//DeregisterProxyClient deregisters proxy client from the hub
func (hub *ProxyHub) DeregisterProxyClient(pClient *ProxyClient) {
	log.Debugf("deregister proxy client[%s]", pClient.name)
	delete(hub.proxys, pClient.name)
	log.Debugf(hub.Status())
}

//GetProxyClient returns ProxyClient instance by name
func (hub *ProxyHub) GetProxyClient(name string) *ProxyClient {
	if pClient, ok := hub.proxys[name]; ok {
		return pClient
	}
	return nil
}

//Status returns status of the hub
func (hub *ProxyHub) Status() string {
	statusTemplate := `proxy hub status:
	ongoing websocket proxy connections = %d`

	status := fmt.Sprintf(statusTemplate, len(hub.proxys))
	return status
}
