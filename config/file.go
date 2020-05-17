package config

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/bmatcuk/doublestar"
	"github.com/shurcooL/httpfs/vfsutil"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const configRootConfig = "config"

var configFileExtensions = []string{".yaml", ".yml", ".ini", ".json", ".json5", ".properties", ".toml"}

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

func NewHttpFileProvidersFromGlob(name string, fs http.FileSystem, glob string) []Provider {
	var configFiles = make(types.StringSet)
	_ = vfsutil.Walk(fs, "/", func(path string, info os.FileInfo, err error) error {
		for _, ext := range configFileExtensions {
			fileGlob := glob + ext
			if inc, err2 := doublestar.Match(fileGlob, path); err2 != nil {
				continue
			} else if inc {
				configFiles.Add(path)
			}
		}
		return nil
	})

	var providers []Provider
	for fileName := range configFiles {
		provider := NewProvider(name, fileName, HttpFileContentReader(fs, fileName))
		providers = append(providers, provider)
	}

	return providers
}

type configConfig struct {
	Path []string
}

func AddConfigFoldersFromPathConfig(cfg *Config) {
	var pathConfig configConfig
	if err := cfg.Populate(&pathConfig, configRootConfig); err != nil {
		return
	}

	AddConfigFolders(pathConfig.Path...)
}

var configFolders = types.StringSet{
	".": {},
}

func AddConfigFolders(folders ...string) {
	for _, folder := range folders {
		absFolder, err := filepath.Abs(folder)
		if err == nil {
			configFolders.Add(absFolder)
		}
	}
}

func ConfigFolders() []string {
	return configFolders.Values()
}

func NewFileProvidersFromBaseName(name, baseName string) []Provider {
	var results []Provider
	for _, folder := range configFolders.Values() {
		for _, ext := range configFileExtensions {
			fullPath := path.Join(folder, baseName+ext)
			info, err := os.Stat(fullPath)
			if os.IsNotExist(err) || info.IsDir() {
				continue
			}

			provider := NewFileProvider(name, fullPath)
			results = append(results, provider)
		}
	}

	if len(results) == 0 {
		logger.Warnf("Could not find %s.{yaml,yml,ini,json,json5,properties,toml}", baseName)
	}

	return results
}

func NewFileProvidersFromGlob(name, glob string) []Provider {
	var results []Provider
	for _, folder := range configFolders.Values() {
		folderGlob := path.Join(folder, glob)
		for _, ext := range configFileExtensions {
			fileGlob := folderGlob + ext
			files, err := doublestar.Glob(fileGlob)
			if err != nil {
				continue
			}

			for _, file := range files {
				provider := NewFileProvider(name, file)
				results = append(results, provider)
			}
		}
	}

	return results
}
