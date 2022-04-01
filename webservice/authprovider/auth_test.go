// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package authprovider

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestNewResourcePatternAuthenticationConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *ResourcePatternAuthenticationConfig
		wantErr bool
	}{
		{
			name: "Defaults",
			args: args{
				cfg: configtest.NewInMemoryConfig(nil),
			},
			want: &ResourcePatternAuthenticationConfig{
				Blacklist: []string{
					"/api/**",
					"/admin",
					"/admin/**",
				},
				Whitelist: []string{
					"/admin/health",
					"/admin/info",
					"/admin/alive",
				},
			},
		},
		{
			name: "Custom",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"security.resources.patterns.blacklist": "/a,/b,/c",
					"security.resources.patterns.whitelist": "/d,/e,/f",
				}),
			},
			want: &ResourcePatternAuthenticationConfig{
				Blacklist: []string{
					"/a",
					"/b",
					"/c",
				},
				Whitelist: []string{
					"/d",
					"/e",
					"/f",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewResourcePatternAuthenticationConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewResourcePatternAuthenticationConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResourcePatternAuthenticationConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
