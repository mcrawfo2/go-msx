// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package config

// Defaults holds application-wide defaults
var Defaults = map[string]string{}
var DefaultsProvider = NewInMemoryProvider("Default", Defaults)
var DefaultsCache = NewCacheProvider(DefaultsProvider)

// EmbeddedDefaultsProvider presents the defaults from go-msx
var EmbeddedDefaultsProviders = NewHttpFileProvidersFromGlob("EmbeddedDefaults", EmbeddedDefaultsFileSystem, "**/defaults-*")
