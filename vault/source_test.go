package vault

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestNewTokenSource(t *testing.T) {
	emptyCfg := configtest.NewInMemoryConfig(nil)
	conn := new(Connection)

	type args struct {
		source string
		cfg    *config.Config
		conn   *Connection
	}
	tests := []struct {
		name            string
		args            args
		wantTokenSource TokenSource
		wantErr         bool
	}{
		{
			name: "Config",
			args: args{
				source: "config",
				cfg:    emptyCfg,
				conn:   nil,
			},
			wantTokenSource: ConfigSource{
				cfg: emptyCfg,
			},
		},
		{
			name: "AppRole",
			args: args{
				source: "approle",
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.cloud.vault.token-source.approle.role-id":   "role-id",
					"spring.cloud.vault.token-source.approle.secret-id": "secret-id",
				}),
				conn: conn,
			},
			wantTokenSource: &AppRoleSource{
				cfg: &AppRoleConfig{
					RoleId:   "role-id",
					SecretId: "secret-id",
				},
				conn: conn,
			},
			wantErr: false,
		},
		{
			name: "Kubernetes",
			args: args{
				source: "kubernetes",
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.cloud.vault.token-source.kubernetes.role": "role",
				}),
				conn: conn,
			},
			wantTokenSource: &KubernetesSource{
				cfg: &KubernetesConfig{
					JWTPath: "/run/secrets/kubernetes.io/serviceaccount/token",
					Role:    "role",
				},
				conn: conn,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTokenSource, err := NewTokenSource(tt.args.source, tt.args.cfg, tt.args.conn)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTokenSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTokenSource, tt.wantTokenSource) {
				t.Errorf("NewTokenSource() gotTokenSource = %v, want %v", gotTokenSource, tt.wantTokenSource)
			}
		})
	}
}
