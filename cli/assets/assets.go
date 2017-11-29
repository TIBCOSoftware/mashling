// Code generated by go-bindata.
// sources:
// assets/banner.txt
// assets/default_manifest
// schema/mashling_schema-0.2.json
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

var _assetsBannerTxt = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x55\xcf\xee\xa4\x30\x08\x3e\x8f\x89\xef\xc0\x8d\xf9\xa5\x0b\xdc\x79\x96\x49\xb8\xec\xbd\x87\x5e\xfb\xf0\x1b\x68\xeb\x38\x4e\xd5\xdd\x45\x63\x94\xf2\x7d\x50\xfe\x54\xd8\x04\x11\x91\x55\x04\xd7\xe5\xad\x93\x26\xd4\x75\xec\x1f\xaa\x7e\x8b\xb8\x8e\x99\x49\xbb\xb8\x36\x0c\xb5\xa3\x5c\xa3\x4a\x24\xbc\x2e\x80\x44\xd4\x4c\x45\xd8\xb9\x44\xd7\x05\xa8\x19\xfa\x12\xe1\xd0\x36\xca\x6e\x06\x4d\x0b\xb7\x82\xb8\x2e\xa4\x42\x8c\x88\xf0\x4f\x48\x16\x8d\x4d\xb3\x0a\xee\xf5\x42\x7b\x7a\xed\x11\x1e\xbc\x6a\xec\x03\x00\x69\x17\x2f\xe0\x07\xd6\x17\x65\x82\x1e\x58\x77\xad\xf8\x0e\x87\xbe\x7c\xc8\x05\x16\x50\x14\xa3\x7c\x88\x74\xc4\x4e\xe8\xbc\xaa\xbc\x15\x19\x49\xc9\xf3\x4f\x44\xa2\x5f\x21\xa2\x7a\x3a\x99\x79\x64\x95\xb6\x5e\x38\x52\x7e\xe4\xce\x3b\x29\xbd\xa5\x75\x8d\xcc\x90\x90\x4a\xd9\xa0\xac\x29\xe5\x34\x11\x65\x9c\x82\xd1\xd1\x7d\x9b\x39\x89\xcc\xb0\xee\x9f\xe6\x70\xe0\x5c\xb2\x7a\xf8\x39\x27\xe2\xac\x07\x02\x25\x49\xd2\x08\xe6\x78\xe0\x52\x12\x61\xca\xd9\x8d\x89\x8e\x51\x73\x12\x4a\x97\x04\xa0\xa5\xa4\xc0\x4f\xf7\xed\x65\xe9\x29\xe4\x33\x06\x90\x92\xf3\x19\xc3\xbe\x08\xa7\x31\x00\xe6\x72\xcf\x20\x57\x31\x78\x26\xce\x19\x64\x30\xe8\x39\x01\x48\x29\x53\x02\x3f\x22\x06\xcd\x15\x1e\x67\x11\xc8\xbb\x29\xe4\x2a\x03\xbd\x12\x74\xda\x42\x72\xda\x44\x9b\xf7\xa4\x3e\xe9\xff\x81\xf5\x0a\x26\x8e\xb9\xa1\x5d\xc0\x31\x37\x17\x29\xef\xfd\x37\xc6\x16\x91\x49\xa5\xf7\xab\x0f\xec\x25\xd2\x4b\x1e\xfd\xc5\x1c\x03\xce\xa3\xd1\x94\xd9\x8f\xa3\x4b\xb0\x4f\x6d\xa2\xf6\xd7\x08\xb0\xf4\xc4\xfb\x28\xb3\xdc\x80\xb5\x94\xbc\x3f\xd1\xf4\x77\xdf\x2d\x21\xdf\xec\x37\xda\x2c\xed\x8e\x1a\x4c\x2c\xad\x45\x24\xc9\x5d\xd4\x90\x7b\xd8\x3b\xa1\xf1\xb7\x22\x46\x4f\xc6\x0d\x83\x84\x7b\xdf\x76\x3b\x70\x3d\x71\xcc\x8e\x26\x62\x65\xff\x91\x31\xe1\x90\x1b\x32\x14\xf7\x78\x14\x9d\xe8\x9a\xdc\x87\x37\x2a\x12\x14\x34\x63\xdf\xd4\xd4\x4e\x75\x5c\x97\x75\x79\x3c\x1e\x53\x32\xeb\x4f\x3b\xea\x4f\x21\x15\x6a\x3c\x9f\xf6\xf3\xa1\x7f\xda\x4f\xc3\x18\x98\xdf\xce\xe8\x2f\xfe\x5a\xa1\x9a\x39\xb2\x7d\x8f\x95\x70\x6f\xd6\x7c\x55\x40\x03\x04\x83\x17\x08\x18\xfa\xc3\x42\xf7\xaa\xe1\x13\xb7\x95\xf8\x16\xb7\xec\xc0\xf7\xf5\xb4\x0a\x2f\x73\xcb\x7a\xb8\x9e\x11\x44\xed\x81\xd7\x86\x1c\x3a\xab\x2f\xb3\x5f\x56\xcd\x4c\xda\x77\xbb\xc6\x8a\x63\xe2\xcd\xe4\x3c\x95\xdf\xb9\x35\xe9\x8e\xfe\xca\xbc\x36\xfa\x3f\x01\x00\x00\xff\xff\x0c\x3c\x9d\x92\xb1\x09\x00\x00")

func assetsBannerTxtBytes() ([]byte, error) {
	return bindataRead(
		_assetsBannerTxt,
		"assets/banner.txt",
	)
}

func assetsBannerTxt() (*asset, error) {
	bytes, err := assetsBannerTxtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/banner.txt", size: 2481, mode: os.FileMode(438), modTime: time.Unix(1510825690, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsDefault_manifest = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x9a\xcd\x72\xdb\x48\x92\xc7\xcf\x76\x84\xdf\xa1\x43\xe7\x96\x51\xdf\x55\xe9\xe3\xee\x65\xf7\xb2\x3b\x11\x3d\xb7\x89\x39\x64\x66\x65\x91\xe8\x26\x01\x0c\x00\xda\x52\x4f\xf4\xbb\x4f\x80\xb4\x2c\xca\x4d\xb6\x00\x8a\x61\xd9\x17\x91\x84\x98\xc9\xff\x2f\x8b\x55\xf9\x41\xfc\xfb\xdd\xdb\x37\x37\x1f\xa5\x1f\xea\xb6\xb9\xf9\xf0\x93\xfa\x79\x7a\x9d\xa5\x93\x26\x4b\xc3\xb5\x0c\x37\x1f\x7e\xfa\xc7\xbb\xb7\x6f\xde\x4c\x6f\x7c\xf3\xe6\xa6\xde\x76\x6d\x3f\x76\x38\xae\x6f\x3e\xfc\x74\xb3\xaa\xc7\xf5\x8e\xde\x73\xbb\xad\xfe\xaf\xcd\xf2\xb7\xbe\xde\x4a\xf5\xeb\xd0\x36\xfb\x37\xfc\x7c\x30\xe9\xa5\x6b\x87\x7a\x6c\xfb\xfb\xc9\x64\x3d\x8e\xdd\xf0\xa1\xaa\x66\x9a\x7e\xac\x3f\x2b\xbb\x49\xce\x29\x9b\x25\x7b\xe7\xac\x49\x0c\x40\x62\x5d\x34\xd1\xc4\x10\xa2\x8e\x20\x52\x8c\x45\xd0\x0f\xb6\xd4\x63\xc3\x7b\x95\x5b\x1c\x46\xe9\x6f\xa6\xcb\x7f\xfc\x3c\x87\xe5\x97\x75\xdb\xd5\xe5\xbe\x1a\xb0\xc7\x2d\xce\xe6\x38\x67\xf6\xc8\x60\x9d\x0a\xc9\x06\x28\x25\xd8\xa0\x49\x1b\x03\x2a\x9b\x14\xb4\x2f\xc5\x3a\x8c\xa9\x58\x12\xf0\xfe\x1a\x0c\x75\xbf\xeb\x06\x69\xaa\x4d\xbb\xea\x77\xc3\x7c\x88\x73\x76\x47\x2b\x01\xd1\x19\x94\x82\x8e\x8c\x0a\x99\x8b\x53\x2a\x82\x29\x96\xb2\xf5\xe4\x9d\x01\x48\x42\x96\xae\x40\xf1\xf7\xff\xfd\xaf\xff\xfe\xff\x5f\xda\x32\x7e\xc2\x5e\xaa\xb2\x69\x57\xed\x2d\x6f\xea\x4a\x9a\x8f\xb3\x89\xce\xf8\x38\x45\x16\x33\x19\x27\x4a\x48\x1b\x0c\x08\x36\xa5\xe2\x38\x06\x97\x21\x81\xcf\x4e\x19\xc8\xc1\x19\x5d\xce\x92\x7d\xbe\xfe\x40\xb2\xd7\xf9\x72\xda\xdd\x58\x6f\x7e\x08\xdc\xbd\xd0\x17\xf2\xb6\xcd\xd8\xd7\x54\x21\x8f\x75\xdb\x4c\x17\x3f\xbd\x0c\xfd\xe0\xef\x14\x3e\xc4\xa0\x43\x71\xd6\x38\x65\x1c\x50\xf1\x86\xb3\x21\xed\x75\x30\x99\x92\x51\x9a\x74\x16\xad\xdc\x6c\xfc\x63\xcd\xd7\x8b\xc2\xc7\x7a\xbc\x9f\x76\xe3\x0f\x15\x86\x2f\xa2\xaf\x1c\x87\x5e\xba\xcd\xfd\x8f\x17\x89\x83\xec\xab\xc7\x62\x18\x7f\xc4\x50\x0c\xe3\x95\x22\xb1\x6d\xb3\x6c\xaa\xa1\xde\x76\x1b\xf9\x61\x02\xf1\x44\xf4\x75\xe2\x30\xf6\xf5\x6a\x25\xfd\x3e\xb4\x15\xb7\xfd\xfc\x6c\xff\xda\xc1\xf8\xb3\xf2\x97\x45\x64\x73\xac\xfb\x02\xf6\xcd\x69\x6e\x16\x2b\x56\x01\x29\xe7\x43\x09\x9a\xac\x37\xc5\x05\x76\x21\x70\x32\x4e\x5b\xf2\x2a\x3a\x91\xf3\x89\xf2\x52\xaa\x2d\x0e\xeb\x4d\xdd\xac\x2a\xb9\x1b\x0f\x12\xbf\x84\x6c\xd5\xf6\xf5\x66\x83\xdb\xdd\xdd\xe7\x2b\x17\x82\x3f\x7c\xc4\x29\x6e\x1f\x0d\x26\x15\x9c\x0f\xc9\x46\x0d\x9c\x82\xf3\x05\x8a\x28\x54\x49\x0c\x23\x10\xfa\xe8\xc3\xfc\xf5\x9e\x83\xf1\xe2\x58\x6d\x6a\xaa\xb8\x6d\x72\x3d\x65\xe3\x4b\xf7\xc2\x37\x0d\xcb\x57\x8a\xaf\x12\x81\x17\x14\x8d\xdf\x9c\x7d\x61\xdd\x88\x1d\xf2\x5a\xaa\x71\xdd\xd7\x65\xdc\x3b\x98\xbe\x50\xfb\x57\xb3\x89\x9f\xf8\x38\x05\xea\x4a\x8c\x48\x49\x0c\x84\x1c\x1c\xfb\x28\x41\x50\x73\x10\x5b\x54\xd1\xde\x10\xc7\x1c\xd0\xe2\x22\xd0\x47\x9d\xf3\x69\x33\x7e\x14\x5e\xad\xab\x55\x7b\x3b\x74\xf2\xa9\x9a\xfe\xcc\xc6\xfc\xca\xf8\x14\xa8\x70\x16\x24\x0e\xde\x81\xcf\xc5\x64\x61\xf0\x39\xb2\x43\xc7\x56\x94\xd1\xa0\xac\xb2\x5e\xce\x37\xa5\x5f\x83\xee\x3f\x68\x3e\x9f\x7c\x5e\x8a\x55\x7b\xdb\xcb\x50\x6f\x6a\x69\xf8\xbe\xa2\x5e\xf0\xb7\x05\x87\xda\x49\x2f\xa7\x70\x49\x17\x49\x96\x3c\x29\x5b\x82\x71\xce\xab\x64\x2c\x45\xaf\x29\x04\x63\x3c\x94\xc2\x01\x4b\x54\xb3\x71\x1f\x94\x5e\x44\x7c\x27\x7d\x8d\x9b\xdb\xa1\xc1\xae\x9b\x5f\xd7\x3e\xef\xe0\x88\x97\xc0\x7b\x51\x9a\xc0\xba\x80\xac\x21\xb3\x01\x21\x1d\x7c\x0a\x0c\x2a\x4b\x06\x40\x48\x7c\x85\xdc\xf5\x20\xeb\x5f\x3b\xd9\xcd\x2f\xc8\xce\x58\x1d\xed\x44\xc7\x9c\x94\xe7\xa2\xad\x51\x9e\xbc\x2f\x01\x44\x3b\x62\x0a\x90\x42\xcc\x1a\x05\x4c\x81\x74\x05\x82\xa9\x04\x68\x57\x65\x3b\x56\x87\x87\xd9\x10\xe7\x0d\x8f\xa6\x3f\xa0\x90\x22\x58\x2f\x62\x92\x70\x20\x93\x82\x0d\x8e\x08\x81\x5c\xce\xc1\x69\xc5\x64\xf3\x35\xa6\x3f\xab\x29\xc5\x76\x7d\x3b\xb6\xb4\x2b\x55\xdd\x2e\xe0\x38\x32\x3c\x05\x11\x74\xc0\x64\x24\x6b\x93\x63\xca\xc6\x65\x97\x2c\xd8\x60\x25\x15\xcb\x3e\xdb\x62\x54\x30\x91\xe7\xcf\x0c\xea\xf6\x62\xae\xfd\x93\xef\x18\xed\xa0\x6f\x09\xdd\x06\x9b\xd5\xe5\x7c\x4f\xac\x4f\x11\xea\xe0\x6c\x48\x56\x34\x79\x87\x20\x29\x61\x36\x21\x43\x2a\x49\x3b\xa5\x38\x31\x64\x28\xae\xc0\xb7\x25\xe4\xdb\x95\x34\xb7\xab\xb6\xca\x32\x70\x5f\x77\x63\x3b\xff\xd8\x7f\x35\xe8\x53\xa2\x5f\x10\x87\xf1\xbe\x93\xf9\xa5\xea\xab\x50\x1f\x24\xbe\x94\xb1\x1a\xeb\xad\x0c\x23\x6e\xbb\x1f\x80\xf6\x48\xec\x62\xee\x85\x29\xfd\x8c\xd5\x51\xe1\xed\x2d\x06\xa7\x5d\x54\x0e\x02\x19\x6b\x22\x12\x23\x6b\x45\x36\x58\x08\x94\x21\x89\xf3\x7c\x1e\x71\x09\xc1\xbe\x27\xab\xb6\xbb\xbb\x05\xfa\x4f\xd8\x3c\xaa\x37\xd9\x17\x29\x2a\x50\x02\xcd\x10\x35\x69\x87\x18\x4a\xd4\x8c\x1e\x29\xb0\xb2\x68\x83\x28\xb9\x82\xfa\x5f\xb1\x9b\xe4\xfc\xb2\xdf\x95\xb3\xf5\x9f\xb1\x3a\x9a\x01\x20\x12\x28\xed\x3c\x29\x6f\x63\xf0\x36\x10\x21\x59\x6b\xac\x36\xc5\x25\x6d\xb4\xb1\xa2\x21\x5e\x83\x60\x37\x95\xaf\x03\xaf\xb7\x75\x1e\xab\x49\x6c\xdf\xee\xc6\x05\xc5\xf0\xf3\x0e\x8e\xca\x7f\x4d\x90\x4c\x22\x06\xf1\xa0\x1c\xa1\xb0\xf2\x11\xb5\x77\xac\x80\xd1\xa9\x22\xb1\xa0\x9c\xaf\x87\xe7\x73\x6d\xea\xd5\x7a\x1c\x46\xe9\x1e\x9f\xdd\x8e\x3d\xb2\xf4\xb7\xab\xf9\x79\x6e\xa6\x97\x47\x42\x45\x5a\xa9\x94\x75\x02\xa5\xac\xb3\x51\x38\x60\xd0\xae\x24\xe4\xcc\x86\x94\x4f\x84\xa2\xe2\x35\xa6\x37\x6d\x27\xcd\xa4\xa5\x6e\x56\x5f\xa6\x73\xab\xf6\xb6\xa5\x41\xfa\x8f\x0b\x16\x70\xb6\x9f\x47\x4a\xf4\xa6\x18\xeb\x8c\x73\x60\x5c\xc8\x9e\x98\x4d\xb4\x12\x3c\x53\xe6\x82\xbe\xc4\x1c\x38\xd8\x6b\xd4\xf9\x47\xea\x2a\xc2\xa1\xe6\xe5\x8b\x38\xc7\xc5\xd1\xf1\x4e\xd6\x60\x31\x2a\x6a\x0d\xa8\x1d\x69\x32\x56\x67\xe7\x75\x2e\x56\xb2\x72\x18\x8d\x14\xa1\x2b\xaf\xe0\x93\x55\xb8\x90\xed\x9c\x8b\x23\x36\x1b\x34\x01\xe7\xa0\x48\x22\xb0\x63\x8b\xb1\x20\x24\xa7\xc9\xb2\xb6\x46\x9c\x52\x21\xe0\x35\xce\x95\x49\xcc\xef\x75\xf7\x5b\xdd\x54\x87\x87\xdb\xe9\x3b\xf5\xa8\x70\x11\xe1\x3c\x47\x47\x5d\x9c\x17\x50\xa4\x54\xd4\x0a\x1d\x5b\x87\x1a\x63\x4e\x18\x53\x06\x55\x80\x94\x56\xa4\x28\xbb\x7c\x05\xce\xae\x96\xbe\x17\xae\x36\xbf\xbb\xd9\x48\x27\x6d\x8e\xce\x90\xc4\x26\x82\x85\x5c\x34\x65\xf0\x92\x92\x16\xc3\xc6\x86\x88\xd1\x01\x04\x87\x59\x17\x96\x6b\xfc\xe6\xfd\xa0\xe4\xee\xee\x7f\x70\x58\x7f\x7e\xb0\x66\x31\xc9\xc1\xf0\xe4\x51\xa1\x94\x0a\xa4\x2d\x47\x63\x4a\x2c\xda\xd8\x90\x58\x29\xb4\xd9\xb2\x41\x49\x08\x00\xa8\x38\xcc\xae\x96\xbe\x68\x5c\x40\xf9\xdb\xaa\x92\xbe\x5f\xf2\x63\xc5\x29\x93\x47\xa8\xa2\x3d\x43\x54\x59\x3c\xc5\x50\x90\x15\x79\x40\x22\x6b\x72\x30\xac\x23\x73\x24\xc9\x26\x5c\xa3\xbb\xee\xb9\x6f\x3f\x6d\xe4\x7e\x3a\x92\xb7\x32\xf6\x35\xcf\x87\xf8\x4b\xdb\xa3\x53\xa1\x58\x55\x04\x14\x38\xf4\xda\xb2\x38\x8e\x4a\x11\xa0\x77\xce\x27\x22\x52\x0c\x01\x20\x5c\x7e\x9a\xef\xcb\xcd\xf7\x6d\xbf\xaa\xee\x2a\xee\xef\xbb\xb1\xad\x86\x61\x5d\x8d\xd2\x6f\xeb\x06\x9f\x9b\x1c\xb7\xef\x57\x6d\xbb\xda\xc8\xd0\xee\x7a\x96\x3d\xd8\xc1\xc9\xc9\xbe\x1a\x0d\xd8\x62\xb2\x23\xed\x28\x49\xc8\xb6\x28\x6f\x41\x6c\xb2\x25\x64\x95\xad\x2a\x2c\xb6\x9c\xdf\xf9\x7f\x1a\x30\x1e\x0b\x5d\x02\xda\xc8\x58\x4d\xe9\x54\xee\x9e\x2b\x0c\x4f\x00\x36\x72\xb2\x2e\x44\x6b\xa3\x02\x4d\xca\x1b\x8f\x25\xf8\x2c\xe0\x72\x31\x42\x51\x20\x51\x06\x08\x26\x33\x89\x99\x4d\xf7\x20\x70\x29\xd8\xa4\xfd\xb9\x53\xe2\x15\xb1\x0e\xf2\x96\x42\xd5\xb9\x79\xee\x06\xac\x57\x64\xda\xab\x5b\x8c\xd4\x8c\xd2\x37\xb8\x39\x34\x97\xd2\xd7\xcf\x36\xfe\xaf\x49\x78\x42\xec\x52\xe0\x8d\xdc\xed\x57\x7f\x23\xcf\xb5\x93\xaf\x08\x7a\x2c\x72\x29\xe0\xbe\x76\xfd\x7e\xd1\x0e\xf2\x96\x40\x0d\xf7\x43\xb5\x6b\xea\x0b\x96\x6b\xb8\x3f\x9d\xce\xc4\x18\x00\xb6\x91\x19\x34\x7a\x05\x45\xd3\x54\x76\x40\x8a\x26\x9b\x48\xa2\x58\x20\xe1\xfc\x8a\x63\xaf\x6e\x29\xd2\xa7\xba\xc9\xed\xa7\x0b\x76\xdb\xb7\xa2\x7a\x10\xb8\x04\x6c\xca\x15\xd5\x20\xbc\xeb\xa5\xa2\x3a\xd7\xfd\xee\xd9\x9b\x61\x4e\x10\x3e\x4d\x89\x47\xf7\x7a\xa6\x12\x7c\x28\x58\x6c\xb1\xb1\x04\x50\xb9\x68\xb4\xc6\x6b\x4f\x2e\x82\xd3\x5e\x34\x86\x18\xe6\x0f\xd6\xbe\x56\xba\x18\x75\xec\xb1\x19\x4a\xdb\x6f\xbf\x63\xc8\x47\x8d\x8b\xf1\x76\x4d\xcd\x6d\x3e\x04\xe8\x3b\x26\x7c\x22\xf3\x62\xc8\xe6\xfb\x5e\xc6\x27\x32\x9f\x85\x9c\x54\xbe\x3f\x62\x5d\x49\xb3\x1f\x59\x57\x87\x7f\x61\x57\x0f\x15\x76\x75\x85\x4d\xd3\x8e\xb8\xec\xde\x94\x83\x8b\xa9\x67\x78\xf0\x7a\xf2\x3c\xd2\x1c\x0b\x88\x77\x31\x63\xc8\x94\x62\x30\x41\xb1\x38\x10\x8c\x7e\x3a\x96\x3c\x28\xc7\x70\xbe\x49\xfd\x3a\x00\x7f\xa1\xfc\x2a\xe1\xe8\x3b\xae\x86\x11\xc7\x05\xf7\xa7\xbf\x7e\x24\x8e\x44\x2f\x0f\x42\xdf\xf1\x7c\xd4\xe9\x93\xa6\x3f\x67\x26\x47\x59\x8c\x51\x80\x10\x52\x76\x49\x12\x44\xe5\x5d\xe0\x84\x51\x2b\x9d\xa0\x44\x17\xb4\x8d\xaa\x5c\x3c\x39\x3a\x6c\xb9\x55\x8f\xdd\xfa\x90\x05\x1f\x5f\x57\xd8\x75\xf9\xc9\x74\xe1\x19\x96\xbf\xb4\x3d\x6a\xe2\x01\x9d\xb5\x81\x95\xf1\xce\x08\x51\x41\xa5\xd9\x46\x70\x40\x39\x6a\x6f\x0a\xa9\x12\x3c\x3c\xd3\xf6\xbe\x7b\xfb\xe6\x9f\xef\xde\xfe\xf1\xee\xed\x7f\x02\x00\x00\xff\xff\xb3\x15\xed\xbd\xde\x32\x00\x00")

func assetsDefault_manifestBytes() ([]byte, error) {
	return bindataRead(
		_assetsDefault_manifest,
		"assets/default_manifest",
	)
}

func assetsDefault_manifest() (*asset, error) {
	bytes, err := assetsDefault_manifestBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/default_manifest", size: 13022, mode: os.FileMode(438), modTime: time.Unix(1510825690, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _schemaMashling_schema02Json = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x58\xbd\x6e\xdb\x30\x10\xde\x03\xe4\x1d\x08\x25\x63\x12\x76\xe8\xe4\xb1\x9d\x3a\xa5\x40\xbb\x15\x41\x70\x96\x4e\x12\x53\x89\x54\x49\x3a\x86\x11\xf8\xdd\x0b\xfd\x5a\x54\x49\x4a\xae\x25\xc3\x01\xe2\xc9\xb8\xff\x3b\x7d\x77\x47\xf2\xed\xfa\x8a\x10\x42\x82\x5b\x15\xa6\x98\x43\xb0\x22\x41\xaa\x75\xb1\xa2\xf4\x45\x09\x7e\x5f\x53\x1f\x84\x4c\x68\x24\x21\xd6\xf7\x9f\x3e\xd3\x9a\x76\x13\xdc\x35\xaa\x7a\x57\x60\xa9\x27\xd6\x2f\x18\xea\x8e\x2c\xf1\xcf\x86\x49\x8c\x82\x15\xf9\x55\x93\x2a\x72\x0e\x2a\xcd\x18\x4f\x9e\x1b\x87\x77\x3d\x5e\x02\x1a\xb7\xb0\x0b\x6a\xd2\x53\x6b\xa9\x90\xa2\x40\xa9\x19\xaa\x60\x45\xde\x3c\xb6\x0c\xae\x11\x9b\xd2\x92\xf1\x24\x38\x70\xf7\x36\xbf\x6e\x7d\x33\xb7\x8e\x6d\xcf\xb1\x63\x73\xc8\x71\xa8\x53\x31\x5e\x51\x2a\x26\xb8\x95\x17\x0a\x1e\xb3\x64\x23\x41\x33\xc1\x95\x55\x44\x4b\x96\x24\x28\xed\x4c\x7c\x45\xae\x9f\x53\xe0\x51\xe6\x17\xc9\x18\xff\xad\x02\x93\xfd\x34\xcc\xd0\x55\x7b\x33\x47\x2b\x8b\x78\x3f\x40\xfb\xdb\xfb\xea\x33\xb7\xdd\x88\xa9\x22\x83\xdd\xf3\x22\x41\xb7\xc6\x59\x0e\xc9\x02\xd6\x51\x85\x92\x15\x7a\x89\xb2\x0c\x20\xe7\x36\xcf\x34\xe6\x1e\x7e\x25\x73\x2b\x31\x2e\x43\xb8\xa1\x11\xc6\x8c\xb3\xca\x26\x35\x5c\x58\xe2\x72\xc5\x66\x64\x05\x52\x76\xd3\x61\x34\xa9\x41\x1f\x2c\x91\x94\xe1\xe2\x8c\x49\xd5\x9d\xbb\x5c\x46\xa5\xfd\xb3\xa4\xd3\x0d\xb2\x25\x72\x69\x8c\xcf\x9e\x88\x49\x1a\xda\x09\x20\x8a\xaa\x00\x20\xfb\xde\x9f\x9d\x31\x64\x0a\x7b\xfb\xa7\xfe\xdb\x6a\x07\xbd\xc0\x07\x4b\xce\xec\x9c\xb3\xad\xa8\xca\xac\x8d\xa1\x50\x6b\xc6\x93\x0b\xdd\x1c\x8d\xd6\xfb\x99\xbd\x5d\x39\xc7\x0d\xdb\xbf\x70\x27\xe6\x00\x9e\xaf\x71\x46\x43\x26\xff\x02\xde\x42\xfa\xdf\x1e\xe8\x9f\xc1\xda\x6e\xfd\x00\x38\xb1\xd5\x67\xa8\xf5\x01\xf0\xf7\x06\x70\xf3\xb0\xb0\x08\xcc\x47\x00\x2b\x38\x3e\xc6\x76\x7d\xd7\x97\xf1\x7b\xed\x89\xc5\x28\x91\x87\xc3\x10\xba\x50\xa6\x61\xe5\xc4\x28\x0e\x3b\xf4\x88\x30\x2e\xb2\xc9\x97\xec\xc7\x02\x24\xf8\x4e\x54\x97\xd8\x8d\xb6\x3c\x0e\x98\x9b\xbf\xfa\x1d\x8e\x4e\xae\xd2\x54\xe8\x82\x06\x97\x8d\x4a\x80\x45\x5e\x76\x79\x00\x76\x60\xde\x15\xd9\x18\xb2\xcd\xd8\xbc\x22\xc4\x52\x10\xb7\xb4\xeb\xe4\x4d\x9a\x3c\xa7\xbb\xf2\xe2\x68\xd4\x55\x7d\x69\x38\xd2\x97\xc7\x20\x69\xb0\x03\x9b\x4c\xb7\xef\x68\x6a\x45\x69\xc2\x74\xba\x59\x3f\x84\x22\xa7\x3f\xbf\x7d\xf9\xfa\xf8\x43\xc4\x7a\x0b\x12\x69\x9c\x89\x44\xdc\x87\x82\x6b\xc9\xd6\x74\x9d\x89\x35\xcd\x41\x69\x94\x14\xc2\x12\x7d\xa5\xc0\xb6\x77\xa5\x69\x5e\xde\x1e\x5e\x94\x73\xc2\x11\x7b\x03\xf9\x6a\x31\x71\x97\x39\xac\xcf\xb7\x19\xab\x4b\xe7\xdc\x6b\xd1\xfb\x52\x16\x31\x55\x80\x0e\x53\x3c\xfd\xa8\x37\xe1\x22\x6b\xde\x2b\x5d\x1f\x63\xca\x75\x77\xd6\x31\xda\x2b\xc2\x39\x63\xf7\x8f\x4d\x32\x7d\x74\x56\xa2\xfe\x27\x18\xe2\x19\x83\xe4\x98\x51\x58\x09\xb3\x09\x53\x83\x1c\x35\xa5\xc8\xc8\xa4\x22\xfd\x14\x97\xf0\x7d\xfc\x2c\x59\x6c\x2e\xd4\x7f\xf7\xd7\x57\xfb\xeb\xab\xbf\x01\x00\x00\xff\xff\x0e\xdd\x1a\x2c\x9c\x18\x00\x00")

func schemaMashling_schema02JsonBytes() ([]byte, error) {
	return bindataRead(
		_schemaMashling_schema02Json,
		"schema/mashling_schema-0.2.json",
	)
}

func schemaMashling_schema02Json() (*asset, error) {
	bytes, err := schemaMashling_schema02JsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema/mashling_schema-0.2.json", size: 6300, mode: os.FileMode(438), modTime: time.Unix(1510825690, 0)}
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
	"assets/banner.txt": assetsBannerTxt,
	"assets/default_manifest": assetsDefault_manifest,
	"schema/mashling_schema-0.2.json": schemaMashling_schema02Json,
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
	"assets": &bintree{nil, map[string]*bintree{
		"banner.txt": &bintree{assetsBannerTxt, map[string]*bintree{}},
		"default_manifest": &bintree{assetsDefault_manifest, map[string]*bintree{}},
	}},
	"schema": &bintree{nil, map[string]*bintree{
		"mashling_schema-0.2.json": &bintree{schemaMashling_schema02Json, map[string]*bintree{}},
	}},
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

