package grpc

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/imdario/mergo"
	"google.golang.org/grpc"
)

//GRPC is a grpc service
type GRPC struct {
	Request  GRPCRequest  `json:"request"`
	Response GRPCResponse `json:"response"`
}

//GRPCRequest is a grpc service request
type GRPCRequest struct {
	HostURL          string                 `json:"hosturl"`
	GrpcMthdParamtrs map[string]interface{} `json:"grpcMthdParamtrs"`
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
	opts := []grpc.DialOption{grpc.WithInsecure()}

	conn, err := grpc.Dial(g.Request.HostURL, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var clientInterfaceObj interface{}

	servicename := g.Request.GrpcMthdParamtrs["servicename"].(string)
	protoname := g.Request.GrpcMthdParamtrs["protoname"].(string)
	protoname = strings.Split(protoname, ".")[0]

	if len(servicename) == 0 && len(protoname) == 0 {
		return errors.New("Service name and Proto name required")
	}

	if len(ClientServiceRegistery.ClientServices) != 0 {
		for k, service := range ClientServiceRegistery.ClientServices {
			if strings.Compare(k, protoname+servicename) == 0 {
				fmt.Println("*********** getting client service for proto ************", protoname, servicename)
				clientInterfaceObj = service.GetRegisteredClientService(conn)
			}
		}
	} else {
		log.Println("***********client Services not available for proto ***********", protoname, servicename)
		log.Fatal(errors.New("client stubs not registered"))
	}

	inputs := make([]reflect.Value, 2)

	inputs[0] = reflect.ValueOf(g.Request.GrpcMthdParamtrs["contextdata"])
	inputs[1] = reflect.ValueOf(g.Request.GrpcMthdParamtrs["reqdata"])

	resultArr := reflect.ValueOf(clientInterfaceObj).MethodByName(g.Request.GrpcMthdParamtrs["methodname"].(string)).Call(inputs)

	res := resultArr[0]
	grpcErr := resultArr[1]
	if !grpcErr.IsNil() {
		return errors.New("Error occured in grpc service :" + fmt.Sprintf("%v", grpcErr))
	}

	g.Response.Body = res.Interface()

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