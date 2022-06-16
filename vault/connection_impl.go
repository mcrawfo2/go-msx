// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package vault

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type connectionImpl struct {
	client *api.Client
	cfg    *ConnectionConfig
}

func (c connectionImpl) LoginWithKubernetes(ctx context.Context, jwt, role string) (token string, err error) {
	payload := map[string]interface{}{
		"jwt":  jwt,
		"role": role,
	}

	result, err := c.write(ctx, c.cfg.Kubernetes.LoginPath, payload)
	if err != nil {
		return "", err
	}

	return result.Auth.ClientToken, nil
}

func (c connectionImpl) LoginWithAppRole(ctx context.Context, roleId, secretId string) (token string, err error) {
	payload := map[string]interface{}{
		"role_id":   roleId,
		"secret_id": secretId,
	}

	result, err := c.write(ctx, c.cfg.AppRole.LoginPath, payload)
	if err != nil {
		return "", err
	}

	return result.Auth.ClientToken, nil
}

func (c connectionImpl) Host() string {
	return c.cfg.Host
}

func (c connectionImpl) Client() *api.Client {
	return c.client
}

func (c connectionImpl) ListSecrets(ctx context.Context, path string) (results map[string]string, err error) {
	results = make(map[string]string)
	if secrets, err := c.read(ctx, path, nil); err != nil {
		return nil, errors.Wrap(err, "Failed to list vault secrets")
	} else if secrets != nil {
		for key, val := range secrets.Data {
			results[key] = val.(string)
		}
	}

	return
}

func (c connectionImpl) ListV2Secrets(ctx context.Context, path string) (results []string, err error) {
	if secret, err := c.list(ctx, "/v1/v2secret/metadata"+path, nil); err != nil {
		return nil, err
	} else if secret != nil {
		tempData := secret["data"].(map[string]interface{})["keys"].([]interface{})
		for _, v := range tempData {
			results = append(results, v.(string))
		}
	}
	return
}

func (c connectionImpl) callVaultAPI(ctx context.Context, method string, path string, query url.Values) (*api.Response, error) {
	r := c.client.NewRequest(method, path)
	r.Params = query
	return c.client.RawRequestWithContext(ctx, r)
}

func (c connectionImpl) list(ctx context.Context, path string, query url.Values) (map[string]interface{}, error) {
	resp, err := c.callVaultAPI(ctx, "LIST", path, query)
	var secrets map[string]interface{}
	if resp != nil && resp.StatusCode == 404 {
		secret, parseErr := api.ParseSecret(resp.Body)
		switch parseErr {
		case nil:
		case io.EOF:
			return nil, nil
		default:
			return nil, err
		}
		if secret != nil && len(secret.Warnings) > 0 {
			return nil, errors.New(strings.Join(secret.Warnings[:], ","))
		}
		return nil, nil
	} else if err != nil {
		logger.Errorf("An error occurred calling vault %+v", err)
		return nil, err
	} else {
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Errorf("An error occurred parsing the response from vault ")
			return nil, err
		}
		parseErr := json.Unmarshal(bytes, &secrets)
		if parseErr != nil {
			logger.Errorf("Error occurred when unmarshalling list of keys ")
			return nil, err
		}
	}
	return secrets, err
}

// Copied from vault/logical to allow custom context
func (c connectionImpl) read(ctx context.Context, path string, query url.Values) (*api.Secret, error) {
	r := c.client.NewRequest("GET", "/v1/"+path)
	r.Params = query
	resp, err := c.client.RawRequestWithContext(ctx, r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		secret, parseErr := api.ParseSecret(resp.Body)
		switch parseErr {
		case nil:
		case io.EOF:
			return nil, nil
		default:
			return nil, err
		}
		if secret != nil && (len(secret.Warnings) > 0 || len(secret.Data) > 0) {
			return secret, nil
		}
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return api.ParseSecret(resp.Body)
}

func (c connectionImpl) StoreSecrets(ctx context.Context, path string, secrets map[string]string) (err error) {
	if _, err = c.write(ctx, path, secrets); err != nil {
		err = errors.Wrap(err, "Failed to store vault secrets")
	}
	return
}

func (c connectionImpl) DeleteSecrets(ctx context.Context, p string, request VersionRequest) (err error) {
	if p[0] == '/' {
		p = p[1:]
	}
	pathParts := strings.Split(p, "/")
	pathParts = append([]string{pathParts[0], "delete"})

	if _, err = c.write(ctx, p, request); err != nil {
		err = errors.Wrap(err, "Failed to delete vault secret versions")
	}

	return
}

func (c connectionImpl) write(ctx context.Context, path string, requestBody interface{}) (*api.Secret, error) {
	r := c.client.NewRequest("POST", "/v1/"+path)
	if err := r.SetJSONBody(requestBody); err != nil {
		return nil, err
	}

	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()
	resp, err := c.client.RawRequestWithContext(ctx, r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		secret, parseErr := api.ParseSecret(resp.Body)
		switch parseErr {
		case nil:
		case io.EOF:
			return nil, nil
		default:
			return nil, err
		}
		if secret != nil && (len(secret.Warnings) > 0 || len(secret.Data) > 0) {
			return secret, err
		}
	}
	if err != nil {
		return nil, err
	}

	return api.ParseSecret(resp.Body)
}

func (c connectionImpl) action(ctx context.Context, method string, path string, requestBody interface{}) (*api.Secret, error) {
	r := c.client.NewRequest(method, "/v1/"+path)
	if err := r.SetJSONBody(requestBody); err != nil {
		return nil, err
	}

	if method == http.MethodPatch {
		if r.Headers == nil {
			r.Headers = make(http.Header)
		}
		r.Headers.Set("Content-Type", "application/merge-patch+json")
	}

	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()
	resp, err := c.client.RawRequestWithContext(ctx, r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		secret, parseErr := api.ParseSecret(resp.Body)
		switch parseErr {
		case nil:
		case io.EOF:
			return nil, nil
		default:
			return nil, err
		}
		if secret != nil && (len(secret.Warnings) > 0 || len(secret.Data) > 0) {
			return secret, err
		}
	}
	if err != nil {
		return nil, err
	}

	return api.ParseSecret(resp.Body)
}

func (c connectionImpl) RemoveSecrets(ctx context.Context, path string) (err error) {
	if _, err = c.delete(ctx, path); err != nil {
		err = errors.Wrap(err, "Failed to remove vault secrets")
	}

	return
}

func (c connectionImpl) delete(ctx context.Context, path string) (*api.Secret, error) {
	r := c.client.NewRequest("DELETE", "/v1/"+path)
	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()
	resp, err := c.client.RawRequestWithContext(ctx, r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		secret, parseErr := api.ParseSecret(resp.Body)
		switch parseErr {
		case nil:
		case io.EOF:
			return nil, nil
		default:
			return nil, err
		}
		if secret != nil && (len(secret.Warnings) > 0 || len(secret.Data) > 0) {
			return secret, err
		}
	}
	if err != nil {
		return nil, err
	}
	return api.ParseSecret(resp.Body)
}

func (c connectionImpl) GetVersionedSecrets(ctx context.Context, p string, version *int) (results map[string]interface{}, err error) {
	versionKvDataPath := path.Join(c.cfg.KV2.Mount, "data", p)
	var query url.Values
	if version != nil {
		query = make(url.Values)
		query.Set("version", strconv.Itoa(*version))
	}

	results = make(types.Pojo)
	if secrets, err := c.read(ctx, versionKvDataPath, query); err != nil {
		return nil, errors.Wrap(err, "Failed to list vault secrets")
	} else if secrets != nil {
		for key, val := range secrets.Data {
			results[key] = val
		}
	}

	return
}

func (c connectionImpl) StoreVersionedSecrets(ctx context.Context, p string, request VersionedWriteRequest) (err error) {
	versionKvDataPath := path.Join(c.cfg.KV2.Mount, "data", p)

	if _, err = c.action(ctx, http.MethodPost, versionKvDataPath, request); err != nil {
		err = errors.Wrap(err, "Failed to store vault secrets")
	}
	return
}

func (c connectionImpl) PatchVersionedSecrets(ctx context.Context, p string, request VersionedWriteRequest) (err error) {
	versionKvDataPath := path.Join(c.cfg.KV2.Mount, "data", p)

	if _, err = c.action(ctx, http.MethodPatch, versionKvDataPath, request); err != nil {
		err = errors.Wrap(err, "Failed to patch vault secrets")
	}
	return
}

func (c connectionImpl) DeleteVersionedSecretsLatest(ctx context.Context, p string) (err error) {
	versionKvDataPath := path.Join(c.cfg.KV2.Mount, "data", p)

	if _, err = c.delete(ctx, versionKvDataPath); err != nil {
		err = errors.Wrap(err, "Failed to delete vault secrets")
	}
	return
}

func (c connectionImpl) DeleteVersionedSecrets(ctx context.Context, p string, request VersionRequest) (err error) {
	versionKvDeletePath := path.Join(c.cfg.KV2.Mount, "delete", p)

	if _, err = c.action(ctx, http.MethodPost, versionKvDeletePath, request); err != nil {
		err = errors.Wrap(err, "Failed to soft-delete vault secrets")
	}
	return
}

func (c connectionImpl) UndeleteVersionedSecrets(ctx context.Context, p string, request VersionRequest) (err error) {
	versionKvUndeletePath := path.Join(c.cfg.KV2.Mount, "undelete", p)

	if _, err = c.action(ctx, http.MethodPost, versionKvUndeletePath, request); err != nil {
		err = errors.Wrap(err, "Failed to soft-undelete vault secrets")
	}
	return
}

func (c connectionImpl) DestroyVersionedSecrets(ctx context.Context, p string, request VersionRequest) (err error) {
	versionKvDestroyPath := path.Join(c.cfg.KV2.Mount, "destroy", p)

	if _, err = c.action(ctx, http.MethodPost, versionKvDestroyPath, request); err != nil {
		err = errors.Wrap(err, "Failed to destroy vault secrets")
	}
	return
}

func (c connectionImpl) GetVersionedMetadata(ctx context.Context, p string) (results VersionedMetadata, err error) {
	versionKvMetadataPath := path.Join(c.cfg.KV2.Mount, "metadata", p)

	secrets, err := c.read(ctx, versionKvMetadataPath, nil)
	if err != nil {
		err = errors.Wrap(err, "Failed to read vault metadata")
		return
	}

	if secrets != nil {
		err = mapstructure.Decode(secrets.Data, &results)
		if err != nil {
			err = errors.Wrap(err, "Failed to decode vault metadata")
		}
	}

	return
}

func (c connectionImpl) StoreVersionedMetadata(ctx context.Context, p string, request VersionedMetadataRequest) (err error) {
	versionKvMetadataPath := path.Join(c.cfg.KV2.Mount, "metadata", p)

	if _, err = c.action(ctx, http.MethodPost, versionKvMetadataPath, request); err != nil {
		err = errors.Wrap(err, "Failed to store vault secrets metadata")
	}
	return
}

func (c connectionImpl) DeleteVersionedMetadata(ctx context.Context, p string) (err error) {
	versionKvMetadataPath := path.Join(c.cfg.KV2.Mount, "metadata", p)

	if _, err = c.delete(ctx, versionKvMetadataPath); err != nil {
		err = errors.Wrap(err, "Failed to delete vault secrets metadata")
	}
	return
}

func (c connectionImpl) CreateTransitKey(ctx context.Context, keyName string, request CreateTransitKeyRequest) (err error) {
	p := "transit/keys/" + keyName
	if _, err = c.write(ctx, p, request); err != nil {
		err = errors.Wrap(err, "Failed to create transit key")
	}
	return
}

func (c connectionImpl) GetTransitKeys(ctx context.Context) (results []string, err error) {
	p := "transit/keys"
	params := url.Values{"list": []string{"true"}}

	secrets, err := c.read(ctx, p, params)
	if err != nil {
		return results, errors.Wrap(err, "Failed to get transit keys")
	}

	if secrets != nil {
		keys, ok := secrets.Data["keys"].([]interface{})
		if ok {
			for _, key := range keys {
				results = append(results, fmt.Sprintf("%v", key))
			}
		}
	}
	return
}

func (c connectionImpl) TransitEncrypt(ctx context.Context, keyName string, plaintext string) (ciphertext string, err error) {
	p := "/transit/encrypt/" + keyName

	data := map[string]interface{}{
		"plaintext": base64.StdEncoding.EncodeToString([]byte(plaintext)),
	}

	result, err := c.write(ctx, p, data)
	if err != nil {
		return "", errors.Wrap(err, "Failed to encrypt data")
	}

	ciphertext = result.Data["ciphertext"].(string)
	return
}

func (c connectionImpl) TransitBulkDecrypt(ctx context.Context, keyName string, ciphertexts ...string) (plaintext []string, err error) {
	path := "/transit/decrypt/" + keyName

	var entries []types.Pojo
	for _, ciphertext := range ciphertexts {
		entries = append(entries, types.Pojo{
			"ciphertext": ciphertext,
		})
	}

	data := types.Pojo{
		"batch_input": entries,
	}

	result, err := c.write(ctx, path, data)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to decrypt data")
	}

	batchResultsData, err := types.Pojo(result.Data).ArrayValue("batch_results")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse decrypt response")
	}

	var batchResults []types.Pojo
	for i := range batchResultsData {
		br, err := batchResultsData.ObjectValue(i)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to parse decrypt response")
		}
		batchResults = append(batchResults, br)
	}

	for _, batchResult := range batchResults {
		plaintextBytes, err := base64.StdEncoding.DecodeString(batchResult["plaintext"].(string))
		if err != nil {
			return nil, err
		}

		plaintext = append(plaintext, string(plaintextBytes))
	}

	return
}

func (c connectionImpl) TransitDecrypt(ctx context.Context, keyName string, ciphertext string) (plaintext string, err error) {
	p := "/transit/decrypt/" + keyName

	data := map[string]interface{}{
		"ciphertext": ciphertext,
	}

	result, err := c.write(ctx, p, data)
	if err != nil {
		return "", errors.Wrap(err, "Failed to decrypt data")
	}

	plaintextBytes, err := base64.StdEncoding.DecodeString(result.Data["plaintext"].(string))
	if err != nil {
		return "", err
	}

	plaintext = string(plaintextBytes)
	return
}

func (c connectionImpl) IssueCertificate(ctx context.Context, role string, request IssueCertificateRequest) (cert *tls.Certificate, err error) {
	cert, _, err = c.IssueCustomCertificate(ctx, c.cfg.Issuer.Mount, role, request)
	return
}

func (c connectionImpl) IssueCustomCertificate(ctx context.Context, mount string, role string, request IssueCertificateRequest) (cert *tls.Certificate, ca *x509.Certificate, err error) {
	fullPath := path.Join(mount, "issue", role)

	if secret, err := c.write(ctx, fullPath, request.Data()); err != nil {
		return nil, nil, errors.Wrap(err, "Failed to issue certificate")
	} else {
		crtPEM := []byte(secret.Data["certificate"].(string))
		keyPEM := []byte(secret.Data["private_key"].(string))
		caPEM := []byte(secret.Data["issuing_ca"].(string))

		var converted tls.Certificate
		converted, err = tls.X509KeyPair(crtPEM, keyPEM)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Failed to parse X.509 keypair")
		} else {
			cert = &converted
		}

		// Decode the first block (CA cert)
		pemBlock, _ := pem.Decode(caPEM)
		if pemBlock == nil || pemBlock.Type != "CERTIFICATE" {
			return nil, nil, errors.Errorf("Could not decode certificate: found %q", pemBlock.Type)
		}

		var convertedAuthority *x509.Certificate
		convertedAuthority, err = x509.ParseCertificate(pemBlock.Bytes)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Failed to parse X509 ca certificate")
		} else {
			ca = convertedAuthority
		}
	}

	return
}

func (c connectionImpl) ReadCaCertificate(ctx context.Context) (cert *x509.Certificate, err error) {
	return c.ReadCustomCaCertificate(ctx, c.cfg.Issuer.Mount)
}

func (c connectionImpl) ReadCustomCaCertificate(ctx context.Context, mount string) (cert *x509.Certificate, err error) {
	fullPath := path.Join(mount, "ca", "pem")

	pemBytes, err := c.readRaw(ctx, fullPath)
	if err != nil {
		return nil, err
	}

	// Decode the first block (CA cert)
	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock == nil || pemBlock.Type != "CERTIFICATE" {
		return nil, errors.Errorf("Could not decode certificate: found %q", pemBlock.Type)
	}

	return x509.ParseCertificate(pemBlock.Bytes)
}

// Health returns a health check of the Vault server.
// Copied from vault/api to allow custom context
func (c connectionImpl) Health(ctx context.Context) (response *api.HealthResponse, err error) {
	r := c.client.NewRequest("GET", "/v1/sys/health")
	// If the code is 400 or above it will automatically turn into an error,
	// but the sys/health API defaults to returning 5xx when not sealed or
	// inited, so we force this code to be something else so we parse correctly
	r.Params.Add("uninitcode", "299")
	r.Params.Add("sealedcode", "299")
	r.Params.Add("standbycode", "299")
	r.Params.Add("drsecondarycode", "299")
	r.Params.Add("performancestandbycode", "299")

	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()
	resp, err := c.client.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result api.HealthResponse
	err = resp.DecodeJSON(&result)
	if err != nil {
		return nil, err
	}

	response = &result
	return
}

func (c connectionImpl) GenerateRandomBytes(ctx context.Context, length int) (data []byte, err error) {
	body := map[string]interface{}{
		"format": "hex",
		"bytes":  length,
	}

	secret, err := c.write(ctx, "/sys/tools/random", body)
	if err != nil {
		return nil, err
	}

	dataString := secret.Data["random_bytes"]
	data, err = hex.DecodeString(dataString.(string))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c connectionImpl) readRaw(ctx context.Context, path string) ([]byte, error) {
	r := c.client.NewRequest("GET", "/v1/"+path)

	resp, err := c.client.RawRequestWithContext(ctx, r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

func newConnectionImpl(cfg *ConnectionConfig, client *api.Client) *connectionImpl {
	return &connectionImpl{
		cfg:    cfg,
		client: client,
	}
}
