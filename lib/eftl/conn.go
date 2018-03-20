/*
 * Copyright (c) 2001-2017 TIBCO Software Inc.
 * All Rights Reserved. Confidential & Proprietary.
 * For more information, please contact:
 * TIBCO Software Inc., Palo Alto, California, USA
 *
 * $Id: conn.go 95090 2017-08-07 21:40:08Z $
 */

package eftl

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Errors
var (
	ErrTimeout          = errors.New("operation timed out")
	ErrNotConnected     = errors.New("not connected")
	ErrInvalidResponse  = errors.New("received invalid response from server")
	ErrMessageTooBig    = errors.New("message too big")
	ErrNotAuthenticated = errors.New("not authenticated")
	ErrShuttingDown     = errors.New("server is shutting down")
	ErrForceClose       = errors.New("server has forcibly closed the connection")
	ErrNotAuthorized    = errors.New("not authorized for the operation")
)

// Options available to configure the connection.
type Options struct {
	// Username for authenticating with the server if not specified with
	// the URL.
	Username string

	// Password for authenticating with the server if not specified with
	// the URL.
	Password string

	// ClientID specifies an optional client identifier if not specified with
	// the URL. The server will generate a client identifier if one is not
	// specified.
	ClientID string

	// TLSConfig specifies the TLS configuration to use when creating a secure
	// connection to the server.
	TLSConfig *tls.Config

	// Timeout specifies the duration for a synchronous operation with the
	// server to complete. The default is 2 seconds.
	Timeout time.Duration

	// HandshakeTimeout specifies the duration for the websocket handshake
	// with the server to complete. The default is 10 seconds.
	HandshakeTimeout time.Duration
}

// Connection represents a connection to the server.
type Connection struct {
	URL         *url.URL
	Options     Options
	ErrorChan   chan error
	reconnectID string
	wg          sync.WaitGroup
	mu          sync.Mutex
	ws          *websocket.Conn
	connected   bool
	pubs        []*Completion
	pubSeqNum   int64
	subs        map[string]*Subscription
	subSeqNum   int64
	lastSeqNum  int64
}

// Subscription represents an interest in application messages.
// When returned from an asynchronous subscribe operation a non-nil
// Error indicates a subscription failure.
type Subscription struct {
	Matcher          string
	Durable          string
	MessageChan      chan Message
	Error            error
	subscriptionID   string
	subscriptionChan chan *Subscription
}

// Completion represents a completed publish operation. When returned
// from an asynchronous publish operation a non-nil Error indicates
// a publish failure.
type Completion struct {
	Message        Message
	Error          error
	seqNum         int64
	completionChan chan *Completion
}

// subprotocol used for websocket communications.
const subprotocol = "v1.eftl.tibco.com"

// op codes
const (
	opHeartbeat    = 0
	opLogin        = 1
	opWelcome      = 2
	opSubscribe    = 3
	opSubscribed   = 4
	opUnsubscribe  = 5
	opUnsubscribed = 6
	opEvent        = 7
	opPublish      = 8
	opAck          = 9
	opError        = 10
	opDisconnect   = 11
)

// defaults
const (
	defaultHandshakeTimeout = 10 * time.Second
	defaultTimeout          = 2 * time.Second
)

// Connect establishes a connection to the server at the specified url.
//
// The url can be in either of these forms:
//   ws://host:port/channel
//   wss://host:port/channel
//
// Optionally, the url can contain the username, password, and/or client identifier:
//   ws://username:password@host:port/channel?clientId=<identifier>
//   wss://username:password@host:port/channel?clientId=<identifier>
//
// Connect is a synchronous operation and will block until either a
// connection has been establish or an error occurs.
//
// The errorChan is used to receive asynchronous connection errors
// once the connection has been established.
func Connect(urlStr string, opts *Options, errorChan chan error) (*Connection, error) {
	if opts == nil {
		opts = &Options{}
	}
	url, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	// initialize the connection
	conn := &Connection{
		URL:       url,
		Options:   *opts,
		ErrorChan: errorChan,
		pubs:      make([]*Completion, 0),
		subs:      make(map[string]*Subscription),
	}
	// set default values
	if conn.Options.HandshakeTimeout == 0 {
		conn.Options.HandshakeTimeout = defaultHandshakeTimeout
	}
	if conn.Options.Timeout == 0 {
		conn.Options.Timeout = defaultTimeout
	}
	// connect to the server
	err = conn.connect()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Reconnect re-establishes the connection to the server following a
// connection error or a disconnect. Reconnect is a synchronous operation
// and will block until either a connection has been established or an
// error occurs. Upon success subscriptions are re-established.
func (conn *Connection) Reconnect() error {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if conn.connected {
		return nil
	}
	// connect to the server
	return conn.connect()
}

// Disconnect closes the connection to the server.
func (conn *Connection) Disconnect() {
	conn.mu.Lock()
	if !conn.connected {
		conn.mu.Unlock()
		return
	}
	// send disconnect message
	conn.sendMessage(Message{
		"op": opDisconnect,
	})
	// close the connection to the server
	conn.disconnect()
	conn.mu.Unlock()
	// wait for the dispatcher go routine to end
	conn.wg.Wait()
}

// IsConnected test if the connection is connected to the server.
func (conn *Connection) IsConnected() bool {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	return conn.connected
}

// Publish an application message.
//
// It is recommended to publish messages to a specific destination
// by including the string field "_dest":
//
//     conn.Publish(Message{
//         "_dest": "sample",
//         "text": "Hello, World!",
//     })
//
func (conn *Connection) Publish(msg Message) error {
	completionChan := make(chan *Completion, 1)
	if err := conn.PublishAsync(msg, completionChan); err != nil {
		return err
	}
	select {
	case completion := <-completionChan:
		return completion.Error
	case <-time.After(conn.Options.Timeout):
		return ErrTimeout
	}
}

// PublishAsync an application message asynchronously. The optional
// completionChan will receive notification once the publish
// operation completes.
//
// It is recommended to publish messages to a specific destination
// by including the string field "_dest":
//
//     conn.Publish(Message{
//         "_dest": "sample",
//         "text": "Hello, World!",
//     })
//
func (conn *Connection) PublishAsync(msg Message, completionChan chan *Completion) error {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if !conn.connected {
		return ErrNotConnected
	}
	// register the publish
	conn.pubSeqNum++
	conn.pubs = append(conn.pubs, &Completion{
		Message:        msg,
		seqNum:         conn.pubSeqNum,
		completionChan: completionChan,
	})
	// send publish message
	return conn.sendMessage(Message{
		"op":   opPublish,
		"seq":  conn.pubSeqNum,
		"body": msg,
	})
}

// Subscribe registers interest in application messages.
// A content matcher can be used to register interest in certain messages.
// A durable name can be specified to create a durable subscription.
// Messages are received on the messageChan.
//
// It is recommended to subscribe to messages published to a specific
// destination by creating a content matcher with the string field "_dest":
//
//     conn.Subscribe("{\"_dest\": \"sample\"}", "", messageChan)
//
func (conn *Connection) Subscribe(matcher string, durable string, messageChan chan Message) (*Subscription, error) {
	subscriptionChan := make(chan *Subscription, 1)
	if err := conn.SubscribeAsync(matcher, durable, messageChan, subscriptionChan); err != nil {
		return nil, err
	}
	select {
	case sub := <-subscriptionChan:
		return sub, sub.Error
	case <-time.After(conn.Options.Timeout):
		return nil, ErrTimeout
	}
}

// SubscribeAsync registers interest in application messages asynchronously.
// A content matcher can be used to register interest in certain messages.
// A durable name can be specified to create a durable subscription.
// Messages are received on the messageChan. The subscriptionChan
// will receive notification once the subscribe operation completes.
//
// It is recommended to subscribe to messages published to a specific
// destination by creating a content matcher with the string field "_dest":
//
//     conn.Subscribe("{\"_dest\": \"sample\"}", "", messageChan)
//
func (conn *Connection) SubscribeAsync(matcher string, durable string, messageChan chan Message, subscriptionChan chan *Subscription) error {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if !conn.connected {
		return ErrNotConnected
	}
	// register the subscription
	conn.subSeqNum++
	sid := strconv.FormatInt(conn.subSeqNum, 10)
	sub := &Subscription{
		Matcher:          matcher,
		Durable:          durable,
		MessageChan:      messageChan,
		subscriptionID:   sid,
		subscriptionChan: subscriptionChan,
	}
	conn.subs[sid] = sub
	// send subscribe protocol
	msg := Message{
		"op": opSubscribe,
		"id": sid,
	}
	if matcher != "" {
		msg["matcher"] = matcher
	}
	if durable != "" {
		msg["durable"] = durable
	}
	return conn.sendMessage(msg)
}

// Unsubscribe unregisters the subscription.
func (conn *Connection) Unsubscribe(sub *Subscription) error {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if !conn.connected {
		return ErrNotConnected
	}
	// send unsubscribe protocol
	conn.sendMessage(Message{
		"op": opUnsubscribe,
		"id": sub.subscriptionID,
	})
	// unregister the subscription
	delete(conn.subs, sub.subscriptionID)
	return nil
}

// UnsubscribeAll unregisters all subscriptions.
func (conn *Connection) UnsubscribeAll() error {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if !conn.connected {
		return ErrNotConnected
	}
	for _, sub := range conn.subs {
		// send unsubscribe protocol
		conn.sendMessage(Message{
			"op": opUnsubscribe,
			"id": sub.subscriptionID,
		})
		// unregister the subscription
		delete(conn.subs, sub.subscriptionID)
	}
	return nil
}

func (conn *Connection) connect() error {
	// create websocket connection
	d := &websocket.Dialer{
		HandshakeTimeout: conn.Options.HandshakeTimeout,
		Subprotocols:     []string{subprotocol},
		TLSClientConfig:  conn.Options.TLSConfig,
	}
	u := &url.URL{
		Scheme: conn.URL.Scheme,
		Host:   conn.URL.Host,
		Path:   conn.URL.Path,
	}
	ws, _, err := d.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	conn.ws = ws
	// send login message
	msg := Message{
		"op":             opLogin,
		"client_type":    "golang",
		"client_version": Version,
		"login_options": Message{
			"_qos":    "true",
			"_resume": "true",
		},
	}
	if conn.URL.User != nil {
		msg["user"] = conn.URL.User.Username()
	} else if conn.Options.Username != "" {
		msg["user"] = conn.Options.Username
	}
	if conn.URL.User != nil {
		msg["password"], _ = conn.URL.User.Password()
	} else if conn.Options.Password != "" {
		msg["password"] = conn.Options.Password
	}
	if conn.URL.Query().Get("clientId") != "" {
		msg["client_id"] = conn.URL.Query().Get("clientId")
	} else if conn.Options.ClientID != "" {
		msg["client_id"] = conn.Options.ClientID
	}
	if conn.reconnectID != "" {
		msg["id_token"] = conn.reconnectID
	}
	err = conn.sendMessage(msg)
	if err != nil {
		conn.ws.Close()
		return err
	}
	// set a read deadline
	conn.ws.SetReadDeadline(time.Now().Add(conn.Options.Timeout))
	defer conn.ws.SetReadDeadline(time.Time{})
	// receive welcome message
	msg, err = conn.nextMessage()
	if err != nil {
		conn.ws.Close()
		return err
	}
	// op code
	if op, ok := msg["op"].(int64); !ok || op != opWelcome {
		conn.ws.Close()
		return ErrInvalidResponse
	}
	// client id
	if val, ok := msg["client_id"].(string); ok {
		conn.Options.ClientID = val
	}
	// token id
	if val, ok := msg["id_token"].(string); ok {
		conn.reconnectID = val
	}
	// resume
	resume := false
	if val, ok := msg["_resume"].(string); ok {
		resume, _ = strconv.ParseBool(val)
	}
	// mark the connection as connected
	conn.connected = true
	// re-establish subscriptions
	for _, sub := range conn.subs {
		msg := Message{
			"op": opSubscribe,
			"id": sub.subscriptionID,
		}
		if sub.Matcher != "" {
			msg["matcher"] = sub.Matcher
		}
		if sub.Durable != "" {
			msg["durable"] = sub.Durable
		}
		conn.sendMessage(msg)
	}
	if resume {
		// re-send unacknowledged messages
		for _, comp := range conn.pubs {
			conn.sendMessage(Message{
				"op":   opPublish,
				"seq":  comp.seqNum,
				"body": comp.Message,
			})
		}
	} else {
		conn.lastSeqNum = 0
	}
	// process incoming messages
	conn.wg.Add(1)
	go conn.dispatch()
	return nil
}

func (conn *Connection) disconnect() error {
	// mark the connection as not connected
	conn.connected = false
	// disconnect from the server
	return conn.ws.Close()
}

func (conn *Connection) dispatch() {
	defer conn.wg.Done()
	for {
		// read the next message
		msg, err := conn.nextMessage()
		if err != nil {
			conn.handleDisconnect(err)
			break
		}
		// process the message
		if op, ok := msg["op"].(int64); ok {
			switch op {
			case opHeartbeat:
				conn.handleHeartbeat(msg)
			case opEvent:
				conn.handleMessage(msg)
			case opSubscribed:
				conn.handleSubscribed(msg)
			case opUnsubscribed:
				conn.handleUnsubscribed(msg)
			case opAck:
				conn.handleAck(msg)
			case opError:
			}
		}
	}
}

func (conn *Connection) handleDisconnect(err error) {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if conn.connected {
		conn.disconnect()
		// send notification to the error channel
		if conn.ErrorChan != nil {
			conn.ErrorChan <- err
		}
		// clear pending completions
		for i := 0; i < len(conn.pubs); i++ {
			comp := conn.pubs[i]
			conn.pubs[i] = nil
			comp.Error = err
			if comp.completionChan != nil {
				comp.completionChan <- comp
			}
		}
		conn.pubs = conn.pubs[:0]
	}
}

func (conn *Connection) handleHeartbeat(msg Message) {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	conn.sendMessage(msg)
}

func (conn *Connection) handleMessage(msg Message) {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	seq, _ := msg["seq"].(int64)
	body, _ := msg["body"].(Message)
	if sid, ok := msg["to"].(string); ok {
		if seq == 0 || seq > conn.lastSeqNum {
			if sub, ok := conn.subs[sid]; ok {
				if sub.MessageChan != nil {
					sub.MessageChan <- body
				}
			}
			if seq > 0 {
				conn.lastSeqNum = seq
			}
		}
	}
	if seq > 0 {
		// acknowledge message receipt
		conn.sendMessage(Message{
			"op":  opAck,
			"seq": seq,
		})
	}
}

func (conn *Connection) handleSubscribed(msg Message) {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if sid, ok := msg["id"].(string); ok {
		if sub, ok := conn.subs[sid]; ok {
			if sub.subscriptionChan != nil {
				sub.subscriptionChan <- sub
			}
		}
	}
}

func (conn *Connection) handleUnsubscribed(msg Message) {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if sid, ok := msg["id"].(string); ok {
		if sub, ok := conn.subs[sid]; ok {
			errCode, _ := msg["err"].(int64)
			if errCode == 12 {
				sub.Error = ErrNotAuthorized
			} else {
				reason, _ := msg["reason"].(string)
				sub.Error = fmt.Errorf("%d: %s", errCode, reason)
			}
			delete(conn.subs, sid)
			if sub.subscriptionChan != nil {
				sub.subscriptionChan <- sub
			}
		}
	}
}

func (conn *Connection) handleAck(msg Message) {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if seq, ok := msg["seq"].(int64); ok {
		var err error
		if errCode, ok := msg["err"].(int64); ok {
			if errCode == 12 {
				err = ErrNotAuthorized
			} else {
				reason, _ := msg["reason"].(string)
				err = fmt.Errorf("%d: %s", errCode, reason)
			}
		}
		k := 0
		for _, comp := range conn.pubs {
			if comp.seqNum <= seq {
				comp.Error = err
				if comp.completionChan != nil {
					comp.completionChan <- comp
				}
			} else {
				conn.pubs[k] = comp
				k++
			}
		}
		conn.pubs = conn.pubs[:k]
	}
}

func (conn *Connection) sendMessage(msg Message) error {
	return conn.ws.WriteJSON(msg)
}

func (conn *Connection) nextMessage() (msg Message, err error) {
	msg = make(Message)
	err = conn.ws.ReadJSON(&msg)
	// translate a websocket.CloseError
	if closeErr, ok := err.(*websocket.CloseError); ok {
		switch closeErr.Code {
		case 1001:
			err = ErrShuttingDown
		case 1009:
			err = ErrMessageTooBig
		case 4000:
			err = ErrForceClose
		case 4002:
			err = ErrNotAuthenticated
		default:
			err = errors.New(closeErr.Text)
		}
	}
	return
}
