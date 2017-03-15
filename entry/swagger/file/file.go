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

var _swagger_json = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xc4\x52\xbd\xce\x9b\x40\x10\xec\x79\x8a\xd5\x26\xe5\xa7\x0f\xc7\xa5\xbb\x14\x2e\x52\x44\x4a\x1f\xb9\x38\xc3\x02\x67\xc1\xdd\x65\x77\x71\x6c\x45\xbc\x7b\x74\xe7\x1f\x0e\x9c\x48\xe9\x42\x81\x8e\xdb\xd9\xd1\xcc\x30\xbf\x0a\x00\x94\x9f\xa6\x6d\x89\x71\x07\xb8\x7d\xdf\xe0\x5b\xbc\xb3\xae\xf1\xb8\x83\x38\x07\x40\xb5\xda\x53\x9c\x93\x53\xbe\x96\xe9\xfd\x1e\xd8\xab\x4f\x68\x00\x3c\x13\x8b\xf5\x2e\x62\xee\x47\x70\x5e\x41\x48\xb1\x00\x98\x12\xa7\x54\x1d\x0d\x24\xb8\x83\xef\xb7\xa5\x4e\x35\x3c\x08\xe2\x59\x22\xf6\x90\xb0\x95\x77\x32\x2e\xc0\x26\x84\xde\x56\x46\xad\x77\xe5\x49\xbc\x9b\xb1\x81\x7d\x3d\x56\xff\x88\x35\xda\xc9\x6c\xac\x3c\x7f\x2a\xe9\x62\x86\xd0\x53\x49\x55\x37\x5b\x8e\x50\x2f\x9a\x7d\x03\xa0\x0f\xc4\x89\xf4\x4b\x1d\x8d\xee\xe3\xc2\xdb\x3c\x66\x92\xe0\x9d\x90\x2c\xb6\x00\x70\xbb\xd9\xac\xae\x00\xb0\x26\xa9\xd8\x06\xbd\xa7\x96\x11\xa5\x71\x0a\xcb\xbc\xac\x01\xe0\x47\xa6\x26\x6e\x7c\x28\x6b\x6a\xac\xb3\x91\x41\x6e\xff\x24\x2a\xfa\x4a\x22\xa6\x25\x5c\xac\x4d\xc5\x9f\xce\x53\x26\x3e\x18\x36\x03\x29\xf1\x1c\xe3\xed\x59\xc9\x76\x66\x48\x4d\x38\xfa\xfa\xba\xd6\x6c\xdd\xdf\x26\x4c\x3f\x46\xcb\x14\x63\x53\x1e\xe9\x3f\x78\x3d\x64\x5e\xd5\xb4\x6b\x97\xf8\xcd\x5c\xf6\x91\xf6\x73\xb0\x33\xdd\xa1\xc8\x69\xa6\x67\x93\x33\x31\x73\x97\x5e\x54\x65\x5d\xd2\x6b\x48\xa9\xf9\xe3\x89\x2a\x7d\xa6\x13\xab\x1b\x88\xd5\xae\x3a\x83\x67\xd3\x8f\xb4\xae\xd1\x83\x44\x94\xad\x6b\x17\x11\x63\xe3\x79\x30\x9a\x4d\x8b\x75\x04\x99\x87\x62\x2a\x7e\x07\x00\x00\xff\xff\x58\xd4\x38\xba\xf9\x03\x00\x00")

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
