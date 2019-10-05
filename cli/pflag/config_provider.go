package pflag

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cli/extract"
	msxConfig "cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/support/config"
	"cto-github.cisco.com/NFV-BU/go-msx/support/log"
	"github.com/spf13/pflag"
	"os"
	"sync"
)

var logger = log.NewLogger("msx.cli.pflag")

type ConfigProvider struct {
	prefix  string
	extras  map[string]string
	flagset *pflag.FlagSet
	once    sync.Once
}

func (f *ConfigProvider) Load(ctx context.Context) (settings map[string]string, err error) {
	logger.Info("Loading command line config")

	f.once.Do(func() {
		f.extras = extract.Extras(func(name string) bool {
			return f.flagset.Lookup(name) != nil
		})
	})

	if !f.flagset.Parsed() {
		if err := f.flagset.Parse(os.Args[1:]); err != nil {
			return nil, err
		}
	}

	settings = make(map[string]string)
	f.flagset.VisitAll(func(flag *pflag.Flag) {
		key := config.NormalizeKey(f.prefix + flag.Name)
		settings[key] = flag.Value.String()
	})

	// Apply extras
	for k, v := range f.extras {
		settings[k] = v
	}

	return settings, nil
}

func NewPflagSource(flagset *pflag.FlagSet, prefix string) *ConfigProvider {
	return &ConfigProvider{
		prefix:  prefix,
		flagset: flagset,
	}
}

func RegisterConfigProvider(flagset *pflag.FlagSet, prefix string) {
	if flagset == nil {
		logger.Warning("Invalid CLI flag set.")
	} else {
		msxConfig.RegisterProviderFactory(msxConfig.SourceCommandLine, func(*config.Config) config.Provider {
			return config.NewOnceLoader(NewPflagSource(flagset, prefix))
		})
	}
}
