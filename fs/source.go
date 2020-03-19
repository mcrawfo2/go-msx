package fs

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var errFilesystemUnavailable = errors.New("Filesystem unavailable")
var fs http.FileSystem
var fsConfig *FileSystemConfig
var fsMtx sync.Mutex

func ConfigureFileSystem(cfg *config.Config) (err error) {
	fsMtx.Lock()
	defer fsMtx.Unlock()

	fsConfig, err = NewFileSystemConfig(cfg)
	if err != nil {
		return err
	}

	if fsConfig.Sources == "" {
		fsConfig.Sources, err = getSourceDir()
		if err != nil {
			return err
		}
	}

	fs, err = newVirtualFileSystem()
	return err
}

func FileSystem() (http.FileSystem, error) {
	if fs == nil {
		return nil, errFilesystemUnavailable
	}
	return fs, nil
}

func SourcePath(path string) (string, error) {
	if fsConfig.Sources == "" {
		return path, errFilesystemUnavailable
	}

	path = strings.TrimPrefix(path, fsConfig.Sources)
	return path, nil
}

// For unit testing
func SetSources() error {
	var err error
	if fsConfig == nil {
		fsConfig = new(FileSystemConfig)
	}
	_, file, _, _ := runtime.Caller(1)
	thence := findSourceDir(file)
	fsConfig.Sources, err = filepath.Abs(thence)
	fs, err = newVirtualFileSystem()
	return err
}

func newVirtualFileSystem() (http.FileSystem, error) {
	sourceFileSystem, err := newSourceFileSystem()
	if err == errFilesystemUnavailable {
		logger.Info("Using release filesystem")
		return newReleaseFileSystem(fsConfig.Root), nil
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
	if fsConfig.Sources == "" {
		return nil, errFilesystemUnavailable
	}
	_, err := os.Stat(fsConfig.Sources)
	if os.IsNotExist(err) {
		return nil, errFilesystemUnavailable
	}
	parentFileSystem := newReleaseFileSystem("/")
	return NewPrefixFileSystem(parentFileSystem, fsConfig.Sources)
}

func newStagingFileSystem() (http.FileSystem, error) {
	if fsConfig.Sources == "" {
		return nil, errFilesystemUnavailable
	}
	_, err := os.Stat(fsConfig.Sources)
	if os.IsNotExist(err) {
		return nil, errFilesystemUnavailable
	}
	parentFileSystem := newReleaseFileSystem("/")
	return NewPrefixFileSystem(parentFileSystem, filepath.Join(fsConfig.Sources, "/dist/root"))
}

func newReleaseFileSystem(root string) http.FileSystem {
	return http.Dir(root)
}

func getSourceDir() (string, error) {
	file, ok := getEntryPointFile()
	if !ok {
		return "", errFilesystemUnavailable
	}

	thence := findSourceDir(file)
	if thence == "" {
		return "", errFilesystemUnavailable
	}

	return thence, nil
}

// Hack when fs.sources is missing
func getEntryPointFile() (string, bool) {
	pcs := make([]uintptr, 32)
	frameCount := runtime.Callers(3, pcs)
	frames := runtime.CallersFrames(pcs[:frameCount])
	var lastFrame runtime.Frame
	for {
		frame, more := frames.Next()
		if strings.HasSuffix(frame.Function, "go-msx/app.Run") {
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
