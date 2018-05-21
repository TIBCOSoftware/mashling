package cache

import "errors"

// Enabled determines if the cache has been configured with a valid provider.
var Enabled = false

// Cache is the initialized CacheProvider instance.
var Cache Provider

// Provider dictates the operations a CacheProvider must implement.
type Provider interface {
	Initialize(configuration interface{}) error
	InCache(key string) bool
	LoadFromCache(key string) ([]byte, error)
	WriteToCache(key string, contents []byte) error
}

// Initialize initializes the provided cache provider type with the
// configuration, sets the global provider, and marks the cache as enabled.
func Initialize(providerType string, configuration interface{}) error {
	var err error
	switch providerType {
	case "file":
		Cache = &FileCache{}
		err = Cache.Initialize(configuration)
		if err != nil {
			return err
		}
		Enabled = true
	default:
		return errors.New("unknown cache provider type")
	}
	return err
}
