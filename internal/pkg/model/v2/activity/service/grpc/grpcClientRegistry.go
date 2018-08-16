package grpc

import (
	//used for generated stub files
	_ "github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

// ClientServiceRegistery holds all the client services written in proto file
var ClientServiceRegistery = NewServiceRegistry()

// ClientService interface used by grpc client side
type ClientService interface {
	ServiceInfo() *ServiceInfo
	GetRegisteredClientService(cs *grpc.ClientConn) interface{}
	InvokeMethod(map[string]interface{}) map[string]interface{}
}

// ServiceInfo holds name of service and name of proto
type ServiceInfo struct {
	ServiceName string
	ProtoName   string
}

// ServiceRegistry data structure to hold the client services
type ServiceRegistry struct {
	ClientServices map[string]ClientService
}

// NewServiceRegistry creates a new service registry
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{ClientServices: make(map[string]ClientService)}
}

// RegisterClientService resgisters service
func (sr *ServiceRegistry) RegisterClientService(service ClientService) {
	sr.ClientServices[service.ServiceInfo().ProtoName+service.ServiceInfo().ServiceName] = service
}
