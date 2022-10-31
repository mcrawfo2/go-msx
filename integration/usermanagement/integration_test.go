// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package usermanagement

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/auth"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/idm"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/secrets"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestNewIntegration(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	ctxWithConfig := configtest.ContextWithNewInMemoryConfig(
		context.Background(),
		map[string]string{
			"remoteservice.usermanagementservice.service": "usermanagementservice",
		})

	authIntWithConfig, err := auth.NewIntegration(ctxWithConfig)
	if err != nil {
		t.Errorf("NewIntegration() %v", err)
		return
	}
	idmIntWithConfig, err := idm.NewIntegration(ctxWithConfig)
	if err != nil {
		t.Errorf("NewIntegration() %v", err)
		return
	}
	secretsIntWithConfig, err := secrets.NewIntegration(ctxWithConfig)
	if err != nil {
		t.Errorf("NewIntegration() %v", err)
		return
	}

	ctxWithConfigDifferentName := configtest.ContextWithNewInMemoryConfig(
		context.Background(),
		map[string]string{
			"remoteservice.usermanagementservice.service": "testservice1",
			"remoteservice.authservice.service":           "testservice2",
			"remoteservice.secretsservice.service":        "testservice3",
		})
	authIntWithConfigDifferentName, err := auth.NewIntegration(ctxWithConfigDifferentName)
	if err != nil {
		t.Errorf("NewIntegration() %v", err)
		return
	}
	idmIntWithConfigDifferentName, err := idm.NewIntegration(ctxWithConfigDifferentName)
	if err != nil {
		t.Errorf("NewIntegration() %v", err)
		return
	}
	secretsIntWithConfigDifferentName, err := secrets.NewIntegration(ctxWithConfigDifferentName)
	if err != nil {
		t.Errorf("NewIntegration() %v", err)
		return
	}

	tests := []struct {
		name string
		args args
		want Api
	}{
		{
			name: "NonExisting",
			args: args{
				ctx: ctxWithConfig,
			},
			want: &Integration{
				authIntWithConfig,
				idmIntWithConfig,
				secretsIntWithConfig,
			},
		},
		{
			name: "Existing",
			args: args{
				ctx: ContextWithIntegration(ctxWithConfig, &Integration{}),
			},
			want: &Integration{},
		},
		{
			name: "ServiceName",
			args: args{
				ctx: ctxWithConfig,
			},
			want: &Integration{
				authIntWithConfig,
				idmIntWithConfig,
				secretsIntWithConfig,
			},
		},
		{
			name: "DifferentServiceName",
			args: args{
				ctx: ctxWithConfigDifferentName,
			},
			want: &Integration{
				authIntWithConfigDifferentName,
				idmIntWithConfigDifferentName,
				secretsIntWithConfigDifferentName,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := NewIntegration(tt.args.ctx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIntegration() got = %v, want %v", got, tt.want)
			}
		})
	}
}
