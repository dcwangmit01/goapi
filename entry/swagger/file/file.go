package file

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

var _swagger_json = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xc4\x52\xcd\xce\x9b\x30\x10\xbc\xf3\x14\xab\x6d\x8f\x51\x48\x73\xcc\x3d\x87\x1e\xfa\x04\x55\x0e\x0e\x2c\xe0\x08\x6c\x77\x77\x49\x1b\x55\xbc\x7b\x65\xe7\x07\x43\x5a\xa9\xb7\x8f\x03\x32\xde\xd9\xd1\xcc\x30\xbf\x0b\x00\x94\x9f\xa6\x6d\x89\xf1\x00\xb8\xdf\xee\x70\x13\xef\xac\x6b\x3c\x1e\x20\xce\x01\x50\xad\xf6\x14\xe7\xe4\x94\x6f\x65\x7a\x6f\x03\x7b\xf5\x09\x0d\x80\x57\x62\xb1\xde\x45\xcc\xe3\x08\xce\x2b\x08\x29\x16\x00\x53\xe2\x94\xaa\xa3\x81\x04\x0f\xf0\xfd\xbe\xd4\xa9\x86\x27\x41\x3c\x4b\xc4\x9e\x12\xb6\xf2\x4e\xc6\x05\xd8\x84\xd0\xdb\xca\xa8\xf5\xae\xbc\x88\x77\x33\x36\xb0\xaf\xc7\xea\x3f\xb1\x46\x3b\x99\x8d\x95\xd7\x2f\x25\xfd\x32\x43\xe8\xa9\xa4\xaa\x9b\x2d\x47\xa8\x17\xcd\xbe\x01\xd0\x07\xe2\x44\xfa\xb5\x8e\x46\x8f\x71\x61\x33\x8f\x99\x24\x78\x27\x24\x8b\x2d\x00\xdc\xef\x76\xab\x2b\x00\xac\x49\x2a\xb6\x41\x1f\xa9\x65\x44\x69\x9c\xc2\x32\x6f\x6b\x00\xf8\x99\xa9\x89\x1b\x9f\xca\x9a\x1a\xeb\x6c\x64\x90\xfb\x3f\x89\x8a\xbe\x91\x88\x69\x09\x17\x6b\x53\xf1\xb7\xf3\x94\x89\x0f\x86\xcd\x40\x4a\x3c\xc7\x78\x7f\x56\xb2\x9d\x19\x52\x13\xce\xbe\xbe\xad\x35\x5b\xf7\xaf\x09\xd3\x8f\xd1\x32\xc5\xd8\x94\x47\xfa\x00\xaf\xa7\xcc\xab\x9a\x76\xed\x12\x8f\x91\x73\x26\x3a\x15\x39\xc1\xf4\xea\x70\x26\x63\x6e\xd1\x9b\x9e\xac\x45\x7a\x0b\x29\x2f\x7f\xbe\x50\xa5\xaf\x5c\x62\x69\x03\xb1\xda\x55\x5b\xf0\x6a\xfa\x91\xd6\x05\x7a\x92\x88\xb2\x75\xed\x22\x5c\x6c\x3c\x0f\x46\xb3\x69\xb1\x36\x9f\x79\x28\xa6\xe2\x4f\x00\x00\x00\xff\xff\xd6\x64\x2b\x4c\xf3\x03\x00\x00")

func swagger_json() ([]byte, error) {
	return bindata_read(
		_swagger_json,
		"swagger.json",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
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
var _bindata = map[string]func() ([]byte, error){
	"swagger.json": swagger_json,
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
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() ([]byte, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"swagger.json": &_bintree_t{swagger_json, map[string]*_bintree_t{
	}},
}}
