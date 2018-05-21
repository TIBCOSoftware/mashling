package cache

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// FileCache is a cache provider backed by the local file system.
type FileCache struct {
	directory string
}

// Initialize inializes the FileCache provider.
func (f *FileCache) Initialize(configuration interface{}) error {
	cacheLocation, ok := configuration.(string)
	if !ok {
		return errors.New("invalid FileCache configuration type")
	}
	f.directory = cacheLocation
	if _, err := os.Stat(f.directory); os.IsNotExist(err) {
		err := os.Mkdir(f.directory, 0744)
		if err != nil {
			return err
		}
	}
	return nil
}

// InCache checks the directory for a file matching the provided key.
func (f *FileCache) InCache(key string) bool {
	if _, err := os.Stat(filepath.Join(f.directory, key)); !os.IsNotExist(err) {
		return true
	}
	return false
}

// LoadFromCache loads the contents of the file matching the key from the cache
// directory.
func (f *FileCache) LoadFromCache(key string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(f.directory, key))
}

// WriteToCache writes the contents to a file named with the provided key in the
// cache directory.
func (f *FileCache) WriteToCache(key string, contents []byte) error {
	return ioutil.WriteFile(filepath.Join(f.directory, key), contents, 0744)
}
