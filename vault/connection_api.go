// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

//go:generate mockery --name=ConnectionApi --inpackage --structname=MockConnection --filename=mock_connection.go

package vault

import (
	"context"
	"crypto/tls"
	"crypto/x509"

	"github.com/hashicorp/vault/api"
)

type ConnectionApi interface {
	// KV V1
	ListSecrets(ctx context.Context, path string) (results map[string]string, err error)
	StoreSecrets(ctx context.Context, path string, secrets map[string]string) (err error)
	RemoveSecrets(ctx context.Context, path string) (err error)

	// KV-V2
	GetVersionedSecrets(ctx context.Context, path string, version *int) (results map[string]interface{}, err error)
	StoreVersionedSecrets(ctx context.Context, path string, request VersionedWriteRequest) (err error)
	PatchVersionedSecrets(ctx context.Context, path string, request VersionedWriteRequest) (err error)
	DeleteVersionedSecretsLatest(ctx context.Context, p string) (err error)
	DeleteVersionedSecrets(ctx context.Context, path string, request VersionRequest) (err error)
	UndeleteVersionedSecrets(ctx context.Context, path string, request VersionRequest) (err error)
	DestroyVersionedSecrets(ctx context.Context, path string, request VersionRequest) (err error)
	GetVersionedMetadata(ctx context.Context, path string) (results VersionedMetadata, err error)
	StoreVersionedMetadata(ctx context.Context, path string, request VersionedMetadataRequest) (err error)
	DeleteVersionedMetadata(ctx context.Context, path string) (err error)

	// Transit
	CreateTransitKey(ctx context.Context, keyName string, request CreateTransitKeyRequest) (err error)
	TransitEncrypt(ctx context.Context, keyName string, plaintext string) (ciphertext string, err error)
	TransitDecrypt(ctx context.Context, keyName string, ciphertext string) (plaintext string, err error)
	TransitBulkDecrypt(ctx context.Context, keyName string, ciphertext ...string) (plaintext []string, err error)
	GetTransitKeys(ctx context.Context) ([]string, error)

	// Certificate (Default Mount)
	IssueCertificate(ctx context.Context, role string, request IssueCertificateRequest) (cert *tls.Certificate, err error)
	ReadCaCertificate(ctx context.Context) (cert *x509.Certificate, err error)

	// Certificate (Custom Mount)
	IssueCustomCertificate(ctx context.Context, mount string, role string, request IssueCertificateRequest) (cert *tls.Certificate, ca *x509.Certificate, err error)
	ReadCustomCaCertificate(ctx context.Context, mount string) (cert *x509.Certificate, err error)

	Health(ctx context.Context) (response *api.HealthResponse, err error)

	// Auth
	LoginWithKubernetes(ctx context.Context, jwt, role string) (token string, err error)
	LoginWithAppRole(ctx context.Context, roleId, secretId string) (token string, err error)

	// Sys-Tools
	GenerateRandomBytes(ctx context.Context, length int) (data []byte, err error)
}
