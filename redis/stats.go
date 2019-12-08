package redis

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stats"
	"github.com/go-redis/redis/v7"
	"time"
)

const (
	statsSubsystemRedis         = "redis"
	statsCounterRedisCalls      = "calls"
	statsGaugeRedisCallsActive  = "callsActive"
	statsCounterRedisCallErrors = "callErrors"
	statsHistogramRedisCallTime = "callTime"
)

var (
	countRedisCalls       = stats.NewCounterVec(statsSubsystemRedis, statsCounterRedisCalls, "command")
	countRedisCallErrors  = stats.NewCounterVec(statsSubsystemRedis, statsCounterRedisCallErrors, "command")
	gaugeRedisCallsActive = stats.NewGaugeVec(statsSubsystemRedis, statsGaugeRedisCallsActive, "command")
	histRedisCallTimeVec  = stats.NewHistogramVec(statsSubsystemRedis, statsHistogramRedisCallTime, nil, "command")
)

type statsHook struct{}

func (s *statsHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	gaugeRedisCallsActive.WithLabelValues(cmd.Name()).Inc()
	countRedisCalls.WithLabelValues(cmd.Name()).Inc()
	return contextWithStartTime(ctx), nil
}

func (s *statsHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	startTime := startTimeFromContext(ctx)
	gaugeRedisCallsActive.WithLabelValues(cmd.Name()).Dec()
	histRedisCallTimeVec.WithLabelValues(cmd.Name()).Observe(float64(time.Since(startTime)) / float64(time.Millisecond))
	if cmd.Err() != nil {
		countRedisCallErrors.WithLabelValues(cmd.Name()).Inc()
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
