package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/config/consulprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/config/vaultprovider"
	"github.com/bmatcuk/doublestar"
	"github.com/pkg/errors"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	SourceConsul      = "Consul"
	SourceVault       = "Vault"
	SourceApplication = "Application"
	SourceProfile     = "Profile"
	SourceCommandLine = "CommandLine"
	SourceBuild       = "Build"
	SourceBootStrap   = "Bootstrap"
	SourceEnvironment = "Environment"
	SourceOverride    = "Override"
	SourceDefaults    = "Defaults"
	SourceDefault     = "Default"

	configKeyAppName = "spring.application.name"
	configRootConfig = "config"
)

var (
	applicationConfig    *config.Config
	overrideConfig       = make(map[string]string)
	configFileExtensions = []string{".yaml", ".yml", ".ini", ".json", ".json5", ".properties", ".toml"}
)

type ConfigConfig struct {
	Path []string `config:"default="`
}

var configConfig ConfigConfig

type Sources struct {
	Default          config.Provider
	DefaultsFiles    []config.Provider
	BootstrapFiles   []config.Provider
	ApplicationFiles []config.Provider
	BuildFiles       []config.Provider
	Consul           []config.Provider
	Vault            []config.Provider
	ProfileFiles     []config.Provider
	Environment      config.Provider
	CommandLine      []config.Provider
	Override         config.Provider
}

type SourcesList []config.Provider

func (s *SourcesList) Append(providers ...config.Provider) {
	for _, provider := range providers {
		if provider != nil {
			*s = append(*s, provider)
		}
	}
}

func (c Sources) Providers() []config.Provider {
	sourcesList := &SourcesList{}
	sourcesList.Append(c.Default)
	sourcesList.Append(c.DefaultsFiles...)
	sourcesList.Append(c.BootstrapFiles...)
	sourcesList.Append(c.ApplicationFiles...)
	sourcesList.Append(c.BuildFiles...)
	sourcesList.Append(c.Consul...)
	sourcesList.Append(c.Vault...)
	sourcesList.Append(c.ProfileFiles...)
	sourcesList.Append(c.Environment)
	sourcesList.Append(c.CommandLine...)
	sourcesList.Append(c.Override)
	return *sourcesList
}

type ProviderFactory func(string, *config.Config) ([]config.Provider, error)

var providerFactories = map[string]ProviderFactory{
	SourceCommandLine: nil,
	SourceConsul:      nil,
	SourceVault:       nil,
	SourceApplication: newApplicationProviders,
	SourceProfile:     newProfileProvider,
}

func RegisterProviderFactory(name string, factory ProviderFactory) {
	providerFactories[name] = factory
}

func newDefaultsProvider() config.Provider {
	return config.NewCachedLoader(config.NewStatic(SourceDefault, map[string]string{
		"profile": "default",
	}))
}

func newDefaultsFilesProviders(cfg *config.Config) []config.Provider {
	files := findConfigFilesGlob(cfg, "defaults-*")
	return newFileProviders(SourceDefaults, files)
}

func newBootstrapProviders(cfg *config.Config) []config.Provider {
	files := findConfigFiles(cfg, "bootstrap")
	return newFileProviders(SourceBootStrap, files)
}

func newEnvironmentProvider() config.Provider {
	return config.NewCachedLoader(config.NewEnvironment(SourceEnvironment))
}

func newOverrideProvider(static map[string]string) config.Provider {
	return config.NewCachedLoader(config.NewStatic(SourceOverride, static))
}

func newApplicationProviders(name string, cfg *config.Config) ([]config.Provider, error) {
	files := findConfigFiles(cfg, "application")
	if appName, err := cfg.String(configKeyAppName); err == nil {
		files = append(files, findConfigFiles(cfg, appName)...)
	}
	return newFileProviders(SourceApplication, files), nil
}

func newBuildProvider(cfg *config.Config) []config.Provider {
	files := findConfigFiles(cfg, "buildinfo")
	return newFileProviders(SourceBuild, files)
}

func newProfileProvider(name string, cfg *config.Config) ([]config.Provider, error) {
	var parts []string
	if appName, err := cfg.String(configKeyAppName); err != nil {
		return nil, errors.Wrap(err, "Application name not defined")
	} else {
		parts = []string{appName}
	}

	if profile, err := cfg.StringOr("profile", "default"); err != nil {
		return nil, err
	} else if profile == "default" {
		// don't add it
	} else {
		parts = append(parts, profile)
	}

	files := findConfigFiles(cfg, strings.Join(parts, "."))
	return newFileProviders(SourceProfile, files), nil
}

func newFileProviders(name string, files []string) []config.Provider {
	var providers []config.Provider
	for _, file := range files {
		providers = append(providers, config.NewFileProvider(SourceApplication, file))
	}
	return providers
}

func newProviders(name string, cfg *config.Config) ([]config.Provider, error) {
	if providerFactories[name] == nil {
		return nil, nil
	}
	providerFactory := providerFactories[name]
	return providerFactory(name, cfg)
}

func findConfigFiles(cfg *config.Config, baseName string) []string {
	folders := findConfigFolders(cfg)

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

func findConfigFolders(cfg *config.Config) []string {
	folders := []string{"."}
	folders = append(folders, configConfig.Path...)
	if cfg != nil {
		appName, err := cfg.String(configKeyAppName)
		if err != nil && appName != "" {
			folders = append(folders, path.Join("etc", appName))
		}
	}

	for i, folder := range folders {
		absFolder, err := filepath.Abs(folder)
		if err == nil {
			folders[i] = absFolder
		}
	}

	return folders
}

func findConfigFilesGlob(cfg *config.Config, glob string) []string {
	folders := findConfigFolders(cfg)

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

func init() {
	OnEvent(EventConfigure, PhaseBefore, registerRemoteConfigProviders)
	OnEvent(EventConfigure, PhaseDuring, loadConfig)
	OnEvent(EventStart, PhaseAfter, watchConfig)
}

func registerRemoteConfigProviders(ctx context.Context) error {
	RegisterProviderFactory(SourceConsul, func(name string, cfg *config.Config) (providers []config.Provider, err error) {
		providers, err = consulprovider.NewConfigProvidersFromConfig(name, cfg)
		if err != nil {
			return nil, err
		}

		for i := range providers {
			providers[i] = config.NewCachedLoader(providers[i])
		}

		return providers, nil
	})

	RegisterProviderFactory(SourceVault, func(name string, cfg *config.Config) (providers []config.Provider, err error) {
		providers, err = vaultprovider.NewConfigProvidersFromConfig(name, cfg)
		if err != nil {
			return nil, err
		}

		for i := range providers {
			providers[i] = config.NewCachedLoader(providers[i])
		}

		return providers, nil
	})

	return nil
}

func watchConfig(ctx context.Context) error {
	logger.Info("Watching configuration for changes")
	cfg := applicationConfig
	cfg.Notify = func(keys []string) {
		for _, k := range keys {
			logger.Warnf("Configuration changed: %s", k)
		}
		if err := application.Refresh(); err != nil {
			logger.Errorf("Failed to refresh: ", err)
		}
	}

	return cfg.Watch(ctx)
}

func mustLoadConfig(ctx context.Context, cfg *config.Config) error {
	loadContext, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	return cfg.Load(loadContext)
}

func loadConfig(ctx context.Context) (err error) {
	sources := &Sources{
		Default:     newDefaultsProvider(),
		Environment: newEnvironmentProvider(),
		Override:    newOverrideProvider(overrideConfig),
	}

	if sources.CommandLine, err = newProviders(SourceCommandLine, nil); err != nil {
		return err
	}

	var cfg = config.NewConfig(sources.Providers()...)
	if err := mustLoadConfig(ctx, cfg); err != nil {
		return err
	}

	if err := cfg.Populate(&configConfig, configRootConfig); err != nil {
		return err
	}

	sources.DefaultsFiles = newDefaultsFilesProviders(cfg)
	sources.BootstrapFiles = newBootstrapProviders(cfg)

	if sources.ApplicationFiles, err = newProviders(SourceApplication, cfg); err != nil {
		return err
	}

	if sources.ProfileFiles, err = newProviders(SourceProfile, cfg); err != nil {
		return err
	}

	sources.BuildFiles = newBuildProvider(cfg)

	cfg = config.NewConfig(sources.Providers()...)
	if err := mustLoadConfig(ctx, cfg); err != nil {
		return err
	}

	if sources.Consul, err = newProviders(SourceConsul, cfg); err != nil {
		return err
	}

	if sources.Vault, err = newProviders(SourceVault, cfg); err != nil {
		return err
	}

	cfg = config.NewConfig(sources.Providers()...)
	if err := mustLoadConfig(ctx, cfg); err != nil {
		return err
	}

	applicationConfig = cfg

	contextInjectors.Register(func(ctx context.Context) context.Context {
		return config.ContextWithConfig(ctx, applicationConfig)
	})

	return nil
}

func Config() *config.Config {
	return applicationConfig
}

func OverrideConfig(override map[string]string) {
	for k, v := range override {
		overrideConfig[k] = v
	}
}
