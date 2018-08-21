package grpc

import (
	"errors"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/logger"

	"google.golang.org/grpc/credentials"

	"github.com/imdario/mergo"
	"google.golang.org/grpc"
)

// log is the default package logger
var log = logger.GetLogger("tibco-service-grpc")

// GRPC is grpc service
type GRPC struct {
	Request  GRPCRequest  `json:"request"`
	Response GRPCResponse `json:"response"`
}

// GRPCRequest is grpc service request
type GRPCRequest struct {
	HostURL          string                 `json:"hosturl"`
	GrpcMthdParamtrs map[string]interface{} `json:"grpcMthdParamtrs"`
	EnableTLS        string                 `json:"enableTLS"`
	ClientCert       string                 `json:"clientCert"`
	Header           map[string]string      `json:"header"`
	PathParams       map[string]string      `json:"pathParams"`
	OperatingMode    string                 `json:"operatingMode"`
	ServiceName      string                 `json:"serviceName"`
	ProtoName        string                 `json:"protoName"`
	MethodName       string                 `json:"methodName"`
	Params           map[string]string      `json:"params"`
	QueryParams      map[string]string      `json:"queryParams"`
	Content          interface{}            `json:"content"`
}

// GRPCResponse is grpc service response
type GRPCResponse struct {
	Body interface{} `json:"body"`
}

// InitializeGRPC  initialize GRPC service with provided settings.
func InitializeGRPC(settings map[string]interface{}) (grpcService *GRPC, err error) {
	grpc := &GRPC{}
	req := GRPCRequest{}
	req.PathParams = make(map[string]string)
	req.Header = make(map[string]string)
	req.GrpcMthdParamtrs = make(map[string]interface{})
	req.Params = make(map[string]string)
	req.QueryParams = make(map[string]string)
	grpc.Request = req
	err = grpc.setRequestValues(settings)
	return grpc, err
}

// Execute invokes this GRPC service.
func (g *GRPC) Execute() (err error) {

	g.Response = GRPCResponse{}

	opts := []grpc.DialOption{}
	log.Debug("enableTLS: ", g.Request.EnableTLS)
	if strings.Compare(g.Request.EnableTLS, "true") == 0 {
		log.Debug("ClientCert: ", g.Request.ClientCert)
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

	log.Debug("operating mode: ", g.Request.OperatingMode)

	switch g.Request.OperatingMode {
	case "grpc-to-grpc":
		return gRPCTogRPCHandler(g, conn)
	case "rest-to-grpc":
		return restTogRPCHandler(g, conn)
	}

	log.Error("Invalid use of service , OperatingMode not recognised")
	return errors.New("Invalid use of service , OperatingMode not recognised")
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
			g.Request.OperatingMode = "grpc-to-grpc"
			grpcData, ok := v.(map[string]interface{})
			if !ok {
				return errors.New("invalid type for grpcData")
			}
			if err := mergo.Merge(&g.Request.GrpcMthdParamtrs, grpcData, mergo.WithOverride); err != nil {
				return errors.New("unable to merge params values")
			}
		case "header":
			if len(g.Request.OperatingMode) == 0 {
				g.Request.OperatingMode = "rest-to-grpc"
			}
			header, ok := v.(map[string]string)
			if !ok {
				return errors.New("invalid type for header")
			}
			if err := mergo.Merge(&g.Request.Header, header, mergo.WithOverride); err != nil {
				return errors.New("unable to merge header values")
			}
		case "serviceName":
			name, ok := v.(string)
			if !ok {
				return errors.New("invalid type for serviceName")
			}
			g.Request.ServiceName = name
		case "protoName":
			name, ok := v.(string)
			if !ok {
				return errors.New("invalid type for protoName")
			}
			g.Request.ProtoName = name
		case "methodName":
			name, ok := v.(string)
			if !ok {
				return errors.New("invalid type for methodName")
			}
			g.Request.MethodName = name
		case "params":
			if len(g.Request.OperatingMode) == 0 {
				g.Request.OperatingMode = "rest-to-grpc"
			}
			params, ok := v.(map[string]string)
			if !ok {
				return errors.New("invalid type for params")
			}
			if err := mergo.Merge(&g.Request.Params, params, mergo.WithOverride); err != nil {
				return errors.New("unable to merge params values")
			}
		case "queryParams":
			if len(g.Request.OperatingMode) == 0 {
				g.Request.OperatingMode = "rest-to-grpc"
			}
			queryParams, ok := v.(map[string]string)
			if !ok {
				return errors.New("invalid type for queryParams")
			}
			if err := mergo.Merge(&g.Request.QueryParams, queryParams, mergo.WithOverride); err != nil {
				return errors.New("unable to merge queryParams values")
			}
		case "content":
			if len(g.Request.OperatingMode) == 0 {
				g.Request.OperatingMode = "rest-to-grpc"
			}
			g.Request.Content = v
		case "pathParams":
			if len(g.Request.OperatingMode) == 0 {
				g.Request.OperatingMode = "rest-to-grpc"
			}
			pathParams, ok := v.(map[string]string)
			if !ok {
				return errors.New("invalid type for pathParams")
			}
			if err := mergo.Merge(&g.Request.PathParams, pathParams, mergo.WithOverride); err != nil {
				return errors.New("unable to merge pathParams values")
			}
		case "operatingMode":
			mode, ok := v.(string)
			if !ok {
				return errors.New("invalid type for operatingMode")
			}
			g.Request.OperatingMode = mode
		default:
			// ignore and move on.
		}
	}
	return nil
}
