package redis

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stats"
	"github.com/go-redis/redis/v7"
	"time"
)

const (
	statsCounterRedisCalls            = "redis.calls"
	statsCounterRedisCallErrors       = "redis.callErrors"
	statsTimerRedisCallTime           = "redis.timer"
)

type statsHook struct {}

func (s *statsHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return contextWithStartTime(ctx), nil
}

func (s *statsHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	startTime := startTimeFromContext(ctx)
	stats.Incr(stats.Name(statsCounterRedisCalls, cmd.Name(), ""), 1)
	stats.PrecisionTiming(stats.Name(statsTimerRedisCallTime, cmd.Name(), ""), time.Since(startTime))
	if cmd.Err() != nil {
		stats.Incr(stats.Name(statsCounterRedisCallErrors, cmd.Name(), ""), 1)
	}
	return nil
}

func (s *statsHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return contextWithStartTime(ctx), nil
}

func (s *statsHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}

type contextRedisKey int
const contextRedisStartTime contextRedisKey = iota

func contextWithStartTime(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextRedisStartTime, time.Now())
}

func startTimeFromContext(ctx context.Context) time.Time {
	return ctx.Value(contextRedisStartTime).(time.Time)
}