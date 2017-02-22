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

	info := bindataFileInfo{name: "migrations/0001-CreateGamesTable.sql", size: 329, mode: os.FileMode(420), modTime: time.Unix(1487785213, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrations0002CreateoffertemplatestableSql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\xd0\x3f\x6b\xe3\x40\x10\x05\xf0\x5e\x9f\x62\x70\x63\x09\xee\x9a\xbb\xf3\x15\x76\x25\x3b\x63\xa2\x44\x5e\x05\x69\x8d\xed\x34\x62\xac\x1d\x2b\x02\xfd\xcb\x6a\x64\x08\x21\xdf\x3d\x58\x24\xa4\x30\x06\xd5\xf3\x7b\x6f\x1f\xbb\x8a\xd1\xd7\x08\xb8\xd7\xa8\x92\x20\x52\x10\xac\x41\x45\x1a\x70\x1f\x24\x3a\x81\x49\xdf\x17\xe6\x77\xd3\x75\xed\x64\xe1\x38\x5f\x58\xfb\xcb\x10\xa1\x39\x9d\xd8\xa6\xc2\x55\x5b\x92\x70\x07\xae\x03\x50\x18\xc8\x5e\xc8\xba\x7f\xff\x7b\xf0\x14\x07\x1b\x3f\x3e\xc0\x23\x1e\xe0\x0e\xd7\xfe\x36\xd4\x70\x69\x4b\x73\xae\xd9\x92\x70\x7a\xfe\xe7\x7a\xbf\x1c\x80\x9a\x2a\x86\x33\xd9\x21\xfa\x67\x36\xf3\x86\x05\x6a\x1b\x86\x97\x6b\x6b\x1b\xd3\x67\x92\x16\xe6\xb6\xc9\xa9\xe2\x9b\x00\x62\x5c\x63\x8c\x6a\x85\xc9\x00\x3b\xb7\x30\xc3\xbb\x59\x53\x0b\xd7\xd2\xc1\x43\x12\xa9\xe5\x8f\xff\x5e\x3b\x7d\xff\x98\xc2\x7c\x3e\x5c\x2f\xbe\x62\x21\x43\x42\x63\x7d\xcb\xb6\x68\xcc\x58\x7d\xb2\xfc\xda\x73\x9d\xbd\x8d\x0d\x88\x2d\xf2\x9c\xed\x58\xce\x35\x1d\x4b\x36\x70\x6c\x9a\xf2\x5a\x8b\xed\x79\xd8\x5c\x52\xc6\x15\xd7\x72\xfb\xb3\x33\xcb\x24\x6c\x52\x12\x90\xa2\xe2\x4e\xa8\x6a\x61\x17\xe8\x7b\xd0\xc1\x06\xe1\x39\x52\x78\xdd\xaf\xa2\x9d\xeb\x39\xde\xc2\xf9\x0c\x00\x00\xff\xff\xd2\x36\xac\x14\x71\x02\x00\x00")

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

	info := bindataFileInfo{name: "migrations/0002-CreateOfferTemplatesTable.sql", size: 625, mode: os.FileMode(420), modTime: time.Unix(1487785213, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrations0003CreateofferstableSql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\x92\xcb\x6a\xeb\x30\x10\x86\xf7\x7e\x8a\x21\x2b\x1b\xce\x01\xf7\x92\x6e\xba\x72\x93\x09\x35\x75\xe4\x62\x2b\x24\xe9\x46\xa8\xf6\xc4\x11\xf8\x86\x24\x07\xfa\xf6\x25\xce\xc5\x29\x2d\x24\x5b\xf1\x7d\xff\xaf\x19\x66\x92\x60\xc0\x11\x70\xc5\x91\xa5\x61\xcc\x20\x9c\x01\x8b\x39\xe0\x2a\x4c\x79\x0a\xa3\xae\x53\xf9\xff\xc6\x98\x76\xf4\xec\x38\x47\x98\x07\x2f\x11\x42\xb3\xd9\x90\x36\xe0\x3a\x00\x00\x2a\x87\x6c\x2b\xb5\xfb\xf0\xe4\xc1\x7b\x12\xce\x83\x64\x0d\x6f\xb8\x86\x29\xce\x82\x45\xc4\x61\x1f\x23\x0a\xaa\x49\x4b\x4b\x62\xf7\xe8\x7a\xff\x7a\xaf\x90\x15\x09\x95\xc3\x4e\xea\xde\xbf\x1f\x8f\xbd\xbe\x9f\x2d\xa2\x08\x12\x9c\x61\x82\x6c\x82\x69\x0f\x1a\x57\xe5\x47\xaf\x2f\x17\x96\xaa\xb6\xdc\x07\x5e\xd6\xff\x65\xff\xc4\x2f\x72\xda\x52\x7e\x91\xbe\xfc\xc1\x9d\xef\xfb\x43\xc8\x81\x32\x44\xb5\xc8\x9a\xae\xb6\xa4\x41\xd5\x96\x0a\xd2\x43\xcf\x69\x46\xff\x00\x7f\x36\x5d\xb1\xb5\x37\xe3\x99\x26\x69\x29\x17\xd2\x82\x55\x15\x19\x2b\xab\x16\x96\x21\x7f\x05\x1e\xce\x11\x3e\x62\x86\xbf\x5d\x16\x2f\x4f\x1b\xec\xda\xfc\xba\x7f\x9e\x24\x2b\xa5\xaa\x6e\xa6\x4b\x69\xac\xe8\x87\xbf\xc6\x3b\xde\x70\x1e\x21\x9b\xe2\x6a\xd8\xac\x13\xb3\xf3\xad\x9c\x1f\x3d\xe7\x3b\x00\x00\xff\xff\xac\x6d\x31\xc6\x79\x02\x00\x00")

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

	info := bindataFileInfo{name: "migrations/0003-CreateOffersTable.sql", size: 633, mode: os.FileMode(420), modTime: time.Unix(1487785213, 0)}
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

