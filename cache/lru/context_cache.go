// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package lru

import (
	"context"
	"github.com/pkg/errors"
)

type ContextCache interface {
	Get(ctx context.Context, key string) (any, bool, error)
	Set(ctx context.Context, key string, value any) error
	Clear(ctx context.Context) error
}

type CacheProviderFactory func(ctx context.Context, configRoot string) (ContextCache, error)

var cacheProviderFactoryRegistry = make(map[string]CacheProviderFactory)

func RegisterCacheProvider(providerName string, factory CacheProviderFactory) {
	cacheProviderFactoryRegistry[providerName] = factory
}

func NewContextCache(ctx context.Context, providerName string, configRoot string) (ContextCache, error) {
	factory, ok := cacheProviderFactoryRegistry[providerName]
	if !ok {
		return nil, errors.Errorf("No such provider registered: %s", providerName)
	}

	return factory(ctx, configRoot)
}

type ContextCacheAdapter struct {
	Lru Cache
}

func (r ContextCacheAdapter) Get(ctx context.Context, key string) (any, bool, error) {
	val, exists := r.Lru.Get(key)
	return val, exists, nil
}

func (r ContextCacheAdapter) Set(ctx context.Context, key string, value any) (err error) {
	r.Lru.Set(key, value)
	return nil
}

func (r ContextCacheAdapter) Clear(ctx context.Context) (err error) {
	r.Lru.Clear()
	return nil
}

