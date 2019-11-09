package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/config/consulprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/config/vaultprovider"
	"github.com/pkg/errors"
	"os"
	"path"
	"strings"
	"time"
)

const (
	SourceConsul      = "Consul"
	SourceVault       = "Vault"
	SourceApplication = "Application"
	SourceProfile     = "Profile"
	SourceCommandLine = "CommandLine"
)

var (
	applicationConfig *config.Config
	overrideConfig    = make(map[string]string)
)

type Sources struct {
	Defaults        config.Provider
	BootstrapFile   config.Provider
	ApplicationFile config.Provider
	BuildFile       config.Provider
	Consul          config.Provider
	Vault           config.Provider
	ProfileFile     config.Provider
	Environment     config.Provider
	CommandLine     config.Provider
	Override        config.Provider
}

func (c Sources) Providers() []config.Provider {
	var providers = []config.Provider{
		c.Defaults,
		c.BootstrapFile,
		c.ApplicationFile,
		c.BuildFile,
		c.Consul,
		c.Vault,
		c.ProfileFile,
		c.Environment,
		c.CommandLine,
		c.Override,
	}

	for i := 0; i < len(providers); i++ {
		if providers[i] == nil {
			providers = append(providers[:i], providers[i+1:]...)
			i--
		}
	}

	return providers
}

type ProviderFactory func(*config.Config) (config.Provider, error)

var providerFactories = map[string]ProviderFactory{
	SourceCommandLine: nil,
	SourceConsul:      nil,
	SourceVault:       nil,
	SourceApplication: newApplicationProvider,
	SourceProfile:     newProfileProvider,
}

func RegisterProviderFactory(name string, factory ProviderFactory) {
	providerFactories[name] = factory
}

func newDefaultsProvider() config.Provider {
	return config.NewCachedLoader(config.NewStatic("default", map[string]string{
		"profile": "default",
	}))
}

func newBootstrapProvider() config.Provider {
	if configFile := findConfigFile("bootstrap"); configFile == "" {
		return nil
	} else {
		return newFileProvider(configFile)
	}
}

func newEnvironmentProvider() config.Provider {
	return config.NewCachedLoader(config.NewEnvironment())
}

func newOverrideProvider(static map[string]string) config.Provider {
	return config.NewCachedLoader(config.NewStatic("override", static))
}

func newApplicationProvider(*config.Config) (config.Provider, error) {
	if configFile := findConfigFile("application"); configFile == "" {
		return nil, nil
	} else {
		return newFileProvider(configFile), nil
	}
}

func newBuildProvider() config.Provider {
	if configFile := findConfigFile("build"); configFile == "" {
		return nil
	} else {
		return newFileProvider(configFile)
	}
}

func newProfileProvider(config *config.Config) (config.Provider, error) {
	var parts []string
	if appName, err := config.String("spring.application.name"); err != nil {
		return nil, errors.Wrap(err, "Application name not defined")
	} else {
		parts = []string{appName}
	}

	if profile, err := config.StringOr("profile", "default"); err != nil {
		return nil, err
	} else if profile == "default" {
		// don't add it
	} else {
		parts = append(parts, profile)
	}

	if configFile := findConfigFile(strings.Join(parts, ".")); configFile == "" {
		return nil, nil
	} else {
		return newFileProvider(configFile), nil
	}
}

func newProvider(name string, cfg *config.Config) (config.Provider, error) {
	if providerFactories[name] == nil {
		return nil, nil
	}
	providerFactory := providerFactories[name]
	return providerFactory(cfg)
}

func newFileProvider(fileName string) config.Provider {
	fileExt := strings.ToLower(path.Ext(fileName))
	switch fileExt {
	case ".yml", ".yaml":
		return config.NewCachedLoader(config.NewYAMLFile(fileName))
	case ".ini":
		return config.NewCachedLoader(config.NewINIFile(fileName))
	case ".json", ".json5":
		return config.NewCachedLoader(config.NewJSONFile(fileName))
	case ".properties":
		return config.NewCachedLoader(config.NewPropertiesFile(fileName))
	default:
		logger.Error("Unknown config file extension: ", fileExt)
		return nil
	}
}

func findConfigFile(baseName string) string {
	extensions := []string{".yaml", ".yml", ".ini", ".json", ".json5", ".properties"}
	for _, ext := range extensions {
		fullName := baseName + ext
		info, err := os.Stat(fullName)
		if os.IsNotExist(err) || info.IsDir() {
			continue
		}

		return fullName
	}

	logger.Warnf("Could not find %s.{yaml,yml,ini,json,json5,properties}", baseName)
	return ""
}

func init() {
	OnEvent(EventConfigure, PhaseBefore, registerRemoteConfigProviders)
	OnEvent(EventConfigure, PhaseDuring, loadConfig)
	OnEvent(EventStart, PhaseAfter, watchConfig)
}

func registerRemoteConfigProviders(ctx context.Context) error {
	RegisterProviderFactory(SourceConsul, consulprovider.NewConfigProviderFromConfig)
	RegisterProviderFactory(SourceVault, vaultprovider.NewConfigProviderFromConfig)
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
		Defaults:      newDefaultsProvider(),
		BootstrapFile: newBootstrapProvider(),
		BuildFile:     newBuildProvider(),
		Environment:   newEnvironmentProvider(),
		Override:      newOverrideProvider(overrideConfig),
	}

	if sources.CommandLine, err = newProvider(SourceCommandLine, nil); err != nil {
		return err
	}

	var cfg = config.NewConfig(sources.Providers()...)
	if err := mustLoadConfig(ctx, cfg); err != nil {
		return err
	}

	if sources.ApplicationFile, err = newProvider(SourceApplication, cfg); err != nil {
		return err
	}

	if sources.ProfileFile, err = newProvider(SourceProfile, cfg); err != nil {
		return err
	}

	cfg = config.NewConfig(sources.Providers()...)
	if err := mustLoadConfig(ctx, cfg); err != nil {
		return err
	}

	if sources.Consul, err = newProvider(SourceConsul, cfg); err != nil {
		return err
	}

	if sources.Vault, err = newProvider(SourceVault, cfg); err != nil {
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
