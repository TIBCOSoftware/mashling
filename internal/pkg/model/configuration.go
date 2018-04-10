package model

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/TIBCOSoftware/mashling/internal/pkg/consul"
	gwerrors "github.com/TIBCOSoftware/mashling/internal/pkg/model/errors"
	v1 "github.com/TIBCOSoftware/mashling/internal/pkg/model/v1"
	v2 "github.com/TIBCOSoftware/mashling/internal/pkg/model/v2"
)

// Gateway represents an instance of a specific version of a Gateway.
type Gateway interface {
	Init(string) error
	Start() error
	Stop() error
	Version() string
	AppVersion() string
	Name() string
	Description() string
	Errors() []gwerrors.Error
	Configuration() interface{}
	Swagger(hostname string, triggerName string) ([]byte, error)
	ConsulServiceDefinition() ([]consul.ServiceDefinition, error)
}

type schema struct {
	Version string `json:"mashling_schema"`
}

// LoadFromFile loads a Gateway configuration from the specified file.
func LoadFromFile(file string) (Gateway, error) {
	configuration, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return Load(configuration)
}

// LoadFromEnv loads a Gateway configuration from the specified environment variable.
func LoadFromEnv(envVarName string) (Gateway, error) {
	var configuration []byte
	b64GatewayJSON := os.Getenv(envVarName)
	if b64GatewayJSON == "" {
		return nil, errors.New("configuration environment variable is empty")
	}
	configuration, err := base64.StdEncoding.DecodeString(b64GatewayJSON)
	if err != nil {
		return nil, err
	}
	return Load(configuration)
}

// Load loads a Gateway configuration from the provided byte array.
func Load(configuration []byte) (Gateway, error) {
	schemaVersion, err := version(configuration)
	if err != nil {
		return nil, err
	}
	if schemaVersion == v1.Version {
		return v1.LoadGateway(configuration)
	} else if schemaVersion == v2.Version {
		return v2.LoadGateway(configuration)
	} else {
		return nil, errors.New("schema version not found")
	}
}

func version(gatewayJSON []byte) (string, error) {
	schema := &schema{}
	err := json.Unmarshal(gatewayJSON, schema)
	return schema.Version, err
}
