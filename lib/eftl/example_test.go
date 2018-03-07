/*
 * Copyright (c) 2001-2017 TIBCO Software Inc.
 * All Rights Reserved. Confidential & Proprietary.
 * For more information, please contact:
 * TIBCO Software Inc., Palo Alto, California, USA
 *
 * $Id: example_test.go 92167 2017-03-08 17:49:12Z $
 */

package eftl

import (
	"fmt"
	"log"
	"net"
	"os"
	"testing"
	"time"
)

func skip() bool {
	_, err := net.Dial("tcp", "127.0.0.1:9191")
	return err != nil
}

func TestMain(m *testing.M) {
	if skip() {
		return
	}
	os.Exit(m.Run())
}

// Connect to the server.
func ExampleConnect() {
	errChan := make(chan error)

	opts := &Options{
		Username: "username",
		Password: "password",
	}

	// connect to the server
	conn, err := Connect("ws://localhost:9191/channel", opts, errChan)
	if err != nil {
		log.Printf("connect failed: %s", err)
		return
	}

	// disconnect from the server when done with the connection
	defer conn.Disconnect()

	// listen for asynnchronous errors
	go func() {
		for err := range errChan {
			log.Printf("connection error: %s", err)
		}
	}()
	// Output:
}

// Reconnect to the server.
func ExampleConnection_Reconnect() {
	// connect to the server
	conn, err := Connect("ws://localhost:9191/channel", nil, nil)
	if err != nil {
		log.Printf("connect failed: %s", err)
		return
	}

	// disconnect from the server when done with the connection
	defer conn.Disconnect()

	// disconnect from the server
	conn.Disconnect()

	// reconnect to the server
	err = conn.Reconnect()
	if err != nil {
		log.Printf("reconnect failed: %s", err)
		return
	}
	// Output:
}

// Publish messages.
func ExampleConnection_Publish() {
	// connect to the server
	conn, err := Connect("ws://localhost:9191/channel", nil, nil)
	if err != nil {
		log.Printf("connect failed: %s", err)
		return
	}

	// disconnect from the server when done with the connection
	defer conn.Disconnect()

	// publish a message
	err = conn.Publish(Message{
		"_dest":        "sample",
		"field-int":    99,
		"field-float":  0.99,
		"field-string": "hellow, world!",
		"field-time":   time.Now(),
		"field-message": Message{
			"field-bytes": []byte("this is an embedded message"),
		},
	})
	if err != nil {
		log.Printf("publish failed: %s", err)
		return
	}
	// Output:
}

// Publish messages asynchronously.
func ExampleConnection_PublishAsync() {
	// connect to the server
	conn, err := Connect("ws://localhost:9191/channel", nil, nil)
	if err != nil {
		log.Printf("connect failed: %s", err)
		return
	}

	// disconnect from the server when done with the connection
	defer conn.Disconnect()

	compChan := make(chan *Completion, 1)

	// publish a message
	err = conn.PublishAsync(Message{
		"_dest":        "sample",
		"field-int":    99,
		"field-float":  0.99,
		"field-string": "hellow, world!",
		"field-time":   time.Now(),
		"field-message": Message{
			"field-bytes": []byte("this is an embedded message"),
		},
	}, compChan)
	if err != nil {
		log.Printf("publish failed: %s", err)
		return
	}

	// wait for publish operation to complete
	comp := <-compChan

	if comp.Error != nil {
		log.Printf("publish completion failed: %s", err)
	} else {
		fmt.Println("published message")
	}
	// Output: published message
}

// Subscribe to messages.
func ExampleConnection_Subscribe() {
	errChan := make(chan error, 1)

	// connect to the server
	conn, err := Connect("ws://localhost:9191/channel", nil, errChan)
	if err != nil {
		log.Printf("connect failed: %s", err)
		return
	}

	// disconnect from the server when done with the connection
	defer conn.Disconnect()

	msgChan := make(chan Message, 1)

	// subscribe to messages
	_, err = conn.Subscribe("{\"_dest\": \"sample\"}", "", msgChan)
	if err != nil {
		log.Printf("subscribe failed: %s", err)
		return
	}

	done := false
	go func() {
		for !done {
			err = conn.Publish(Message{
				"_dest": "sample",
				"test":  "hello world",
			})
			if err != nil {
				log.Printf("publish failed: %s", err)
				return
			}
		}
	}()

	// receive messages
	for {
		select {
		case msg := <-msgChan:
			fmt.Println(msg["test"])
			done = true
			return
		case err := <-errChan:
			log.Printf("connection error: %s", err)
			return
		}
	}
	// Output: hello world
}

// Subscribe to messages asynchronously.
func ExampleConnection_SubscribeAsync() {
	errChan := make(chan error)

	// connect to the server
	conn, err := Connect("ws://localhost:9191/channel", nil, errChan)
	if err != nil {
		log.Printf("connect failed: %s", err)
		return
	}

	// disconnect from the server when done with the connection
	defer conn.Disconnect()

	subChan := make(chan *Subscription)
	msgChan := make(chan Message)

	// subscribe to messages
	err = conn.SubscribeAsync("{\"_dest\": \"sample\"}", "", msgChan, subChan)
	if err != nil {
		log.Printf("subscribe failed: %s", err)
		return
	}

	done := false
	go func() {
		for !done {
			err = conn.Publish(Message{
				"_dest": "sample",
				"test":  "hello world",
			})
			if err != nil {
				log.Printf("publish failed: %s", err)
				return
			}
		}
	}()

	// wait for subsribe operation to complete and receive messages
	for {
		select {
		case sub := <-subChan:
			if sub.Error != nil {
				log.Printf("subscription failed: %s", sub.Error)
				return
			}
			fmt.Println("subscription succeeded")
		case msg := <-msgChan:
			fmt.Println(msg["test"])
			done = true
			return
		case err := <-errChan:
			log.Printf("connection error: %s", err)
			return
		}
	}
	// Output: subscription succeeded
	// hello world
}
