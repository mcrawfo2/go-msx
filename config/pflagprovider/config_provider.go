package pflagprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/config/args"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"sync"
)

var logger = log.NewLogger("msx.config.pflagprovider")

type ConfigProvider struct {
	name    string
	prefix  string
	extras  map[string]string
	flagset *pflag.FlagSet
	once    sync.Once
}

func (f *ConfigProvider) Description() string {
	return fmt.Sprintf("%s: pflag", f.name)
}

func (f *ConfigProvider) Load(ctx context.Context) (settings map[string]string, err error) {
	logger.Info("Loading command line config")

	f.once.Do(func() {
		f.extras = args.Extras(func(name string) bool {
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

func NewPflagSource(name string, flagset *pflag.FlagSet, prefix string) *ConfigProvider {
	return &ConfigProvider{
		name:    name,
		prefix:  prefix,
		flagset: flagset,
	}
}
