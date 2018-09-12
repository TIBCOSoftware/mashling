// Code generated by go-bindata.
// sources:
// internal/app/assets/banner.txt
// DO NOT EDIT!

package assets

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

var _bannerTxt = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x95\xbd\x8e\xe4\x30\x08\xc7\xeb\xcd\x53\xd0\x31\x2b\x1f\xd0\xf3\x2c\x23\xd1\x5c\xef\xc2\xad\x1f\xfe\xf4\xf7\x47\x26\x9b\x71\x32\xb7\x24\x8a\x12\xec\x1f\x60\x0c\x0e\xed\xc2\xcc\xac\x6e\xc6\xdb\x4b\x65\x5d\xa4\xab\x14\xef\xee\xb8\xcd\x36\x22\x55\x15\x1f\x02\x25\xa6\xf9\x40\xa0\x70\x17\x31\xdd\x88\x45\xa4\x4f\x34\x53\x18\x32\xdf\x48\xfa\x34\x8c\x08\x0f\x65\x37\x37\x26\x51\x57\xd2\x47\x61\xde\xc4\x4d\x94\x99\xe9\x57\xa0\x9a\x63\xb1\xea\xc6\x47\xb5\xc9\xd1\xb8\x8f\xf0\x4e\x3e\x1d\x6b\x20\x62\x39\x04\x4b\xfc\x03\xc5\xa0\x2d\xe0\x81\xc2\xb1\xf3\x2b\x16\x79\xf3\x60\xd7\x28\xb1\x39\xb7\x2d\x63\x96\x33\xba\xb0\x86\xad\xd4\x6d\x8f\xcb\x05\x89\x17\x11\xf3\xb7\xf8\xd8\x91\x49\x55\x9d\x09\x95\xb9\xff\x67\x83\x3f\xd2\x86\xda\x49\x2f\xe9\x95\x62\x0b\x90\x52\x29\x3b\xa9\x9e\x52\x4e\x0b\x71\xe5\x15\xcb\x80\xc7\x12\x73\x32\x5b\xa1\xf0\x2e\x4b\x9a\x34\x97\xec\x88\x3d\xe7\x24\x9a\xfd\xc4\xbb\x58\xb2\xce\x2f\x71\xd2\x52\x92\x70\xca\x19\x73\x45\xce\x21\x6b\x32\x49\x77\x3c\x79\x29\xa9\xe1\xcb\x35\x63\x43\x46\xf6\xf4\xc2\x00\x59\xc9\xf9\xca\xc0\x31\xfd\x57\x11\x10\xe7\xf2\xd9\x80\xdd\x44\x80\x2c\x5c\x1b\xb0\x69\xc0\x2f\x79\xb2\x52\x96\x3c\x4e\x84\x69\xe5\x06\xe7\x95\x7f\x7b\x15\x83\xdd\xac\x7e\xec\x81\x5c\x56\x8e\x5d\xd5\xce\xee\x3b\x39\x3a\xfb\xf7\x28\xb6\x2e\x69\x6b\x15\x39\x44\xdb\x5a\xe5\x3a\xd9\xa3\xea\x66\x9f\x32\xab\xb8\x8d\x22\x45\x87\xde\x81\xd8\xea\x56\x55\xaa\xad\xa1\x75\x96\x97\xab\xe2\xf0\xb9\x63\xd1\xa6\x49\xfa\x7f\xa1\xb1\x36\x52\x8e\xde\x55\xbb\x67\xbd\x94\x7c\x3c\xbd\xfc\xef\x58\xa9\xb0\xde\xaf\xb5\x15\x57\x3a\x1c\x2c\x9c\xd4\x7a\x65\x58\xb2\x0f\x21\x53\x1e\x31\x1f\x44\xe6\x0f\x49\x94\x91\x88\x7b\x03\xd6\x9c\x63\xc9\xfd\x68\x45\xce\x54\x01\x8b\xa8\x2b\x7e\x56\x2a\x3c\xe5\xde\x16\x1b\xfc\x9d\xc5\x17\xba\x2e\x1f\x83\x9b\x7b\xd1\x2c\xc8\xca\xf8\xae\x96\x7e\x7c\xf3\xb6\x7d\x7d\x7d\x2d\x4d\xc5\x78\xc6\x59\x7f\x45\x54\xaa\xed\xf9\x88\xef\x1f\xfa\x47\x7c\x37\x24\x28\x70\xc3\x1e\x5e\xf0\x5a\xa9\x46\x00\xec\xdf\x73\xa4\x39\x8f\x68\x9e\x2a\x71\x10\x53\xd0\x93\x8c\x82\xf1\x88\xa6\x7b\xd6\xe6\x91\xf7\x91\xf6\x6d\x98\xd9\xb9\xd7\xf5\x88\x4a\xcf\xc0\xc4\x7a\xba\x1e\x2d\x84\x3a\xa2\xae\x0d\x9c\xaa\xa8\xcf\x88\x3f\x51\x23\xc2\xfa\x77\xbf\xe6\x08\x90\xf6\x16\x76\x99\xc4\xf7\xac\x86\x75\x37\xff\x35\xbb\xf9\xde\xfe\x05\x00\x00\xff\xff\x37\x0c\xc8\x3c\x83\x09\x00\x00")

func bannerTxtBytes() ([]byte, error) {
	return bindataRead(
		_bannerTxt,
		"banner.txt",
	)
}

func bannerTxt() (*asset, error) {
	bytes, err := bannerTxtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "banner.txt", size: 2435, mode: os.FileMode(420), modTime: time.Unix(1528122438, 0)}
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
	"banner.txt": bannerTxt,
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
	"banner.txt": &bintree{bannerTxt, map[string]*bintree{}},
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
