// Code generated by go-bindata.
// sources:
// migrations/0001-CreateGamesTable.sql
// migrations/0002-CreateOfferTemplatesTable.sql
// migrations/0003-CreateOffersTable.sql
// DO NOT EDIT!

package migrations

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

var _migrations0001CreategamestableSql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x6c\x8d\x41\x4b\xc3\x30\x18\x86\xef\xf9\x15\x2f\xbb\xac\x05\xbd\x08\xbb\x6c\xa7\x4e\xbf\x41\x35\x26\xd2\xa4\xb8\x9d\xc6\xe7\x12\xb4\x60\x66\x69\x12\x2f\xe2\x7f\x17\x57\x11\xaa\x3b\x3f\xcf\xf3\xbe\xd7\x0d\x55\x96\x40\x5b\x4b\xca\xd4\x5a\xa1\xde\x40\x69\x0b\xda\xd6\xc6\x1a\xcc\x72\xee\xdc\xe5\x5b\x8c\xfd\x6c\x25\xc4\x8f\x6c\xab\xb5\x24\x3c\x73\xf0\x11\x85\x00\x80\xce\xe1\x9d\x87\xc3\x0b\x0f\xc5\xd5\x62\x51\xe2\xa1\xa9\xef\xab\x66\x87\x3b\xda\x5d\x9c\x84\x23\x07\x3f\x55\xbe\x4f\x54\x2b\xe5\xc8\x83\x4f\xec\x38\x31\x6e\x8d\x56\xeb\x5f\x88\x1b\xda\x54\xad\xb4\x98\x7f\x7c\xce\x97\xcb\x13\x1c\x83\xa7\x7c\x74\xaf\x7e\xff\xf7\x78\xba\x7a\x18\x3c\x27\xef\xf6\x9c\x90\xba\xe0\x63\xe2\xd0\xff\xdf\x56\xfa\xb1\x28\xc7\x20\xf7\xee\x4c\xd0\x4a\x29\xca\x95\xf8\x0a\x00\x00\xff\xff\x4f\x78\x61\xe4\x2b\x01\x00\x00")

func migrations0001CreategamestableSqlBytes() ([]byte, error) {
	return bindataRead(
		_migrations0001CreategamestableSql,
		"migrations/0001-CreateGamesTable.sql",
	)
}

func migrations0001CreategamestableSql() (*asset, error) {
	bytes, err := migrations0001CreategamestableSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/0001-CreateGamesTable.sql", size: 299, mode: os.FileMode(420), modTime: time.Unix(1487285165, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrations0002CreateoffertemplatestableSql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\xcd\x4d\x4b\xc3\x40\x10\xc6\xf1\xfb\x7e\x8a\x87\x5e\xda\x82\x5e\x84\x5e\xda\x53\xaa\x5b\x88\xc6\x44\x92\x2d\xb4\xa7\xb2\xcd\x4e\x62\x20\xd9\x8d\x9b\x59\x41\xc4\xef\x2e\x11\x11\x5f\x40\x72\xfe\xff\x66\x9e\xeb\x5c\x46\x4a\x42\x1e\x94\x4c\x8b\x38\x4b\x11\xef\x90\x66\x0a\xf2\x10\x17\xaa\xc0\x2c\x84\xc6\x5c\xba\x61\xe8\x67\x1b\x21\x3e\xb1\x8a\xb6\x89\x84\xab\x2a\xf2\x27\xa6\xae\x6f\x35\xd3\x80\x85\x00\x1a\x83\x67\xed\xcb\x47\xed\x17\x57\xab\xd5\x12\x0f\x79\x7c\x1f\xe5\x47\xdc\xc9\xe3\x85\x00\xac\xee\xe8\x27\x18\xa7\xd2\x7d\x92\x8c\xb5\xf7\xce\x84\x92\x4f\xbf\x9f\x7c\x37\xb5\xee\xe8\x5f\x50\x3a\xcb\x64\x79\xc0\x6d\x91\xa5\xdb\xaf\x84\x1b\xb9\x8b\xf6\x89\xc2\xfc\xf5\x6d\x8e\xf5\xfa\xa3\x8e\xbe\x23\xd6\x46\xb3\x9e\xea\x7b\xf2\x8d\x33\x53\x75\xe5\xe9\x29\x90\x2d\x5f\xa6\x1e\xb0\x6f\xea\x9a\xfc\x54\x4e\x56\x9f\x5b\x32\x38\x3b\xd7\xfe\xd5\xec\x03\x89\xe5\x46\xbc\x07\x00\x00\xff\xff\x88\x5b\x9a\x79\xe4\x01\x00\x00")

func migrations0002CreateoffertemplatestableSqlBytes() ([]byte, error) {
	return bindataRead(
		_migrations0002CreateoffertemplatestableSql,
		"migrations/0002-CreateOfferTemplatesTable.sql",
	)
}

func migrations0002CreateoffertemplatestableSql() (*asset, error) {
	bytes, err := migrations0002CreateoffertemplatestableSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/0002-CreateOfferTemplatesTable.sql", size: 484, mode: os.FileMode(420), modTime: time.Unix(1487442237, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrations0003CreateofferstableSql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\x8f\xc1\x4b\xc3\x30\x14\x87\xef\xfd\x2b\x1e\x3b\xb5\xa0\x50\xc5\x9d\x76\xaa\xf3\x15\x8a\x35\x95\x34\xc3\xed\x14\xc2\xf2\x36\x03\x8d\x0b\x49\x3a\xf0\xbf\x97\xc5\x8d\x4e\xd8\xc1\xf3\xfb\xbe\xef\xc7\x5b\x72\xac\x04\x02\xae\x05\xb2\xbe\xe9\x18\x34\x35\xb0\x4e\x00\xae\x9b\x5e\xf4\x30\x1b\x47\xa3\xef\x0f\x21\xb8\xd9\x22\xcb\xce\xb0\xa8\x9e\x5b\x84\xc3\x6e\x47\x3e\x40\x9e\x01\x00\x18\x0d\x27\x12\xde\x79\xf3\x56\xf1\x0d\xbc\xe2\x06\x5e\xb0\xae\x56\xad\x48\x07\xb9\xa7\x2f\xf2\x2a\x92\x3c\x3e\xe5\xc5\x5d\x72\xf6\xca\x92\x34\x1a\x8e\xca\x6f\x3f\x95\xcf\x1f\xe7\xf3\x22\x6d\xb3\x55\xdb\x02\xc7\x1a\x39\xb2\x25\xf6\x09\x0c\xb9\xd1\x67\x2f\x0d\xcb\x48\xd6\x0d\xa7\xe0\x7f\x0a\x7f\x95\xab\x96\x1b\xd4\x37\xf9\xeb\xc6\x43\x59\x96\x53\xe4\x97\xda\x7a\x52\x91\xb4\x54\x11\xa2\xb1\x14\xa2\xb2\x6e\xda\xb9\xfc\xc9\xba\x8f\xcb\x6b\xa3\xd3\x37\x84\xa9\x37\x28\x63\x6f\x9d\xb3\x62\x91\xfd\x04\x00\x00\xff\xff\x6e\xd9\x76\x0e\x92\x01\x00\x00")

func migrations0003CreateofferstableSqlBytes() ([]byte, error) {
	return bindataRead(
		_migrations0003CreateofferstableSql,
		"migrations/0003-CreateOffersTable.sql",
	)
}

func migrations0003CreateofferstableSql() (*asset, error) {
	bytes, err := migrations0003CreateofferstableSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/0003-CreateOffersTable.sql", size: 402, mode: os.FileMode(420), modTime: time.Unix(1487442580, 0)}
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
	"migrations/0001-CreateGamesTable.sql": migrations0001CreategamestableSql,
	"migrations/0002-CreateOfferTemplatesTable.sql": migrations0002CreateoffertemplatestableSql,
	"migrations/0003-CreateOffersTable.sql": migrations0003CreateofferstableSql,
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
	"migrations": &bintree{nil, map[string]*bintree{
		"0001-CreateGamesTable.sql": &bintree{migrations0001CreategamestableSql, map[string]*bintree{}},
		"0002-CreateOfferTemplatesTable.sql": &bintree{migrations0002CreateoffertemplatestableSql, map[string]*bintree{}},
		"0003-CreateOffersTable.sql": &bintree{migrations0003CreateofferstableSql, map[string]*bintree{}},
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

