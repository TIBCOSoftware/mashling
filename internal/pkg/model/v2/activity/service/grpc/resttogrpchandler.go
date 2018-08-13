package grpc

import (
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/mashling/lib/util"
	"google.golang.org/grpc"
)

func restTogRPCHandler(g *GRPC, conn *grpc.ClientConn) error {

	fmt.Println("@@@@@@@in grpc service ", g.Request.Headers)
	fmt.Println("@@@@@@@@in grpc service ", g.Request.PathParams)

	clServFlag := false
	if len(ClientServiceRegistery.ClientServices) != 0 {
		for k, service := range ClientServiceRegistery.ClientServices {
			if strings.Compare(k, g.Request.ProtoName+g.Request.ServiceName) == 0 {
				log.Debugf("client service object found for proto [%v] and service [%v]", g.Request.ProtoName, g.Request.ServiceName)

				InvokeMethodData := make(map[string]interface{})
				InvokeMethodData["ClientObject"] = service.GetRegisteredClientService(conn)
				InvokeMethodData["MethodName"] = g.Request.MethodName
				InvokeMethodData["PathParams"] = g.Request.PathParams

				values := service.InvokeMethod(InvokeMethodData)
				if strings.Compare(string(values[0].([]byte)), "null") != 0 {
					err := util.Unmarshal("application/json", values[0].([]byte), &g.Response.Body)
					if err != nil {
						return err
					}
				} else {
					erroString := fmt.Sprintf("%v", values[1])
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
