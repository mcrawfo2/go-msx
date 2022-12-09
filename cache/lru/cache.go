// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

// Package lru implements a cache for arbitrary values
// It provides an interface type Cache and a concrete type HeapMapCache
// Entries are added with a key, a value and an individual TTL - time to live.
// New uses should call NewCache2 to instantiate a cache with the deAgeOnAccess setting.
// The cache will expire entries after the TTL has passed. It checks every
// ExpireFrequency for expired entries and expires them in batches
// of at most ExpireLimit at once.
// The cache has no size limit. It will grow until the process runs out of memory, unless entries are expired.
// The cache is safe for concurrent access.
// The setting DeAgeOnAccess being true will cause the cache to reset the TTL of an entry when it is accessed or updated in true LRU fashion.
// When this setting is false (default) it behaves like a simple TTL cache.
// When the setting metrics is true (default false), the cache will emit metrics to the stats package thus:
// - entries: the number of entries in the cache
// - hits: the number of cache hits
// - misses: the number of cache misses
// - sets: the number of times set or setWithTTL were called
// - evictions: the number of times an entry was evicted
// - gcRuns: the number of times the garbage collector was run
// - gcSizes: a histogram of the number of entries evicted in each garbage collection run
// - deAgedAt: a histogram of the remaining time to live of entries when they are deaged
// The metricsPrefix setting is used to prefix the metrics names.
// The timeSource setting is used to provide a clock for testing purposes.
package lru

import (
	"github.com/benbjohnson/clock"
	"sync"
	"time"
)

type Cache interface {
	Get(key string) (any, bool)
	Set(key string, value any)
	Clear()
}

type HeapMapCache struct {
	sync.RWMutex
	ttl             time.Duration     // period before expiry
	index           map[string]*entry // entries indexed by key
	collected       []string          // temp garbage collection space
	expireFrequency time.Duration     // garbage collection frequency
	deAgeOnAccess   bool              // reset TTL on access or update
	timeSource      clock.Clock       // for test double
	metrics         bool              // enable metrics collection
	metricsPrefix   string            // prefix for metrics
	metricsObs      metricsObserver   // the metrics for the cache
}

// entry is an individual entry in the cache.
type entry struct {
	key     string
	expires int64
	value   any
}

// NewCache2 chains NewCache (which is retained for backwards compatibility) and adds the deageonaccess,
// metrics and metricsPrefix settings
// New uses should call this constructor instead of NewCache
func NewCache2(ttl time.Duration, expireLimit int, expireFrequency time.Duration,
	deageonaccess bool, timeSource clock.Clock,
	metrics bool, metricsPrefix string) *HeapMapCache {

	c := &HeapMapCache{
		ttl:             ttl,
		index:           make(map[string]*entry),
		collected:       make([]string, expireLimit),
		expireFrequency: expireFrequency,
		timeSource:      timeSource,
		deAgeOnAccess:   deageonaccess,
		metrics:         metrics,
		metricsPrefix:   metricsPrefix,
	}

	if metrics {
		name := metricsPrefix + subsystemName
		c.metricsObs = newPrometheusMetrics(name)
	} else {
		c.metricsObs = &nullMetricsObserver{}
	}

	go c.tick()
	return c
}

// Get fetches a value from cache and updates its expiry time if it exists and DeAgeOnAccess is true.
// Returns the value and a bool indicating if it was found.
func (c *HeapMapCache) Get(key string) (value any, found bool) {
	c.RLock()
	defer c.RUnlock()

	e, exists := c.index[key]
	if exists {
		c.metricsObs.OnHit()
	}

	if !exists {
		c.metricsObs.OnMiss()
		return nil, false
	}

	now := c.timeSource.Now().UnixNano()

	if e.expires < now {
		c.metricsObs.OnMiss()
		return nil, false
	}
	if c.deAgeOnAccess {
		c.metricsObs.OnDeAge((e.expires - now) / int64(time.Millisecond))
		e.expires = now + int64(c.ttl)
	}

	return e.value, true
}

// Set adds a value to the cache with the given key and current default TTL
func (c *HeapMapCache) Set(key string, value any) {
	c.SetWithTtl(key, value, c.ttl)
}

// SetWithTtl adds a value to the cache with the given key and TTL
func (c *HeapMapCache) SetWithTtl(key string, value any, ttl time.Duration) {
	c.Lock()
	defer c.Unlock()
	expires := c.timeSource.Now().Add(ttl).UnixNano()
	e, exists := c.index[key]
	if exists {
		e.value = value
		if c.deAgeOnAccess {
			c.metricsObs.OnDeAge((e.expires - time.Now().UnixNano()) / int64(time.Millisecond))
			e.expires = expires
		}
	} else {
		e = &entry{
			key:     key,
			value:   value,
			expires: expires,
		}
		c.metricsObs.OnEntriesInc()
	}
	c.index[key] = e
	c.metricsObs.OnSet()
}

// Clear removes all entries from the cache
func (c *HeapMapCache) Clear() {
	c.Lock()
	defer c.Unlock()
	c.index = make(map[string]*entry)
	c.metricsObs.OnEntriesResize(0)
}

func (c *HeapMapCache) tick() {
	for {
		c.timeSource.Sleep(c.expireFrequency)
		c.expire()
	}
}

func (c *HeapMapCache) expire() {
	c.Lock()
	defer c.Unlock()

	expiredCount := c.collect()
	c.metricsObs.OnGC(expiredCount)
	c.metricsObs.OnEvict(expiredCount)
	for i := 0; i < expiredCount; i++ {
		delete(c.index, c.collected[i])
	}
	c.metricsObs.OnEntriesResize(len(c.index))
}

// collect identifies expired entries and returns their count
func (c *HeapMapCache) collect() int {
	now := c.timeSource.Now().UnixNano()
	j := 0
	for k, e := range c.index {
		if e.expires < now {
			c.collected[j] = k
			j++
		}
		if j == len(c.collected) {
			break
		}
	}
	return j
}

// NewCache creates a new cache with the given TTL, ExpireLimit & ExpireFrequency
func NewCache(ttl time.Duration, expireLimit int,
	expireFrequency time.Duration, timeSource clock.Clock) *HeapMapCache {
	return NewCache2(ttl, expireLimit, expireFrequency, false, timeSource, false, "")
}

// NewCacheFromConfig creates a new cache with the given config.
func NewCacheFromConfig(cfg *CacheConfig) *HeapMapCache {
	return NewCache2(cfg.Ttl, cfg.ExpireLimit, cfg.ExpireFrequency,
		cfg.DeAgeOnAccess, clock.New(), cfg.Metrics, cfg.MetricsPrefix)
}
