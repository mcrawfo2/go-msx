package configtest

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
)

func NewStaticConfig(values map[string]string) *config.Config {
	provider := config.NewStatic("static", values)
	cfg := config.NewConfig(provider)
	_ = cfg.Load(context.Background())
	return cfg
}

func ContextWithNewStaticConfig(ctx context.Context, values map[string]string) context.Context {
	cfg := NewStaticConfig(values)
	return config.ContextWithConfig(ctx, cfg)
}