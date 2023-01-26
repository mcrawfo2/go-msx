// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package cache

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"time"
)

type ContextCacheConfig struct {
	Ttl    time.Duration `config:"default=300s"`
	Prefix string        `config:"default=rc:"`
}

func NewContextCacheConfig(cfg *config.Config, root string) (*ContextCacheConfig, error) {
	var cacheConfig ContextCacheConfig
	if err := cfg.Populate(&cacheConfig, root); err != nil {
		return nil, err
	}
	return &cacheConfig, nil
}
