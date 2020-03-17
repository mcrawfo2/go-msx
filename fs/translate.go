package fs

import (
	"github.com/pkg/errors"
	"net/http"
	"path"
)

type TranslateFs struct {
	fs   http.FileSystem
	root string
}

func (t TranslateFs) Open(name string) (http.File, error) {
	return t.fs.Open(path.Clean(path.Join(t.root, name)))
}

func NewTranslateFs(fs http.FileSystem, root string) (http.FileSystem, error) {
	if !path.IsAbs(root) {
		return nil, errors.New("Expected absolute root path")
	}

	return TranslateFs{
		fs:   fs,
		root: root,
	}, nil
}
