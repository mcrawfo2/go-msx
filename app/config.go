package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/config/consulprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/config/vaultprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/resource"
	"github.com/pkg/errors"
	"path"
	"strings"
	"time"
)

const (
	SourceConsul           = "Consul"
	SourceVault            = "Vault"
	SourceApplication      = "Application"
	SourceProfile          = "Profile"
	SourceCommandLine      = "CommandLine"
	SourceBuild            = "Build"
	SourceBootStrap        = "Bootstrap"
	SourceEnvironment      = "Environment"
	SourceOverride         = "Override"
	SourceDefaultResources = "DefaultResources"
	SourceDefaults         = "Defaults"
	SourceDefault          = "Default"

	configKeyAppName     = "spring.application.name"
	configKeyAppInstance = "spring.application.instance"
	configKeyFsConfigs   = "fs.configs"
)

var (
	applicationConfig *config.Config
	overrideConfig    = make(map[string]string)
)

type Sources struct {
	Default           config.Provider
	DefaultsFiles     []config.Provider
	DefaultsResources []config.Provider
	BootstrapFiles    []config.Provider
	ApplicationFiles  []config.Provider
	BuildFiles        []config.Provider
	Consul            []config.Provider
	Vault             []config.Provider
	ProfileFiles      []config.Provider
	Environment       config.Provider
	CommandLine       []config.Provider
	Override          config.Provider
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
	sourcesList.Append(c.DefaultsResources...)
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
	return config.NewCachedLoader(config.Defaults)
}

func newDefaultsResourcesProviders() []config.Provider {
	return config.NewHttpFileProvidersFromGlob(SourceDefaultResources, resource.Defaults, "**/defaults-*")
}

func newDefaultsFilesProviders(_ *config.Config) []config.Provider {
	return config.NewFileProvidersFromGlob(SourceDefaults, "defaults-*")
}

func newBootstrapProviders(_ *config.Config) []config.Provider {
	return config.NewFileProvidersFromBaseName(SourceBootStrap, "bootstrap")
}

func newEnvironmentProvider() config.Provider {
	return config.NewCachedLoader(config.NewEnvironment(SourceEnvironment))
}

func newOverrideProvider(static map[string]string) config.Provider {
	return config.NewCachedLoader(config.NewStatic(SourceOverride, static))
}

func newApplicationProviders(name string, cfg *config.Config) ([]config.Provider, error) {
	providers := config.NewFileProvidersFromBaseName(name, "application")
	if appName, err := cfg.String(configKeyAppName); err == nil {
		namedProviders := config.NewFileProvidersFromBaseName(name, appName)
		providers = append(providers, namedProviders...)
	}
	return providers, nil
}

func newBuildProvider(_ *config.Config) []config.Provider {
	return config.NewFileProvidersFromBaseName(SourceBuild, "buildinfo")
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

	return config.NewFileProvidersFromBaseName(SourceProfile, strings.Join(parts, ".")), nil
}

func newProviders(name string, cfg *config.Config) ([]config.Provider, error) {
	if providerFactories[name] == nil {
		return nil, nil
	}
	providerFactory := providerFactories[name]
	return providerFactory(name, cfg)
}

func init() {
	OnEvent(EventConfigure, PhaseBefore, registerRemoteConfigProviders)
	OnEvent(EventConfigure, PhaseDuring, loadConfig)
	OnEvent(EventConfigure, PhaseDuring, applyLoggingConfig)
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
			logger.WithError(err).Error("Failed to refresh")
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
		Default:           newDefaultsProvider(),
		DefaultsResources: newDefaultsResourcesProviders(),
		Environment:       newEnvironmentProvider(),
		Override:          newOverrideProvider(overrideConfig),
	}

	if sources.CommandLine, err = newProviders(SourceCommandLine, nil); err != nil {
		return err
	}

	var cfg = config.NewConfig(sources.Providers()...)
	if err := mustLoadConfig(ctx, cfg); err != nil {
		return err
	}

	// Add any config paths from the command line
	config.AddConfigFoldersFromPathConfig(cfg)
	// Support single-folder installs
	config.AddConfigFolders("./etc")
	// Support global installs
	appName, err := cfg.String(configKeyAppName)
	if err == nil && appName != "" {
		config.AddConfigFolders(path.Join("/etc", appName))
	}
	// Allow the build to specify the config folder in fs.configs
	sources.BuildFiles = newBuildProvider(cfg)

	// Reload the config
	cfg = config.NewConfig(sources.Providers()...)
	if err := mustLoadConfig(ctx, cfg); err != nil {
		return err
	}

	fsConfigs, err := cfg.String(configKeyFsConfigs)
	if err == nil && fsConfigs != "" {
		config.AddConfigFolders(fsConfigs)
	}

	logger.WithContext(ctx).Infof("Config Search Path: %v", config.ConfigFolders())

	sources.DefaultsFiles = newDefaultsFilesProviders(cfg)
	sources.BootstrapFiles = newBootstrapProviders(cfg)

	if sources.ApplicationFiles, err = newProviders(SourceApplication, cfg); err != nil {
		return err
	}

	if sources.ProfileFiles, err = newProviders(SourceProfile, cfg); err != nil {
		return err
	}

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

func applyLoggingConfig(ctx context.Context) error {
	settings := config.FromContext(ctx).Settings()
	prefix := "logger."
	n := len(prefix)
	for k, v := range settings {
		if len(k) <= len(prefix) || !strings.HasPrefix(k, prefix) {
			continue
		}
		loggerName := k[n:]
		loggerLevel := log.LevelFromName(strings.ToUpper(v))
		log.SetLoggerLevel(loggerName, loggerLevel)
	}

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
