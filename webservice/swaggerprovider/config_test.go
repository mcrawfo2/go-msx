// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package swaggerprovider

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestDocumentationConfigFromConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *DocumentationConfig
		wantErr bool
	}{
		{
			name: "Defaults",
			args: args{
				cfg: configtest.NewInMemoryConfig(nil),
			},
			want: &DocumentationConfig{
				Enabled:     false,
				ApiPath:     "/apidocs.json",
				SwaggerPath: "/swagger-resources",
				ApiYamlPath: "/apidocs.yml",
				Version:     "2.0",
				Security: DocumentationSecurityConfig{
					Sso: DocumentationSecuritySsoConfig{
						BaseUrl:       "http://localhost:9103/idm",
						TokenPath:     "/v2/token",
						AuthorizePath: "/v2/authorize",
						ClientId:      "nfv-client",
						ClientSecret:  "",
					},
				},
				Ui: DocumentationUiConfig{
					Enabled:  true,
					Endpoint: "/swagger",
					View:     "/swagger-ui.html",
				},
				Server: DocumentationServerConfig{
					Host:        "localhost",
					Port:        0,
					ContextPath: "/",
				},
			},
		},
		{
			name: "Embedded",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"swagger.enabled": "true",
				}),
			},
			want: &DocumentationConfig{
				Enabled:     true,
				ApiPath:     "/apidocs.json",
				SwaggerPath: "/swagger-resources",
				ApiYamlPath: "/apidocs.yml",
				Version:     "2.0",
				Security: DocumentationSecurityConfig{
					Sso: DocumentationSecuritySsoConfig{
						BaseUrl:       "http://localhost:9103/idm",
						TokenPath:     "/v2/token",
						AuthorizePath: "/v2/authorize",
						ClientId:      "nfv-client",
						ClientSecret:  "",
					},
				},
				Ui: DocumentationUiConfig{
					Enabled:  true,
					Endpoint: "/swagger",
					View:     "/swagger-ui.html",
				},
				Server: DocumentationServerConfig{
					Host:        "localhost",
					Port:        0,
					ContextPath: "/",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DocumentationConfigFromConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("DocumentationConfigFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(testhelpers.Diff(tt.want, got))
			}
		})
	}
}
