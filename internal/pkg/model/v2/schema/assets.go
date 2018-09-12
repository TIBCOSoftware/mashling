// Code generated by go-bindata.
// sources:
// schema.json
// DO NOT EDIT!

package schema

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _schemaJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x59\x4d\x4f\xdc\x3c\x10\xbe\xf3\x2b\xa2\xf0\x9e\x5e\x01\xdb\x4a\x3d\xed\xb9\x12\xed\xa9\x15\xdb\x5b\x85\x90\x49\x66\xb3\x46\x89\x1d\xc6\x13\xaa\x55\xc5\x7f\xaf\x36\x8e\xa3\x5d\x62\xe7\x3b\x90\x95\xc2\x81\x83\x67\x3c\xf6\x8e\x9f\x67\xbe\xf2\xf7\xc2\xf3\x3c\xcf\xff\x4f\x05\x3b\x48\x98\xbf\xf6\xfc\x1d\x51\xba\x5e\xad\x9e\x94\x14\xd7\x7a\xf5\x46\x62\xb4\x0a\x91\x6d\xe9\xfa\xd3\x97\x95\x5e\xbb\xf4\xaf\x8a\x9d\x08\xdb\xc3\xb6\xcb\x55\x08\x5b\x2e\x38\x71\x29\xd4\x6a\xa3\xcd\x15\x3a\x47\x12\x7f\xed\xe9\x23\x73\xc1\x57\xae\x52\x46\xc1\xee\x64\x35\x97\x20\x3c\x67\x1c\x21\xf4\xd7\xde\xef\x13\x49\x2e\x15\x2c\x01\xff\x64\xf9\xfe\xea\x74\x7f\x8a\x32\x05\x24\x0e\xaa\x62\x3b\x97\x07\x52\x6c\x79\x94\x21\x3b\xdc\xca\xaa\xa2\xcd\x30\x22\x40\xf1\xb3\xde\x5a\xa9\x7e\xf3\x7f\xad\x3c\xd7\x61\x61\x98\xbb\x82\xc5\x27\x56\x09\x33\xb8\xaa\xdf\x49\xfb\x14\xac\xfe\xa8\x9e\x81\xc8\xf6\x7e\xbd\xb9\x5c\xf1\x51\xca\x18\x98\x68\xa3\xca\x05\x41\x04\xd8\x46\x55\x64\xc9\x63\x5b\xcd\x38\x6e\xa3\x27\x1f\x9f\x20\xa0\x36\x9a\x8a\x90\x8b\xc8\xaf\x55\xbc\x77\x4a\x5f\xad\x92\x57\xfb\xb9\xe6\x41\xcc\xed\x2a\x4a\x96\x7d\x1a\xba\x4e\xbc\x19\x8b\xae\x5f\x61\xb3\x58\x60\x74\x5c\xa3\x28\x33\xaa\x41\xbb\xcf\x09\x92\x06\x32\xf4\x8e\x29\x76\x6b\xd6\x38\x73\x77\xb8\xa6\xfd\xb1\x5d\x8f\x96\x09\xfe\x9c\xc1\xf7\xe2\xfe\x6e\xda\x95\x6e\xd3\x64\xaa\x7a\xed\xa2\xe6\x34\x17\xcb\xb7\x2c\x56\x6f\xce\x73\xa2\xe8\xc8\xa4\x7f\xcb\x08\xfe\xb0\x7d\xaf\x38\x69\x79\xdd\x17\x40\x75\x88\x7b\x83\x42\x68\x08\x2a\x40\x9e\xd6\x07\xd0\x3e\xd8\x0b\x8b\xac\x70\x06\xf8\x2b\x13\x58\x27\x08\x26\x5c\x18\xfc\x7d\x9e\x1a\xa5\x0e\xff\xc6\x6c\xff\xc0\x13\x16\x8d\x1c\x8b\x8c\xe9\xf1\xa3\xdc\xf8\x16\x15\xe0\x0b\x0f\xce\x00\x64\x1b\x7d\xd1\x0f\x0a\x73\x16\xcf\x11\xf2\x28\x02\x9c\xbf\xe7\x7e\xe9\x8b\xce\xc7\x73\x26\xf0\xf6\x07\xf2\x7b\x66\x9d\x6f\x4c\x84\x31\x60\xd7\xac\x13\x5a\x83\x62\xe7\xf4\xe2\xea\x0d\xda\xfb\xca\x4a\x7a\x22\x2e\xa2\x1a\xe8\x2e\x25\x7f\xa9\xba\x94\xfc\x95\xb3\x1a\x4b\xfe\xf7\xe4\xe7\x8f\x8c\xd2\x8c\x3a\xd3\x93\x11\x1b\xda\x3c\x87\x2d\x72\xb1\x41\x4f\xbb\xc2\xe1\x70\x29\xa7\xc9\xce\x04\x6a\x26\x4e\x23\x61\x5a\x10\xa5\x05\x41\x9a\x89\xd1\x44\x88\x66\x22\xd4\x12\xa0\x0a\xfc\x77\x85\xe8\x1d\xa8\x54\x0a\x55\x85\x4b\x03\x48\x01\x51\xe2\x30\x94\x6a\x13\x8d\x30\x35\xcf\xdc\x0a\xa6\x7c\x3b\x6e\x3a\x92\x76\x06\x97\xf2\xd1\x8a\x24\x47\x81\x54\x44\x90\x0f\xc6\x48\xde\xc4\x77\x04\x88\x22\x48\xd5\x30\x80\x30\xb5\x17\xc1\xcc\x01\x82\x05\x7f\xe6\x5f\x6b\x97\x4c\xef\x95\x4f\x3b\xd4\xd0\xfa\xe5\xe7\xee\x8e\x0d\x41\x3a\xfa\x54\x60\x4e\x23\xa9\x8d\xf1\x61\x27\xda\x26\x4c\xed\x62\x2e\xa2\x07\x75\xfc\x85\xe0\x44\x25\x2a\x66\x5d\x83\xb8\x1d\x39\x06\x66\xa5\xc2\xd4\x91\xf5\xd6\xf6\x2b\x3c\x07\xa4\xdf\x7a\xe5\x3c\xda\x43\x33\x97\x18\x69\x28\x99\x1f\x39\xac\x65\x9c\x6a\x22\x39\xc5\xf0\x69\xe9\x43\x97\x3e\xf4\xf8\x6f\xfa\x4f\x4f\xc5\x9e\xb3\x08\x2d\x87\xec\xd9\xb5\x22\xb4\x8d\x49\xbb\x86\x90\xb1\xeb\x37\x2e\x6a\xeb\xfb\x85\xe4\xa5\xea\x42\xf2\xca\x59\x7d\x48\xae\x1c\x49\xb9\x62\x74\x16\x3c\x37\x03\xfa\x31\x4b\x08\xcb\xfa\x4e\x0f\xb2\x07\xf6\x8b\x93\x95\x17\xe5\xf5\xe6\xde\xd4\x98\x0f\x02\x63\xf7\x35\x53\x7e\x72\x59\x6a\xb7\x25\xac\x9f\x7f\x58\x9f\x7b\xed\x76\xa1\xff\xbf\xfe\x0b\x00\x00\xff\xff\xe0\x18\x76\x84\x50\x28\x00\x00")

func schemaJsonBytes() ([]byte, error) {
	return bindataRead(
		_schemaJson,
		"schema.json",
	)
}

func schemaJson() (*asset, error) {
	bytes, err := schemaJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema.json", size: 10320, mode: os.FileMode(420), modTime: time.Unix(1536788500, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"schema.json": schemaJson,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"schema.json": &bintree{schemaJson, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
