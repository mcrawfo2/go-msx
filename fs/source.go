package fs

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var errFilesystemUnavailable = errors.New("Filesystem unavailable: not running from source directory")

func NewVirtualFileSystem(cfg *config.Config) (http.FileSystem, error) {
	fsConfig, err := NewFileSystemConfig(cfg)
	if err != nil {
		return nil, err
	}

	return NewVirtualFileSystemFromConfig(fsConfig)
}

func NewVirtualFileSystemFromConfig(cfg *FileSystemConfig) (http.FileSystem, error) {
	sourceFileSystem, err := newSourceFileSystem()
	if err == errFilesystemUnavailable {
		logger.Info("Using release filesystem")
		return newReleaseFileSystem(cfg.Root), nil
	}

	stagingFileSystem, err := newStagingFileSystem()
	if err == errFilesystemUnavailable {
		logger.Info("Using source filesystem")
		return sourceFileSystem, nil
	}

	logger.Info("Using source and staging overlay filesystem")
	return NewOverlayFileSystem(stagingFileSystem, sourceFileSystem), nil
}

func newSourceFileSystem() (http.FileSystem, error) {
	file, ok := getEntryPointFile()
	if !ok {
		return nil, errFilesystemUnavailable
	}

	thence := findSourceDir(file)
	if thence == "" {
		return nil, errFilesystemUnavailable
	}

	parentFileSystem := newReleaseFileSystem("/")
	return NewPrefixFileSystem(parentFileSystem, thence)
}

func newStagingFileSystem() (http.FileSystem, error) {
	file, ok := getEntryPointFile()
	if !ok {
		return nil, errFilesystemUnavailable
	}

	thence := findSourceDir(file)
	if thence == "" {
		return nil, errFilesystemUnavailable
	}

	parentFileSystem := newReleaseFileSystem("/")
	return NewPrefixFileSystem(parentFileSystem, filepath.Join(thence, "/dist/root"))
}

func newReleaseFileSystem(root string) http.FileSystem {
	return http.Dir(root)
}

// Hack
func getEntryPointFile() (string, bool) {
	pcs := make([]uintptr, 32)
	frameCount := runtime.Callers(2, pcs)
	frames := runtime.CallersFrames(pcs[:frameCount])
	var lastFrame runtime.Frame
	for {
		frame, more := frames.Next()
		if strings.HasSuffix(frame.File, "go-msx/app/command.go") {
			lastFrame, more = frames.Next()
		}
		if !more {
			break
		}
	}
	if lastFrame.File != "" {
		return lastFrame.File, true
	}
	return "", false
}

func findSourceDir(whence string) string {
	for whence != "/" {
		whence = filepath.Dir(whence)
		gomod := filepath.Join(whence, "go.mod")
		_, err := os.Stat(gomod)
		if !os.IsNotExist(err) {
			return whence
		}
	}

	return ""
}
