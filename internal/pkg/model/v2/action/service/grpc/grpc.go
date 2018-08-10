package grpc

import (
	"errors"
	"reflect"
	"strings"

	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/logger"

	"google.golang.org/grpc/credentials"

	"github.com/imdario/mergo"
	"google.golang.org/grpc"
)

// log is the default package logger
var log = logger.GetLogger("tibco-service-grpc")

//GRPC is a grpc service
type GRPC struct {
	Request  GRPCRequest  `json:"request"`
	Response GRPCResponse `json:"response"`
}

//GRPCRequest is a grpc service request
type GRPCRequest struct {
	HostURL          string                 `json:"hosturl"`
	GrpcMthdParamtrs map[string]interface{} `json:"grpcMthdParamtrs"`
	EnableTLS        string                 `json:"enableTLS"`
	ClientCert       string                 `json:"clientCert"`
}

//GRPCResponse is a grpc service response
type GRPCResponse struct {
	Body interface{} `json:"body"`
}

//InitializeGRPC  initializes an GRPC service with provided settings.
func InitializeGRPC(settings map[string]interface{}) (grpcService *GRPC, err error) {
	grpc := &GRPC{}
	req := GRPCRequest{}
	req.GrpcMthdParamtrs = make(map[string]interface{})
	grpc.Request = req
	err = grpc.setRequestValues(settings)
	return grpc, err
}

// Execute invokes this GRPC service.
func (g *GRPC) Execute() (err error) {

	opts := []grpc.DialOption{}
	log.Debug("enableTLS: ", g.Request.EnableTLS)
	if strings.Compare(g.Request.EnableTLS, "true") == 0 {
		log.Debug("ClientCert: ", g.Request.ClientCert)
		fmt.Println("g.Request.ClientCert", g.Request.ClientCert)
		creds, err := credentials.NewClientTLSFromFile(g.Request.ClientCert, "")
		if err != nil {
			log.Error(err)
		}

		opts = []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	} else {
		opts = []grpc.DialOption{grpc.WithInsecure()}
	}

	conn, err := grpc.Dial(g.Request.HostURL, opts...)
	if err != nil {
		log.Error(err)
	}
	defer conn.Close()

	var clientInterfaceObj interface{}

	servicename := g.Request.GrpcMthdParamtrs["servicename"].(string)
	protoname := g.Request.GrpcMthdParamtrs["protoname"].(string)
	protoname = strings.Split(protoname, ".")[0]

	if len(servicename) == 0 && len(protoname) == 0 {
		return errors.New("Service name and Proto name required")
	}

	clServFlag := false
	if len(ClientServiceRegistery.ClientServices) != 0 {
		for k, service := range ClientServiceRegistery.ClientServices {
			if strings.Compare(k, protoname+servicename) == 0 {
				log.Debugf("client service object found for proto [%v] and service [%v]", protoname, servicename)
				clientInterfaceObj = service.GetRegisteredClientService(conn)
				clServFlag = true
			}
		}
		if !clServFlag {
			log.Errorf("client service object not found for proto [%v] and service [%v]", protoname, servicename)
		}
	} else {
		log.Errorf("gRPC Client services not registered")
	}

	inputs := make([]reflect.Value, 2)

	inputs[0] = reflect.ValueOf(g.Request.GrpcMthdParamtrs["contextdata"])
	inputs[1] = reflect.ValueOf(g.Request.GrpcMthdParamtrs["reqdata"])

	resultArr := reflect.ValueOf(clientInterfaceObj).MethodByName(g.Request.GrpcMthdParamtrs["methodname"].(string)).Call(inputs)

	res := resultArr[0]
	grpcErr := resultArr[1]
	if !grpcErr.IsNil() {
		log.Error("Propagating error to calling function")
		log.Error("Error Details: ", grpcErr.Interface())
		g.Response.Body = grpcErr.Interface()
	} else {
		g.Response.Body = res.Interface()
	}
	return nil
}

// UpdateRequest updates a request on an existing GRPC service instance with new values.
func (g *GRPC) UpdateRequest(values map[string]interface{}) (err error) {
	return g.setRequestValues(values)
}

func (g *GRPC) setRequestValues(settings map[string]interface{}) (err error) {
	for k, v := range settings {
		switch k {
		case "hosturl":
			url, ok := v.(string)
			if !ok {
				return errors.New("invalid type for url")
			}
			g.Request.HostURL = url
		case "enableTLS":
			enableTLS, ok := v.(string)
			if !ok {
				return errors.New("invalid type for enableTLS")
			}
			g.Request.EnableTLS = enableTLS
		case "clientCert":
			clientCert, ok := v.(string)
			if !ok {
				return errors.New("invalid type for clientCert")
			}
			g.Request.ClientCert = clientCert
		case "grpcMthdParamtrs":
			grpcData, ok := v.(map[string]interface{})
			if !ok {
				return errors.New("invalid type for grpcData")
			}
			if err := mergo.Merge(&g.Request.GrpcMthdParamtrs, grpcData, mergo.WithOverride); err != nil {
				return errors.New("unable to merge params values")
			}
		default:
			// ignore and move on.
		}
	}
	return nil
}
