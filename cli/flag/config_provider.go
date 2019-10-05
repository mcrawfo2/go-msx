package flag

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cli/extract"
	msxConfig "cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/support/config"
	"cto-github.cisco.com/NFV-BU/go-msx/support/log"
	"flag"
	"os"
	"sync"
)

var logger = log.NewLogger("msx.cli.flag")

type ConfigProvider struct {
	prefix  string
	extras  map[string]string
	flagSet *flag.FlagSet
	once    sync.Once
}

func (f *ConfigProvider) Load(ctx context.Context) (settings map[string]string, err error) {
	logger.Info("Loading command line config")

	f.once.Do(func() {
		f.extras = extract.Extras(func(name string) bool {
			return f.flagSet.Lookup(name) != nil
		})
	})

	if !f.flagSet.Parsed() {
		if err := f.flagSet.Parse(os.Args[1:]); err != nil {
			return nil, err
		}
	}

	settings = make(map[string]string)
	f.flagSet.VisitAll(func(flag *flag.Flag) {
		key := config.NormalizeKey(f.prefix + flag.Name)
		settings[key] = flag.Value.String()
	})

	// Apply extras
	for k, v := range f.extras {
		settings[k] = v
	}

	return settings, nil
}

func NewFlagSource(flagSet *flag.FlagSet, prefix string) *ConfigProvider {
	return &ConfigProvider{
		prefix:  prefix,
		flagSet: flagSet,
	}
}

func RegisterConfigProvider(flagSet *flag.FlagSet, prefix string) {
	if flagSet == nil {
		logger.Warning("Invalid CLI flag set.")
	} else {
		msxConfig.RegisterProviderFactory(msxConfig.SourceCommandLine, func(*config.Config) config.Provider {
			return config.NewOnceLoader(NewFlagSource(flagSet, prefix))
		})
	}
}
