package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config/consulprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/config/vaultprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/support/config"
	"fmt"
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

type ProviderFactory func(*config.Config) config.Provider

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
	if configFile, err := findConfigFile("bootstrap"); err != nil {
		logger.Warn("Failed to load bootstrap config file: ", err.Error())
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

func newApplicationProvider(*config.Config) config.Provider {
	if configFile, err := findConfigFile("application"); err != nil {
		logger.Warn("Failed to load application config file: ", err.Error())
		return nil
	} else {
		return newFileProvider(configFile)
	}
}

func newBuildProvider() config.Provider {
	if configFile, err := findConfigFile("build"); err != nil {
		logger.Warn("Failed to load build info file: ", err.Error())
		return nil
	} else {
		return newFileProvider(configFile)
	}
}

func newProfileProvider(config *config.Config) config.Provider {
	var parts []string
	if appName, err := config.String("spring.application.name"); err != nil {
		logger.Warn("Application name not found: ", err)
		return nil
	} else {
		parts = []string{appName}
	}

	if profile, err := config.StringOr("profile", "default"); err != nil {
		logger.Warn(err)
		return nil
	} else if profile == "default" {
		// don't add it
	} else {
		parts = append(parts, profile)
	}

	if configFile, err := findConfigFile(strings.Join(parts, ".")); err != nil {
		logger.Warn("Failed to locate profile configuration: ", err.Error())
		return nil
	} else {
		return newFileProvider(configFile)
	}
}

func newProvider(name string, cfg *config.Config) config.Provider {
	if providerFactories[name] == nil {
		return nil
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
	case ".json":
		return config.NewCachedLoader(config.NewJSONFile(fileName))
	case ".properties":
		return config.NewCachedLoader(config.NewPropertiesFile(fileName))
	default:
		logger.Error("Unknown config file extension: ", fileExt)
		return nil
	}
}

func findConfigFile(baseName string) (string, error) {
	extensions := []string{".yaml", ".yml", ".ini", ".json", ".properties"}
	for _, ext := range extensions {
		fullName := baseName + ext
		info, err := os.Stat(fullName)
		if os.IsNotExist(err) || info.IsDir() {
			continue
		}

		return fullName, nil
	}

	return "", errors.New(fmt.Sprintf("Could not find %s.{yaml,yml,ini,json,properties}", baseName))
}

func init() {
	OnEvent(EventConfigure, PhaseBefore, registerRemoteConfigProviders)
	OnEvent(EventConfigure, PhaseDuring, loadConfig)
	OnEvent(EventStart, PhaseAfter, watchConfig)
}

func registerRemoteConfigProviders() {
	RegisterProviderFactory(SourceConsul, consulprovider.NewConfigProviderFromConfig)
	RegisterProviderFactory(SourceVault, vaultprovider.NewConfigProviderFromConfig)
}

func watchConfig() {
	logger.Info("Watching configuration for changes")
	cfg := Config()
	cfg.Notify = func(keys []string) {
		for _, k := range keys {
			logger.Warnf("Configuration changed: %s", k)
		}
	}
	cfg.Watch(Context())
}

func mustLoadConfig(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(Context(), time.Second*15)
	defer cancel()

	var err error
	if err = cfg.Load(ctx); err != nil {
		Shutdown()
	}

	return err
}

func loadConfig() {
	sources := &Sources{
		Defaults:      newDefaultsProvider(),
		BootstrapFile: newBootstrapProvider(),
		BuildFile:     newBuildProvider(),
		Environment:   newEnvironmentProvider(),
		CommandLine:   newProvider(SourceCommandLine, nil),
		Override:      newOverrideProvider(overrideConfig),
	}

	config1 := config.NewConfig(sources.Providers()...)
	if err := mustLoadConfig(config1); err != nil {
		logger.Error(err)
		return
	}

	sources.ApplicationFile = newProvider(SourceApplication, config1)
	sources.ProfileFile = newProvider(SourceProfile, config1)
	config2 := config.NewConfig(sources.Providers()...)
	if err := mustLoadConfig(config2); err != nil {
		logger.Error(err)
		return
	}

	sources.Consul = newProvider(SourceConsul, config2)
	sources.Vault = newProvider(SourceVault, config2)
	applicationConfig = config.NewConfig(sources.Providers()...)
	if err := mustLoadConfig(applicationConfig); err != nil {
		logger.Error(err)
		return
	}
}

func Config() *config.Config {
	return applicationConfig
}

func SetOverrideConfig(override map[string]string) {
	if override != nil {
		overrideConfig = override
	}
}
