// Code generated for package utils by go-bindata DO NOT EDIT. (@generated)
// sources:
// test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-hostport-approved.yaml
// test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-hostport-unapproved.yaml
// test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-seccomp-approved.yaml
// test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-seccomp-unapproved.yaml
// test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-seccomp-undefined.yaml
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

var _testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpSeccompApprovedYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x90\x3f\x4f\xc3\x30\x10\xc5\xf7\x48\xf9\x0e\xa7\xcc\x24\x4d\x60\xf3\x86\x10\x62\xed\x02\x4b\xd5\xe1\xb0\xaf\xad\x55\xff\x93\xef\x5c\xca\xb7\x47\x6e\x48\xe8\x40\x86\x0c\xef\xbd\x7b\xd6\xef\x61\xb2\x1f\x94\xd9\xc6\xa0\xe0\x32\xb5\xcd\xd9\x06\xa3\x60\x1b\x4d\xdb\x78\x12\x34\x28\xa8\xda\x06\x20\xa0\x27\x05\x4c\xba\x64\x2b\xdf\xbd\x8e\x41\xe8\x2a\xbd\x21\x1f\xab\x8d\x21\x44\x41\xb1\x31\xf0\x2d\x0e\x35\xaa\xa3\x4f\xc3\x72\x32\xa0\x4b\x27\x1c\xce\xe5\x93\x72\x20\x21\x1e\x6c\xdc\xa4\x68\x14\x74\xb9\x04\xb1\x9e\x36\x86\x0e\x58\x9c\x74\x6d\xc3\x89\xb4\x02\xa8\x4d\xcb\xfd\xcb\xfc\xe2\x6f\x7b\x2e\xe1\x99\xdf\x99\xb2\x82\x69\x1c\xc7\x2a\xdd\x19\x6f\x39\x96\xa4\xe0\x69\x1c\xc7\x59\x3d\x2c\xd2\xe3\x2a\x71\x49\xc9\x91\xa7\x20\xe8\x6e\x26\x2b\xd8\xc1\x04\x7b\x80\xa5\xad\x42\xa2\x0d\x94\x67\xa6\xfe\x6f\x84\x5e\xcb\x75\x65\x07\xb0\x1e\x8f\xa4\xc0\x44\x7d\xa6\x5c\xb9\xd0\x25\x1b\x48\x39\x14\x62\x99\x33\x3a\x7a\x8f\x75\xdb\x1d\x74\x7c\xea\x1e\xa0\xeb\x75\xfd\xb3\x23\x4a\x30\x9d\x3a\xd8\xaf\xc3\xfd\x03\x0c\x80\xce\xc5\xaf\x6d\xb6\x17\xeb\xe8\x48\xaf\xac\xd1\xdd\x06\x57\x70\x40\xc7\x04\xf7\xdf\x4f\x00\x00\x00\xff\xff\x53\x20\xad\x9d\xd6\x01\x00\x00")

func testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpSeccompApprovedYamlBytes() ([]byte, error) {
	return bindataRead(
		_testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpSeccompApprovedYaml,
		"test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-seccomp-approved.yaml",
	)
}

func testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpSeccompApprovedYaml() (*asset, error) {
	bytes, err := testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpSeccompApprovedYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-seccomp-approved.yaml", size: 470, mode: os.FileMode(438), modTime: time.Unix(1598346737, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpSeccompUnapprovedYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x90\xb1\x6e\xf3\x30\x0c\x84\x77\x03\x7e\x07\xc2\xf3\x6f\xc7\xfe\xbb\x69\x2b\x8a\xa2\x6b\x96\x76\x09\x32\x30\x32\x13\xab\x91\x44\x41\x94\xd2\xe4\xed\x0b\xc5\xb5\x9b\xa1\x1a\x34\xdc\x51\x27\xde\x87\xc1\x7c\x50\x14\xc3\x5e\xc1\x65\xa8\xab\xb3\xf1\xa3\x82\x2d\x8f\x75\xe5\x28\xe1\x88\x09\x55\x5d\x01\x78\x74\xa4\x40\x48\xe7\x68\xd2\xad\xd5\xec\x13\x5d\x53\x3b\x92\xe3\x62\xa3\xf7\x9c\x30\x19\xf6\x72\x1f\x87\x32\xaa\xd9\x85\x6e\x79\xd2\xa1\x0d\x13\x76\xe7\x7c\xa0\xe8\x29\x91\x74\x86\x37\x81\x47\x05\x8d\x65\x8d\x76\x62\x49\x9b\xec\x31\x84\xc8\x17\x1a\xdb\x10\xf9\x68\x2c\x75\x9f\xc2\xbe\xa9\x2b\x09\xa4\x15\x40\xc9\x5e\x12\x5f\xe6\x1d\x7e\xfe\x8b\xd9\x3f\xcb\xbb\x50\x54\x30\xf4\x7d\x5f\xa4\x07\xe3\x2d\x72\x0e\x0a\x9e\xfa\xbe\x9f\xd5\xe3\x22\xfd\x5f\x25\xc9\x21\x58\x72\xe4\x13\xda\xbb\x29\x0a\x76\x30\xc0\x1e\x60\x49\x2b\xb5\xd1\x78\x8a\x73\xcb\xf6\x17\x4b\xab\xd3\x75\xa5\x01\x60\x1c\x9e\x48\xc1\x21\xcb\xed\xc0\xd7\x59\xd3\xec\x1c\x16\xba\x3b\x68\x64\x6a\xfe\x41\xd3\xea\x72\x8b\x25\x0a\x30\x4c\x0d\xec\x57\x74\x7f\x14\x04\x40\x6b\xf9\x6b\x1b\xcd\xc5\x58\x3a\xd1\xab\x68\xb4\x77\xe4\x0a\x8e\x68\x85\xe0\xe1\x7c\x07\x00\x00\xff\xff\x95\x72\x7e\x36\xd7\x01\x00\x00")

func testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpSeccompUnapprovedYamlBytes() ([]byte, error) {
	return bindataRead(
		_testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpSeccompUnapprovedYaml,
		"test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-seccomp-unapproved.yaml",
	)
}

func testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpSeccompUnapprovedYaml() (*asset, error) {
	bytes, err := testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpSeccompUnapprovedYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-seccomp-unapproved.yaml", size: 471, mode: os.FileMode(438), modTime: time.Unix(1598290256, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpSeccompUndefinedYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x90\xb1\xae\xdb\x30\x0c\x45\x77\x03\xfe\x07\xc2\x6f\xad\x1d\xbb\xdd\xb4\x15\x45\xd1\x35\x4b\xbb\x04\x19\x68\x89\x4e\x84\x48\xa4\x20\x4a\x69\xf2\xf7\x85\xe3\xa4\xc8\xf0\x34\x68\x38\x24\x0f\xc8\x8b\xc9\xff\xa1\xac\x5e\xd8\xc0\x75\x6a\x9b\x8b\x67\x67\x60\x2f\xae\x6d\x22\x15\x74\x58\xd0\xb4\x0d\x00\x63\x24\x03\x29\xcb\x9c\x7b\x25\x6b\x25\xa6\xbe\xb2\xa3\xc5\x33\xb9\xb5\xe1\x03\x99\xa5\x60\xf1\xc2\xfa\x98\x80\x0f\x78\x36\x0e\x4a\xb6\x66\x5f\xee\x03\x86\x74\xc6\xe1\x52\x67\xca\x4c\x85\x74\xf0\xb2\x4b\xe2\x0c\x74\xb9\x72\xf1\x91\x76\x8e\x16\xac\xa1\x74\x6d\xa3\x89\xac\x01\x58\x55\xaf\xf9\x1f\xc2\x85\x6e\x65\xd3\x43\xae\xfc\x5d\x7f\x2b\x65\x03\xd3\x38\x8e\x2b\x7a\x2b\xfc\xca\x52\x93\x81\x6f\xe3\x38\x6e\x74\x79\xa1\xaf\xff\x91\xd6\x94\x02\x45\xe2\x82\xe1\x51\x54\x03\x07\x98\xe0\x08\xf0\xb2\x59\xe1\x82\x9e\x29\x6f\x47\xf5\xcf\x20\x94\x6c\x6f\xcb\xad\x77\x14\x65\x73\xf9\x88\x27\x32\x30\x57\xbd\xcf\x72\xdb\x98\x95\x18\x71\xcd\xf3\x00\x9d\x9e\xbb\x2f\xd0\xf5\x76\xfd\x35\x10\x25\x98\xce\x1d\x1c\x9f\x8b\x7c\x7a\x20\x00\x86\x20\x7f\xf7\xd9\x5f\x7d\xa0\x13\xfd\x54\x8b\xe1\x91\xb0\x81\x05\x83\x12\xbc\xbd\x7f\x01\x00\x00\xff\xff\x10\xfc\x21\x5f\xc9\x01\x00\x00")

func testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpSeccompUndefinedYamlBytes() ([]byte, error) {
	return bindataRead(
		_testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpSeccompUndefinedYaml,
		"test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-seccomp-undefined.yaml",
	)
}

func testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpSeccompUndefinedYaml() (*asset, error) {
	bytes, err := testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpSeccompUndefinedYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-seccomp-undefined.yaml", size: 457, mode: os.FileMode(438), modTime: time.Unix(1598276904, 0)}
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
	"test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-seccomp-approved.yaml":       testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpSeccompApprovedYaml,
	"test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-seccomp-unapproved.yaml":     testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpSeccompUnapprovedYaml,
	"test/features/kubernetes/podsecuritypolicy/features/yaml/psp-azp-seccomp-undefined.yaml":      testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpSeccompUndefinedYaml,
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
							"psp-azp-seccomp-approved.yaml":       &bintree{testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpSeccompApprovedYaml, map[string]*bintree{}},
							"psp-azp-seccomp-unapproved.yaml":     &bintree{testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpSeccompUnapprovedYaml, map[string]*bintree{}},
							"psp-azp-seccomp-undefined.yaml":      &bintree{testFeaturesKubernetesPodsecuritypolicyFeaturesYamlPspAzpSeccompUndefinedYaml, map[string]*bintree{}},
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
