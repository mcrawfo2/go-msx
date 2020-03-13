package lru

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"time"
)

type CacheConfig struct {
	Ttl             time.Duration `config:"default=300s"`
	ExpireLimit     int           `config:"default=100"`
	ExpireFrequency time.Duration `config:"default=30s"`
}

func NewCacheConfig(cfg *config.Config, root string) (*CacheConfig, error) {
	var cacheConfig CacheConfig
	if err := cfg.Populate(&cacheConfig, root); err != nil {
		return nil, err
	}
	return &cacheConfig, nil
}
