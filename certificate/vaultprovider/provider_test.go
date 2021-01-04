package vaultprovider

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
	"time"
)

func TestNewProviderConfig(t *testing.T) {
	type args struct {
		cfg        *config.Config
		configRoot string
	}
	tests := []struct {
		name    string
		args    args
		want    ProviderConfig
		wantErr bool
	}{
		{
			name: "Simple",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"certificate.source.alpha.provider": "vault",
					"certificate.source.alpha.role":     "vault-role-alpha",
					"certificate.source.alpha.cn":       "key=value, key2=value2",
				}),
				configRoot: "certificate.source.alpha",
			},
			want: ProviderConfig{
				Role:     "vault-role-alpha",
				TTL:      730 * time.Hour,
				CN:       "key=value, key2=value2",
				AltNames: []string{"localhost"},
				IPSans:   []string{"127.0.0.1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got ProviderConfig
			err := tt.args.cfg.Populate(&got, tt.args.configRoot)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewProviderConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProviderConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
