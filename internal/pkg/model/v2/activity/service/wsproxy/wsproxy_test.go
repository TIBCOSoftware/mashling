package wsproxy

import (
	"net/http"
	"testing"
	"time"

	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/types"
	"github.com/gorilla/websocket"
)

func TestWSProxy(t *testing.T) {
	fini := make(chan bool, 1)
	wsHandler := func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		for {
			mt, message, err := conn.ReadMessage()
			next := string(message)
			t.Log(next)
			if len(next) > 5 {
				fini <- true
				break
			}
			if err != nil {
				break
			}
			conn.WriteMessage(mt, []byte(next+"."))
			if err != nil {
				break
			}
		}
	}

	middleware := http.NewServeMux()
	middleware.HandleFunc("/ws", wsHandler)
	server := http.Server{
		Addr:    "localhost:8282",
		Handler: middleware,
	}
	done := make(chan bool, 1)
	go func() {
		server.ListenAndServe()
		done <- true
	}()
	defer func() {
		err := server.Shutdown(nil)
		if err != nil {
			t.Fatal(err)
		}
		<-done
	}()

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8282/ws", http.Header{})
	for err != nil {
		conn, _, err = websocket.DefaultDialer.Dial("ws://localhost:8282/ws", http.Header{})
	}
	defer conn.Close()

	service := types.Service{
		Type: "ws",
		Settings: map[string]interface{}{
			"url":            "ws://localhost:8282/ws",
			"maxConnections": 2.0,
		},
	}
	instance, err := InitializeWSProxy("test", service.Settings)
	if err != nil {
		t.Fatal(err)
	}
	err = instance.UpdateRequest(map[string]interface{}{
		"wsconnection": conn,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = instance.Execute()
	if err != nil {
		t.Fatal(err)
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte("test"))
	if err != nil {
		t.Fatal(err)
	}

	select {
	case <-fini:
	case <-time.After(30 * time.Second):
		t.Fatal("test failed: timed out")
	}
}
