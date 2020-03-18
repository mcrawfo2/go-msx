package fs

import (
	"github.com/pkg/errors"
	"net/http"
	"path"
)

type PrefixFileSystem struct {
	fs   http.FileSystem
	root string
}

func (t PrefixFileSystem) Open(name string) (http.File, error) {
	prefixedName := path.Clean(path.Join(t.root, name))
	return t.fs.Open(prefixedName)
}

func NewPrefixFileSystem(fs http.FileSystem, root string) (http.FileSystem, error) {
	if !path.IsAbs(root) {
		return nil, errors.New("Expected absolute root path")
	}

	return PrefixFileSystem{
		fs:   fs,
		root: root,
	}, nil
}
