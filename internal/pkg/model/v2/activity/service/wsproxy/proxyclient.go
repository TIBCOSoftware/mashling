package wsproxy

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// ProxyClient is proxy between client websocket connection and server websocket connection
type ProxyClient struct {
	clientConn *websocket.Conn
	serverConn *websocket.Conn
}

// start creates new ProxyClient instance from client connection and server URL
func start(clientConn *websocket.Conn, serverURL string) (err error) {
	//create new proxy client
	pClient := &ProxyClient{
		clientConn: clientConn,
	}
	//establish server connection
	fmt.Println("connecting to ", serverURL)
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		fmt.Println("Error while connecting to server", err)
		return err
	}
	pClient.serverConn = conn
	//run upstreamPump
	go pClient.upstreamPump()
	//run downstreamPump
	go pClient.downstreamPump()

	return nil
}

// upstreamPump pumps message from client connection to server connection
func (pc *ProxyClient) upstreamPump() (err error) {
	defer pc.closeConnections()
	for {
		mt, message, err := pc.clientConn.ReadMessage()
		if err != nil {
			fmt.Println("Read error", err)
			return err
		}
		// fmt.Println("message:", string(message))
		err = pc.serverConn.WriteMessage(mt, []byte(message))
		if err != nil {
			fmt.Println("Write error", err)
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
			fmt.Println("Read error", err)
			return err
		}
		// fmt.Println("message:", string(message))
		err = pc.clientConn.WriteMessage(mt, []byte(message))
		if err != nil {
			fmt.Println("Write error", err)
			return err
		}
	}
}

// closeConnections closes clinet & sever connections
func (pc *ProxyClient) closeConnections() {
	pc.clientConn.Close()
	pc.serverConn.Close()
}
