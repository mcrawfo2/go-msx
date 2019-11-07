package cassandracheck

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra"
	"cto-github.cisco.com/NFV-BU/go-msx/health"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"github.com/gocql/gocql"
)

func Check(ctx context.Context) health.CheckResult {
	cassandraPool := cassandra.PoolFromContext(ctx)
	if cassandraPool == nil {
		return health.CheckResult{
			Status: health.StatusDown,
			Details: map[string]interface{}{
				"error": "Cassandra pool not found in context",
			},
		}
	}

	healthResult := health.CheckResult{Details: make(map[string]interface{})}

	err := trace.Operation(ctx, "kafka.healthCheck", func() error {
		return cassandraPool.WithSession(func(session *gocql.Session) error {
			var version *string

			if err := session.Query("SELECT release_version FROM system.local").
				WithContext(ctx).
				Consistency(gocql.One).Scan(&version); err != nil {
				return err
			}

			healthResult.Details["version"] = *version
			healthResult.Status = health.StatusUp
			return nil
		})
	})

	if err != nil {
		healthResult = health.CheckResult{
			Status: health.StatusDown,
			Details: map[string]interface{}{
				"error": err.Error(),
			},
		}
	}

	return healthResult
}
