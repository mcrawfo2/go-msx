package vault

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestNewConnectionConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *ConnectionConfig
		wantErr bool
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
			},
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConnectionConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConnectionConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConnectionConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
