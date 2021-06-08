package vault

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
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

	if secrets, err := c.read(ctx, path); err != nil {
		return nil, errors.Wrap(err, "Failed to list vault secrets")
	} else if secrets != nil {
		for key, val := range secrets.Data {
			results[key] = val.(string)
		}
	}

	return
}

// Copied from vault/logical to allow custom context
func (c connectionImpl) read(ctx context.Context, path string) (*api.Secret, error) {
	r := c.client.NewRequest("GET", "/v1/"+path)

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

func (c connectionImpl) StoreSecrets(ctx context.Context, path string, secrets map[string]string) (err error) {
	if _, err = c.write(ctx, path, secrets); err != nil {
		err = errors.Wrap(err, "Failed to store vault secrets")
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

func (c connectionImpl) CreateTransitKey(ctx context.Context, keyName string, request CreateTransitKeyRequest) (err error) {
	path := "transit/keys/" + keyName
	if _, err = c.write(ctx, path, request); err != nil {
		err = errors.Wrap(err, "Failed to create transit key")
	}
	return
}

func (c connectionImpl) TransitEncrypt(ctx context.Context, keyName string, plaintext string) (ciphertext string, err error) {
	path := "/transit/encrypt/" + keyName

	data := map[string]interface{}{
		"plaintext": base64.StdEncoding.EncodeToString([]byte(plaintext)),
	}

	result, err := c.write(ctx, path, data)
	if err != nil {
		return "", errors.Wrap(err, "Failed to encrypt data")
	}

	ciphertext = result.Data["ciphertext"].(string)
	return
}

func (c connectionImpl) TransitDecrypt(ctx context.Context, keyName string, ciphertext string) (plaintext string, err error) {
	path := "/transit/decrypt/" + keyName

	data := map[string]interface{}{
		"ciphertext": ciphertext,
	}

	result, err := c.write(ctx, path, data)
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
	path := c.cfg.Issuer.Mount + "/issue/" + role

	if secret, err := c.write(ctx, path, request.Data()); err != nil {
		return nil, errors.Wrap(err, "Failed to issue certificate")
	} else {
		crtPEM := []byte(secret.Data["certificate"].(string))
		keyPEM := []byte(secret.Data["private_key"].(string))

		var converted tls.Certificate
		converted, err = tls.X509KeyPair(crtPEM, keyPEM)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to parse X.509 keypair")
		} else {
			cert = &converted
		}
	}

	return
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

func (c connectionImpl) ReadCaCertificate(ctx context.Context) (cert *x509.Certificate, err error) {
	pemBytes, err := c.readRaw(ctx, c.cfg.Issuer.Mount+"/ca/pem")
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

func newConnectionImpl(cfg *ConnectionConfig, client *api.Client) *connectionImpl {
	return &connectionImpl{
		cfg:    cfg,
		client: client,
	}
}
