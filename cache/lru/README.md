# MSX LRU Cache

An LRU cache implementation which expires key/value pairs based on a TTL duration.
Inspired by [rcache](https://github.com/karlseguin/rcache).

- Entries are added with a key, a value and an individual TTL - time to live.  
- New uses should call NewCache2
- The cache will expire entries after the TTL has passed.
- The cache checks every `ExpireFrequency` for expired entries and expires them in batches of at most `ExpireLimit` at once.
- The cache has no size limit. It will grow until the process runs out of memory, unless entries are expired.
- The cache is safe for concurrent access.
- The setting `DeAgeOnAccess` being true will cause the cache to reset the TTL of an entry when it is accessed or updated, in true LRU fashion.
- When this setting is `false` (default for backwards compatibility) it behaves like a simple TTL cache. New uses should probably set this to `true`.
- When the setting `metrics` is `true` (default `false`), the cache will emit metrics.
- The timeSource setting is used to provide a clock for testing purposes.

## Usage

### Instantiation

To create a new cache with 120 second retention, which expires up to 100 keys every 15 seconds with de-aging switched on, metrics on, with prefix "cat", and a normal time source:

```go
myCache := lru.NewCache2(120 * time.Second, 100, 15 * time.Second, true,
	clock.New(), true, "cat_")
```

lru provides an interface type Cache and a concrete type HeapMapCache; NewCache2 returns an instance of HeapMapCache which implements the former.

### Storage

To store a key/value pair:

```go
myCache.Set("somekey", "myvalue")
```

### Retrieval

To retrieve a key/value pair:

```go
value, exists := myCache.get("somekey")
if !exists { 
  // fill cache for "somekey"
}
```

## Metrics

When initialized with `metrics` set true, the cache will emit metric events to the stats package thus:  

- `entries`: the number of entries in the cache
- `hits`: the number of cache hits
- `misses`: the number of cache misses
- `sets`: the number of times set or setWithTTL were called
- `evictions`: the number of times an entry was evicted
- `gcRuns`: the number of times the garbage collector was run
- `gcSizes`: a histogram of the number of entries evicted in each garbage collection run
- `deAgedAt`: a histogram of the remaining time to live of entries when they are deaged

The `metricsPrefix` setting is used to prefix the metrics names in the output system.