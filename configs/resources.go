package configs

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
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
		return nil, fmt.Errorf("read %q: %w", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
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

var _configsResourcesXml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x58\xcd\x6e\xdb\x46\x17\xdd\xfb\x29\xf8\x6d\x05\xd8\x7c\x01\x41\x8b\xaf\x36\x8a\x00\x0d\x5a\xb4\xee\x5a\xa0\xc8\x49\xcc\x54\x24\x95\x19\x52\x69\x76\x2e\xfc\x17\xb5\x92\xad\x36\x76\x6d\xa8\x0a\x6c\xa5\x16\x2c\xa4\x88\x68\xe7\xaf\xac\x28\xc5\x2f\x33\x33\xa4\x57\x7a\x85\x82\x1c\xd9\x92\xc8\xa1\x24\xab\x8d\xd1\xe5\x70\xee\xcc\xdc\x73\xee\xf5\x3d\x47\x4e\x7f\x0d\x90\x61\x41\x19\x64\x16\x04\x21\xfd\xbf\xc5\xc5\x05\x41\x10\x88\xfd\x17\xee\xec\x92\xd3\x1a\x76\x5e\x61\x67\xdd\xeb\x9c\xf9\x7b\x75\x52\x39\xb8\x3a\x7c\x27\xe0\xde\xa5\xb7\xdf\xc2\x9d\x6d\x52\x2d\x63\x67\x1d\x3b\xaf\x2c\x98\xef\x77\xcb\xb8\xb7\x4d\x8f\xec\x7b\xcb\xd8\xa9\x90\x6a\xd9\xab\x6d\xd2\x8d\x2d\xb2\xfd\xbe\xdf\x2d\x93\xea\xae\xf7\x87\x4d\x8f\x6c\x76\xd4\x82\x79\xfa\xa2\x49\x8f\x9b\xb4\xf4\x51\xf0\xcf\x7e\xf1\x4b\x6f\xee\x2d\xf7\xbb\xe5\xe0\x61\xdf\x6e\x7a\xfb\x2d\xfa\xec\x4f\x5a\x2f\xf9\x76\x93\xed\x92\xf3\x0b\x7a\x64\xfb\xed\x53\xec\x74\xc8\x9e\x8d\xdd\xa6\xdf\xbe\xbc\x3a\x6c\xf7\xbb\x65\xfa\xa6\x41\xeb\x25\x72\x7e\xc1\x76\xb1\x53\xf1\x37\x7a\x6c\x77\x41\x10\x16\x17\xaf\x41\x09\xa3\x28\xfa\xdd\x1a\x7d\xf6\xb3\x57\xdb\xbc\x0e\x80\x23\x14\x08\x42\x5a\x55\x32\x9a\xaa\x83\xb4\xa8\x2a\x83\x2f\xba\xa4\x81\x0c\x3b\x93\x16\xc3\x05\xfb\xae\x01\x73\xcd\x50\x32\x5f\x7d\xf9\xcd\x6a\x5a\x1c\x2c\xc6\x77\xbe\x8d\x6e\x40\xc3\x32\x41\x46\x2c\xa6\x44\x4b\x57\x8b\x00\x22\x29\x2f\xa6\x52\x69\x91\x7d\x8f\xc5\x20\x00\x45\x04\xf2\x0f\xe2\x31\x63\x39\x0f\xb3\x5e\x52\x75\x64\x42\x4b\x36\x55\x43\x1f\x22\xb8\xc1\x50\xde\xc1\xee\xe9\x28\x86\x49\x28\x26\xe0\x18\x6e\x2d\xaf\x7c\xb1\xb2\xba\x12\xdb\x1d\x62\x18\x49\x28\x86\x42\x8c\x50\x9f\x88\x4a\x53\x11\xe2\x22\xc2\xae\x4b\x7e\x6c\xdc\x05\xa2\xc1\xf2\xf3\x95\xf8\xc1\x21\xd8\x41\x9e\x51\xa0\x63\x31\x0f\xd4\x3c\x10\xad\x42\xde\x90\x94\xf9\x09\x31\x81\x56\xc8\x4b\x26\x48\x62\x84\xb6\x1a\xf4\xc5\xe5\xdd\x56\x7a\x00\x3e\x7b\x9d\xdb\xfc\xe8\x18\x3b\x01\x51\x9c\x1e\xfe\x75\x07\xbb\x1f\x70\xf7\xc4\xff\xf8\xfc\x3f\x53\xf7\x20\xd5\xac\x09\x25\x1d\x15\x0c\x68\xce\x0f\x5c\x91\x4c\x29\x27\x21\x1e\xec\x70\xfc\xd0\x83\x73\x5a\x69\x93\xce\x9d\x20\x1f\xc2\x43\x8f\xf3\x22\xf8\x1e\xc8\x13\xfb\x3a\x08\xd2\x8d\x78\xd9\x63\x41\x4f\xa0\xca\x8f\x8a\xb0\xc2\xe1\x65\x49\x32\xb3\x05\x68\x28\x96\x6c\x8e\x52\x34\xec\xfd\x2d\x52\x6f\x79\xfb\xc7\xb8\x73\xe6\xed\xda\xe4\xe5\x06\x2d\x9d\xf9\x8d\x72\x84\x2d\x31\xc6\xff\x8c\x4f\x9b\x00\x71\xdf\x25\x7b\x36\xa9\xb7\xb0\xeb\xe2\xde\xc1\xec\xef\xce\xde\x16\xb2\xa1\x69\x40\x37\xe3\x5d\x71\xb5\x5e\x0b\x04\xf3\x79\x05\xf7\xea\xfc\x96\x48\xee\xdb\x89\x4d\x3a\xba\xe2\x68\x64\xa0\xa9\x10\x14\x55\xf0\x44\x20\x5b\x6f\xaf\xf6\xdb\xe4\xf5\x61\xc8\xc0\xef\xde\xc9\x0f\xd8\x79\x4d\xeb\x25\xa6\xf6\x81\x33\x08\x55\x99\xfe\xf6\x8e\x56\x5e\x0e\x24\x97\x41\x63\x17\x44\x75\x96\xb4\x1b\xf4\xc4\xf9\xf7\x74\x96\xbd\x32\xde\x6e\x51\x74\x1c\x8b\xe0\xbd\x75\x3d\xf7\xd8\x6b\x37\xbc\xea\x76\xb2\x51\x40\x4f\xd1\x7d\x49\x97\x1e\xc6\xdc\xc2\xe8\xf1\x31\x2c\xbc\x1a\xdf\xdc\xb2\x14\x08\x7e\xbc\xca\xcc\x13\xc5\x2f\x1b\x12\x90\x9a\xf0\xc7\x1b\x9a\x08\x2d\xbc\x1f\xce\x35\x99\x86\xe9\x41\x83\x37\x91\x99\x4d\x9b\x37\x3d\x79\x4d\x82\x92\x6c\x4e\xcb\x6d\x4a\x43\x06\x69\x86\xeb\x68\x21\xfc\xf7\x9b\xb4\x53\x9d\xb5\x10\xe1\x97\x25\xf4\x14\x71\x8a\x10\x56\x34\x61\xee\x26\xcf\xd6\x49\x23\x79\xea\xdc\x2d\x40\xe3\x11\x90\xe7\xd3\x93\x6b\x2c\x00\x16\x55\x99\x27\x28\xf5\x4a\xb2\x73\xfa\x44\x78\x90\x95\xcb\x26\x60\x4a\x8a\xcb\x2a\x39\x34\x39\x96\x01\xcc\x9a\xaa\x36\x9f\xe3\x18\x10\xb5\x66\x20\xce\x80\xc5\x8e\x4b\xeb\x9d\xbb\x65\x49\x93\xe4\x35\x55\xff\x47\x60\x26\xf8\x88\x29\x0e\xe2\x53\x55\xfe\x71\x7e\xb6\x19\x34\xc3\x68\x26\xed\x86\xdf\x6e\x78\x17\xee\xb4\xe9\x2c\x59\x8a\x6a\x72\xf4\x25\x38\x3d\xdb\x40\x00\x45\xae\xec\xe2\xce\x4f\xd8\xfd\xc0\x6e\xba\xfd\xc4\x43\xd9\xc2\xa3\x6c\x78\xf3\x5c\x25\x36\x0a\x00\x72\x9d\x3f\xb3\x01\xf1\xac\x6e\x4d\x30\x2d\xef\xd0\x93\x1d\xf6\x43\x22\x99\x5d\x53\x82\x0f\x81\xb9\x2a\xa1\xef\xa2\x14\x0f\x7e\x82\x84\xb7\xf0\x84\x3c\x95\x28\xd6\x9a\xa1\x6b\xaa\x72\x6b\xb1\xfe\xec\xfe\xf2\xff\x93\x13\x95\x35\x25\x17\x4d\x31\x38\xc1\xf5\x18\x31\x27\x31\xd5\x7c\xf0\x7a\x7f\x54\xdf\x34\x25\x97\x0c\x28\x2d\xde\xfc\x03\xe6\xef\x00\x00\x00\xff\xff\xed\xfe\x6c\x8b\x8b\x11\x00\x00")

func configsResourcesXmlBytes() ([]byte, error) {
	return bindataRead(
		_configsResourcesXml,
		"configs/resources.xml",
	)
}

func configsResourcesXml() (*asset, error) {
	bytes, err := configsResourcesXmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "configs/resources.xml", size: 4491, mode: os.FileMode(0644), modTime: time.Unix(1632736751, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x4e, 0x2c, 0xdb, 0xa4, 0xf2, 0x23, 0x94, 0xd5, 0xe0, 0x70, 0x1, 0x99, 0x22, 0x84, 0x27, 0x6f, 0x84, 0x5a, 0x4e, 0xcd, 0xdc, 0x9d, 0xee, 0x51, 0xd1, 0xd8, 0x5, 0xa1, 0x47, 0xa3, 0x3c, 0x42}}
	return a, nil
}

var _configsEmptyJsonExampleJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func configsEmptyJsonExampleJsonBytes() ([]byte, error) {
	return bindataRead(
		_configsEmptyJsonExampleJson,
		"configs/empty-json-example.json",
	)
}

func configsEmptyJsonExampleJson() (*asset, error) {
	bytes, err := configsEmptyJsonExampleJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "configs/empty-json-example.json", size: 0, mode: os.FileMode(0644), modTime: time.Unix(1635147716, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xe3, 0xb0, 0xc4, 0x42, 0x98, 0xfc, 0x1c, 0x14, 0x9a, 0xfb, 0xf4, 0xc8, 0x99, 0x6f, 0xb9, 0x24, 0x27, 0xae, 0x41, 0xe4, 0x64, 0x9b, 0x93, 0x4c, 0xa4, 0x95, 0x99, 0x1b, 0x78, 0x52, 0xb8, 0x55}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
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

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
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
	"configs/resources.xml":           configsResourcesXml,
	"configs/empty-json-example.json": configsEmptyJsonExampleJson,
}

// AssetDebug is true if the assets were built with the debug flag enabled.
const AssetDebug = false

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
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
	"configs": {nil, map[string]*bintree{
		"configs/empty-json-example.json": {configsEmptyJsonExampleJson, map[string]*bintree{}},
		"resources.xml":                   {configsResourcesXml, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory.
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
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
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
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
