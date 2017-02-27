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

var _migrations0001CreategamestableSql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x84\x8d\xb1\x4e\xc3\x30\x14\x45\x77\x7f\xc5\x55\x97\x26\x12\x2c\x48\x5d\xda\x29\x85\x57\x61\x70\x6d\x94\x38\x6a\xcb\x52\x3d\x6a\x0b\x22\xe1\x12\x25\x36\x0b\xe2\xdf\x11\x0d\x42\x0a\x0c\x9d\xef\x39\xf7\x5c\x97\x54\x58\x02\x6d\x2d\xe9\x4a\x1a\x0d\xb9\x82\x36\x16\xb4\x95\x95\xad\x30\x49\xa9\x71\x97\x6f\x7d\xdf\x4e\x16\x42\xfc\xc0\xb6\x58\x2a\xc2\x33\x07\xdf\x23\x13\x00\xd0\x38\xbc\x73\x77\x78\xe1\x2e\xbb\x9a\xcd\x72\x3c\x94\x72\x5d\x94\x3b\xdc\xd3\xee\xe2\x04\x1c\x39\xf8\x31\xf2\x1d\xd1\xb5\x52\xc3\x1e\x7c\x64\xc7\x91\x71\x57\x19\xbd\xfc\x1d\x71\x43\xab\xa2\x56\x16\xd3\x8f\xcf\xe9\x7c\x7e\x1a\x07\xe1\x29\x1d\xdd\xab\xdf\xff\x0d\x8f\x5f\x0f\x9d\xe7\xe8\xdd\x9e\x23\x62\x13\x7c\x1f\x39\xb4\xd8\x48\x7b\x0b\x2b\xd7\x84\x47\xa3\xe9\x7f\x4a\x9b\x4d\x96\x0f\x7e\x6a\xdd\x79\xbf\x56\x4a\xe4\x0b\xf1\x15\x00\x00\xff\xff\x2e\xa6\x18\x71\x49\x01\x00\x00")

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

	info := bindataFileInfo{name: "migrations/0001-CreateGamesTable.sql", size: 329, mode: os.FileMode(420), modTime: time.Unix(1487781584, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrations0002CreateoffertemplatestableSql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\xd0\xc1\x6f\xa2\x40\x14\x06\xf0\x3b\x7f\xc5\x8b\x17\x21\xd9\xbd\xec\xae\x7b\xd0\x13\xda\x67\x4a\x8b\x43\x0b\x43\xd4\x5e\xc8\xc8\x3c\x29\x09\x30\x74\x78\x98\x34\x4d\xff\xf7\x46\x62\xd3\x36\xc6\xc4\xf3\xfc\xbe\x6f\xbe\xbc\x45\x8c\xbe\x44\xc0\x8d\x44\x91\x04\x91\x80\x60\x09\x22\x92\x80\x9b\x20\x91\x09\x8c\xfa\xbe\xd4\xbf\x4d\xd7\xb5\xa3\x99\xe3\x9c\xb0\xf4\xe7\x21\x82\xd9\xef\xc9\x66\x4c\x75\x5b\x29\xa6\x0e\x5c\x07\xa0\xd4\x90\x3f\x2b\xeb\xfe\xfd\xef\xc1\x43\x1c\xac\xfc\x78\x0b\xf7\xb8\x85\x1b\x5c\xfa\x69\x28\xe1\xd8\x96\x15\xd4\x90\x55\x4c\xd9\xe1\x9f\xeb\xfd\x72\x00\x1a\x55\x13\x1c\x94\x1d\xa2\x7f\x26\x13\x0f\x52\x11\x3c\xa6\x38\x0c\x11\x69\x18\x1e\x51\x6b\x8d\xee\x73\xce\x4a\xfd\x93\x7e\x37\x85\xaa\xe9\x22\x80\x18\x97\x18\xa3\x58\x60\x32\xc0\xce\x2d\xf5\xf0\x7d\x6e\x1a\xa6\x86\x3b\xb8\x4b\x22\x31\xff\xf2\x9f\xa3\xc7\x6f\xef\x63\x98\x4e\x87\xd7\xa3\xaf\x89\x95\x56\xac\x4e\xfe\x12\x6b\xc9\x96\x46\x5f\x5b\xba\xb7\xf4\xd2\x53\x93\xbf\x5e\x1b\x60\x5b\x16\x05\xd9\x6b\x39\x35\x6a\x57\x91\x86\x9d\x31\xd5\xb9\x66\xdb\xd3\xb0\xb9\x52\x39\xd5\xd4\xf0\xe5\x1b\xe7\x96\x14\x93\xce\x14\x03\x97\x35\x75\xac\xea\x16\xd6\x81\xbc\x05\x19\xac\x10\x9e\x22\x81\xe7\xfd\x22\x5a\xbb\x9e\xe3\xcd\x9c\x8f\x00\x00\x00\xff\xff\x18\xae\xf2\xad\x6f\x02\x00\x00")

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

	info := bindataFileInfo{name: "migrations/0002-CreateOfferTemplatesTable.sql", size: 623, mode: os.FileMode(420), modTime: time.Unix(1488149873, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrations0003CreateofferstableSql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\x92\x4f\x6f\xaa\x40\x14\xc5\xf7\x7c\x8a\x1b\x57\x90\xf8\x12\xde\x7b\xb5\x9b\xae\xa8\x5e\xd3\x49\x71\x68\x61\x88\xda\xcd\x64\x0a\x57\x24\xe1\x5f\x86\xc1\xa4\xdf\xbe\x11\xff\x60\x63\x13\xdd\xff\xce\x39\x73\xe6\xdc\x69\x88\x9e\x40\xc0\x95\x40\x1e\xb1\x80\x03\x9b\x03\x0f\x04\xe0\x8a\x45\x22\x82\x51\xd7\xe5\xe9\x9f\xba\x6d\x9b\xd1\x93\x65\x1d\x61\xe1\x3d\xfb\x08\xf5\x66\x43\xba\x05\xdb\x02\x00\xc8\x53\x48\xb6\x4a\xdb\xff\x1f\x1d\x78\x0b\xd9\xc2\x0b\xd7\xf0\x8a\x6b\x98\xe1\xdc\x8b\x7d\x01\x7b\x1b\x99\x51\x45\x5a\x19\x92\xbb\x07\xdb\x19\xf7\xba\x4c\x95\x24\xf3\x14\x76\x4a\xf7\xfa\x7f\x93\x89\xd3\xe7\xf3\xd8\xf7\x21\xc4\x39\x86\xc8\xa7\x18\xf5\x60\x6b\xe7\xe9\x51\xd7\x87\x4b\x43\x65\x53\xec\x0d\x2f\xe3\x7f\x53\xff\xc4\x2f\x7c\x9a\x42\x7d\x91\xbe\x7c\xc1\x5f\xd7\x75\x07\x93\x03\xd5\x12\x55\x32\xa9\xbb\xca\x90\x86\xbc\x32\x94\x91\x1e\x72\x4e\x1d\xdd\x03\xfc\x59\x77\xd9\xd6\xdc\x8d\x27\x9a\x94\xa1\x54\x2a\x03\x26\x2f\xa9\x35\xaa\x6c\x60\xc9\xc4\x0b\x08\xb6\x40\xf8\x08\x38\x5e\x6b\x79\xb0\x3c\xfd\x60\xd7\xa4\xb7\xf5\xe7\x26\x49\xa1\xf2\xf2\x6e\xba\x50\xad\x91\x7d\xf9\x5b\xbc\xe5\x0c\xe7\x11\x73\xf6\x1e\x23\x30\x3e\xc3\xd5\x69\x60\x79\xfe\x68\x79\x35\x9d\x15\xf0\xf3\x31\x1d\xf1\xf1\x30\xcc\xf8\x7a\x6b\xc7\xfa\x0e\x00\x00\xff\xff\xdf\x6f\x19\x8d\xb6\x02\x00\x00")

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

	info := bindataFileInfo{name: "migrations/0003-CreateOffersTable.sql", size: 694, mode: os.FileMode(420), modTime: time.Unix(1487876695, 0)}
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

