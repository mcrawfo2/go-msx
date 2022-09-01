package vault

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/hashicorp/vault/api"
)

type DisConnection struct{}

func (d DisConnection) ListSecrets(ctx context.Context, path string) (results map[string]string, err error) {
	return map[string]string{}, nil
}

func (d DisConnection) StoreSecrets(ctx context.Context, path string, secrets map[string]string) (err error) {
	return nil
}

func (d DisConnection) RemoveSecrets(ctx context.Context, path string) (err error) {
	return nil
}

func (d DisConnection) ListV2Secrets(ctx context.Context, path string) (keys []string, err error) {
	return nil, nil
}

func (d DisConnection) GetVersionedSecrets(ctx context.Context, path string, version *int) (results map[string]interface{}, err error) {
	return map[string]interface{}{}, nil
}

func (d DisConnection) StoreVersionedSecrets(ctx context.Context, path string, request VersionedWriteRequest) (err error) {
	return nil
}

func (d DisConnection) PatchVersionedSecrets(ctx context.Context, path string, request VersionedWriteRequest) (err error) {
	return nil
}

func (d DisConnection) DeleteVersionedSecretsLatest(ctx context.Context, p string) (err error) {
	return nil
}

func (d DisConnection) DeleteVersionedSecrets(ctx context.Context, path string, request VersionRequest) (err error) {
	return nil
}

func (d DisConnection) UndeleteVersionedSecrets(ctx context.Context, path string, request VersionRequest) (err error) {
	return nil
}

func (d DisConnection) DestroyVersionedSecrets(ctx context.Context, path string, request VersionRequest) (err error) {
	return nil
}

func (d DisConnection) GetVersionedMetadata(ctx context.Context, path string) (results VersionedMetadata, err error) {
	return VersionedMetadata{}, nil
}

func (d DisConnection) StoreVersionedMetadata(ctx context.Context, path string, request VersionedMetadataRequest) (err error) {
	return nil
}

func (d DisConnection) DeleteVersionedMetadata(ctx context.Context, path string) (err error) {
	return nil
}

func (d DisConnection) CreateTransitKey(ctx context.Context, keyName string, request CreateTransitKeyRequest) (err error) {
	return nil
}

func (d DisConnection) TransitEncrypt(ctx context.Context, keyName string, plaintext string) (ciphertext string, err error) {
	return "ciphertext", nil
}

func (d DisConnection) TransitDecrypt(ctx context.Context, keyName string, ciphertext string) (plaintext string, err error) {
	return "plaintext", nil
}

func (d DisConnection) TransitBulkDecrypt(ctx context.Context, keyName string, ciphertext ...string) (plaintext []string, err error) {
	for range ciphertext {
		plaintext = append(plaintext, "plaintext")
	}
	return
}

func (d DisConnection) GetTransitKeys(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (d DisConnection) IssueCertificate(ctx context.Context, role string, request IssueCertificateRequest) (cert *tls.Certificate, err error) {
	return nil, nil
}

func (d DisConnection) ReadCaCertificate(ctx context.Context) (cert *x509.Certificate, err error) {
	return nil, nil
}

func (d DisConnection) IssueCustomCertificate(ctx context.Context, mount string, role string, request IssueCertificateRequest) (cert *tls.Certificate, ca *x509.Certificate, err error) {
	return nil, nil, nil
}

func (d DisConnection) ReadCustomCaCertificate(ctx context.Context, mount string) (cert *x509.Certificate, err error) {
	return nil, nil
}

func (d DisConnection) Health(ctx context.Context) (response *api.HealthResponse, err error) {
	return nil, nil
}

func (d DisConnection) LoginWithKubernetes(ctx context.Context, jwt, role string) (token string, err error) {
	return "token", nil
}

func (d DisConnection) LoginWithAppRole(ctx context.Context, roleId, secretId string) (token string, err error) {
	return "token", nil
}

func (d DisConnection) GenerateRandomBytes(ctx context.Context, length int) (data []byte, err error) {
	return []byte{}, nil
}
