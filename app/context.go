package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/consul"
	"cto-github.cisco.com/NFV-BU/go-msx/vault"
	"github.com/pkg/errors"
)

func init() {
	OnEvent(EventConfigure, PhaseAfter, withConfig(configureConsulPool))
	OnEvent(EventConfigure, PhaseAfter, withConfig(configureVaultPool))
	OnEvent(EventConfigure, PhaseAfter, withConfig(configureCassandraPool))
}

type configHandler func(cfg *config.Config) error

func withConfig(handler configHandler) Observer {
	return func(ctx context.Context) error {
		var cfg *config.Config
		if cfg = config.FromContext(ctx); cfg == nil {
			return errors.New("Failed to retrieve config from context")
		}

		return handler(cfg)
	}
}

func configureConsulPool(cfg *config.Config) error {
	if err := consul.ConfigurePool(cfg); err != nil && err != consul.ErrDisabled {
		return err
	} else if err != consul.ErrDisabled {
		RegisterInjector(consul.ContextWithPool)
	}

	return nil
}

func configureVaultPool(cfg *config.Config) error {
	if err := vault.ConfigurePool(cfg); err != nil && err != vault.ErrDisabled {
		return err
	} else if err != vault.ErrDisabled {
		RegisterInjector(vault.ContextWithPool)
	}

	return nil
}

func configureCassandraPool(cfg *config.Config) error {
	if err := cassandra.ConfigurePool(cfg); err != nil && err != cassandra.ErrDisabled {
		return err
	} else if err != cassandra.ErrDisabled {
		RegisterInjector(cassandra.ContextWithPool)
	}

	return nil
}

type ContextInjector func(ctx context.Context) context.Context

var contextInjectors []ContextInjector

func RegisterInjector(injector ContextInjector) {
	contextInjectors = append(contextInjectors, injector)
}

func injectContextValues(ctx context.Context) context.Context {
	for _, contextInjector := range contextInjectors {
		ctx = contextInjector(ctx)
	}
	return ctx
}
