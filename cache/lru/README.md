# MSX LRU Cache

An LRU cache implementation which expires key/value pairs based on a TTL duration.
Inspired by [rcache](https://github.com/karlseguin/rcache).

## Usage

### Instantiation

To create a new cache with 120 second retention, which expires up to 100 keys every 15 seconds:

```go
myCache := lru.NewCache(120 * time.Second, 100, 15 * time.Second)
```

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
