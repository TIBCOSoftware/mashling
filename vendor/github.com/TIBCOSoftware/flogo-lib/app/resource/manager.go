package resource

import (
	"errors"
	"strings"
)

// ResourceManager interface
type Manager interface {
	// LoadResources tells the manager to load the specified resource set
	LoadResource(config *Config) error

	// GetResource get the resource that corresponds to the specified id
	GetResource(id string) interface{}
}

var managers = make(map[string]Manager)

// RegisterManager registers a resource manager for the specified type
func RegisterManager(resourceType string, manager Manager) error {

	_, exists := managers[resourceType]

	if exists {
		return errors.New("Resource Manager already registered for type: " + resourceType)
	}

	managers[resourceType] = manager
	return nil
}

// GetManager gets the manager for the specified resource type
func GetManager(resourceType string) Manager {
	return managers[resourceType]
}

// Load specified resource into its corresponding Resource Manager
func Load(config *Config) error {
	resType, err := GetTypeFromID(config.ID)
	if err != nil {
		return err
	}

	manager := GetManager(resType)

	if manager == nil {
		return errors.New("unsupported resource type: " + resType)
	}

	return manager.LoadResource(config)
}

// Get gets the specified resource, id format is {type}:{id}"
func Get(id string) (interface{}, error) {

	resType, err := GetTypeFromID(id)
	if err != nil {
		return nil, err
	}

	manager := GetManager(resType)

	if manager == nil {
		return nil, errors.New("unsupported resource type: " + resType)
	}

	return manager.GetResource(id), nil
}

func GetTypeFromID(id string) (string, error) {

	idx := strings.Index(id, ":")

	if idx < 0 {
		return "", errors.New("Invalid resource id: " + id)
	}

	return id[:idx], nil
}
