// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package lru

import (
	"github.com/benbjohnson/clock"
	"sync"
	"time"
)

type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Clear()
}

type HeapMapCache struct {
	sync.RWMutex
	ttl             time.Duration     // period before expiry
	index           map[string]*entry // entries indexed by key
	collected       []string          // temp garbage collection space
	expireFrequency time.Duration     // garbage collection frequency
	timeSource      clock.Clock       // for test double
}

func (c *HeapMapCache) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()

	e, exists := c.index[key]

	if !exists {
		return nil, false
	}

	if e.expires < c.timeSource.Now().UnixNano() {
		return nil, false
	}

	return e.value, true
}

func (c *HeapMapCache) Set(key string, value interface{}) {
	c.SetWithTtl(key, value, c.ttl)
}

func (c *HeapMapCache) SetWithTtl(key string, value interface{}, ttl time.Duration) {
	c.Lock()
	defer c.Unlock()
	expires := c.timeSource.Now().Add(ttl).UnixNano()
	e, exists := c.index[key]
	if exists {
		e.value = value
	} else {
		e = &entry{
			key:     key,
			value:   value,
			expires: expires,
			index:   -1,
		}
	}
	c.index[key] = e
}

func (c *HeapMapCache) Clear() {
	c.Lock()
	defer c.Unlock()
	c.index = make(map[string]*entry)
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
	for i := 0; i < expiredCount; i++ {
		delete(c.index, c.collected[i])
	}
}

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

func NewCache(ttl time.Duration, expireLimit int, expireFrequency time.Duration, timeSource clock.Clock) *HeapMapCache {
	c := &HeapMapCache{
		ttl:             ttl,
		index:           make(map[string]*entry),
		collected:       make([]string, expireLimit),
		expireFrequency: expireFrequency,
		timeSource:      timeSource,
	}
	go c.tick()
	return c
}

func NewCacheFromConfig(cfg *CacheConfig) *HeapMapCache {
	return NewCache(cfg.Ttl, cfg.ExpireLimit, cfg.ExpireFrequency, clock.New())
}
