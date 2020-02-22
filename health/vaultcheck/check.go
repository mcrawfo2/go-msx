package vaultcheck

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/health"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"cto-github.cisco.com/NFV-BU/go-msx/vault"
)

func Check(ctx context.Context) health.CheckResult {
	vaultPool := vault.PoolFromContext(ctx)
	if vaultPool == nil {
		return health.CheckResult{
			Status: health.StatusDown,
			Details: map[string]interface{}{
				"error": "Vault pool not found in context",
			},
		}
	}

	var healthResult health.CheckResult
	_ = trace.Operation(ctx, "vault.healthCheck", func(ctx context.Context) error {
		return vaultPool.WithConnection(func(connection *vault.Connection) error {
			healthResponse, err := connection.Health(ctx)
			if err != nil {
				healthResult = health.CheckResult{
					Status: health.StatusDown,
					Details: map[string]interface{}{
						"error": err.Error(),
					},
				}
				return err
			}

			version := healthResponse.Version
			healthResult = health.CheckResult{
				Status: health.StatusUp,
				Details: map[string]interface{}{
					"version": version,
				},
			}

			return nil
		})
	})

	return healthResult

}
