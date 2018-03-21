package awsiot

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// log is the default package logger
var log = logger.GetLogger("activity-tibco-rest")

const (
	ivThingName   = "thingName"
	ivAwsEndpoint = "awsEndpoint"
	ivDesired     = "desired"
	ivReported    = "reported"

	ovResult = "result"
)

// AwsIoT is an Activity that is used to update an Aws IoT device shadow
// inputs : {method,uri,params}
// outputs: {result}
type AwsIoT struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AwsIoT activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &AwsIoT{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *AwsIoT) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Invokes a Aws Iot Shadow Update
func (a *AwsIoT) Eval(context activity.Context) (done bool, err error) {

	thingName := context.GetInput(ivThingName).(string)
	awsEndpoint := context.GetInput(ivAwsEndpoint).(string)

	req := &ShadowRequest{State: &ShadowState{}}

	if context.GetInput(ivDesired) != nil {
		desired := context.GetInput(ivDesired).(map[string]string)
		req.State.Desired = desired
	}

	if context.GetInput(ivReported) != nil {
		reported := context.GetInput(ivReported).(map[string]string)
		req.State.Reported = reported
	}

	reqJSON, err := json.Marshal(req)

	if err != nil {
		return false, activity.NewError(err.Error(), "", nil)
	}

	log.Debugf("Shadow Request: %s", string(reqJSON))

	brokerURI := fmt.Sprintf("ssl://%s:%d", awsEndpoint, 8883)
	log.Debugf("Broker URI: %s", brokerURI)

	tlsConfig := NewTLSConfig(thingName)

	opts := MQTT.NewClientOptions()
	opts.AddBroker(brokerURI)
	opts.SetClientID(context.FlowDetails().ID())
	opts.SetTLSConfig(tlsConfig)

	// Start the connection
	client := MQTT.NewClient(opts)
	defer client.Disconnect(250)

	token := client.Connect()

	if token.Wait() && token.Error() != nil {
		log.Errorf("Error connecting to '%s': %s", brokerURI, token.Error().Error())
		return false, activity.NewError(token.Error().Error(), "", nil)
	}

	thingUpdate := fmt.Sprintf("$aws/things/%s/shadow/update", thingName)
	Publish(client, thingUpdate, 1, string(reqJSON))

	return true, nil
}

////////////////////////////////////////////////////////////////////////////////////////
// Utils

// Publish publishes a client message
func Publish(client MQTT.Client, topic string, qos int, input string) error {
	token := client.Publish(topic, byte(qos), false, input)
	if token.Wait() && token.Error() != nil {
		log.Error(token.Error())
		return token.Error()
	}
	return nil
}

// NewTLSConfig creates a TLS configuration for the specified 'thing'
func NewTLSConfig(thingName string) *tls.Config {
	// Import root CA
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile("things/root-CA.pem.crt")
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	thingDir := "things/" + thingName + "/"

	// Import client certificate/key pair for the specified 'thing'
	cert, err := tls.LoadX509KeyPair(thingDir+"device.pem.crt", thingDir+"device.pem.key")
	if err != nil {
		panic(err)
	}

	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		panic(err)
	}

	return &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	}
}

// ShadowRequest is a simple structure representing a Aws Shadow Update Request
type ShadowRequest struct {
	State *ShadowState `json:"state"`
}

// ShadowState is the state to be updated
type ShadowState struct {
	Desired  map[string]string `json:"desired,omitempty"`
	Reported map[string]string `json:"reported,omitempty"`
}
