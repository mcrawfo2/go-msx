package tokensource

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestNewAppRoleConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *AppRoleConfig
		wantErr bool
	}{
		{
			name: "Custom",
			args: args{
				cfg: configtest.NewStaticConfig(map[string]string{
					"spring.cloud.vault.token-source.app-role.role-id":   "some-role",
					"spring.cloud.vault.token-source.app-role.secret-id": "some-role-secret",
				}),
			},
			want: &AppRoleConfig{
				Path:     "/auth/approle/login",
				RoleId:   "some-role",
				SecretId: "some-role-secret",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAppRoleConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAppRoleConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAppRoleConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
