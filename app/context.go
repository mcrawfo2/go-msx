// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/consul"
	"cto-github.cisco.com/NFV-BU/go-msx/fs"
	"cto-github.cisco.com/NFV-BU/go-msx/health"
	"cto-github.cisco.com/NFV-BU/go-msx/health/consulcheck"
	"cto-github.cisco.com/NFV-BU/go-msx/health/kafkacheck"
	"cto-github.cisco.com/NFV-BU/go-msx/health/redischeck"
	"cto-github.cisco.com/NFV-BU/go-msx/health/sqldbcheck"
	"cto-github.cisco.com/NFV-BU/go-msx/health/vaultcheck"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/kafka"
	"cto-github.cisco.com/NFV-BU/go-msx/redis"
	"cto-github.cisco.com/NFV-BU/go-msx/sqldb"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/vault"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/pkg/errors"
)

var contextInjectors = new(types.ContextInjectors)

func init() {
	OnEvent(EventConfigure, PhaseAfter, configureHttpClientFactory)
	OnEvent(EventConfigure, PhaseAfter, withConfig(configureConsulPool))
	OnEvent(EventConfigure, PhaseAfter, configureVaultPool)
	OnEvent(EventConfigure, PhaseAfter, configureSqlDbPool)
	OnEvent(EventConfigure, PhaseAfter, configureSqlDbCrudRepositoryFactory)
	OnEvent(EventConfigure, PhaseAfter, configureRedisPool)
	OnEvent(EventConfigure, PhaseAfter, configureKafkaPool)
	OnEvent(EventConfigure, PhaseAfter, withConfig(fs.ConfigureFileSystem))
	OnEvent(EventConfigure, PhaseAfter, configureWebService)
}

func RegisterContextInjector(injector types.ContextInjector) {
	contextInjectors.Register(injector)
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
	RegisterContextInjector(func(ctx context.Context) context.Context {
		return httpclient.ContextWithFactory(ctx, httpClientFactory)
	})
	return nil
}

func configureSqlDbCrudRepositoryFactory(context.Context) error {
	crudRepositoryFactory := sqldb.NewProductionCrudRepositoryFactory()
	RegisterContextInjector(func(ctx context.Context) context.Context {
		return sqldb.ContextWithCrudRepositoryFactory(ctx, crudRepositoryFactory)
	})
	return nil
}

func configureConsulPool(cfg *config.Config) error {
	if err := consul.ConfigurePool(cfg); err != nil && err != consul.ErrDisabled {
		return err
	} else if err == nil {
		RegisterContextInjector(consul.ContextWithPool)
		health.RegisterCheck("consul", consulcheck.Check)
	}

	return nil
}

func configureVaultPool(ctx context.Context) error {
	if err := vault.ConfigurePool(ctx); err != nil && err != vault.ErrDisabled {
		return err
	} else if err == nil {
		RegisterContextInjector(vault.ContextWithPool)
		RegisterContextInjector(func(ctx context.Context) context.Context {
			return vault.ContextWithConnection(ctx, vault.PoolFromContext(ctx).Connection())
		})
		health.RegisterCheck("vault", vaultcheck.Check)
	}

	return nil
}

func configureSqlDbPool(ctx context.Context) error {
	if err := sqldb.ConfigurePool(ctx); err != nil && err != sqldb.ErrDisabled {
		return err
	} else if err == nil {
		RegisterContextInjector(sqldb.ContextWithPool)
		health.RegisterCheck("sqldb", sqldbcheck.Check)
	}

	return nil
}

func configureRedisPool(ctx context.Context) error {
	if err := redis.ConfigurePool(ctx); err != nil && err != redis.ErrDisabled {
		return err
	} else if err == nil {
		RegisterContextInjector(redis.ContextWithPool)
		health.RegisterCheck("redis", redischeck.Check)
	}

	return nil
}

func configureKafkaPool(ctx context.Context) error {
	if err := kafka.ConfigurePool(ctx); err != nil && err != kafka.ErrDisabled {
		return err
	} else if err == nil {
		RegisterContextInjector(kafka.ContextWithPool)
		health.RegisterCheck("kafka", kafkacheck.Check)
	}

	return nil
}

func configureWebService(ctx context.Context) error {
	return withConfig(func(cfg *config.Config) error {
		if err := webservice.ConfigureWebServer(cfg, ctx); err != nil && err != webservice.ErrDisabled {
			return err
		} else if err == nil {
			RegisterContextInjector(webservice.ContextWithWebServer)
		} else {
			logger.Warn(err.Error())
		}

		return nil
	})(ctx)
}
