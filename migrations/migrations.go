// Code generated by go-bindata.
// sources:
// migrations/0001-CreateGamesTable.sql
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

var _migrations0001CreategamestableSql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x64\x8e\x31\x4f\xc3\x30\x10\x46\x77\xff\x8a\x4f\x5d\x9a\x48\x30\x80\xd4\xa5\x9d\x52\xea\x4a\x86\xe0\xa0\xc4\x11\xed\x14\x1d\xf5\x29\x78\x70\x88\x12\xa7\x0b\xe2\xbf\x23\x5c\x9a\xa5\xf3\xfb\xde\xbb\x7b\x2a\x65\x66\x24\xe4\xc1\x48\x5d\xa9\x42\x43\xed\xa1\x0b\x03\x79\x50\x95\xa9\xb0\x98\x26\x67\xef\xbf\xc6\xb1\x5f\x6c\x84\xf8\x1f\x9b\x6c\x9b\x4b\xb4\xe4\x79\x44\x22\x00\xc0\x59\xd4\xb5\xda\xe1\xad\x54\xaf\x59\x79\xc4\x8b\x3c\x62\x27\xf7\x59\x9d\x1b\xfc\x15\x9a\x96\x3b\x1e\x28\x70\x73\x7e\x48\xd2\xbb\xe8\x74\xe4\x19\x67\x1a\x4e\x9f\x34\x24\x8f\xab\x55\x1a\xef\xea\x3a\xcf\x2f\xdc\x73\x20\x4b\x81\xf0\x5c\x15\x7a\x3b\xc3\xb9\xbb\xfc\xfe\x59\xae\xd7\x11\x5e\x84\xd3\xc0\x14\xd8\x36\x14\x10\x9c\xe7\x31\x90\xef\x6f\x35\x5d\xbc\x5f\x3f\x98\x7a\x7b\x15\x3e\x5c\xeb\xba\x10\x97\x22\xdd\x88\xdf\x00\x00\x00\xff\xff\xc6\xd2\xf7\xe5\x16\x01\x00\x00")

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

	info := bindataFileInfo{name: "migrations/0001-CreateGamesTable.sql", size: 278, mode: os.FileMode(420), modTime: time.Unix(1487098708, 0)}
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

