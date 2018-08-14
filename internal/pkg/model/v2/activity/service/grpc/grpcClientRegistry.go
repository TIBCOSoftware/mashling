package grpc

import (
	//used for generated stub files
	_ "github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

var ClientServiceRegistery = NewServiceRegistry()

type ClientService interface {
	ServiceInfo() *ServiceInfo
	GetRegisteredClientService(cs *grpc.ClientConn) interface{}
	InvokeMethod(map[string]interface{}) map[string]interface{}
}

type ServiceInfo struct {
	ServiceName string
	ProtoName   string
}

type ServiceRegistry struct {
	ClientServices map[string]ClientService
}

// NewServiceRegistry creates a new service registry
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{ClientServices: make(map[string]ClientService)}
}

//RegisterClientService resgisters service
func (sr *ServiceRegistry) RegisterClientService(service ClientService) {
	sr.ClientServices[service.ServiceInfo().ProtoName+service.ServiceInfo().ServiceName] = service
}
