package redischeck

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/health"
	"cto-github.cisco.com/NFV-BU/go-msx/redis"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"github.com/pkg/errors"
)

func Check(ctx context.Context) health.CheckResult {
	redisPool := redis.PoolFromContext(ctx)
	if redisPool == nil {
		return health.CheckResult{
			Status: health.StatusDown,
			Details: map[string]interface{}{
				"error": "Redis pool not found in context",
			},
		}
	}

	var healthResult health.CheckResult
	err := trace.Operation(ctx, "redis.healthCheck", func(ctx context.Context) error {
		conn := redisPool.Connection()
		pingCmd := conn.Client(ctx).Ping()
		pingResult, err := pingCmd.Result()
		healthResult = health.CheckResult{
			Status: health.StatusUp,
			Details: map[string]interface{}{
				"version": conn.Version(),
			},
		}

		if err != nil {
			return err
		} else if pingResult != "PONG" {
			return errors.New("Redis returned invalid PING response: " + pingResult)
		}
		return nil
	})

	if err != nil {
		healthResult.Details["error"] = err.Error()
		healthResult.Status = health.StatusDown
	}

	return healthResult
}
