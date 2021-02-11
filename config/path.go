package config

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/bmatcuk/doublestar"
	"github.com/shurcooL/httpfs/vfsutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

const configRootConfig = "config"

var configFileExtensions = []string{".yaml", ".yml", ".ini", ".json", ".json5", ".properties"}

func NewFileProvider(name, fileName string) Provider {
	fileWatcher := NewFileNotifier(fileName)
	reader := FileContentReader(fileName)
	return newFileProvider(name, fileName, reader, fileWatcher)
}

func NewHttpFileProvider(name string, fs http.FileSystem, fileName string) Provider {
	return newFileProvider(name, fileName, HttpFileContentReader(fs, fileName), nil)
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
		provider := NewCacheProvider(NewHttpFileProvider(name, fs, fileName))
		providers = append(providers, provider)
	}

	return providers
}

type fsConfig struct {
	Configs string
	Local   string
	Roots   struct {
		Sources string
		Staging string
		Release string
		Command string
	}
}

func AddConfigFoldersFromFsConfig(cfg *Config) {
	var fs fsConfig
	if err := cfg.Populate(&fs, "fs"); err != nil {
		return
	}

	if fs.Roots.Sources != "" {
		if fs.Roots.Command != "" {
			AddConfigFolders(
				fs.Roots.Command)
		}

		AddConfigFolders(
			fs.Roots.Sources+fs.Local,
			fs.Roots.Staging+fs.Configs)
	}

	if fs.Roots.Release != "" {
		AddConfigFolders(
			fs.Roots.Release + fs.Configs)
	}
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

var configFolders = types.StringStack{
	".",
}

func AddConfigFolders(folders ...string) {
	for _, folder := range folders {
		absFolder, err := filepath.Abs(folder)
		if err == nil {
			configFolders = append(configFolders, absFolder)
		}
	}
}

func Folders() []string {
	return append([]string{}, configFolders...)
}

func NewFileProvidersFromBaseName(name, baseName string) []Provider {
	var found = map[string]string{}
	var results []Provider
	for _, folder := range configFolders {
		for _, ext := range configFileExtensions {
			fullPath := path.Join(folder, baseName+ext)
			info, err := os.Stat(fullPath)
			if os.IsNotExist(err) || info.IsDir() {
				continue
			}

			fileName := filepath.Base(fullPath)
			if previous, ok := found[fileName]; ok {
				logger.Warnf("Skipping %q due to previously found %q", fullPath, previous)
				continue
			}
			found[fileName] = fullPath

			provider := NewCacheProvider(NewFileProvider(name, fullPath))
			results = append(results, provider)
		}
	}

	if len(results) == 0 {
		logger.Warnf("Could not find %s.{yaml,yml,ini,json,json5,properties}", baseName)
	}

	return results
}

func NewFileProvidersFromGlob(name, glob string) []Provider {
	var found = map[string]string{}
	var results []Provider
	for _, folder := range configFolders {
		folderGlob := path.Join(folder, glob)
		for _, ext := range configFileExtensions {
			fileGlob := folderGlob + ext
			files, err := doublestar.Glob(fileGlob)
			if err != nil {
				continue
			}

			for _, file := range files {
				fileName := filepath.Base(file)
				if previous, ok := found[fileName]; ok {
					logger.Warnf("Skipping %q due to previously found %q", file, previous)
					continue
				}
				found[fileName] = file

				provider := NewCacheProvider(NewFileProvider(name, file))
				results = append(results, provider)
			}
		}
	}

	return results
}
