package vaultcheck

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/health"
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

	client := vaultPool.Connection().Client()

	healthResponse, err := client.Sys().Health()
	if err != nil {
		return health.CheckResult{
			Status: health.StatusDown,
			Details: map[string]interface{}{
				"error": err.Error(),
			},
		}
	}

	version := healthResponse.Version
	return health.CheckResult{
		Status: health.StatusUp,
		Details: map[string]interface{}{
			"version": version,
		},
	}
}
