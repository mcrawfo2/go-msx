// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package config

import "context"

type Provided interface {
	Provider() Provider
}

type Cache struct {
	provider Provider
	entries  ProviderEntries
	cached   bool
}

func (c *Cache) Provider() Provider {
	if p, ok := c.provider.(Provided); ok {
		return p.Provider()
	}
	return c.provider
}

func (c *Cache) Description() string {
	return c.provider.Description()
}

func (c *Cache) Run(ctx context.Context) {
	c.provider.Run(ctx)
}

func (c *Cache) Notify() <-chan struct{} {
	return c.provider.Notify()
}

func (c *Cache) Load(ctx context.Context) (ProviderEntries, error) {
	if c.cached {
		return c.entries, nil
	}

	entries, err := c.provider.Load(ctx)
	if err != nil {
		return nil, err
	}

	if len(entries) > 0 {
		logger.Infof("Retrieved %d configs from %s", len(entries), c.Description())
	} else {
		logger.Warningf("No configs retrieved from %s", c.Description())
	}

	entries.SortByNormalizedName()

	c.cached = true
	c.entries = entries
	return entries, nil
}

func (c *Cache) Invalidate() {
	c.cached = false
}

func NewCacheProvider(provider Provider) *Cache {
	return &Cache{
		provider: provider,
	}
}
