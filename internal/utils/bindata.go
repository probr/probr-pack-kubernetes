// Code generated for package utils by go-bindata DO NOT EDIT. (@generated)
// sources:
// test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-hostport-approved.yaml
// test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-hostport-unapproved.yaml
package utils

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

var _testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpHostportApprovedYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x91\xbb\x6e\xf3\x30\x0c\x85\x77\x03\x7e\x07\xc2\xff\xfa\xfb\xd6\x6e\xda\x8a\xb6\xe8\x54\x20\x4b\xbb\x04\x1d\x18\x99\x89\x85\x48\xa2\xa0\x8b\x9b\xa0\xe8\xbb\x17\x8a\xeb\x24\x43\xea\xc1\xc3\x47\xf2\xe8\xf0\x10\x9d\x7a\x27\x1f\x14\x5b\x01\x53\x5f\x16\x7b\x65\x07\x01\x2b\x1e\xca\xc2\x50\xc4\x01\x23\x8a\xb2\x00\xf8\x67\xd1\x90\x80\x40\x32\x79\x15\x8f\xb5\x64\x1b\xe9\x10\xeb\x81\x0c\xe7\x3a\x5a\xcb\x11\xa3\x62\x1b\x4e\xfd\x90\x5b\x25\x1b\xd7\x2c\x23\x0d\x6a\x37\x62\xb3\x4f\x1b\xf2\x96\x22\x85\x46\x71\xeb\x78\x10\x50\xf9\x64\xa3\x32\xd4\x0e\xb4\xc5\xa4\x63\x55\x16\xc1\x91\x14\x00\x59\x69\x99\x7f\x9c\x5f\xfc\x55\xf7\xc9\x3e\x84\xb7\x40\x5e\x40\xdf\x75\x5d\x46\x57\x85\x17\xcf\xc9\x09\xb8\xef\xba\x6e\xa6\xdb\x05\xdd\x9d\x51\x48\xce\x69\x32\x64\x23\xea\x53\x31\x08\x58\x43\x0f\x1f\x8b\xd6\xc4\x3a\x19\x9a\xd7\xa9\xe1\xbc\x7f\x2d\xe3\xa1\x9e\x58\xcf\x2a\x64\x5c\x3c\x3e\x29\x2f\xe0\xeb\x3b\x93\x9c\x0b\x2a\x4b\xfe\xf6\xdc\x12\x17\x80\x32\xb8\x23\x01\x9b\x14\x8e\x1b\x3e\xcc\x4c\xb2\x31\x98\xf3\x5f\x43\x15\xc6\xea\x3f\x54\xb5\xcc\xff\xa0\x89\x1c\xf4\x63\x75\x31\xb7\xd8\x7b\xe5\x64\xe3\x12\xf9\x9f\x2e\x01\x4c\xee\x5b\x61\x1c\x05\xb4\xf9\xa8\xed\xc5\xc9\xed\x80\x01\x50\x6b\xfe\x5c\x79\x35\x29\x4d\x3b\x7a\x0e\x12\xf5\xe9\xc0\x02\xb6\xa8\x03\xc1\xd5\xf7\x13\x00\x00\xff\xff\xbd\xec\xbb\x0d\x46\x02\x00\x00")

func testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpHostportApprovedYamlBytes() ([]byte, error) {
	return bindataRead(
		_testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpHostportApprovedYaml,
		"test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-hostport-approved.yaml",
	)
}

func testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpHostportApprovedYaml() (*asset, error) {
	bytes, err := testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpHostportApprovedYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-hostport-approved.yaml", size: 582, mode: os.FileMode(438), modTime: time.Unix(1597863384, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpHostportUnapprovedYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x52\xb1\x8e\x13\x31\x10\xed\x23\xe5\x1f\x9e\x96\x96\x4d\x36\x20\xa1\x93\x3b\x04\x88\x0a\x29\x0d\x34\x27\x8a\xc9\xee\xe4\xd6\x3a\xdb\x63\xd9\xe3\x25\x11\xe2\xdf\x91\xb3\xec\x25\x42\x39\x17\x2e\xde\x9b\x79\xe3\xf7\x3c\x14\xed\x0f\x4e\xd9\x4a\x30\x98\x76\xeb\xd5\xb3\x0d\x83\xc1\x5e\x86\xf5\xca\xb3\xd2\x40\x4a\x66\xbd\x02\xde\x04\xf2\x6c\x90\xb9\x2f\xc9\xea\xb9\xed\x25\x28\x9f\xb4\x1d\xd8\x4b\xe5\x29\x04\x51\x52\x2b\x21\x5f\xea\x51\x4b\x7b\xf1\x71\xb3\xb4\x6c\xc8\xc5\x91\x36\xcf\xe5\xc0\x29\xb0\x72\xde\x58\xd9\x46\x19\x0c\x9a\x54\x82\x5a\xcf\xdb\x81\x8f\x54\x9c\x36\xeb\x55\x8e\xdc\x1b\xa0\x2a\x2d\xfd\x9f\xe6\x89\xff\xd4\x53\x09\x1f\xf3\xf7\xcc\xc9\x60\xd7\x75\x5d\x85\x6e\x88\xaf\x49\x4a\x34\x78\xdf\x75\xdd\x8c\x1e\x17\xe8\xdd\x0b\x94\x4b\x8c\x8e\x3d\x07\x25\x77\x21\xb3\xc1\x23\x76\xf8\xb9\x68\x4d\xe2\x8a\xe7\xd9\x4e\x8b\x17\xff\x6d\xaf\xa7\x76\x12\x37\xab\xb0\x8f\x7a\xfe\x6c\x93\xc1\xef\x3f\x15\xa9\xb9\x90\x0d\x9c\xee\xf7\x2d\x71\x01\xd6\xd3\x13\x1b\x1c\x4a\x3e\x1f\xe4\x34\x63\xbd\x78\x4f\x35\xff\x47\x34\x79\x6c\xde\xa2\x69\xfb\x7a\x67\xc7\x1c\xb1\x1b\x9b\xeb\xe3\x80\x28\x49\x97\xac\xdb\xeb\xdc\xbd\x24\x35\x78\xe8\x1e\x3e\xcc\x14\x30\x4a\xd6\xff\xd1\xd9\xdb\x37\x29\xe1\x46\xe3\x15\x8b\x80\xaf\x75\x7b\xd2\xd1\x60\x5b\x37\x62\x7b\xb5\x71\xff\x77\x00\x72\x4e\x7e\xed\x93\x9d\xac\xe3\x27\xfe\x92\x7b\x72\x97\xed\x30\x38\x92\xcb\x8c\x9b\xf3\x37\x00\x00\xff\xff\x48\xd0\xb0\x8c\x83\x02\x00\x00")

func testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpHostportUnapprovedYamlBytes() ([]byte, error) {
	return bindataRead(
		_testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpHostportUnapprovedYaml,
		"test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-hostport-unapproved.yaml",
	)
}

func testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpHostportUnapprovedYaml() (*asset, error) {
	bytes, err := testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpHostportUnapprovedYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-hostport-unapproved.yaml", size: 643, mode: os.FileMode(438), modTime: time.Unix(1597863427, 0)}
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
	"test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-hostport-approved.yaml":   testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpHostportApprovedYaml,
	"test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-hostport-unapproved.yaml": testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpHostportUnapprovedYaml,
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
	"test": &bintree{nil, map[string]*bintree{
		"features": &bintree{nil, map[string]*bintree{
			"kubernetes": &bintree{nil, map[string]*bintree{
				"podsecuritypolicy": &bintree{nil, map[string]*bintree{
					"features": &bintree{nil, map[string]*bintree{
						"yaml": &bintree{nil, map[string]*bintree{
							"psp-azp-hostport-approved.yaml":   &bintree{testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpHostportApprovedYaml, map[string]*bintree{}},
							"psp-azp-hostport-unapproved.yaml": &bintree{testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpHostportUnapprovedYaml, map[string]*bintree{}},
						}},
					}},
				}},
			}},
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
