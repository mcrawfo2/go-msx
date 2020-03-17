package config

import (
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

type ContentReader func() ([]byte, error)

func FileContentReader(fileName string) ContentReader {
	return func() (bytes []byte, err error) {
		return ioutil.ReadFile(fileName)
	}
}

func HttpFileContentReader(fs http.FileSystem, fileName string) ContentReader {
	return func() (bytes []byte, err error) {
		file, err := fs.Open(fileName)
		if err != nil {
			return nil, err
		}
		return ioutil.ReadAll(file)
	}
}

func NewFileProvider(name, fileName string) Provider {
	return NewProvider(name, fileName, FileContentReader(fileName))
}

func NewHttpFileProvider(name string, fs http.FileSystem, fileName string) Provider {
	return NewProvider(name, fileName, HttpFileContentReader(fs, fileName))
}

func NewProvider(name, fileName string, reader ContentReader) Provider {
	fileExt := strings.ToLower(path.Ext(fileName))
	switch fileExt {
	case ".yml", ".yaml":
		return NewCachedLoader(NewYAMLFile(name, fileName, reader))
	case ".ini":
		return NewCachedLoader(NewINIFile(name, fileName, reader))
	case ".json", ".json5":
		return NewCachedLoader(NewJSONFile(name, fileName, reader))
	case ".toml":
		return NewCachedLoader(NewTOMLFile(name, fileName, reader))
	case ".properties":
		return NewCachedLoader(NewPropertiesFile(name, fileName, reader))
	default:
		logger.Error("Unknown config file extension: ", fileExt)
		return nil
	}
}
