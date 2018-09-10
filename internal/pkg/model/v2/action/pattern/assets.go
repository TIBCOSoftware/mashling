// Code generated by go-bindata.
// sources:
// DefaultHttpPattern.json
// DefaultChannelPattern.json
// DO NOT EDIT!

package pattern

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

var _defaulthttppatternJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x55\xcb\x6e\xdb\x3a\x10\xdd\xfb\x2b\xe6\x12\x46\xb2\x09\x04\xdf\xb6\x2b\x03\x5e\x34\x69\x81\x22\x45\x81\x2e\xd2\x66\x51\x74\x41\x4b\xe3\x98\xb1\x4c\x12\xc3\x51\x52\xd7\xd1\xbf\x17\xa4\x24\xeb\x61\xd9\x55\xdd\x34\x3b\x51\x9c\xc7\x39\x67\x1e\xdc\x8e\x00\x44\xa2\x9c\x95\x1c\x2f\xc5\x14\xfc\x19\x40\x68\xb9\x46\x31\x05\xf1\x0e\x17\x32\x4b\xf9\x03\xb3\xfd\x2c\x99\x91\xb4\xb8\x28\x2c\xc8\x64\x8c\x4e\x4c\xe1\x5b\x38\x43\xe9\x19\xee\x1c\xa3\x6d\x5e\xb5\xaf\x83\x89\x5a\xf8\xf0\xb1\xd1\x8b\x28\x73\x78\x7d\x7b\x03\xb3\x19\x30\x65\x58\xc6\xaf\x63\x21\x3d\xa8\x38\x80\xb9\xbe\xbd\xf9\x2a\x53\x95\x48\x36\xd4\x35\x53\xda\x66\xbc\xc3\xdf\xb8\x60\xb3\x42\xed\xbd\xc7\x5b\x2b\x37\xa9\x91\x49\xb4\x44\x99\x20\x45\x6f\x33\x5e\x1a\x52\x3f\x25\x2b\xa3\xf3\x4e\x40\x00\xb1\xc2\x4d\xe1\x17\x50\xde\x3f\xf2\x47\xdc\xe4\xa2\x65\x95\x37\x4e\xf9\xc5\x50\xb6\x57\x8a\xe2\x4c\xf1\x25\xa1\x5c\x21\x0d\x20\xde\x76\x10\x27\x25\x2d\x25\x5e\xc8\xd4\x21\x3c\x3d\x41\x53\xcc\x88\xd0\x59\xa3\x1d\x46\x0f\xfe\xd7\x00\x44\xbe\x21\x2e\x65\xbc\x42\x9d\x0c\xae\x44\x46\x69\x43\xcf\x79\xe1\xfd\x85\xd2\x1e\xe9\x35\xf2\x7b\x22\x43\x62\x1a\x80\xfc\x33\xcd\xe1\xec\x0c\x1a\x54\x6a\x1d\x2a\x00\xf0\xdf\x0c\xce\xcf\x87\x97\x66\xa8\x16\xc6\x22\x85\xbe\x2b\x30\x66\x9a\xdb\x75\x7d\x69\x9a\xb3\x17\xa0\x49\xe8\x90\x8f\x90\xdc\x7d\x7f\xaf\xc3\x8b\x0a\xea\xa9\xeb\xc4\x53\xff\x4d\xab\x87\x91\xe8\x52\xc2\x46\xfb\x75\xae\x4c\xc6\x07\xe8\xc6\x26\xf1\x6a\xbd\x99\xfc\xbf\xd7\xd1\x89\x64\xd9\xe3\xd2\xc8\x24\xc6\xdb\x5e\xa0\x9d\x95\xd3\x16\xed\xf9\xfb\xa4\x7d\x13\x31\x29\x6b\xf1\xe0\x4a\xf8\x3b\x99\x5e\x9f\x22\x53\x5c\x00\x84\x79\x89\xbd\x44\xf8\x3c\x2a\x55\x59\x42\x4f\xfc\x29\xa1\x57\x93\xc9\x21\x42\x62\xbc\xed\x1d\xc0\xb9\x49\x8e\xbe\x29\xf5\x4c\x8c\x9a\x7f\xfc\x29\xd0\xa8\x66\xb4\x1e\x8f\x0a\xdc\xee\x09\xef\x7d\x35\x45\x82\x2e\x26\x65\xab\xd9\x2c\x2d\x10\x9c\x59\x23\x84\x37\xd3\xd5\xc6\xbc\xb1\x21\xd4\xfd\x63\x39\xbf\xa5\x84\x7b\xb9\x0e\xac\x89\x6e\xb6\xab\x4e\x09\xab\x45\xb3\x97\x30\x3e\x10\xcf\x21\xb3\xd2\x77\xae\x55\x0c\xb1\x2e\xca\x20\x64\xa3\x4b\x7d\xc7\x33\xfe\xf0\x65\x13\x77\xf5\xfa\xc9\x8f\xd2\xe8\x7b\xdd\xba\x1c\x3e\xc9\x15\x82\xd4\xb0\x64\xb6\x10\xcb\x34\x05\x36\xb0\x31\x19\xc1\xbc\xeb\x59\xb1\xf1\xa6\xa5\x7e\x23\x5f\xc3\x7c\xf4\x2b\x00\x00\xff\xff\x20\xea\x80\xa2\x7e\x09\x00\x00")

func defaulthttppatternJsonBytes() ([]byte, error) {
	return bindataRead(
		_defaulthttppatternJson,
		"DefaultHttpPattern.json",
	)
}

func defaulthttppatternJson() (*asset, error) {
	bytes, err := defaulthttppatternJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "DefaultHttpPattern.json", size: 2430, mode: os.FileMode(420), modTime: time.Unix(1536238166, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _defaultchannelpatternJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x56\x4d\x6f\x1a\x31\x10\xbd\xf3\x2b\x26\x56\x94\x5c\xd2\x25\xfd\x38\x21\x71\x48\xe8\xa5\xa9\xaa\x56\x4a\xd4\x1c\xaa\x1e\x06\xef\x00\x0e\x8b\xbd\xb2\xc7\xa4\x34\xd9\xff\x5e\xd9\xec\xb2\xec\x02\x09\x49\x5b\x6e\x78\x67\xfc\x66\xde\x1b\x3f\x9b\x87\x0e\x80\x48\x95\xcb\x91\xe5\x44\xf4\x20\xac\x01\x84\xc6\x19\x89\x1e\x88\x8f\x34\x42\x9f\xf1\x60\x82\x5a\x53\xf6\x0d\x99\xc9\x6a\x71\xb6\x4c\xb2\xc6\x33\x39\xd1\x83\x1f\x71\x0d\xe5\xe6\x18\x73\x4c\xf9\x7a\xa8\x19\x8e\x29\x6a\x14\x2a\x48\xa3\x47\x89\x77\x74\x75\x7b\x03\xfd\x3e\xb0\xf5\x54\xe2\xd7\x58\x64\xe7\x4a\xc6\x7e\xae\x6e\x6f\xbe\x63\xa6\x52\x64\x63\xdb\x69\x4a\xe7\x9e\x57\x14\xd6\x02\x6c\xa6\xa4\xc3\xee\xe3\x87\x1c\x17\x99\xc1\x34\x99\x10\xa6\x64\x93\x0b\xcf\x13\x63\xd5\x6f\x64\x65\x74\xd1\x02\x04\x10\x53\x5a\x2c\xf7\xc5\x2e\xef\xee\xf9\x33\x2d\x0a\xd1\xc8\x2a\xd6\x56\xc5\xd9\xbe\x6c\x07\xca\x4a\xaf\xf8\xd2\x12\x4e\xc9\xee\x41\xbc\xb9\x41\xbc\xaa\x68\x29\xf1\x08\x33\x47\xf0\xf8\x08\xeb\x62\x26\x96\x5c\x6e\xb4\xa3\x64\x1e\x3e\xed\xd1\xd1\x85\x0c\xa2\x5d\xa2\x9c\x92\x4e\xf7\x9e\x45\x0c\xb8\x2d\x11\x00\x21\x97\xa7\x2c\x80\x33\x39\xde\x18\x07\x80\x98\x63\xe6\x69\x95\xd0\x8a\x17\xff\x6d\x30\x70\x72\x02\x0d\xbe\xb5\x5c\x64\xad\xb1\x70\xd4\x87\xd3\xd3\xfd\xc7\xb7\xaf\x5a\x26\x27\x1b\xcf\xe6\xb2\x45\xaf\xb9\x39\xfb\x03\xb3\xec\x1f\x80\xa5\x25\x47\xfc\x04\xc7\xd5\xef\x9f\x35\xbc\xa8\xfa\x7c\xed\x8d\x13\x98\x3f\xe3\x86\xe8\x9a\x36\xa5\x28\x8b\xe8\x45\x8c\x56\xc8\x78\xde\x41\x57\x9a\x34\xa8\xf5\xe1\xfc\xed\xc6\x7d\x93\x22\xe3\x76\x6f\x54\x95\xc4\xf1\xc3\xd6\x46\x8b\xc3\x9a\xa1\x19\x49\xd8\xaa\x3c\xa7\x9d\xb7\xc6\xdf\xc9\xf4\xfe\x35\x32\xc9\x65\x83\x30\x2c\x7b\x2f\x3b\xfc\x37\x2a\x55\x55\xe2\x99\x78\x29\xa1\x77\xe7\xe7\x2f\x20\x54\xcd\x37\x70\xba\xf6\x52\x92\x73\x47\x2f\x20\x51\xbb\xa5\xb3\xfe\x25\xac\x22\xc1\xca\xbd\xb5\x71\xaa\x1e\x56\x7f\x01\xb6\x3e\xb9\x22\x25\x27\xad\xca\x2b\xd7\x96\x19\x04\xce\xcc\x08\xe2\x83\xeb\xea\x64\x5e\xe4\x11\xea\xee\xbe\x74\x76\x29\xee\x46\xad\x1d\x17\x48\xbb\xda\xa0\x35\xdc\xea\x0a\xda\x28\x28\x77\xe0\x39\x62\x56\x7a\xdc\x7c\x87\xc4\x6c\x39\x20\x81\x6b\xe7\x37\x78\x81\xe9\x57\x18\xa8\x18\xd7\x17\x53\xf1\x24\x8d\xed\x4f\x63\x9b\xc5\x17\x9c\x12\xa0\x06\x8c\xd9\x20\x31\xcb\x80\x0d\x2c\x8c\xb7\x30\x6c\xef\xad\x18\x8d\x32\x33\x36\x01\x7f\xae\x78\xf1\x1c\x21\x4b\xd1\xce\x63\xc5\x13\x3f\x4c\xa4\x99\x75\x6f\x3e\x5d\x0e\xbe\x5e\x9b\x11\xdf\xa3\xa5\x6e\x04\x7b\x13\x18\x5a\x35\xec\x62\x89\xda\xad\xde\xe1\x26\xd7\x4e\x38\x36\x45\xe7\x4f\x00\x00\x00\xff\xff\xc3\xe9\x1a\x88\x31\x0a\x00\x00")

func defaultchannelpatternJsonBytes() ([]byte, error) {
	return bindataRead(
		_defaultchannelpatternJson,
		"DefaultChannelPattern.json",
	)
}

func defaultchannelpatternJson() (*asset, error) {
	bytes, err := defaultchannelpatternJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "DefaultChannelPattern.json", size: 2609, mode: os.FileMode(420), modTime: time.Unix(1536238166, 0)}
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
	"DefaultHttpPattern.json":    defaulthttppatternJson,
	"DefaultChannelPattern.json": defaultchannelpatternJson,
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
	"DefaultChannelPattern.json": &bintree{defaultchannelpatternJson, map[string]*bintree{}},
	"DefaultHttpPattern.json":    &bintree{defaulthttppatternJson, map[string]*bintree{}},
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