package grpc

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/TIBCOSoftware/mashling/lib/util"
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

				if g.Request.GrpcMthdParamtrs["contextdata"] != nil {

					inputs := make([]reflect.Value, 2)

					inputs[0] = reflect.ValueOf(g.Request.GrpcMthdParamtrs["contextdata"])
					inputs[1] = reflect.ValueOf(g.Request.GrpcMthdParamtrs["reqdata"])

					resultArr := reflect.ValueOf(clientInterfaceObj).MethodByName(g.Request.GrpcMthdParamtrs["methodName"].(string)).Call(inputs)

					res := resultArr[0]
					grpcErr := resultArr[1]
					if !grpcErr.IsNil() {
						erroString := fmt.Sprintf("%v", grpcErr.Interface())
						log.Error("Propagating error to calling function:", erroString)
						erroString = "{\"error\":\"true\",\"details\":{\"error\":\"" + erroString + "\"}}"
						err := util.Unmarshal("application/json", []byte(erroString), &g.Response.Body)
						if err != nil {
							return err
						}
					} else {
						g.Response.Body = res.Interface()
					}
				} else {
					InvokeMethodData := make(map[string]interface{})
					InvokeMethodData["ClientObject"] = clientInterfaceObj
					InvokeMethodData["MethodName"] = g.Request.GrpcMthdParamtrs["methodName"]
					InvokeMethodData["reqdata"] = g.Request.GrpcMthdParamtrs["reqdata"]
					InvokeMethodData["strmReq"] = g.Request.GrpcMthdParamtrs["strmReq"]

					resMap := service.InvokeMethod(InvokeMethodData)

					if resMap["Error"] != nil {
						log.Errorf("Error occured:%v", resMap["Error"])
						erroString := fmt.Sprintf("%v", resMap["Error"])
						erroString = "{\"error\":\"true\",\"details\":{\"error\":\"" + erroString + "\"}}"
						err := util.Unmarshal("application/json", []byte(erroString), &g.Response.Body)
						if err != nil {
							return err
						}
					}

				}

			}
		}
		if !clServFlag {
			log.Errorf("client service object not found for proto [%v] and service [%v]", protoName, serviceName)
		}
	} else {
		log.Errorf("gRPC Client services not registered")
	}
	return nil
}
