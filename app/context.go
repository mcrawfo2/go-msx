package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/consul"
	"cto-github.cisco.com/NFV-BU/go-msx/health"
	"cto-github.cisco.com/NFV-BU/go-msx/health/cassandracheck"
	"cto-github.cisco.com/NFV-BU/go-msx/health/consulcheck"
	"cto-github.cisco.com/NFV-BU/go-msx/health/kafkacheck"
	"cto-github.cisco.com/NFV-BU/go-msx/health/redischeck"
	"cto-github.cisco.com/NFV-BU/go-msx/health/vaultcheck"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/kafka"
	"cto-github.cisco.com/NFV-BU/go-msx/redis"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/vault"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/pkg/errors"
)

var contextInjectors = new(types.ContextInjectors)

func init() {
	OnEvent(EventConfigure, PhaseAfter, configureHttpClientFactory)
	OnEvent(EventConfigure, PhaseAfter, withConfig(configureConsulPool))
	OnEvent(EventConfigure, PhaseAfter, withConfig(configureVaultPool))
	OnEvent(EventConfigure, PhaseAfter, withConfig(configureCassandraPool))
	OnEvent(EventConfigure, PhaseAfter, configureCassandraCrudRepositoryFactory)
	OnEvent(EventConfigure, PhaseAfter, withConfig(configureRedisPool))
	OnEvent(EventConfigure, PhaseAfter, withConfig(configureKafkaPool))
	OnEvent(EventConfigure, PhaseAfter, configureWebService)
	OnEvent(EventCommand, CommandMigrate, func(ctx context.Context) error {
		// Only during migrate command
		OnEvent(EventConfigure, PhaseAfter, createCassandraKeyspace)
		return nil
	})
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

func configureHttpClientFactory(ctx context.Context) error {
	httpClientFactory, err := httpclient.NewProductionHttpClientFactory(ctx)
	if err != nil {
		return err
	}
	contextInjectors.Register(func(ctx context.Context) context.Context {
		return httpclient.ContextWithFactory(ctx, httpClientFactory)
	})
	return nil
}

func configureCassandraCrudRepositoryFactory(context.Context) error {
	crudRepositoryFactory := cassandra.NewProductionCrudRepositoryFactory()
	contextInjectors.Register(func(ctx context.Context) context.Context {
		return cassandra.ContextWithCrudRepositoryFactory(ctx, crudRepositoryFactory)
	})
	return nil
}

func configureConsulPool(cfg *config.Config) error {
	if err := consul.ConfigurePool(cfg); err != nil && err != consul.ErrDisabled {
		return err
	} else if err != consul.ErrDisabled {
		contextInjectors.Register(consul.ContextWithPool)
		health.RegisterCheck("consul", consulcheck.Check)
	}

	return nil
}

func configureVaultPool(cfg *config.Config) error {
	if err := vault.ConfigurePool(cfg); err != nil && err != vault.ErrDisabled {
		return err
	} else if err != vault.ErrDisabled {
		contextInjectors.Register(vault.ContextWithPool)
		health.RegisterCheck("vault", vaultcheck.Check)
	}

	return nil
}

func configureCassandraPool(cfg *config.Config) error {
	if err := cassandra.ConfigurePool(cfg); err != nil && err != cassandra.ErrDisabled {
		return err
	} else if err != cassandra.ErrDisabled {
		contextInjectors.Register(cassandra.ContextWithPool)
		health.RegisterCheck("cassandra", cassandracheck.Check)
	}

	return nil
}

func configureRedisPool(cfg *config.Config) error {
	if err := redis.ConfigurePool(cfg); err != nil && err != redis.ErrDisabled {
		return err
	} else if err != redis.ErrDisabled {
		contextInjectors.Register(redis.ContextWithPool)
		health.RegisterCheck("redis", redischeck.Check)
	}

	return nil
}

func configureKafkaPool(cfg *config.Config) error {
	if err := kafka.ConfigurePool(cfg); err != nil && err != kafka.ErrDisabled {
		return err
	} else if err != kafka.ErrDisabled {
		contextInjectors.Register(kafka.ContextWithPool)
		health.RegisterCheck("kafka", kafkacheck.Check)
	}

	return nil
}

func configureWebService(ctx context.Context) error {
	return withConfig(func(cfg *config.Config) error {
		if err := webservice.ConfigureWebServer(cfg, ctx); err != nil && err != webservice.ErrDisabled {
			return err
		} else if err == nil {
			contextInjectors.Register(webservice.ContextWithWebServer)
		} else {
			logger.Warn(err.Error())
		}

		return nil
	})(ctx)
}

func createCassandraKeyspace(ctx context.Context) error {
	return cassandra.CreateKeyspaceForPool(ctx)
}
