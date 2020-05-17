package appconfig

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/bmatcuk/doublestar"
	"github.com/shurcooL/httpfs/vfsutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

const (
	configKeyAppName = "spring.application.name"
)

var (
	logger               = log.NewLogger("msx.app.appconfig")
	configFileExtensions = []string{".yaml", ".yml", ".ini", ".json", ".json5", ".properties", ".toml"}
)

// Deprecated
func FindConfigFiles(cfg *config.Config, baseName string) []string {
	folders := FindConfigFolders(cfg)

	var results []string
	for _, folder := range folders {
		for _, ext := range configFileExtensions {
			fullPath := path.Join(folder, baseName+ext)
			info, err := os.Stat(fullPath)
			if os.IsNotExist(err) || info.IsDir() {
				continue
			}

			results = append(results, fullPath)
		}
	}

	if len(results) == 0 {
		logger.Warnf("Could not find %s.{yaml,yml,ini,json,json5,properties}", baseName)
	}

	return results
}

// Deprecated
func FindConfigFolders(cfg *config.Config) []string {
	folders := []string{"."}
	if cfg != nil {
		appName, err := cfg.String(configKeyAppName)
		if err == nil && appName != "" {
			folders = append(folders, path.Join("/etc", appName))
		}
	}
	folders = append(folders, Config.Path...)

	for i, folder := range folders {
		absFolder, err := filepath.Abs(folder)
		if err == nil {
			folders[i] = absFolder
		}
	}

	return folders
}

// Deprecated
func FindConfigFilesGlob(cfg *config.Config, glob string) []string {
	folders := FindConfigFolders(cfg)

	var results []string
	for _, folder := range folders {
		folderGlob := path.Join(folder, glob)
		for _, ext := range configFileExtensions {
			fileGlob := folderGlob + ext
			files, err := doublestar.Glob(fileGlob)
			if err != nil {
				continue
			}
			results = append(results, files...)
		}
	}

	return results
}

// Deprecated
func FindConfigHttpFilesGlob(fs http.FileSystem, glob string) []string {
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
	return configFiles.Values()
}
