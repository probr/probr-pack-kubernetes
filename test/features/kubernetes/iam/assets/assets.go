// Code generated for package iam by go-bindata DO NOT EDIT. (@generated)
// sources:
// assets/yaml/iam-azi-test-aib.yaml
package iam

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

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _assetsYamlIamAziTestAibYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x93\xcf\x6b\x1b\x3f\x10\xc5\xef\x06\xff\x0f\x43\x72\xd6\x7a\xd7\xd9\xac\x6d\x41\x0e\x26\xc9\xf7\x4b\x0e\x6e\x96\x26\x2d\x94\x52\x8c\x56\x9a\xb5\x55\x6b\x35\x42\x3f\xdc\xe4\xbf\x2f\x6b\x7b\x89\x03\x6d\x2f\xed\xf5\x33\xbc\x37\x7a\x8f\x91\x70\xfa\x33\xfa\xa0\xc9\x72\xd8\x17\xe3\xd1\x4e\x5b\xc5\xa1\x26\x35\x1e\x75\x18\x85\x12\x51\xf0\xf1\x08\xc0\x8a\x0e\x39\x5c\x3e\x2c\x57\xcf\xf7\x4f\xcf\x97\x3d\x32\xa2\x41\x13\x0e\x53\x00\x21\x94\x23\xa5\x55\xa3\xad\xd2\x76\xc3\x41\x61\x47\x00\x97\xcf\x8f\x77\x8f\x1c\x2c\xa2\x82\x48\xa0\xed\x77\x94\x11\xe2\x56\x87\x5e\x26\xac\xa5\x28\xa2\x26\x3b\xd8\x04\x94\x92\x3a\x97\x05\x94\xc9\xeb\xf8\x9a\x09\xe3\xb6\x22\xdb\xa5\x06\xbd\xc5\x88\x21\xd3\x34\x71\xa4\x38\x5c\xf8\x64\xa3\xee\x70\xa2\xb0\x15\xc9\xc4\x8b\xf1\x28\x38\x94\x07\x9f\x41\x7d\x4b\x36\xe2\x4b\x3c\x7a\xfb\x64\x97\xe1\x53\x40\xcf\xa1\xc8\xf3\xfc\x8d\xfd\xef\x29\x39\x0e\x57\x03\x6c\x07\x32\x1d\x48\x48\xce\x19\xec\xd0\x46\x61\x0e\xb3\xc0\xe1\x2b\x14\xf0\xad\x9f\x4a\xb2\x51\x68\x8b\xfe\x98\x81\x9d\xaa\xea\xf3\x1f\x33\xe9\x4e\x6c\x90\x43\x27\x7d\xd6\x69\xe9\x29\x50\x1b\x33\x49\xdd\x64\x37\x0f\x13\x21\x14\x73\xa4\x98\x56\x68\xa3\x8e\xaf\x93\x5e\xc7\x8b\x6c\x7a\xaa\xd5\x6f\x86\x6a\x7a\x6b\xc6\x42\x6a\x82\xf4\xda\xf5\xa5\x69\x75\x93\xcf\x67\x65\xb5\x98\x21\xab\xaa\xab\x82\x95\xb2\xa9\x58\xd3\xcc\x5a\xa6\xaa\xab\xf9\xb4\x5a\x48\xac\xd4\xec\x5c\x2e\x8d\xee\xf7\xa8\x9b\xb2\x44\xd5\x54\x62\xc6\xf2\x36\x9f\xb1\xb2\x91\x05\x5b\x5c\x97\x73\xb6\x28\x64\xbb\xc8\x0b\x5c\x94\xd3\xeb\x73\xa1\xc7\x40\xc9\x4b\xdc\xf4\xf1\x6f\x56\xb7\x6b\xbf\x61\xce\x53\xe3\x99\x30\x86\x39\x32\x5a\x6a\x0c\xeb\xba\x47\x4b\x63\xea\x01\xa4\x5d\xa0\x14\xb7\x47\x2b\xb4\xfb\xb3\x30\xc7\x9e\x56\x5f\xd6\xf5\xe3\xdd\xfa\xc3\x72\x75\x3f\x8c\x00\xf6\xc2\x24\xfc\xcf\x53\xc7\xdf\x18\x40\xab\xd1\xa8\x8f\xd8\xbe\x83\x27\x5c\x8b\xb8\xe5\x30\x9c\x6c\xd6\x7b\xff\x61\xd3\x53\xbd\xbc\xfd\xd7\xeb\x82\x13\xf2\x77\x3b\x1f\xea\xbf\x5f\x16\xa2\x88\x29\x64\x8e\xd4\xe0\xf6\xcb\x33\xef\x8f\xc6\x18\xfa\x51\x7b\xbd\xd7\x06\x37\x78\x1f\xa4\x30\x87\x4f\xc6\xa1\x15\x26\x1c\xde\x68\x49\xe1\x13\x1a\x94\x91\xfc\x49\xf7\xfe\x8f\x51\xe0\x60\xb4\x4d\x2f\x3f\x03\x00\x00\xff\xff\xb7\x8f\xeb\xb0\x21\x04\x00\x00")

func assetsYamlIamAziTestAibYamlBytes() ([]byte, error) {
	return bindataRead(
		_assetsYamlIamAziTestAibYaml,
		"assets/yaml/iam-azi-test-aib.yaml",
	)
}

func assetsYamlIamAziTestAibYaml() (*asset, error) {
	bytes, err := assetsYamlIamAziTestAibYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/yaml/iam-azi-test-aib.yaml", size: 1057, mode: os.FileMode(438), modTime: time.Unix(1600092251, 0)}
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
	"assets/yaml/iam-azi-test-aib.yaml": assetsYamlIamAziTestAibYaml,
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
		"yaml": &bintree{nil, map[string]*bintree{
			"iam-azi-test-aib.yaml": &bintree{assetsYamlIamAziTestAibYaml, map[string]*bintree{}},
		}},
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
