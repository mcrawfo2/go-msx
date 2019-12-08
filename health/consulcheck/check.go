package consulcheck

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/consul"
	"cto-github.cisco.com/NFV-BU/go-msx/health"
)

func Check(ctx context.Context) health.CheckResult {
	consulPool := consul.PoolFromContext(ctx)
	if consulPool == nil {
		return health.CheckResult{
			Status: health.StatusDown,
			Details: map[string]interface{}{
				"error": "Consul pool not found in context",
			},
		}
	}

	conn := consulPool.Connection()
	checks, err := conn.NodeHealth(ctx)
	if err != nil {
		return health.CheckResult{
			Status: health.StatusDown,
			Details: map[string]interface{}{
				"error": err.Error(),
			},
		}
	}

	for _, check := range checks {
		if check.CheckID == "serfHealth" {
			if check.Status == "passing" {
				return health.CheckResult{Status: health.StatusUp}
			}
		}
	}

	return health.CheckResult{
		Status: health.StatusDown,
		Details: map[string]interface{}{
			"error": "Consul serfHealth check missing or failed",
		},
	}
}
