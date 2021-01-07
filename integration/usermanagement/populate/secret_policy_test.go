package populate

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestNewSecretPolicyPopulatorConfigFromConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *SecretPolicyPopulatorConfig
		wantErr bool
	}{
		{
			name: "StructDefaults",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"populate.root": "/platform-common",
				}),
			},
			want: &SecretPolicyPopulatorConfig{
				Enabled: false,
				Root:    "/platform-common/usermanagement",
			},
		},
		{
			name: "CustomOptions",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"populate.root": "/platform-common",
					"populate.usermanagement.secret-policy.enabled": "true",
				}),
			},
			want: &SecretPolicyPopulatorConfig{
				Enabled: true,
				Root:    "/platform-common/usermanagement",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSecretPolicyPopulatorConfigFromConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSecretPolicyPopulatorConfigFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSecretPolicyPopulatorConfigFromConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
