// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

// This file loads static files from an embed.FS and returns them in the structure
// provided by bou.ke/staticfiles (staticFiles) and expected by the rest of skel

package skel

import (
	"crypto/sha256"
	"embed"
	"fmt"
	"io/fs"
	"mime"
	"strings"
	"time"

	"github.com/pkg/errors"
)

//go:embed all:_templates/*
var statics embed.FS

const (
	staticRoot = "_templates"
)

var StaticFileReadError = errors.New("static file read error")

type staticFilesFile struct {
	data  string
	mime  string
	mtime time.Time
	size  int    // size is the size before compression. If 0, it means the data is uncompressed
	hash  string // hash is a sha256 hash of the file contents. Used for the Etag, and useful for caching
}

// provideStaticFiles returns a map of static files derived from the
// _templates directory, or the first error encountered
func provideStaticFiles() (fls map[string]*staticFilesFile, err error) {

	fls = make(map[string]*staticFilesFile)

	err = fs.WalkDir(statics, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		key := strings.TrimPrefix(path, staticRoot+"/")
		if _, already := fls[key]; !already {
			if !d.IsDir() {
				data, err := fs.ReadFile(statics, path)
				if err != nil {
					return fmt.Errorf("file %s: %s, %w", key, err, StaticFileReadError)
				}
				info, err := d.Info()
				if err != nil {
					return fmt.Errorf("file info %s: %s, %w", key, err, StaticFileReadError)
				}
				fls[key] = &staticFilesFile{
					mime:  mime.TypeByExtension(key),
					data:  string(data),
					mtime: info.ModTime(),
					size:  0, // report as always uncompressed
					hash:  fmt.Sprintf("%x", sha256.Sum256(data)),
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walkdir %s: %w", staticRoot, StaticFileReadError)
	}

	return fls, nil
}
