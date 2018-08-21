package grpc

import (
	"errors"
	"reflect"
	"strings"

	"google.golang.org/grpc"
)

var clientInterfaceObj interface{}

func gRPCTogRPCHandler(g *GRPC, conn *grpc.ClientConn) error {

	serviceName := g.Request.GrpcMthdParamtrs["serviceName"].(string)
	protoName := g.Request.GrpcMthdParamtrs["protoName"].(string)
	protoName = strings.Split(protoName, ".")[0]

	if len(serviceName) == 0 && len(protoName) == 0 {
		return errors.New("Service name and Proto name required")
	}

	clServFlag := false
	if len(ClientServiceRegistery.ClientServices) != 0 {
		for k, service := range ClientServiceRegistery.ClientServices {
			if strings.Compare(k, protoName+serviceName) == 0 {
				log.Debugf("client service object found for proto [%v] and service [%v]", protoName, serviceName)
				clientInterfaceObj = service.GetRegisteredClientService(conn)
				clServFlag = true
			}
		}
		if !clServFlag {
			log.Errorf("client service object not found for proto [%v] and service [%v]", protoName, serviceName)
		}
	} else {
		log.Errorf("gRPC Client services not registered")
	}

	inputs := make([]reflect.Value, 2)

	inputs[0] = reflect.ValueOf(g.Request.GrpcMthdParamtrs["contextdata"])
	inputs[1] = reflect.ValueOf(g.Request.GrpcMthdParamtrs["reqdata"])

	resultArr := reflect.ValueOf(clientInterfaceObj).MethodByName(g.Request.GrpcMthdParamtrs["methodName"].(string)).Call(inputs)

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
