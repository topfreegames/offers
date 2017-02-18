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

var _migrations0002CreateoffertemplatestableSql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\xcd\xc1\x4b\xc3\x30\x1c\xc5\xf1\x7b\xfe\x8a\xc7\x2e\xdb\x40\x2f\xc2\x2e\xdb\xa9\xd3\x0c\xaa\x35\x95\x36\x83\xed\x34\xb2\xe6\xd7\x1a\x68\x93\x9a\x26\x82\x88\xff\xbb\x74\x88\x17\x71\xf4\xfc\xfd\x3c\xde\x7d\xc1\x13\xc9\xc1\x0f\x92\x8b\x32\xcd\x05\xd2\x1d\x44\x2e\xc1\x0f\x69\x29\x4b\xcc\x62\x34\xfa\xd6\x0d\x43\x3f\xdb\x30\xf6\x83\x65\xb2\xcd\x38\x5c\x5d\x93\x3f\x05\xea\xfa\x56\x05\x1a\xb0\x60\x80\xd1\x18\x07\x78\x29\xd2\xe7\xa4\x38\xe2\x89\x1f\x6f\x18\x60\x55\x47\x78\x57\xbe\x7a\x55\x7e\x71\xb7\x5a\x2d\x2f\x17\x62\x9f\x65\x63\xed\xbd\xd3\xb1\x0a\x27\xa3\xff\x37\x8d\xea\xe8\x2a\xa8\x9c\x0d\x64\xc3\x80\xc7\x32\x17\xdb\xdf\x84\x07\xbe\x4b\xf6\x99\xc4\xfc\xf3\x6b\x8e\xf5\xfa\x52\x47\xdf\x51\x50\x5a\x05\x35\xd5\xf7\xe4\x8d\xd3\x53\x75\xed\xe9\x2d\x92\xad\x3e\xa6\x0e\x82\x37\x4d\x43\x7e\x2a\x27\xab\xce\x2d\x69\x9c\x9d\x6b\xff\xea\xe0\x23\xb1\xe5\x86\x7d\x07\x00\x00\xff\xff\x57\x87\x8c\x06\xdc\x01\x00\x00")

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

	info := bindataFileInfo{name: "migrations/0002-CreateOfferTemplatesTable.sql", size: 476, mode: os.FileMode(420), modTime: time.Unix(1487364390, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrations0003CreateofferstableSql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x6c\x8f\x41\x4b\xc3\x40\x14\x84\xef\xf9\x15\x73\x4c\xc0\x43\x14\x7b\xf2\x14\xeb\x2b\x88\x31\x95\x75\x8b\xf4\xb4\x3c\xba\xaf\x75\x21\xab\xcb\x66\x53\xf0\xdf\x8b\xab\x21\x15\x72\x9e\xf9\xe6\x63\xd6\x8a\x1a\x4d\xd0\xcd\x7d\x4b\xf8\x3c\x1e\x25\x0e\x28\x0b\x00\x70\x16\xe3\xe8\x2c\x5e\xd4\xe3\x73\xa3\xf6\x78\xa2\x3d\x1e\x68\xd3\xec\x5a\x9d\x03\x73\x92\x0f\x89\x9c\xc4\x9c\x6f\xcb\xea\x2a\x33\x27\xf6\x62\x9c\xc5\x99\xe3\xe1\x9d\x63\x79\xb3\x5a\x55\xe8\xb6\x1a\xdd\xae\x6d\xa1\x68\x43\x8a\xba\x35\xbd\xe6\xe2\x50\x3a\xfb\xc7\x65\xb1\x49\xe2\x43\xff\x33\x38\xa9\x97\xc8\xff\xd5\x8b\x8d\xd0\xf3\x97\xc4\x4b\xfb\x75\x5d\xd7\xb3\xfe\xb7\x75\x88\xc2\x49\xac\xe1\x84\xe4\xbc\x0c\x89\x7d\x98\x3d\xd3\xbf\x6e\xfb\x36\x5d\x1a\x83\x5d\x00\xe6\xbd\x9e\x9d\x5f\x8a\x8b\xea\xae\xf8\x0e\x00\x00\xff\xff\x81\xe4\x6c\x3c\x5d\x01\x00\x00")

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

	info := bindataFileInfo{name: "migrations/0003-CreateOffersTable.sql", size: 349, mode: os.FileMode(420), modTime: time.Unix(1487344379, 0)}
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

