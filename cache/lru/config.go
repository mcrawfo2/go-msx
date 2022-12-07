// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package lru

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"time"
)

type CacheConfig struct {
	Ttl             time.Duration `config:"default=300s"`
	ExpireLimit     int           `config:"default=100"`
	ExpireFrequency time.Duration `config:"default=30s"`
	DeAgeOnAccess   bool          `config:"default=false"`
	Metrics         bool          `config:"default=false"`
	MetricsPrefix   string        `config:"default=cache"`
}

func NewCacheConfig(cfg *config.Config, root string) (*CacheConfig, error) {
	var cacheConfig CacheConfig
	if err := cfg.Populate(&cacheConfig, root); err != nil {
		return nil, err
	}
	return &cacheConfig, nil
}
