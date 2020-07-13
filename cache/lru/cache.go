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
	heap            entryHeap         // entries ordered by expiry
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
	c.Lock()
	defer c.Unlock()
	expires := c.timeSource.Now().Add(c.ttl).UnixNano()
	e, exists := c.index[key]
	if exists {
		e.value = value
		c.heap.update(e, expires)
	} else {
		e = &entry{
			key:     key,
			value:   value,
			expires: expires,
			index:   -1,
		}

		c.heap.Push(e)
	}
	c.index[key] = e
}

func (c *HeapMapCache) Clear() {
	c.Lock()
	defer c.Unlock()
	c.index = make(map[string]*entry)
	c.heap = make(entryHeap, 0)
}

func (c *HeapMapCache) collect() {
	for {
		c.timeSource.Sleep(c.expireFrequency)
		c.expire()
	}
}

func (c *HeapMapCache) expire() {
	c.Lock()
	defer c.Unlock()
	now := c.timeSource.Now().UnixNano()
	expiredCount := c.heap.expire(now, c.collected)
	for i := 0; i < expiredCount; i++ {
		delete(c.index, c.collected[i])
	}
}

func NewCache(ttl time.Duration, expireLimit int, expireFrequency time.Duration, timeSource clock.Clock) *HeapMapCache {
	c := &HeapMapCache{
		ttl:             ttl,
		index:           make(map[string]*entry),
		heap:            make(entryHeap, 0),
		collected:       make([]string, expireLimit),
		expireFrequency: expireFrequency,
		timeSource:      timeSource,
	}
	go c.collect()
	return c
}

func NewCacheFromConfig(cfg *CacheConfig) *HeapMapCache {
	return NewCache(cfg.Ttl, cfg.ExpireLimit, cfg.ExpireFrequency, clock.New())
}
