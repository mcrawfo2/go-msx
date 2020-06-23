package sqldbcheck

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/health"
	"cto-github.cisco.com/NFV-BU/go-msx/sqldb"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"github.com/jmoiron/sqlx"
	"regexp"
)

var versionParser = regexp.MustCompile(`^(\w+).*\sv?(\d+\.\d+\.\d+)\s.*$`)

func Check(ctx context.Context) health.CheckResult {
	pool, err := sqldb.PoolFromContext(ctx)
	if err != nil && err != sqldb.ErrDisabled {
		return health.CheckResult{
			Status: health.StatusDown,
			Details: map[string]interface{}{
				"error": "SQL database pool not found in context",
			},
		}
	} else if err != nil && err == sqldb.ErrDisabled {
		return health.CheckResult{
			Status: health.StatusDown,
			Details: map[string]interface{}{
				"error": err.Error(),
			},
		}
	}

	healthResult := health.CheckResult{Details: make(map[string]interface{})}

	err = trace.Operation(ctx, "sqldb.healthCheck", func(ctx context.Context) error {
		return pool.WithSqlxConnection(ctx, func(ctx context.Context, session *sqlx.DB) error {
			var version *string
			var server *string

			if pool.Config().Driver == "postgres" {
				// Request server version
				if err := session.GetContext(ctx, &version, `SELECT version()`); err != nil {
					return err
				}

				if version != nil {
					parts := versionParser.FindStringSubmatch(*version)
					if len(parts) == 3 {
						server = &parts[1]
						version = &parts[2]
					} else {
						version = nil
						server = nil
					}
				}
			} else if err := session.PingContext(ctx); err != nil {
				return err
			}

			if version == nil || server == nil {
				v := "Unknown"
				version = &v
				server = &v
			}

			healthResult.Details["version"] = *version
			healthResult.Details["server"] = *server
			healthResult.Details["driver"] = pool.Config().Driver
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
