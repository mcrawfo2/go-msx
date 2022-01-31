package vault

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"github.com/hashicorp/vault/api"
)

const tracePrefixVault = "vault."

type traceConnection struct {
	ConnectionApi
}

func (s traceConnection) ListSecrets(ctx context.Context, path string) (results map[string]string, err error) {
	err = trace.Operation(ctx, tracePrefixVault+statsApiListSecrets, func(ctx context.Context) error {
		results, err = s.ConnectionApi.ListSecrets(ctx, path)
		return err
	})
	return
}

func (s traceConnection) StoreSecrets(ctx context.Context, path string, secrets map[string]string) (err error) {
	err = trace.Operation(ctx, tracePrefixVault+statsApiStoreSecrets, func(ctx context.Context) error {
		return s.ConnectionApi.StoreSecrets(ctx, path, secrets)
	})
	return
}

func (s traceConnection) RemoveSecrets(ctx context.Context, path string) (err error) {
	err = trace.Operation(ctx, tracePrefixVault+statsApiRemoveSecrets, func(ctx context.Context) error {
		return s.ConnectionApi.RemoveSecrets(ctx, path)
	})
	return
}

func (s traceConnection) GetVersionedSecrets(ctx context.Context, path string, version *int) (results map[string]interface{}, err error) {
	err = trace.Operation(ctx, tracePrefixVault+statsApiGetVersionedSecrets, func(ctx context.Context) error {
		results, err = s.ConnectionApi.GetVersionedSecrets(ctx, path, version)
		return err
	})
	return
}

func (s traceConnection) StoreVersionedSecrets(ctx context.Context, path string, request VersionedWriteRequest) (err error) {
	err = trace.Operation(ctx, tracePrefixVault+statsApiStoreVersionedSecrets, func(ctx context.Context) error {
		return s.ConnectionApi.StoreVersionedSecrets(ctx, path, request)
	})
	return
}

func (s traceConnection) PatchVersionedSecrets(ctx context.Context, path string, request VersionedWriteRequest) (err error) {
	err = trace.Operation(ctx, tracePrefixVault+statsApiPatchVersionedSecrets, func(ctx context.Context) error {
		return s.ConnectionApi.PatchVersionedSecrets(ctx, path, request)
	})
	return
}

func (s traceConnection) DeleteVersionedSecrets(ctx context.Context, path string, request VersionRequest) (err error) {
	err = trace.Operation(ctx, tracePrefixVault+statsApiDeleteVersionedSecrets, func(ctx context.Context) error {
		return s.ConnectionApi.DeleteVersionedSecrets(ctx, path, request)
	})
	return
}

func (s traceConnection) UndeleteVersionedSecrets(ctx context.Context, path string, request VersionRequest) (err error) {
	err = trace.Operation(ctx, tracePrefixVault+statsApiUndeleteVersionedSecrets, func(ctx context.Context) error {
		return s.ConnectionApi.UndeleteVersionedSecrets(ctx, path, request)
	})
	return
}

func (s traceConnection) DestroyVersionedSecrets(ctx context.Context, path string, request VersionRequest) (err error) {
	err = trace.Operation(ctx, tracePrefixVault+statsApiDestroyVersionedSecrets, func(ctx context.Context) error {
		return s.ConnectionApi.DestroyVersionedSecrets(ctx, path, request)
	})
	return
}

func (s traceConnection) GetVersionedMetadata(ctx context.Context, path string) (results VersionedMetadata, err error) {
	err = trace.Operation(ctx, tracePrefixVault+statsApiGetVersionedMetadata, func(ctx context.Context) error {
		results, err = s.ConnectionApi.GetVersionedMetadata(ctx, path)
		return err
	})
	return
}

func (s traceConnection) StoreVersionedMetadata(ctx context.Context, path string, request VersionedMetadataRequest) (err error) {
	err = trace.Operation(ctx, tracePrefixVault+statsApiStoreVersionedMetadata, func(ctx context.Context) error {
		return s.ConnectionApi.StoreVersionedMetadata(ctx, path, request)
	})
	return
}

func (s traceConnection) DeleteVersionedMetadata(ctx context.Context, path string) (err error) {
	err = trace.Operation(ctx, tracePrefixVault+statsApiDeleteVersionedMetadata, func(ctx context.Context) error {
		return s.ConnectionApi.DeleteVersionedMetadata(ctx, path)
	})
	return
}

func (s traceConnection) CreateTransitKey(ctx context.Context, keyName string, request CreateTransitKeyRequest) (err error) {
	err = trace.Operation(ctx, tracePrefixVault+statsApiCreateTransitKey, func(ctx context.Context) error {
		return s.ConnectionApi.CreateTransitKey(ctx, keyName, request)
	})
	return
}

func (s traceConnection) TransitEncrypt(ctx context.Context, keyName string, plaintext string) (ciphertext string, err error) {
	err = trace.Operation(ctx, tracePrefixVault+statsApiTransitEncrypt, func(ctx context.Context) error {
		ciphertext, err = s.ConnectionApi.TransitEncrypt(ctx, keyName, plaintext)
		return err
	})
	return
}

func (s traceConnection) TransitDecrypt(ctx context.Context, keyName string, ciphertext string) (plaintext string, err error) {
	err = trace.Operation(ctx, tracePrefixVault+statsApiTransitDecrypt, func(ctx context.Context) error {
		plaintext, err = s.ConnectionApi.TransitDecrypt(ctx, keyName, ciphertext)
		return err
	})
	return
}

func (s traceConnection) TransitBulkDecrypt(ctx context.Context, keyName string, ciphertexts ...string) (plaintext []string, err error) {
	err = trace.Operation(ctx, tracePrefixVault+statsApiTransitDecrypt, func(ctx context.Context) error {
		plaintext, err = s.ConnectionApi.TransitBulkDecrypt(ctx, keyName, ciphertexts...)
		return err
	})
	return
}

func (s traceConnection) IssueCertificate(ctx context.Context, role string, request IssueCertificateRequest) (cert *tls.Certificate, err error) {
	err = trace.Operation(ctx, tracePrefixVault+statsApiIssueCertificate, func(ctx context.Context) error {
		cert, err = s.ConnectionApi.IssueCertificate(ctx, role, request)
		return err
	})
	return
}

func (s traceConnection) GenerateRandomBytes(ctx context.Context, length int) (data []byte, err error) {
	err = trace.Operation(ctx, tracePrefixVault+statsApiGenerateRandomBytes, func(ctx context.Context) error {
		data, err = s.ConnectionApi.GenerateRandomBytes(ctx, length)
		return err
	})
	return
}

func (s traceConnection) Health(ctx context.Context) (response *api.HealthResponse, err error) {
	err = trace.Operation(ctx, tracePrefixVault+statsApiHealth, func(ctx context.Context) error {
		response, err = s.ConnectionApi.Health(ctx)
		return err
	})
	return
}

func (s traceConnection) ReadCaCertificate(ctx context.Context) (cert *x509.Certificate, err error) {
	err = trace.Operation(ctx, tracePrefixVault+statsApiReadCaCertificate, func(ctx context.Context) error {
		cert, err = s.ConnectionApi.ReadCaCertificate(ctx)
		return err
	})
	return
}

func newTraceConnection(api ConnectionApi) traceConnection {
	return traceConnection{
		ConnectionApi: api,
	}
}
