// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package cache

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/redis"
	"encoding/json"
	"errors"
	goredis "github.com/go-redis/redis/v8"
	"time"
)

type ContextCache[I any] struct {
	Ttl    time.Duration
	Prefix string
}

func NewContextCache[I any](ttl time.Duration, prefix string) ContextCache[I] {
	return ContextCache[I]{Ttl: ttl, Prefix: prefix}
}

func NewContextCacheFromConfig[I any](cfg *ContextCacheConfig) ContextCache[I] {
	return ContextCache[I]{Ttl: cfg.Ttl, Prefix: cfg.Prefix}
}

func (r ContextCache[I]) Get(ctx context.Context, key string) (any, bool, error) {
	redisPool := redis.PoolFromContext(ctx)
	redisClient := redisPool.Connection().Client(ctx)

	val, err := redisClient.Get(ctx, r.Prefix+key).Result()

	if errors.Is(err, goredis.Nil) { // redis: nil is expected when key does not exist
		return nil, false, nil
	}
	if err != nil {
		logger.Debug(err)
		return nil, false, err
	}

	var ret I
	err = json.Unmarshal([]byte(val), &ret)
	if err != nil {
		logger.Debug(err)
		return nil, true, err
	}

	return ret, true, nil
}

func (r ContextCache[I]) Set(ctx context.Context, key string, value any) (err error) {
	redisPool := redis.PoolFromContext(ctx)
	redisClient := redisPool.Connection().Client(ctx)

	strVal, err := json.Marshal(value)
	if err != nil {
		logger.Debug(err)
		return
	}

	err = redisClient.Set(ctx, r.Prefix+key, strVal, r.Ttl).Err()
	if err != nil {
		logger.Debug(err)
		return
	}

	return
}

func (r ContextCache[I]) Clear(ctx context.Context) (err error) {
	redisPool := redis.PoolFromContext(ctx)
	redisClient := redisPool.Connection().Client(ctx)

	iter := redisClient.Scan(ctx, 0, r.Prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		err = redisClient.Del(ctx, iter.Val()).Err()
		if err != nil {
			logger.Debug(err)
		}
	}
	if err = iter.Err(); err != nil {
		logger.Debug(err)
	}

	return
}
