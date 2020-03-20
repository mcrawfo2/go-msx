package vault

import (
	"cto-github.cisco.com/NFV-BU/go-msx/stats"
	"time"
)

const (
	statsSubsystemVault         = "vault"
	statsGaugeVaultCalls        = "calls"
	statsHistogramVaultCallTime = "call_time"
	statsCounterVaultCallErrors = "call_errors"

	statsApiListSecrets      = "listSecrets"
	statsApiStoreSecrets     = "storeSecrets"
	statsApiHealth           = "health"
	statsApiCreateTransitKey = "createTransitKey"
	statsApiTransitEncrypt   = "transitEncrypt"
	statsApiTransitDecrypt   = "transitDecrypt"
)

var (
	histVecVaultCallTime    = stats.NewHistogramVec(statsSubsystemVault, statsHistogramVaultCallTime, nil, "api", "param")
	gaugeVecVaultCalls      = stats.NewGaugeVec(statsSubsystemVault, statsGaugeVaultCalls, "api", "param")
	countVecVaultCallErrors = stats.NewCounterVec(statsSubsystemVault, statsCounterVaultCallErrors, "api", "param")
)

type queryFunc func() error

type statsObserver struct{}

func (o *statsObserver) Observe(api, param string, queryFunc queryFunc) (err error) {
	start := time.Now()
	gaugeVecVaultCalls.WithLabelValues(api, param).Inc()

	defer func() {
		gaugeVecVaultCalls.WithLabelValues(api, param).Dec()
		histVecVaultCallTime.WithLabelValues(api, param).Observe(float64(time.Since(start)) / float64(time.Millisecond))
		if err != nil {
			countVecVaultCallErrors.WithLabelValues(api, param).Inc()
		}
	}()

	err = queryFunc()
	return err
}
