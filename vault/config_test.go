package vault

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewConnectionConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name             string
		args             args
		want             *ConnectionConfig
		wantAddress      string
		wantErr          bool
		wantClientConfig bool
	}{
		{
			name: "Defaults",
			args: args{
				cfg: configtest.NewInMemoryConfig(nil),
			},
			want: &ConnectionConfig{
				Enabled: true,
				Host:    "localhost",
				Port:    8200,
				Scheme:  "http",
				TokenSource: ConnectionTokenSourceConfig{
					Source: "config",
				},
				Ssl: ConnectionSslConfig{
					Cacert:     "",
					ClientCert: "",
					ClientKey:  "",
					Insecure:   true,
				},
				Issuer: ConnectionIssuerConfig{
					Mount: "/pki",
				},
				Kubernetes: ConnectionKubernetesConfig{
					LoginPath: "/auth/kubernetes/login",
				},
				AppRole: ConnectionAppRoleConfig{
					LoginPath: "/auth/approle/login",
				},
				KV: ConnectionKvConfig{
					Mount: "/secret",
				},
				KV2: ConnectionKv2Config{
					Mount: "/v2secret",
				},
			},
			wantAddress:      "http://localhost:8200",
			wantClientConfig: true,
		},
		{
			name: "Invalid",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.cloud.vault.enabled": "foo",
				}),
			},
			wantErr:          true,
			wantClientConfig: true,
		},
		{
			name: "Custom",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.cloud.vault.enabled":             "false",
					"spring.cloud.vault.host":                "remote-vm",
					"spring.cloud.vault.port":                "9999",
					"spring.cloud.vault.scheme":              "https",
					"spring.cloud.vault.token-source.source": "kubernetes",
					"spring.cloud.vault.ssl.ca-cert":         "ca.crt",
					"spring.cloud.vault.ssl.client-cert":     "client.crt",
					"spring.cloud.vault.ssl.client-key":      "client.key",
					"spring.cloud.vault.ssl.insecure":        "false",
					"spring.cloud.vault.issuer.mount":        "/pki/vms",
					"spring.cloud.vault.kv.mount":            "/secret",
					"spring.cloud.vault.kv2.mount":           "/v2secret",
				}),
			},
			want: &ConnectionConfig{
				Enabled: false,
				Host:    "remote-vm",
				Port:    9999,
				Scheme:  "https",
				TokenSource: ConnectionTokenSourceConfig{
					Source: "kubernetes",
				},
				Ssl: ConnectionSslConfig{
					Cacert:     "ca.crt",
					ClientCert: "client.crt",
					ClientKey:  "client.key",
					Insecure:   false,
				},
				Issuer: ConnectionIssuerConfig{
					Mount: "/pki/vms",
				},
				Kubernetes: ConnectionKubernetesConfig{
					LoginPath: "/auth/kubernetes/login",
				},
				AppRole: ConnectionAppRoleConfig{
					LoginPath: "/auth/approle/login",
				},
				KV: ConnectionKvConfig{
					Mount: "/secret",
				},
				KV2: ConnectionKv2Config{
					Mount: "/v2secret",
				},
			},
			wantClientConfig: false,
			wantAddress:      "https://remote-vm:9999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newConnectionConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("newConnectionConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(testhelpers.Diff(tt.want, got))
			}
			if got != nil {
				assert.Equal(t, got.Address(), tt.wantAddress)

				clientConfig, err := got.ClientConfig()
				if tt.wantClientConfig {
					assert.NotNil(t, clientConfig)
					assert.NoError(t, err)
				} else {
					assert.Nil(t, clientConfig)
					assert.Error(t, err)
				}
			}
		})
	}
}
