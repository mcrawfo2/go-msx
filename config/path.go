// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package config

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/bmatcuk/doublestar"
	"github.com/shurcooL/httpfs/vfsutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
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

func NewHttpFileProvidersFromGlob(name string, fsys http.FileSystem, glob string) []Provider {
	var configFiles = make(types.StringSet)
	reg := `.*\.(` +
		strings.ReplaceAll(strings.Join(configFileExtensions, "|"), ".", "") +
		`)$`
	exts, err := regexp.Compile(reg)
	if err != nil {
		logger.Warnf("File extension list caused a regex compile fail %s", reg)
		return nil
	}

	err = vfsutil.Walk(fsys, "/", func(path string, info os.FileInfo, err error) error {
		inc, err2 := doublestar.Match(glob, path)
		if err2 != nil {
			logger.Warnf("Malformed glob (eww) %s", glob)
			return err2
		}
		if inc && exts.MatchString(path) {
			configFiles.Add(path)
		}
		return nil
	})
	if err != nil {
		return nil
	}

	var providers []Provider
	for fileName := range configFiles {
		provider := NewCacheProvider(NewHttpFileProvider(name, fsys, fileName))
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
	var fsys fsConfig
	if err := cfg.Populate(&fsys, "fs"); err != nil {
		return
	}

	if fsys.Roots.Sources != "" {
		if fsys.Roots.Command != "" {
			AddConfigFolders(
				fsys.Roots.Command)
		}

		AddConfigFolders(
			filepath.Join(fsys.Roots.Sources, fsys.Local),
			filepath.Join(fsys.Roots.Staging, fsys.Configs))
	}

	if fsys.Roots.Release != "" {
		AddConfigFolders(
			filepath.Join(fsys.Roots.Release, fsys.Configs))
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
			fullPath := filepath.Join(folder, baseName+ext)
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
