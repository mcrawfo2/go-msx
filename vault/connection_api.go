//go:generate mockery --name=ConnectionApi --inpackage --structname=MockConnection --filename=mock_connection.go

package vault

import (
	"context"
	"crypto/tls"
	"github.com/hashicorp/vault/api"
)

type ConnectionApi interface {
	ListSecrets(ctx context.Context, path string) (results map[string]string, err error)
	StoreSecrets(ctx context.Context, path string, secrets map[string]string) (err error)
	RemoveSecrets(ctx context.Context, path string) (err error)
	CreateTransitKey(ctx context.Context, keyName string, request CreateTransitKeyRequest) (err error)
	TransitEncrypt(ctx context.Context, keyName string, plaintext string) (ciphertext string, err error)
	TransitDecrypt(ctx context.Context, keyName string, ciphertext string) (plaintext string, err error)
	IssueCertificate(ctx context.Context, role string, request IssueCertificateRequest) (cert *tls.Certificate, err error)
	Health(ctx context.Context) (response *api.HealthResponse, err error)
	LoginWithKubernetes(ctx context.Context, jwt, role string) (token string, err error)
	LoginWithAppRole(ctx context.Context, roleId, secretId string) (token string, err error)
	GenerateRandomBytes(ctx context.Context, length int) (data []byte, err error)

	// Deprecated
	Host() string
	// Deprecated
	Client() *api.Client
}
