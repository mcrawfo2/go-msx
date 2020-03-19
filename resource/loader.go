package resource

import (
	"cto-github.cisco.com/NFV-BU/go-msx/fs"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
)

func Load(resourceName string) (data []byte, err error) {
	absPath, err := abs(resourceName)
	if err != nil {
		return nil, err
	}

	return load(absPath)
}

func Unmarshal(resourceName string, target interface{}) (err error) {
	absPath, err := abs(resourceName)
	if err != nil {
		return err
	}

	bytes, err := load(absPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, target)
}

func abs(filename string) (string, error) {
	if strings.HasPrefix(filename, "/") {
		return filename, nil
	}

	_, file, _, ok := runtime.Caller(2)
	if !ok {
		return "", errors.New("Failed to identify source file of caller")
	}

	base := filepath.Dir(file)
	full := filepath.Join(base, filename)
	return fs.SourcePath(full)
}

func load(resourcePath string) ([]byte, error) {
	fileSystem, err := fs.FileSystem()
	if err != nil {
		return nil, err
	}

	reader, err := fileSystem.Open(resourcePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}
