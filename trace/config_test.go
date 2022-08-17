// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package trace

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestNewTracingConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *TracingConfig
		wantErr bool
	}{
		{
			name: "Default",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"info.app.name": "TestNewTracingConfig",
				}),
			},
			want: &TracingConfig{
				Enabled:     true,
				ServiceName: "TestNewTracingConfig",
				Collector:   "jaeger",
				Reporter: TracingReporterConfig{
					Enabled: false,
					Name:    "jaeger",
					Host:    "localhost",
					Port:    6831,
					Url:     "http://localhost:9411/api/v1/spans",
				},
			},
		},
		{
			name: "Zipkin",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"info.app.name":          "TestNewTracingConfig",
					"remote.service.address": "remote-vm",
					"spring.zipkin.base-url": "http://${remote.service.address}:9411/",
					"trace.enabled":          "true",
					"trace.reporter.name":    "zipkin",
					"trace.reporter.url":     "${spring.zipkin.base-url}api/v1/spans",
				}),
			},
			want: &TracingConfig{
				Enabled:     true,
				ServiceName: "TestNewTracingConfig",
				Collector:   "jaeger",
				Reporter: TracingReporterConfig{
					Name: "zipkin",
					Host: "localhost",
					Port: 6831,
					Url:  "http://remote-vm:9411/api/v1/spans",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTracingConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTracingConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTracingConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
