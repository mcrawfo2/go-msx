package config

import (
	"context"
	"github.com/pkg/errors"
)

type configContextKey int

const contextKeyConfig configContextKey = iota

func ContextWithConfig(ctx context.Context, cfg *Config) context.Context {
	return context.WithValue(ctx, contextKeyConfig, cfg)
}

func FromContext(ctx context.Context) *Config {
	configInterface := ctx.Value(contextKeyConfig)
	if configInterface == nil {
		return nil
	}
	if cfg, ok := configInterface.(*Config); !ok {
		logger.Warn("Context config value wrong type")
		return nil
	} else {
		return cfg
	}
}

func MustFromContext(ctx context.Context) *Config {
	cfg := FromContext(ctx)
	if cfg == nil {
		panic(errors.New("Config missing from context"))
	}
	return cfg
}
