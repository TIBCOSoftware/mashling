package grpc

import (

	//used for generated stub files
	_ "github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

//ServiceRegistery holds all the services written in proto file
var ServiceRegistery = NewServiceRegistry()

//ServerService methods to invoke registartion of service
type ServerService interface {
	ServiceInfo() *ServiceInfo
	RunRegisterServerService(s *grpc.Server, t *GRPCTrigger)
}

//ServiceInfo holds name of service and name of proto
type ServiceInfo struct {
	ServiceName string
	ProtoName   string
}

//ServiceRegistry data structure to hold the services
type ServiceRegistry struct {
	ServerServices map[string]ServerService
}

// NewServiceRegistry creates a new service registry
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{ServerServices: make(map[string]ServerService)}
}

//RegisterServerService resgisters server services
func (sr *ServiceRegistry) RegisterServerService(service ServerService) {
	sr.ServerServices[service.ServiceInfo().ProtoName+service.ServiceInfo().ServiceName] = service
}
