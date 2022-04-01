// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package consul

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestNewConnectionConfigFromConfig(t *testing.T) {
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
				cfg: configtest.NewInMemoryConfig(map[string]string{}),
			},
			want: &ConnectionConfig{
				Enabled: true,
				Host:    "localhost",
				Port:    8500,
				Scheme:  "http",
				Config: struct {
					AclToken string `config:"default="`
				}{
					AclToken: "",
				},
				Watch: struct {
					Enabled  bool `config:"default=true"`
					WaitTime int  `config:"default=55"`
				}{
					Enabled:  true,
					WaitTime: 55,
				},
			},
		},
		{
			name: "Custom",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.cloud.consul.enabled":          "false",
					"spring.cloud.consul.host":             "remote-vm",
					"spring.cloud.consul.port":             "9999",
					"spring.cloud.consul.scheme":           "https",
					"spring.cloud.consul.config.acl-token": "replace_with_token_value",
					"spring.cloud.consul.watch.enabled":    "false",
					"spring.cloud.consul.watch.wait-time":  "35",
				}),
			},
			want: &ConnectionConfig{
				Enabled: false,
				Host:    "remote-vm",
				Port:    9999,
				Scheme:  "https",
				Config: struct {
					AclToken string `config:"default="`
				}{
					AclToken: "replace_with_token_value",
				},
				Watch: struct {
					Enabled  bool `config:"default=true"`
					WaitTime int  `config:"default=55"`
				}{
					Enabled:  false,
					WaitTime: 35,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConnectionConfigFromConfig(tt.args.cfg)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewConnectionConfigFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConnectionConfigFromConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
