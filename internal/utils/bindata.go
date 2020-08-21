// Code generated for package utils by go-bindata DO NOT EDIT. (@generated)
// sources:
// test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-hostport-approved.yaml
// test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-hostport-unapproved.yaml
// test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-volumetypes-approved.yaml
// test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-volumetypes-unapproved.yaml
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

var _testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpVolumetypesApprovedYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x92\xbd\xae\x9c\x30\x10\x85\xfb\x95\xf6\x1d\x46\xd4\x61\x17\x92\xce\x5d\x14\x45\xa9\xae\xb4\x4d\xd2\x5c\x45\x57\xb3\x66\x58\x2c\x6c\x8f\xe5\x1f\xb2\xbc\x7d\x64\x0c\x0b\x45\x14\x0a\x84\xe6\x3b\x3e\xcc\x9c\x31\x3a\xf5\x8b\x7c\x50\x6c\x05\x4c\xed\xf9\x34\x2a\xdb\x09\xb8\x71\x77\x3e\x19\x8a\xd8\x61\x44\x71\x3e\x01\x58\x34\x24\x20\x90\x4c\x5e\xc5\xb9\x96\x6c\x23\x3d\x63\xdd\x91\xe1\x8c\xd1\x5a\x8e\x18\x15\xdb\xb0\xc8\x21\x4b\x25\x1b\x77\xd9\x8e\x5c\x50\xbb\x01\x2f\x63\xba\x93\xb7\x14\x29\x5c\x14\x5f\x1d\x77\x02\x2a\x9f\x6c\x54\x86\xae\x1d\xf5\x98\x74\xac\xce\xa7\xe0\x48\x0a\x80\xec\xb4\x9d\xff\x56\xfe\xb8\xba\xfb\x64\xbf\x86\x9f\x81\xbc\x80\xb6\x69\x9a\x5c\x3a\x80\x1f\x9e\x93\x13\xf0\xa5\x69\x9a\x52\xed\xb7\xd2\xe7\x57\x29\x24\xe7\x34\x19\xb2\x11\xf5\x02\x83\x80\x77\x68\xe1\xf7\xe6\x35\xb1\x4e\x86\xd6\x71\xea\x35\x00\x74\xce\xf3\x44\x5d\x5d\x68\xb1\x02\xc9\xb6\x57\x8f\x37\x74\x6b\x77\xf9\x29\xfa\x42\x3e\x0c\xba\x9d\xa8\x48\x26\x1c\x94\xd9\x7d\xa4\x59\x80\x99\x3f\x46\x9a\x8f\x00\xc0\x61\x1c\x16\x92\x3f\x32\xca\xc9\xa3\xb2\xe4\x8b\x45\xbd\x6f\xa6\x96\xf1\xf9\x5a\x08\x80\x32\xf8\x20\x01\xf7\x14\xe6\x3b\x3f\x4b\x4d\xb2\x31\x98\x17\xfc\x0e\x55\x18\xaa\x4f\x50\xd5\x32\xbf\x83\x26\x72\xd0\x0e\x55\x19\x7f\x8f\xb3\x8c\xf9\xc6\xc9\xc6\xad\xe5\xff\x47\x01\x26\x6b\x6f\x4b\xd7\xd7\x7c\x7b\xae\x7b\x47\x9b\xe9\xbf\x57\x0a\x80\x5a\xf3\x9f\x9b\x57\x93\xd2\xf4\xa0\xef\x41\xa2\x5e\xae\x94\x80\x1e\x75\xa0\x63\x2c\x7f\x03\x00\x00\xff\xff\xb9\x6c\x7f\xf1\xb7\x02\x00\x00")

func testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpVolumetypesApprovedYamlBytes() ([]byte, error) {
	return bindataRead(
		_testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpVolumetypesApprovedYaml,
		"test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-volumetypes-approved.yaml",
	)
}

func testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpVolumetypesApprovedYaml() (*asset, error) {
	bytes, err := testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpVolumetypesApprovedYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-volumetypes-approved.yaml", size: 695, mode: os.FileMode(438), modTime: time.Unix(1598032429, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpVolumetypesUnapprovedYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x52\x3d\xab\xdc\x30\x10\xec\x0f\xfc\x1f\x16\xd7\xf1\x57\xd2\xa9\x0b\x21\xa4\x0a\xb8\x49\x9a\x47\x8a\x3d\x79\xef\x2c\x9e\xa4\x15\x5a\xc9\xb9\xf7\xef\x83\x4e\xf6\xbb\x2b\x02\xcf\x85\xc1\x33\xab\xf1\xec\x8c\x30\x98\xdf\x14\xc5\xb0\x57\xb0\x4d\xcd\xe9\xd5\xf8\x45\xc1\xcc\x4b\x73\x72\x94\x70\xc1\x84\xaa\x39\x01\x78\x74\xa4\x40\x48\xe7\x68\xd2\x5b\xa7\xd9\x27\xba\xa5\x6e\x21\xc7\x85\x46\xef\x39\x61\x32\xec\xe5\x3e\x0e\x65\x54\xb3\x0b\xfd\x71\xa4\x47\x1b\x56\xec\x5f\xf3\x99\xa2\xa7\x44\xd2\x1b\x1e\x02\x2f\x0a\xda\x98\x7d\x32\x8e\x86\x85\x2e\x98\x6d\x6a\x9b\x93\x04\xd2\x0a\xa0\x28\x1d\xe7\xbf\xd5\x3f\xee\xea\x31\xfb\xaf\xf2\x4b\x28\x2a\x98\xc6\x71\x2c\xd0\x13\xf1\x23\x72\x0e\x0a\xbe\x8c\xe3\x58\xd1\xcb\x01\x7d\x7e\x87\x24\x87\x60\xc9\x91\x4f\x68\xef\xa4\x28\x78\x81\x09\xfe\x1c\x5a\x1b\xdb\xec\x68\x5f\xa7\xdb\x03\xc8\x1e\x43\x88\xbc\xd1\xd2\x55\xbe\x8a\xc1\xca\x92\x66\x4c\xeb\x6e\x0f\x20\x94\x0f\x18\x36\x8c\x83\xe5\x6b\x59\x54\x0a\x55\x72\x43\xe3\x29\x56\xdd\xee\x91\x6b\xa7\xd3\xed\x3d\x4e\x00\xe3\xf0\x4a\x0a\xce\x59\xde\xce\x7c\xab\x98\x66\xe7\xb0\xd4\xf3\x02\xad\xac\xed\x27\x68\x3b\x5d\xde\x62\x89\x02\x4c\x6b\x5b\xcd\x3f\xc2\xa8\x16\x7f\x72\xf6\xe9\xa8\xe5\xa3\x45\xc0\x95\xe9\xb9\xba\x2f\xed\x0f\x0f\x4f\x87\xec\xff\x2b\x01\x40\x6b\xf9\xef\x1c\xcd\x66\x2c\x5d\xe9\xbb\x68\xb4\xf7\x2b\xa1\xe0\x82\x56\x08\x9e\x9e\x7f\x01\x00\x00\xff\xff\xd4\x7f\x20\xcb\x77\x02\x00\x00")

func testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpVolumetypesUnapprovedYamlBytes() ([]byte, error) {
	return bindataRead(
		_testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpVolumetypesUnapprovedYaml,
		"test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-volumetypes-unapproved.yaml",
	)
}

func testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpVolumetypesUnapprovedYaml() (*asset, error) {
	bytes, err := testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpVolumetypesUnapprovedYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-volumetypes-unapproved.yaml", size: 631, mode: os.FileMode(438), modTime: time.Unix(1598032435, 0)}
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
	"test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-hostport-approved.yaml":      testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpHostportApprovedYaml,
	"test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-hostport-unapproved.yaml":    testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpHostportUnapprovedYaml,
	"test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-volumetypes-approved.yaml":   testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpVolumetypesApprovedYaml,
	"test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-volumetypes-unapproved.yaml": testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpVolumetypesUnapprovedYaml,
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
							"psp-azp-hostport-approved.yaml":      &bintree{testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpHostportApprovedYaml, map[string]*bintree{}},
							"psp-azp-hostport-unapproved.yaml":    &bintree{testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpHostportUnapprovedYaml, map[string]*bintree{}},
							"psp-azp-volumetypes-approved.yaml":   &bintree{testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpVolumetypesApprovedYaml, map[string]*bintree{}},
							"psp-azp-volumetypes-unapproved.yaml": &bintree{testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpVolumetypesUnapprovedYaml, map[string]*bintree{}},
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
