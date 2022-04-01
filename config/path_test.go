// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package config

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewFileProvider(t *testing.T) {
	vals := map[string]string{
		"server.context-path": "www",
		"server.port":         "8882",
		"client.refresh":      "5.05",
		"client.enabled":      "true",
		"client.items[0]":     "abc",
	}

	tests := []struct {
		name     string
		fileName string
		want     map[string]string
		wantErr  bool
	}{
		{
			name:     "Ini",
			fileName: "testdata/config.ini",
			want:     vals,
		},
		{
			name:     "Json",
			fileName: "testdata/config.json",
			want:     vals,
		},
		{
			name:     "Yaml",
			fileName: "testdata/config.yaml",
			want:     vals,
		},
		{
			name:     "Properties",
			fileName: "testdata/config.properties",
			want:     vals,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewFileProvider(tt.name, tt.fileName)
			entries, err := provider.Load(context.Background())
			if tt.wantErr != (err != nil) {
				t.Errorf("FileProvider.Load() got err = %v, wantErr = %v", err, tt.wantErr)
			}

			settings := MapFromEntries(entries)
			if !reflect.DeepEqual(settings, tt.want) {
				t.Errorf("FileProvider.Load() got = %v, want %v", settings, tt.want)
			}
		})
	}
}

func TestAddConfigFoldersFromFsConfig_WithSources(t *testing.T) {
	configFolders = types.StringStack{"."}

	cfg := NewInMemoryConfig(map[string]string{
		"spring.application.name": "TestAddConfigFoldersFromFsConfig",
		"fs.root":                 "/",
		"fs.sources":              "/home/ubuntu/go-msx",
		"fs.resources":            "/var/lib/${spring.application.name}",
		"fs.configs":              "/etc/${spring.application.name}",
		"fs.binaries":             "/usr/bin",
		"fs.local":                "/local",
		"fs.staging":              "/dist/root",
		"fs.command":              "/cmd/app",
		"fs.roots.command":        "${fs.sources}${fs.command}",
		"fs.roots.release":        "${fs.root}",
		"fs.roots.sources":        "${fs.sources}",
		"fs.roots.staging":        "${fs.roots.sources}${fs.staging}",
	})

	AddConfigFoldersFromFsConfig(cfg)

	assert.Equal(t,
		types.StringStack([]string{
			".",
			"/home/ubuntu/go-msx/cmd/app",
			"/home/ubuntu/go-msx/local",
			"/home/ubuntu/go-msx/dist/root/etc/TestAddConfigFoldersFromFsConfig",
			"/etc/TestAddConfigFoldersFromFsConfig",
		}),
		configFolders)
}

func TestAddConfigFoldersFromFsConfig_NoSources(t *testing.T) {
	configFolders = types.StringStack{"."}

	cfg := NewInMemoryConfig(map[string]string{
		"spring.application.name": "TestAddConfigFoldersFromFsConfig",
		"fs.root":                 "/",
		"fs.resources":            "/var/lib/${spring.application.name}",
		"fs.configs":              "/etc/${spring.application.name}",
		"fs.binaries":             "/usr/bin",
		"fs.local":                "/local",
		"fs.staging":              "/dist/root",
		"fs.command":              "/cmd/app",
		"fs.roots.command":        "${fs.sources}${fs.command}",
		"fs.roots.release":        "${fs.root}",
		"fs.roots.sources":        "${fs.sources}",
		"fs.roots.staging":        "${fs.roots.sources}${fs.staging}",
	})

	AddConfigFoldersFromFsConfig(cfg)

	assert.Equal(t,
		types.StringStack([]string{
			".",
			"/etc/TestAddConfigFoldersFromFsConfig",
		}),
		configFolders)
}

func TestAddAddConfigFoldersFromPathConfig_Single(t *testing.T) {
	t.Run("Single", func(t *testing.T) {
		configFolders = types.StringStack{"."}

		cfg := NewInMemoryConfig(map[string]string{
			"config.path": "/home/ubuntu/go-msx/local",
		})

		AddConfigFoldersFromPathConfig(cfg)

		assert.Equal(t,
			types.StringStack([]string{
				".",
				"/home/ubuntu/go-msx/local",
			}),
			configFolders)
	})

	t.Run("Multi", func(t *testing.T) {
		configFolders = types.StringStack{"."}

		cfg := NewInMemoryConfig(map[string]string{
			"config.path[0]": "/home/ubuntu/go-msx/local",
			"config.path[1]": "/home/ubuntu/go-msx/remote",
		})

		AddConfigFoldersFromPathConfig(cfg)

		assert.Equal(t,
			types.StringStack([]string{
				".",
				"/home/ubuntu/go-msx/local",
				"/home/ubuntu/go-msx/remote",
			}),
			configFolders)
	})
}

func TestFolders(t *testing.T) {
	configFolders = []string{"alpha"}
	assert.Equal(t, []string(configFolders), Folders())
}

func TestNewFileProvidersFromBaseName(t *testing.T) {
	configFolder, _ := filepath.Abs("testdata")
	configFolders = []string{configFolder, configFolder}

	providers := NewFileProvidersFromBaseName("testdata", "config")
	assert.Len(t, providers, 4)
}

func TestNewFileProvidersFromGlob(t *testing.T) {
	configFolder, _ := filepath.Abs(".")
	configFolders = []string{configFolder, configFolder}

	providers := NewFileProvidersFromGlob("testdata", "testdata/config")
	assert.Len(t, providers, 4)
}
