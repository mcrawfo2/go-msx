// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package config

import (
	"embed"
	"io/fs"
	"net/http"
)

// Defaults holds application-wide defaults
var Defaults = map[string]string{}
var DefaultsProvider = NewInMemoryProvider("Default", Defaults)
var DefaultsCache = NewCacheProvider(DefaultsProvider)

// EmbeddedDefaultsProviders presents the defaults from go-msx

var EmbeddedDefaultsProviders []Provider
var EmbeddedConfigs fs.FS

// Embeds config files in the executable.
// Due to limitations in go:embed, all config files that need to be
// embedded should be placed in the config/embed directory. Place symlinks
// in the original locations, pointing into the latter, if you wish

//go:embed embed/*
var embeddedConfigs embed.FS

func init() {
	var err error
	EmbeddedConfigs, err = fs.Sub(embeddedConfigs, "embed")
	if err != nil {
		logger.Warnf("Could not strip embed dir from embedded config files")
	}
	EmbeddedDefaultsProviders = NewHttpFileProvidersFromGlob("EmbeddedDefaults",
		http.FS(EmbeddedConfigs), "**/*")
}
