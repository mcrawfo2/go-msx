package config

import "context"

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
