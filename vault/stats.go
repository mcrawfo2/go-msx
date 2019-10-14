package vault

import (
	"cto-github.cisco.com/NFV-BU/go-msx/stats"
	"time"
)

const (
	statsCounterVaultCalls      = "vault.calls"
	statsCounterVaultCallErrors = "vault.callErrors"
	statsTimerVaultCallTime     = "vault.timer"

	statsApiListSecrets = "list-secrets"
)

type queryFunc func() error

type statsObserver struct{}

func (o *statsObserver) Observe(api, param string, queryFunc queryFunc) (err error) {
	start := time.Now()
	defer func() {
		stats.Incr(stats.Name(statsCounterVaultCalls, api, param), 1)
		stats.PrecisionTiming(stats.Name(statsTimerVaultCallTime, api, param), time.Since(start))
		if err != nil {
			stats.Incr(statsCounterVaultCallErrors, 1)
		}
	}()

	err = queryFunc()
	return err
}
