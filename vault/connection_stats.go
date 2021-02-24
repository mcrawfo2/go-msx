package vault

import (
	"context"
	"crypto/tls"
	"cto-github.cisco.com/NFV-BU/go-msx/stats"
	"github.com/hashicorp/vault/api"
	"time"
)

const (
	statsSubsystemVault         = "vault"
	statsGaugeVaultCalls        = "calls"
	statsHistogramVaultCallTime = "call_time"
	statsCounterVaultCallErrors = "call_errors"

	statsApiListSecrets      = "listSecrets"
	statsApiStoreSecrets     = "storeSecrets"
	statsApiRemoveSecrets    = "removeSecrets"
	statsApiHealth           = "health"
	statsApiCreateTransitKey = "createTransitKey"
	statsApiTransitEncrypt   = "transitEncrypt"
	statsApiTransitDecrypt   = "transitDecrypt"
	statsApiIssueCertificate = "issueCertificate"
)

var (
	histVecVaultCallTime    = stats.NewHistogramVec(statsSubsystemVault, statsHistogramVaultCallTime, nil, "api", "param")
	gaugeVecVaultCalls      = stats.NewGaugeVec(statsSubsystemVault, statsGaugeVaultCalls, "api", "param")
	countVecVaultCallErrors = stats.NewCounterVec(statsSubsystemVault, statsCounterVaultCallErrors, "api", "param")
)

type queryFunc func() error

type statsConnection struct {
	ConnectionApi
}

func (s statsConnection) Observe(api, param string, fn queryFunc) (err error) {
	start := time.Now()
	gaugeVecVaultCalls.WithLabelValues(api, param).Inc()

	defer func() {
		gaugeVecVaultCalls.WithLabelValues(api, param).Dec()
		histVecVaultCallTime.WithLabelValues(api, param).Observe(float64(time.Since(start)) / float64(time.Millisecond))
		if err != nil {
			countVecVaultCallErrors.WithLabelValues(api, param).Inc()
		}
	}()

	err = fn()
	return err
}

func (s statsConnection) ListSecrets(ctx context.Context, path string) (results map[string]string, err error) {
	err = s.Observe(statsApiListSecrets, path, func() error {
		results, err = s.ConnectionApi.ListSecrets(ctx, path)
		return err
	})
	return
}

func (s statsConnection) StoreSecrets(ctx context.Context, path string, secrets map[string]string) (err error) {
	err = s.Observe(statsApiStoreSecrets, path, func() error {
		return s.ConnectionApi.StoreSecrets(ctx, path, secrets)
	})
	return
}

func (s statsConnection) RemoveSecrets(ctx context.Context, path string) (err error) {
	err = s.Observe(statsApiRemoveSecrets, path, func() error {
		return s.ConnectionApi.RemoveSecrets(ctx, path)
	})
	return
}

func (s statsConnection) CreateTransitKey(ctx context.Context, keyName string, request CreateTransitKeyRequest) (err error) {
	err = s.Observe(statsApiCreateTransitKey, keyName, func() error {
		return s.ConnectionApi.CreateTransitKey(ctx, keyName, request)
	})
	return
}

func (s statsConnection) TransitEncrypt(ctx context.Context, keyName string, plaintext string) (ciphertext string, err error) {
	err = s.Observe(statsApiTransitEncrypt, keyName, func() error {
		ciphertext, err = s.ConnectionApi.TransitEncrypt(ctx, keyName, plaintext)
		return err
	})
	return
}

func (s statsConnection) TransitDecrypt(ctx context.Context, keyName string, ciphertext string) (plaintext string, err error) {
	err = s.Observe(statsApiTransitDecrypt, keyName, func() error {
		plaintext, err = s.ConnectionApi.TransitDecrypt(ctx, keyName, ciphertext)
		return err
	})
	return
}

func (s statsConnection) IssueCertificate(ctx context.Context, role string, request IssueCertificateRequest) (cert *tls.Certificate, err error) {
	err = s.Observe(statsApiIssueCertificate, role, func() error {
		cert, err = s.ConnectionApi.IssueCertificate(ctx, role, request)
		return err
	})
	return
}

func (s statsConnection) Health(ctx context.Context) (response *api.HealthResponse, err error) {
	err = s.Observe(statsApiHealth, "", func() error {
		response, err = s.ConnectionApi.Health(ctx)
		return err
	})
	return
}

func newStatsConnection(api ConnectionApi) statsConnection {
	return statsConnection{
		ConnectionApi: api,
	}
}
