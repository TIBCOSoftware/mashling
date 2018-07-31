package wsproxy

//ProxyHub holds all ongoing websocket proxy clients instances information
type ProxyHub struct {
	proxys map[string]*ProxyClient
}

var pHub *ProxyHub

func init() {
	log.Infof("initialize proxy hub")
	pHub = &ProxyHub{
		proxys: make(map[string]*ProxyClient),
	}
}

func GetProxyHub() *ProxyHub {
	return pHub
}

func (hub *ProxyHub) RegisterProxyClient(pClient *ProxyClient) {
	log.Infof("register proxy client: %s", pClient.Name)
	hub.proxys[pClient.Name] = pClient
	hub.GetConcurentConnNumber()
}

func (hub *ProxyHub) UnregisterProxyClient(pClient *ProxyClient) {
	log.Infof("unregister proxy client: %s", pClient.Name)
	delete(hub.proxys, pClient.Name)
	hub.GetConcurentConnNumber()
}

func (hub *ProxyHub) GetConcurentConnNumber() {
	log.Infof("number of ongoing proxy connections = %d", len(hub.proxys))
}
