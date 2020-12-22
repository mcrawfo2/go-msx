package vault

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.cloud.vault.token-source.app-role.role-id":   "some-role",
					"spring.cloud.vault.token-source.app-role.secret-id": "some-role-secret",
				}),
			},
			want: &AppRoleConfig{
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

func TestAppRoleSource_Renewable(t *testing.T) {
	appRoleSource := AppRoleSource{
		cfg: &AppRoleConfig{
			RoleId:   "some-role",
			SecretId: "some-role-secret",
		},
		conn: new(MockConnection),
	}

	assert.True(t, appRoleSource.Renewable())
}

func TestAppRoleSource_GetToken(t *testing.T) {
	token := "expected_token"

	conn := new(MockConnection)
	conn.
		On("LoginWithAppRole",
			mock.AnythingOfType("*context.emptyCtx"),
			"some-role",
			"some-role-secret").
		Return(token, nil)

	appRoleSource := AppRoleSource{
		cfg: &AppRoleConfig{
			RoleId:   "some-role",
			SecretId: "some-role-secret",
		},
		conn: conn,
	}

	actualToken, err := appRoleSource.GetToken(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, token, actualToken)

}
