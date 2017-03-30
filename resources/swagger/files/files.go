package files

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

var _swagger_json = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xc4\x52\x31\x6e\xdc\x40\x0c\xec\xf5\x0a\x82\x49\x69\x58\x17\x97\xd7\xbb\x48\x91\x17\x04\x57\xac\x25\x4a\x5a\x43\xda\x65\x48\xea\x92\x43\xa0\xbf\x07\xbb\xf6\x9d\x56\xba\x04\x48\x67\x55\xd4\xce\x70\xc0\x19\xcc\xef\x0a\x00\xf5\xa7\xeb\x7b\x12\x3c\x02\x3e\x3d\x1e\xf0\x21\xbd\xf9\xd0\x45\x3c\x42\xc2\x01\xd0\xbc\x8d\x94\x70\xc7\x5c\x3b\xe6\x47\x96\x68\x31\x33\x01\xf0\x4c\xa2\x3e\x86\x84\xbf\x8f\x10\xa2\x81\x92\x61\x05\xb0\x64\x3d\x6d\x06\x9a\x48\xf1\x08\xdf\xdf\x96\x06\x33\xbe\x0a\xa4\x59\x13\xf7\x94\xb9\x4d\x0c\x3a\x6f\xc8\x8e\x79\xf4\x8d\x33\x1f\x43\xfd\xaa\x31\xac\x5c\x96\xd8\xce\xcd\x7f\x72\x9d\x0d\xba\x9a\xaa\xcf\x5f\x6a\xfa\xe5\x26\x1e\xa9\xa6\x66\x58\xed\x26\x6a\x54\x2b\xfe\x01\x30\x32\x49\x16\xfd\xda\x26\xa3\xcf\x69\xe1\x61\x85\x85\x94\x63\x50\xd2\xcd\x16\x00\x3e\x1d\x0e\xbb\x27\x00\x6c\x49\x1b\xf1\x6c\xef\xa9\x15\x42\x19\xce\x61\xb9\xbb\x35\x00\xfc\x2c\xd4\xa5\x8d\x4f\x75\x4b\x9d\x0f\x3e\x29\x68\x4d\xc1\xe4\x92\x2e\xfa\x46\xaa\xae\x27\xdc\xac\x2d\xd5\xdf\xe6\xa5\x38\x9e\x9d\xb8\x89\x8c\x64\x8d\xf1\xed\xdb\x9d\x1d\xdc\x94\x5b\xf0\x12\xdb\xcb\xfe\x66\x1f\xfe\x85\x08\xfd\x98\xbd\x50\x8a\xcd\x64\xa6\x0f\xf0\x7a\x2a\xbc\x9a\xeb\xf7\x2e\xf1\x39\x69\xae\x42\xa7\xaa\x14\x58\x6e\x1d\x2e\xce\x58\x5b\x74\x77\x4f\xd1\x22\xbb\x70\xce\x2b\xbe\xbc\x52\x63\xb7\x5c\x52\x69\x99\xc4\xfc\xae\x2d\x78\x76\xe3\x4c\xfb\x02\x5d\x45\xd4\xc4\x87\x7e\x13\x2e\x76\x51\x26\x67\x05\x5a\xed\xcd\x17\x1e\xaa\xa5\xfa\x13\x00\x00\xff\xff\xf5\x03\xae\x2c\xef\x03\x00\x00")

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
	Func     func() ([]byte, error)
	Children map[string]*_bintree_t
}

var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"swagger.json": &_bintree_t{swagger_json, map[string]*_bintree_t{}},
}}
