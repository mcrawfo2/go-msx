package cassandracheck

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra"
	"cto-github.cisco.com/NFV-BU/go-msx/health"
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

	result := health.CheckResult{Details:make(map[string]interface{})}

	err := cassandraPool.WithSession(func(session *gocql.Session) error {
		var version *string

		if err := session.Query("SELECT release_version FROM system.local").
			Consistency(gocql.One).Scan(&version); err != nil {
			return err
		}

		result.Details["version"] = *version
		result.Status = health.StatusUp
		return nil
	})

	if err == nil {
		return result
	}

	return health.CheckResult{
		Status: health.StatusDown,
		Details: map[string]interface{}{
			"error": err.Error(),
		},
	}
}
