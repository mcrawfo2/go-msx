// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package httpclient

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestDoFunc_RoundTrip(t *testing.T) {
	resp := new(http.Response)

	tests := []struct {
		name    string
		d       DoFunc
		want    *http.Response
		wantErr bool
	}{
		{
			name: "Success",
			d: func(req *http.Request) (*http.Response, error) {
				return resp, nil
			},
			want:    resp,
			wantErr: false,
		},
		{
			name: "Error",
			d: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("")
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := new(http.Request)
			got, err := tt.d.RoundTrip(request)
			if (err != nil) != tt.wantErr {
				t.Errorf("RoundTrip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RoundTrip() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	cfg := config.NewConfig()
	_ = cfg.Load(context.Background())
	ctx := config.ContextWithConfig(context.Background(), cfg)
	factory, _ := NewProductionHttpClientFactory(ctx)
	ctx = ContextWithFactory(ctx, factory)

	type args struct {
		ctx        context.Context
		configurer Configurer
	}
	tests := []struct {
		name       string
		args       args
		wantClient bool
		wantErr    bool
		validators []ClientValidator
	}{
		{
			name: "NoFactory",
			args: args{
				ctx: context.Background(),
			},
			wantClient: false,
			wantErr:    true,
		},
		{
			name: "NoConfigurer",
			args: args{
				ctx: ctx,
			},
			wantClient: true,
		},
		{
			name: "Configurer",
			args: args{
				ctx: ctx,
				configurer: ClientConfigurer{
					ClientFuncs: []ClientConfigurationFunc{
						ClientTimeout(3 * time.Second),
					},
				},
			},
			wantClient: true,
			validators: []ClientValidator{
				validateClientTimeout(3 * time.Second),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.ctx, tt.args.configurer)

			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got != nil) != tt.wantClient {
				t.Errorf("New() got = %v, wantClient %v", got, tt.wantClient)
			}

			for _, validator := range tt.validators {
				if !validator.Valid(got) {
					t.Errorf("Client validation failed: %s", validator.Description)
				}
			}
		})
	}
}
