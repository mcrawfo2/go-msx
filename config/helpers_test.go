package config

import "context"

func NewInMemoryConfig(values map[string]string) *Config {
	provider := NewInMemoryProvider("testdata", values)
	cfg := NewConfig(provider)
	_ = cfg.Load(context.Background())
	return cfg
}
