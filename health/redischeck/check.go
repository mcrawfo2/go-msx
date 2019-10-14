package redischeck

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/health"
	"cto-github.cisco.com/NFV-BU/go-msx/redis"
)

func Check(ctx context.Context) health.CheckResult {
	redisPool := redis.PoolFromContext(ctx)
	if redisPool == nil {
		return health.CheckResult{
			Status:  health.StatusDown,
			Details: map[string]interface{}{
				"error": "Redis pool not found in context",
			},
		}
	}

	conn := redisPool.Connection()
	pingCmd := conn.Client(ctx).Ping()
	pingResult, err := pingCmd.Result()
	if err != nil {
		return health.CheckResult{
			Status:  health.StatusDown,
			Details: map[string]interface{}{
				"error": err.Error(),
				"version": conn.Version(),
			},
		}
	} else if pingResult != "PONG" {
		return health.CheckResult{
			Status:  health.StatusDown,
			Details: map[string]interface{}{
				"error": "Redis returned invalid PING response: " + pingResult,
				"version": conn.Version(),
			},
		}
	} else {
		return health.CheckResult{
			Status:  health.StatusUp,
			Details: map[string]interface{}{
				"version": conn.Version(),
			},
		}
	}
}

