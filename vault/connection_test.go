// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package vault

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"github.com/stretchr/testify/assert"
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
