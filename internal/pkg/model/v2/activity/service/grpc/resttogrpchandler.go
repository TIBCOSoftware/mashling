package grpc

import (
	"errors"
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/mashling/lib/util"
	"google.golang.org/grpc"
)

func restTogRPCHandler(g *GRPC, conn *grpc.ClientConn) error {
	// check for method name
	if len(g.Request.MethodName) == 0 {
		if len(g.Request.PathParams["grpcMethodName"]) == 0 {
			log.Error("Method name not provided in json/pathParams")
			return errors.New("Method name not provided")
		}
		g.Request.MethodName = g.Request.PathParams["grpcMethodName"]
		log.Debug("Method name: ", g.Request.MethodName)
	}

	clServFlag := false
	if len(ClientServiceRegistery.ClientServices) != 0 {
		for k, service := range ClientServiceRegistery.ClientServices {
			if strings.Compare(k, g.Request.ProtoName+g.Request.ServiceName) == 0 {
				log.Debugf("client service object found for proto [%v] and service [%v]", g.Request.ProtoName, g.Request.ServiceName)

				InvokeMethodData := make(map[string]interface{})
				InvokeMethodData["ClientObject"] = service.GetRegisteredClientService(conn)
				InvokeMethodData["MethodName"] = g.Request.MethodName
				if len(g.Request.PathParams) != 0 {
					InvokeMethodData["PathParams"] = g.Request.PathParams
				}
				if len(g.Request.Params) != 0 {
					InvokeMethodData["Params"] = g.Request.Params
				}
				if len(g.Request.QueryParams) != 0 {
					InvokeMethodData["QueryParams"] = g.Request.QueryParams
				}
				if g.Request.Content != nil {
					InvokeMethodData["Content"] = g.Request.Content
				}

				resMap := service.InvokeMethod(InvokeMethodData)
				if strings.Compare(string(resMap["Response"].([]byte)), "null") != 0 {
					err := util.Unmarshal("application/json", resMap["Response"].([]byte), &g.Response.Body)
					if err != nil {
						return err
					}
				} else {
					erroString := fmt.Sprintf("%v", resMap["Error"])
					erroString = "{\"error\":\"true\",\"details\":{\"error\":\"" + erroString + "\"}}"
					err := util.Unmarshal("application/json", []byte(erroString), &g.Response.Body)
					if err != nil {
						return err
					}
				}
				clServFlag = true
			}
		}
		if !clServFlag {
			log.Errorf("client service object not found for proto [%v] and service [%v]", g.Request.ProtoName, g.Request.ServiceName)
		}
	} else {
		log.Errorf("gRPC Client services not registered")
	}

	return nil
}
