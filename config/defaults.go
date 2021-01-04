package config

// Defaults holds application-wide defaults
var Defaults = map[string]string{}
var DefaultsProvider = NewInMemoryProvider("Default", Defaults)
var DefaultsCache = NewCacheProvider(DefaultsProvider)

// EmbeddedDefaultsProvider presents the defaults from go-msx
var EmbeddedDefaultsProviders = NewHttpFileProvidersFromGlob("EmbeddedDefaults", EmbeddedDefaultsFileSystem, "**/defaults-*")
