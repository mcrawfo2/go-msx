// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

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

func TestNewKubernetesConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *KubernetesConfig
		wantErr bool
	}{
		{
			name: "Defaults",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.cloud.vault.token-source.kubernetes.role": "PlatformMicroservice",
				}),
			},
			want: &KubernetesConfig{
				JWTPath: "/run/secrets/kubernetes.io/serviceaccount/token",
				Role:    "PlatformMicroservice",
			},
		},
		{
			name: "Custom",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.cloud.vault.token-source.kubernetes.jwt-path": "/another/path/to/token",
					"spring.cloud.vault.token-source.kubernetes.role":     "PlatformMicroservice",
				}),
			},
			want: &KubernetesConfig{
				JWTPath: "/another/path/to/token",
				Role:    "PlatformMicroservice",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewKubernetesConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKubernetesConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKubernetesConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKubernetesSource_Renewable(t *testing.T) {
	kubernetesSource := KubernetesSource{
		cfg: &KubernetesConfig{
			JWTPath: "testdata/jwt.txt",
			Role:    "role",
		},
		conn: new(MockConnection),
	}

	assert.True(t, kubernetesSource.Renewable())
}

func TestKubernetesSource_GetToken(t *testing.T) {
	token := "expected_token"

	conn := new(MockConnection)
	conn.
		On("LoginWithKubernetes",
			mock.AnythingOfType("*context.emptyCtx"),
			"jwt",
			"role").
		Return(token, nil)

	kubernetesSource := KubernetesSource{
		cfg: &KubernetesConfig{
			JWTPath: "testdata/jwt.txt",
			Role:    "role",
		},
		conn: conn,
	}

	actualToken, err := kubernetesSource.GetToken(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, token, actualToken)

}
