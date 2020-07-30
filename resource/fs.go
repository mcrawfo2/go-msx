package resource

import (
	"cto-github.cisco.com/NFV-BU/go-msx/fs"
	"github.com/pkg/errors"
	"net/http"
	"os"
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
	if fs.Sources() == "" {
		logger.Info("Using release filesystem")
		return newReleaseFileSystem("resources", filepath.Join(fs.Root(), fs.Resources())), nil
	}

	sourceFileSystem, err := newSourceFileSystem()
	if err == ErrFilesystemUnavailable {
		logger.Info("Using release filesystem")
		return newReleaseFileSystem("resources", filepath.Join(fs.Root(), fs.Resources())), nil
	}

	stagingFileSystem, err := newStagingFileSystem()
	if err == ErrFilesystemUnavailable {
		logger.Info("Using source filesystem")
		return sourceFileSystem, nil
	}

	logger.Info("Using source and staging overlay filesystem")
	return fs.LoggingFilesystem{
		Name: "overlay",
		Fs:   fs.NewOverlayFileSystem(stagingFileSystem, sourceFileSystem),
	}, nil
}

func newSourceFileSystem() (http.FileSystem, error) {
	if fs.Sources() == "" {
		return nil, ErrFilesystemUnavailable
	}
	_, err := os.Stat(fs.Sources())
	if os.IsNotExist(err) {
		return nil, ErrFilesystemUnavailable
	}
	return newReleaseFileSystem("source", fs.Sources()), nil
}

func newStagingFileSystem() (http.FileSystem, error) {
	if fs.Sources() == "" {
		return nil, ErrFilesystemUnavailable
	}
	_, err := os.Stat(fs.Sources())
	if os.IsNotExist(err) {
		return nil, ErrFilesystemUnavailable
	}
	parentFileSystem := newReleaseFileSystem("source", fs.Sources())
	stagingFileSystem, err := fs.NewPrefixFileSystem(parentFileSystem, filepath.Join("/dist/root", fs.Resources()))
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
			Fs: http.Dir(root),
		},
	}
}
