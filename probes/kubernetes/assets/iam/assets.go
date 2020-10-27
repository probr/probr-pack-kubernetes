// Code generated for package iam by go-bindata DO NOT EDIT. (@generated)
// sources:
// assets/yaml/iam-azi-test-aib-curl.yaml
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

var _assetsYamlIamAziProbeAibCurlYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x91\x4f\x6b\x14\x41\x10\xc5\xef\x0b\xfb\x1d\x8a\xd9\xab\xd9\x3f\x7a\xeb\x5b\xd0\x20\x1e\x24\x81\xac\x5e\x82\x87\xda\xee\xda\xdd\x32\xdd\x55\x4d\x57\x75\x8c\xdf\x5e\x66\x26\x83\x78\xf0\x32\x14\xef\xf1\xde\xf0\x7e\x8d\x95\xbf\x53\x33\x56\x09\xf0\x72\x58\xaf\x9e\x59\x52\x80\x07\x4d\xeb\x55\x21\xc7\x84\x8e\x61\xbd\x02\x10\x2c\x14\x60\xf3\xe5\xf6\xeb\xf1\xee\xf1\xb8\x19\xa5\x8c\x27\xca\x36\xb9\x00\x88\xa9\x6a\xe2\x74\x62\x49\x2c\x97\x00\x9b\x44\x45\x01\x36\xc7\xfb\x4f\xf7\x01\x84\x28\x81\x2b\xb0\xfc\xa4\xe8\xe0\x57\xb6\x31\x87\x22\xea\xe8\xac\xb2\xf4\x18\xc5\xa8\xa5\x6e\x8d\x62\x6f\xec\xbf\xb7\x98\xeb\x15\xb7\xcf\xfd\x44\x4d\xc8\xc9\xb6\xac\xbb\xaa\x29\xc0\xd0\xba\x38\x17\xda\x25\x3a\x63\xcf\x3e\xac\x57\x56\x29\x4e\x3d\x4b\xfa\xa3\x8a\xd3\xab\xcf\xdd\xad\xcb\xad\x7d\x33\x6a\x01\x0e\xfb\xfd\xfe\xaf\xf6\xb9\x69\xaf\x01\x3e\x2c\xe2\x79\x51\xde\x2f\x8a\xf5\x5a\x33\x15\x12\xc7\x3c\x79\x16\xe0\x09\x0e\xf0\x63\x74\xa3\x8a\x23\x0b\xb5\x79\xc3\xcd\x1b\xab\x71\xff\xbc\x89\x0b\x5e\x28\x40\xec\x2d\x4f\xa7\xed\xc6\x73\xf6\xa2\x96\x82\x23\xf2\x27\x18\xec\x3a\xbc\x83\xe1\x26\x8e\x5f\xcb\x44\x15\x0e\xd7\x61\xfe\xc7\x7f\x26\x8d\xe0\x73\xd6\x5f\x0f\x8d\x5f\x38\xd3\x85\xee\x2c\x62\x9e\x80\x06\x38\x63\x36\x9a\xde\x4e\x13\x3d\x52\xa6\xe8\xda\xde\x72\xff\xf2\x54\x0b\x90\x59\xfa\xeb\x9f\x00\x00\x00\xff\xff\x86\x9a\xc1\xd4\x0e\x02\x00\x00")

func assetsYamlIamAziProbeAibCurlYamlBytes() ([]byte, error) {
	return bindataRead(
		_assetsYamlIamAziProbeAibCurlYaml,
		"assets/yaml/iam-azi-test-aib-curl.yaml",
	)
}

func assetsYamlIamAziProbeAibCurlYaml() (*asset, error) {
	bytes, err := assetsYamlIamAziProbeAibCurlYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/yaml/iam-azi-test-aib-curl.yaml", size: 526, mode: os.FileMode(438), modTime: time.Unix(1600355666, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsYamlIamAziProbeAibYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x93\xcf\x6b\x1b\x3f\x10\xc5\xef\x06\xff\x0f\x43\x7c\xd6\x7a\xd7\xd9\xac\x6d\x41\x0e\x26\xc9\xf7\x4b\x0e\x6e\x96\x26\x2d\x94\x52\x8c\x56\x9a\xb5\x55\x6b\x35\x42\x3f\xdc\xe4\xbf\x2f\x6b\x7b\x49\x02\x6d\x2f\xed\xf5\x33\xbc\x37\x7a\x8f\x91\x70\xfa\x33\xfa\xa0\xc9\x72\x38\x14\xe3\xd1\x5e\x5b\xc5\xa1\x26\x35\x1e\x75\x18\x85\x12\x51\xf0\xf1\x08\xc0\x8a\x0e\x39\x4c\xee\x57\xeb\xa7\xbb\xc7\xa7\x49\x8f\x8c\x68\xd0\x84\xe3\x14\x40\x08\xe5\x48\x69\xd5\x68\xab\xb4\xdd\x72\x98\x28\xec\x08\x60\xf2\xf4\x70\xfb\xc0\xc1\x22\x2a\x88\x04\xda\x7e\x47\x19\x21\xee\x74\xe8\x75\xc2\x5a\x8a\x22\x6a\xb2\x83\x4f\x40\x29\xa9\x73\x59\x40\x99\xbc\x8e\x2f\x99\x30\x6e\x27\xb2\x7d\x6a\xd0\x5b\x8c\x18\x32\x4d\x53\x47\x8a\xc3\x85\x4f\x36\xea\x0e\xa7\x0a\x5b\x91\x4c\xbc\x18\x8f\x82\x43\x79\xf4\x19\xd4\x37\x64\x23\x3e\xc7\x93\xb7\x4f\x76\x15\x3e\x05\xf4\x1c\x8a\x3c\xcf\x5f\xd9\xff\x9e\x92\xe3\x70\x39\xc0\x76\x20\xb3\x81\x84\xe4\x9c\xc1\x0e\x6d\x14\xe6\x38\x0b\x1c\xbe\x42\x01\xdf\xfa\xa9\x24\x1b\x85\xb6\xe8\x4f\x19\xd8\xb9\xab\x3e\xff\x29\x93\xee\xc4\x16\x39\x74\xd2\x67\x9d\x96\x9e\x02\xb5\x31\x93\xd4\x4d\xf7\x8b\x30\x15\x42\x31\x47\x8a\x69\x85\x36\xea\xf8\x32\xed\x75\xbc\xc8\x66\xe7\x5e\xfd\x76\xa8\xa6\xb7\x66\x2c\xa4\x26\x48\xaf\x5d\x5f\x9a\x56\xd7\xf9\x62\x5e\x56\xcb\x39\xb2\xaa\xba\x2c\x58\x29\x9b\x8a\x35\xcd\xbc\x65\xaa\xba\x5c\xcc\xaa\xa5\xc4\x4a\xcd\xdf\xca\xa5\xd1\xfd\x1e\x75\x5d\x96\xa8\x9a\x4a\xcc\x59\xde\xe6\x73\x56\x36\xb2\x60\xcb\xab\x72\xc1\x96\x85\x6c\x97\x79\x81\xcb\x72\x76\xf5\x56\xe8\x31\x50\xf2\x12\xb7\x7d\xfc\xeb\xf5\xcd\xc6\x6f\x99\xf3\xd4\x78\x26\x8c\x61\x8e\x8c\x96\x1a\xc3\xa6\xee\xd1\xca\x98\x7a\x00\x69\x1f\x28\xc5\xdd\xc9\x0a\xed\xe1\x4d\x98\x53\x4f\xeb\x2f\x9b\xfa\xe1\x76\xf3\x61\xb5\xbe\x1b\x46\x00\x07\x61\x12\xfe\xe7\xa9\xe3\xaf\x0c\xa0\xd5\x68\xd4\x47\x6c\xdf\xc1\x33\xae\x45\xdc\x71\x18\x6e\x36\xeb\xbd\xff\xb0\xe9\xb1\x5e\xdd\xfc\xeb\x75\xc1\x09\xf9\xbb\x9d\xf7\xf5\xdf\x2f\x0b\x51\xc4\x14\x32\x47\x6a\x70\xfb\xe5\x99\xf7\x47\x63\x0c\xfd\xa8\xbd\x3e\x68\x83\x5b\xbc\x0b\x52\x98\xe3\x27\xe3\xd0\x0a\x13\x8e\x6f\xb4\xa4\xf0\x11\x0d\xca\x48\xfe\xac\x7b\xff\xc7\x28\x70\x30\xda\xa6\xe7\x9f\x01\x00\x00\xff\xff\xfc\x4d\xb4\xd8\x22\x04\x00\x00")

func assetsYamlIamAziProbeAibYamlBytes() ([]byte, error) {
	return bindataRead(
		_assetsYamlIamAziProbeAibYaml,
		"assets/yaml/iam-azi-test-aib.yaml",
	)
}

func assetsYamlIamAziProbeAibYaml() (*asset, error) {
	bytes, err := assetsYamlIamAziProbeAibYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/yaml/iam-azi-test-aib.yaml", size: 1058, mode: os.FileMode(438), modTime: time.Unix(1600274827, 0)}
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
	"assets/yaml/iam-azi-test-aib-curl.yaml": assetsYamlIamAziProbeAibCurlYaml,
	"assets/yaml/iam-azi-test-aib.yaml":      assetsYamlIamAziProbeAibYaml,
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
			"iam-azi-test-aib-curl.yaml": &bintree{assetsYamlIamAziProbeAibCurlYaml, map[string]*bintree{}},
			"iam-azi-test-aib.yaml":      &bintree{assetsYamlIamAziProbeAibYaml, map[string]*bintree{}},
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
