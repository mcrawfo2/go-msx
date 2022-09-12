// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package resource

import (
	"cto-github.cisco.com/NFV-BU/go-msx/fs"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

var ErrFilesystemUnavailable = errors.New("Filesystem unavailable")
var rfs http.FileSystem

func FileSystem() (result http.FileSystem, err error) {
	if rfs == nil {
		rfs, err = newFileSystem()
	}
	return rfs, err
}

func newFileSystem() (http.FileSystem, error) {
	fsConfigBytes, _ := json.Marshal(fs.Config())
	logger.Infof("Filesystem Config: %s", string(fsConfigBytes))

	sourceFileSystem, err := newSourceFileSystem()
	if err == ErrFilesystemUnavailable {
		return newReleaseOnlyFileSystem(), nil
	} else if err != nil {
		return nil, err
	}

	stagingFileSystem, err := newStagingFileSystem()
	if err == ErrFilesystemUnavailable {
		logger.Info("Using source filesystem")
		return sourceFileSystem, nil
	} else if err != nil {
		return nil, err
	}

	logger.Info("Using source and staging overlay filesystem")
	return fs.LoggingFilesystem{
		Name: "overlay",
		Fs:   fs.NewOverlayFileSystem(stagingFileSystem, sourceFileSystem),
	}, nil
}

func newReleaseOnlyFileSystem() http.FileSystem {
	logger.Info("Using release filesystem")
	return newReleaseFileSystem("resources", filepath.Join(fs.Root(), fs.Resources()))
}

func newSourceFileSystem() (http.FileSystem, error) {
	if fs.Sources() == "" {
		return nil, ErrFilesystemUnavailable
	}
	_, err := os.Stat(fs.Sources())
	if os.IsNotExist(err) {
		return nil, ErrFilesystemUnavailable
	} else if err != nil {
		return nil, err
	}
	logger.Infof("Located source filesystem: %s", fs.Sources())
	return newReleaseFileSystem("source", fs.Sources()), nil
}

func newStagingFileSystem() (http.FileSystem, error) {
	if fs.Sources() == "" {
		return nil, ErrFilesystemUnavailable
	}
	_, err := os.Stat(fs.Sources())
	if os.IsNotExist(err) {
		return nil, ErrFilesystemUnavailable
	} else if err != nil {
		return nil, err
	}
	parentFileSystem := newReleaseFileSystem("source", fs.Sources())
	stagingFileSystem, err := fs.NewPrefixFileSystem(parentFileSystem, path.Join("/dist/root", fs.Resources()))
	if err != nil {
		return nil, err
	}
	return fs.LoggingFilesystem{
		Name: "staging",
		Fs:   stagingFileSystem,
	}, nil
}

func newReleaseFileSystem(name, root string) http.FileSystem {
	return fs.LoggingFilesystem{
		Name: name,
		Fs: fs.RootLoggingFilesystem{
			Dir: http.Dir(root),
		},
	}
}
