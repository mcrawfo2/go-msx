// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cache/lru"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	redisCache "cto-github.cisco.com/NFV-BU/go-msx/redis/cache"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/adminprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/aliveprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/apilistprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/asyncapiprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/authprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/debugprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/envprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/healthprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/idempotency"
	idempotencyCache "cto-github.cisco.com/NFV-BU/go-msx/webservice/idempotency/cache"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/infoprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/loggersprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/maintenanceprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/metricsprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/swaggerprovider"

	_ "cto-github.cisco.com/NFV-BU/go-msx/ops/restops/httperrors"
)

func init() {
	OnEvent(EventConfigure, PhaseAfter, withConfig(registerRegistrations))
}

func registerRegistrations(cfg *config.Config) error {
	serverEnabled, err := cfg.BoolOr("server.enabled", true)
	if err != nil {
		return err
	}

	if serverEnabled {
		OnEvent(EventStart, PhaseBefore, registerAuthenticationProvider)
		OnEvent(EventStart, PhaseBefore, registerAdminWebServices)
		OnEvent(EventStart, PhaseBefore, registerDebugWebServices)
		OnEvent(EventStart, PhaseBefore, registerSwaggerWebService)
		OnEvent(EventStart, PhaseBefore, registerApiListWebService)
		OnEvent(EventStart, PhaseBefore, registerAsyncApiWebService)

		OnEvent(EventStart, PhaseAfter, webservice.Start)

		OnEvent(EventStart, PhaseAfter, registerIdempotencyCacheRedis)
		OnEvent(EventStart, PhaseAfter, registerIdempotencyCacheInMemory)
		OnEvent(EventStart, PhaseAfter, idempotency.ApplyIdempotencyKeyFilter)

		OnEvent(EventStop, PhaseBefore, webservice.Stop)
	}

	return nil
}

func registerAuthenticationProvider(ctx context.Context) error {
	logger.Info("Registering resource path glob security provider")
	return authprovider.RegisterAuthenticationProvider(ctx)
}

func registerAdminWebServices(ctx context.Context) error {
	logger.Info("Registering admin endpoints")
	err := types.ErrorList{
		adminprovider.RegisterProvider(ctx),
		healthprovider.RegisterProvider(ctx),
		infoprovider.RegisterProvider(ctx),
		aliveprovider.RegisterProvider(ctx),
		metricsprovider.RegisterProvider(ctx),
		envprovider.RegisterProvider(ctx),
		loggersprovider.RegisterProvider(ctx),
		maintenanceprovider.RegisterProvider(ctx),
	}
	return err.Filter()
}

func registerDebugWebServices(ctx context.Context) error {
	debugEnabled, _ := config.FromContext(ctx).BoolOr("server.debug-enabled", false)
	if !debugEnabled {
		logger.Info("Debug endpoints disabled")
		return nil
	}

	logger.Info("Registering debug endpoints")
	return types.ErrorList{
		debugprovider.RegisterProvider(ctx),
	}.Filter()
}

func registerSwaggerWebService(ctx context.Context) error {
	logger.Info("Registering swagger documentation provider")
	if err := swaggerprovider.RegisterSwaggerProvider(ctx); err != nil && err != swaggerprovider.ErrDisabled {
		return err
	} else if err == swaggerprovider.ErrDisabled {
		logger.Info("Swagger documentation provider disabled")
	}

	return nil
}

func registerApiListWebService(ctx context.Context) error {
	logger.Info("Registering apilist documentation provider")
	if err := apilistprovider.RegisterProvider(ctx); err != nil {
		return err
	}

	return nil
}

func registerAsyncApiWebService(ctx context.Context) error {
	logger.WithContext(ctx).Info("Registering AsyncApi documentation provider")
	err := asyncapiprovider.RegisterProvider(ctx)

	switch err {
	case asyncapiprovider.ErrDisabled:
		logger.Info("AsyncApi documentation provider disabled")
		return nil
	default:
		return err
	}
}

func registerIdempotencyCacheRedis(ctx context.Context) error {
	cfg := config.MustFromContext(ctx)

	lru.RegisterCacheProvider(idempotency.CacheProviderInMemory, func(ctx context.Context, configRoot string) (lru.ContextCache, error) {
		lruConfig, err := lru.NewCacheConfig(cfg, idempotency.ConfigRootIdempotencyKeyInMemory)
		if err != nil {
			logger.WithContext(ctx).Error(err)
			return nil, err
		}
		return lru.ContextCacheAdapter{Lru: lru.NewCacheFromConfig(lruConfig)}, nil
	})

	return nil
}

func registerIdempotencyCacheInMemory(ctx context.Context) error {
	cfg := config.MustFromContext(ctx)

	lru.RegisterCacheProvider(idempotency.CacheProviderRedis, func(ctx context.Context, configRoot string) (lru.ContextCache, error) {
		redisConfig, err := redisCache.NewContextCacheConfig(cfg, idempotency.ConfigRootIdempotencyKeyRedis)
		if err != nil {
			logger.WithContext(ctx).Error(err)
			return nil, err
		}
		return redisCache.NewContextCacheFromConfig[idempotencyCache.CachedWebData](redisConfig), nil
	})

	return nil
}


