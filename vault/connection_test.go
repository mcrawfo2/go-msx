// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package vault

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"encoding/json"
	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_NewConnection(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Disabled",
			args: args{
				ctx: configtest.ContextWithNewInMemoryConfig(context.Background(), nil),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Enabled",
			args: args{
				ctx: configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
					"spring.cloud.vault.enabled": "false",
				}),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConnection(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("newConnectionFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want {
				assert.NotNil(t, got)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

func Test_ListV2Secrets_Success(t *testing.T) {
	responseBody := map[string]interface{}{}
	responseBody["data"] = map[string]interface{}{
		"keys": []string{"key1/", "key2/"},
	}
	responseBody["auth"] = nil
	responseBody["renewable"] = false
	responseBody["request_id"] = "20ac948d-2b24-138e-76d5-c3e681d416e9"
	responseBody["warnings"] = nil
	responseBody["wrap_info"] = nil

	ctx := context.Background()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		bodyMarshalled, err := json.Marshal(responseBody)
		if err != nil {
			t.Errorf("Error happened during encoding %+v", err)
		}
		w.Write(bodyMarshalled)
	}))
	defer server.Close()
	cfg := configtest.NewInMemoryConfig(map[string]string{
		"spring.cloud.vault.enabled": "true",
		"spring.cloud.vault.scheme":  "http",
	})
	vaultConfig, _ := newConnectionConfig(cfg)
	newConf, err := vaultConfig.ClientConfig()
	newConf.Address = server.URL
	vaultClient, err := api.NewClient(newConf)

	if err != nil {
		t.Fatal(err)
	}
	vaultConnection := newConnectionImpl(vaultConfig, vaultClient)
	allKeys, err := vaultConnection.ListV2Secrets(ctx, "/secure-choice")
	expectedKeys := []string{"key1/", "key2/"}
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, vaultConnection)
	assert.Equal(t, expectedKeys, allKeys)
}

func Test_ListV2Secrets_Empty(t *testing.T) {
	responseBody := map[string]interface{}{}
	responseBody["data"] = nil
	responseBody["auth"] = nil
	responseBody["renewable"] = false
	responseBody["request_id"] = "20ac948d-2b24-138e-76d5-c3e681d416e9"
	responseBody["warnings"] = []string{"No keys for that path found"}
	responseBody["wrap_info"] = nil

	ctx := context.Background()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		bodyMarshalled, err := json.Marshal(responseBody)
		if err != nil {
			t.Errorf("Error happened during encoding %+v", err)
		}
		w.Write(bodyMarshalled)
	}))
	defer server.Close()
	cfg := configtest.NewInMemoryConfig(map[string]string{
		"spring.cloud.vault.enabled": "true",
		"spring.cloud.vault.scheme":  "http",
	})
	vaultConfig, _ := newConnectionConfig(cfg)
	newConf, err := vaultConfig.ClientConfig()
	newConf.Address = server.URL
	vaultClient, err := api.NewClient(newConf)

	if err != nil {
		t.Fatal(err)
	}
	vaultConnection := newConnectionImpl(vaultConfig, vaultClient)
	_, err = vaultConnection.ListV2Secrets(ctx, "/secure-choice")
	if err != nil {
		assert.Equal(t, "No keys for that path found", err.Error())
	} else {
		assert.Fail(t, "Expected an error object")
	}
}
