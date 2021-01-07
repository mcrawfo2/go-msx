package config

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

type Config struct {
	Values
	caches   []*Cache
	watchers []*Watcher
	layers   []ProviderEntries
	changes  chan WatcherNotification
	notify   chan Snapshot
	values   struct {
		original SnapshotValues // from Load
		latest   SnapshotValues // from Watch
	}
}

func (c *Config) Load(ctx context.Context) error {
	layers := make(Layers, len(c.caches))

	for i, cache := range c.caches {
		entries, err := cache.Load(ctx)
		if err != nil {
			return err
		}

		layers[i] = entries
	}

	merged := layers.Merge()

	resolver := NewResolver(merged)
	entries, err := resolver.Entries()
	if err != nil {
		return err
	}
	values := newSnapshotValues(entries)

	c.layers = layers
	c.values.original, c.values.latest = values, values

	return nil
}

func (c *Config) Watch(ctx context.Context) {
	logger.WithContext(ctx).Info("Watching configuration for changes")

	if len(c.watchers) > 0 || len(c.caches) == 0 {
		return
	}

	for _, cache := range c.caches {
		watcher := NewWatcher(cache)
		c.watchers = append(c.watchers, watcher)
		// *We* watch the watchers
		go c.watch(ctx, watcher)
	}

	for {
		select {
		case <-ctx.Done():
			logger.
				WithContext(ctx).
				WithError(ctx.Err()).
				Info("Configuration watcher context finished.  Exiting watcher.")
			return

		case change := <-c.changes:
			logger.Info("Creating new config snapshot")
			// resolve latest config
			layers := append(Layers{}, c.layers...)
			layerIndex := c.getLayerIndex(change)
			if layerIndex == -1 {
				logger.
					WithContext(ctx).
					Error("Configuration watcher received changes from unregistered provider.")
				continue
			}
			layers[layerIndex] = change.Entries
			merged := layers.Merge()

			snapshot, err := newSnapshot(merged, c.values.latest)
			if err != nil {
				logger.
					WithContext(ctx).
					WithError(err).
					Error("Config snapshot creation failed")
				continue
			}

			c.layers = layers
			c.values.latest = snapshot.SnapshotValues
			c.notify <- snapshot
			logger.
				WithContext(ctx).
				Info("Config snapshot created")

		}
	}
}

func (c *Config) getLayerIndex(change WatcherNotification) int {
	for i, cache := range c.caches {
		if cache.provider == change.Provider {
			return i
		}
	}
	return -1
}

func (c *Config) watch(ctx context.Context, watcher *Watcher) {
	go watcher.Watch(ctx)

	for {
		select {
		case <-ctx.Done():
			return

		case n, ok := <-watcher.Notify():
			if !ok {
				return
			}

			c.changes <- n
		}
	}
}

func (c Config) LatestValues() SnapshotValues {
	if c.layers == nil {
		return emptySnapshotValues
	}

	return c.values.latest
}

func (c Config) OriginalValues() SnapshotValues {
	if c.layers == nil {
		return emptySnapshotValues
	}

	return c.values.original
}

func (c Config) Notify() <-chan Snapshot {
	return c.notify
}

func (c Config) Caches() []Provider {
	var results []Provider
	for _, cache := range c.caches {
		results = append(results, cache)
	}
	return results
}

func NewConfig(providers ...Provider) *Config {
	var caches []*Cache
	for _, provider := range providers {
		var cache *Cache
		// Dont double-cache providers
		if c, ok := provider.(*Cache); ok {
			cache = c
		} else {
			cache = NewCacheProvider(provider)
		}
		caches = append(caches, cache)
	}

	cfg := &Config{
		caches:  caches,
		changes: make(chan WatcherNotification),
		notify:  make(chan Snapshot),
	}
	cfg.Values = configValues{cfg:cfg, original: true}
	return cfg
}

type configValues struct {
	cfg *Config
	original bool
}

func (c configValues) snapshot() (values SnapshotValues, err error) {
	if c.cfg.layers == nil {
		err = errors.Wrap(ErrNotLoaded, "Values not available")
		return emptySnapshotValues, err
	}

	if c.original {
		values = c.cfg.OriginalValues()
	} else {
		values = c.cfg.LatestValues()
	}

	return values, nil
}

func (c configValues) String(key string) (string, error) {
	if values, err := c.snapshot(); err != nil {
		return "", err
	} else {
		return values.String(key)
	}
}

func (c configValues) StringOr(key, alt string) (string, error) {
	if values, err := c.snapshot(); err != nil {
		return "", err
	} else {
		return values.StringOr(key, alt)
	}
}

func (c configValues) Int(key string) (int, error) {
	if values, err := c.snapshot(); err != nil {
		return 0, err
	} else {
		return values.Int(key)
	}
}

func (c configValues) IntOr(key string, alt int) (int, error) {
	if values, err := c.snapshot(); err != nil {
		return 0, err
	} else {
		return values.IntOr(key, alt)
	}
}

func (c configValues) Uint(key string) (uint, error) {
	if values, err := c.snapshot(); err != nil {
		return 0, err
	} else {
		return values.Uint(key)
	}
}

func (c configValues) UintOr(key string, alt uint) (uint, error) {
	if values, err := c.snapshot(); err != nil {
		return 0, err
	} else {
		return values.UintOr(key, alt)
	}
}

func (c configValues) Float(key string) (float64, error) {
	if values, err := c.snapshot(); err != nil {
		return 0, err
	} else {
		return values.Float(key)
	}
}

func (c configValues) FloatOr(key string, alt float64) (float64, error) {
	if values, err := c.snapshot(); err != nil {
		return 0, err
	} else {
		return values.FloatOr(key, alt)
	}
}

func (c configValues) Bool(key string) (bool, error) {
	if values, err := c.snapshot(); err != nil {
		return false, err
	} else {
		return values.Bool(key)
	}
}

func (c configValues) BoolOr(key string, alt bool) (bool, error) {
	if values, err := c.snapshot(); err != nil {
		return false, err
	} else {
		return values.BoolOr(key, alt)
	}
}

func (c configValues) Duration(key string) (time.Duration, error) {
	if values, err := c.snapshot(); err != nil {
		return 0, err
	} else {
		return values.Duration(key)
	}
}

func (c configValues) DurationOr(key string, alt time.Duration) (time.Duration, error) {
	if values, err := c.snapshot(); err != nil {
		return 0, err
	} else {
		return values.DurationOr(key, alt)
	}
}

func (c configValues) Value(key string) (Value, error) {
	values, err := c.snapshot()
	if err != nil {
		return "", err
	}

	entry, err := values.ResolveByName(key)
	if err != nil {
		return "", err
	}

	return entry.ResolvedValue, nil
}

func (c configValues) Populate(target interface{}, prefix string) error {
	if values, err := c.snapshot(); err != nil {
		return nil
	} else {
		return values.Populate(target, prefix)
	}
}

func (c configValues) Settings() map[string]string {
	if values, err := c.snapshot(); err != nil {
		return nil
	} else {
		return values.Settings()
	}
}

func (c configValues) Each(target func(string, string)) {
	if values, err := c.snapshot(); err == nil {
		values.Each(target)
	}
}
