// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package webservice

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/securitytest"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/webservicetest"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestManagementSecurityConfig_EndpointSecurityEnabled(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
		"management.security.endpoint.alice.enabled": "true",
		"management.security.endpoint.bob.enabled":   "false",
	})
	cfg, err := NewManagementSecurityConfig(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	type args struct {
		endpoint string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "EnabledByDefault",
			args: args{endpoint: "charlie"},
			want: true,
		},
		{
			name: "ExplicitlyDisabled",
			args: args{endpoint: "bob"},
			want: false,
		},
		{
			name: "ExplicitlyEnabled",
			args: args{endpoint: "alice"},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, cfg.EndpointSecurityEnabled(tt.args.endpoint))
		})
	}
}

func TestManagementSecurityFilter_Filter(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{})
	cfg, err := NewManagementSecurityConfig(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "NoPermission",
			test: new(webservicetest.RouteBuilderTest).
				WithContextInjector(securitytest.PermissionInjector()).
				WithContextInjector(securitytest.AuthoritiesInjector("ROLE_CLIENT")).
				WithRouteFilter(NewManagementSecurityFilter(cfg)).
				WithResponsePredicate(webservicetest.ResponseHasStatus(401)),
		},
		{
			name: "NoAuthority",
			test: new(webservicetest.RouteBuilderTest).
				WithContextInjector(securitytest.PermissionInjector("IS_API_ADMIN")).
				WithRouteFilter(NewManagementSecurityFilter(cfg)).
				WithResponsePredicate(webservicetest.ResponseHasStatus(401)),
		},
		{
			name: "Authorized",
			test: new(webservicetest.RouteBuilderTest).
				WithContextInjector(securitytest.AuthoritiesInjector("ROLE_CLIENT")).
				WithContextInjector(securitytest.PermissionInjector("IS_API_ADMIN")).
				WithRouteFilter(NewManagementSecurityFilter(cfg)).
				WithRouteTargetReturn(204).
				WithResponsePredicate(webservicetest.ResponseHasStatus(204)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}

}

func TestManagementSecurityFilter_roles(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{})
	ctx = securitytest.AuthoritiesInjector("ROLE_PUBLISHER")(ctx)

	cfg, _ := NewManagementSecurityConfig(ctx)
	filter := ManagementSecurityFilter{cfg: cfg}

	type args struct {
		roles []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "HasRoles",
			args:    args{roles: []string{"ROLE_PUBLISHER"}},
			wantErr: false,
		},
		{
			name:    "NotHasRoles",
			args:    args{roles: []string{"ROLE_CONSUMER"}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg.Roles = tt.args.roles
			err := filter.roles(ctx)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestNewManagementSecurityConfig(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantCfg bool
	}{
		{
			name: "ValidConfig",
			args: args{
				ctx: configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
					"management.security.endpoint.alice.enabled": "true",
					"management.security.endpoint.bob.enabled":   "false",
				}),
			},
			wantErr: false,
			wantCfg: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := NewManagementSecurityConfig(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewManagementSecurityConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (cfg != nil) != tt.wantCfg {
				t.Errorf("NewManagementSecurityConfig() cfg = %v, wantCfg %v", cfg, tt.wantCfg)

			}
		})
	}
}

func TestNewManagementSecurityFilter(t *testing.T) {
	cfg := new(ManagementSecurityConfig)
	filter := NewManagementSecurityFilter(cfg)
	assert.NotNil(t, filter)
}
