// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package fs

import (
	"github.com/pkg/errors"
	"net/http"
	"path"
	"strings"
)

var ErrInvalidFileName = errors.New("Invalid filename")

type PrefixFileSystem struct {
	fs   http.FileSystem
	root string
}

func (t PrefixFileSystem) Open(name string) (http.File, error) {
	prefixedName := path.Clean(path.Join(t.root, name))
	if !strings.HasPrefix(prefixedName, t.root+"/") && prefixedName != t.root {
		return nil, errors.Wrap(ErrInvalidFileName, name)
	}
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
