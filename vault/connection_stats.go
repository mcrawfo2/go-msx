// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package vault

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"time"

	"cto-github.cisco.com/NFV-BU/go-msx/stats"
	"github.com/hashicorp/vault/api"
)

const (
	statsSubsystemVault         = "vault"
	statsGaugeVaultCalls        = "calls"
	statsHistogramVaultCallTime = "call_time"
	statsCounterVaultCallErrors = "call_errors"

	statsApiListSecrets              = "listSecrets"
	statsApiStoreSecrets             = "storeSecrets"
	statsApiRemoveSecrets            = "removeSecrets"
	statsApiGetVersionedSecrets      = "getVersionedSecrets"
	statsApiStoreVersionedSecrets    = "storeVersionedSecrets"
	statsApiPatchVersionedSecrets    = "patchVersionedSecrets"
	statsApiDeleteVersionedSecrets   = "deleteVersionedSecrets"
	statsApiUndeleteVersionedSecrets = "undeleteVersionedSecrets"
	statsApiDestroyVersionedSecrets  = "destroyVersionedSecrets"
	statsApiGetVersionedMetadata     = "getVersionedMetadata"
	statsApiStoreVersionedMetadata   = "storeVersionedMetadata"
	statsApiDeleteVersionedMetadata  = "deleteVersionedMetadata"
	statsApiHealth                   = "health"
	statsApiCreateTransitKey         = "createTransitKey"
	statsApiTransitEncrypt           = "transitEncrypt"
	statsApiTransitDecrypt           = "transitDecrypt"
	statsApiTransitKey               = "transitKey"
	statsApiIssueCertificate         = "issueCertificate"
	statsApiGenerateRandomBytes      = "generateRandomBytes"
	statsApiReadCaCertificate        = "readCaCertificate"
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

func (s statsConnection) GetVersionedSecrets(ctx context.Context, path string, version *int) (results map[string]interface{}, err error) {
	err = s.Observe(statsApiGetVersionedSecrets, path, func() error {
		results, err = s.ConnectionApi.GetVersionedSecrets(ctx, path, version)
		return err
	})
	return
}

func (s statsConnection) StoreVersionedSecrets(ctx context.Context, path string, request VersionedWriteRequest) (err error) {
	err = s.Observe(statsApiStoreVersionedSecrets, path, func() error {
		return s.ConnectionApi.StoreVersionedSecrets(ctx, path, request)
	})
	return
}

func (s statsConnection) PatchVersionedSecrets(ctx context.Context, path string, request VersionedWriteRequest) (err error) {
	err = s.Observe(statsApiPatchVersionedSecrets, path, func() error {
		return s.ConnectionApi.PatchVersionedSecrets(ctx, path, request)
	})
	return
}

func (s statsConnection) DeleteVersionedSecrets(ctx context.Context, path string, request VersionRequest) (err error) {
	err = s.Observe(statsApiDeleteVersionedSecrets, path, func() error {
		return s.ConnectionApi.DeleteVersionedSecrets(ctx, path, request)
	})
	return
}

func (s statsConnection) UndeleteVersionedSecrets(ctx context.Context, path string, request VersionRequest) (err error) {
	err = s.Observe(statsApiUndeleteVersionedSecrets, path, func() error {
		return s.ConnectionApi.UndeleteVersionedSecrets(ctx, path, request)
	})
	return
}

func (s statsConnection) DestroyVersionedSecrets(ctx context.Context, path string, request VersionRequest) (err error) {
	err = s.Observe(statsApiDestroyVersionedSecrets, path, func() error {
		return s.ConnectionApi.DestroyVersionedSecrets(ctx, path, request)
	})
	return
}

func (s statsConnection) GetVersionedMetadata(ctx context.Context, path string) (results VersionedMetadata, err error) {
	err = s.Observe(statsApiGetVersionedMetadata, path, func() error {
		results, err = s.ConnectionApi.GetVersionedMetadata(ctx, path)
		return err
	})
	return
}

func (s statsConnection) StoreVersionedMetadata(ctx context.Context, path string, request VersionedMetadataRequest) (err error) {
	err = s.Observe(statsApiStoreVersionedMetadata, path, func() error {
		return s.ConnectionApi.StoreVersionedMetadata(ctx, path, request)
	})
	return
}

func (s statsConnection) DeleteVersionedMetadata(ctx context.Context, path string) (err error) {
	err = s.Observe(statsApiDeleteVersionedMetadata, path, func() error {
		return s.ConnectionApi.DeleteVersionedMetadata(ctx, path)
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

func (s statsConnection) TransitBulkDecrypt(ctx context.Context, keyName string, ciphertexts ...string) (plaintext []string, err error) {
	err = s.Observe(statsApiTransitDecrypt, keyName, func() error {
		plaintext, err = s.ConnectionApi.TransitBulkDecrypt(ctx, keyName, ciphertexts...)
		return err
	})
	return
}

func (s statsConnection) GetTransitKeys(ctx context.Context) (results []string, err error) {
	err = s.Observe(statsApiTransitKey, "", func() error {
		results, err = s.ConnectionApi.GetTransitKeys(ctx)
		return err
	})
	return
}

func (s statsConnection) IssueCustomCertificate(ctx context.Context, pki string, role string, request IssueCertificateRequest) (cert *tls.Certificate, ca *x509.Certificate, err error) {
	err = s.Observe(statsApiIssueCertificate, pki, func() error {
		cert, ca, err = s.ConnectionApi.IssueCustomCertificate(ctx, pki, role, request)
		return err
	})
	return
}

func (s statsConnection) ReadCustomCaCertificate(ctx context.Context, pki string) (cert *x509.Certificate, err error) {
	err = s.Observe(statsApiReadCaCertificate, pki, func() error {
		cert, err = s.ConnectionApi.ReadCustomCaCertificate(ctx, pki)
		return err
	})
	return
}

func (s statsConnection) GenerateRandomBytes(ctx context.Context, length int) (data []byte, err error) {
	err = s.Observe(statsApiGenerateRandomBytes, "", func() error {
		data, err = s.ConnectionApi.GenerateRandomBytes(ctx, length)
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
