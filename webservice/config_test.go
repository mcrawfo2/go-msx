// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/certificate"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestWebServerConfig_Address(t *testing.T) {
	tests := []struct {
		name string
		cfg  WebServerConfig
		want string
	}{
		{
			name: "Zeros",
			cfg: WebServerConfig{
				Host: "0.0.0.0",
				Port: 80,
			},
			want: "0.0.0.0:80",
		},
		{
			name: "Local",
			cfg: WebServerConfig{
				Host: "127.0.0.1",
				Port: 8080,
			},
			want: "127.0.0.1:8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.Address(); got != tt.want {
				t.Errorf("Address() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebServerConfig_Url(t *testing.T) {
	tests := []struct {
		name string
		cfg  WebServerConfig
		want string
	}{
		{
			name: "Zeros",
			cfg: WebServerConfig{
				Host: "0.0.0.0",
				Port: 80,
				Tls: certificate.TLSConfig{
					Enabled: true,
				},
			},
			want: "https://0.0.0.0:80",
		},
		{
			name: "Local",
			cfg: WebServerConfig{
				Host: "127.0.0.1",
				Port: 8080,
			},
			want: "http://127.0.0.1:8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.Url(); got != tt.want {
				t.Errorf("Address() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewWebServerConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *WebServerConfig
		wantErr bool
	}{
		{
			name: "Defaults",
			args: args{
				cfg: configtest.NewInMemoryConfig(nil),
			},
			want: &WebServerConfig{
				Enabled: false,
				Host:    "0.0.0.0",
				Port:    8080,
				Tls: certificate.TLSConfig{
					Enabled:           false,
					MinVersion:        "tls12",
					CertificateSource: "server",
					CipherSuites: []string{
						"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
						"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
						"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
						"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
						"TLS_RSA_WITH_AES_256_GCM_SHA384",
						"TLS_RSA_WITH_AES_256_CBC_SHA",
					},
				},
				Cors: CorsConfig{
					Enabled:              true,
					CustomAllowedHeaders: []string{},
					CustomExposedHeaders: []string{},
				},
				ContextPath:   "/app",
				StaticPath:    "/www",
				StaticEnabled: true,
				TraceEnabled:  false,
				DebugEnabled:  false,
			},
		},
		{
			name: "Microservice",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					// Embedded
					"server.static-path":            "/www",
					"server.trace-enabled":          "false",
					"server.host":                   "0.0.0.0",
					"server.tls.cert-file":          "server.crt",
					"server.tls.key-file":           "server.key",
					"server.tls.ca-file":            "${server.tls.cert-file}",
					"server.tls.certificate-source": "server",
					"server.tls.cipher-suites":      "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,TLS_RSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_CBC_SHA",
					// Microservice
					"server.port":         "3030",
					"server.context-path": "/dna",
					"server.enabled":      "true",
				}),
			},
			want: &WebServerConfig{
				Enabled: true,
				Host:    "0.0.0.0",
				Port:    3030,
				Tls: certificate.TLSConfig{
					Enabled:           false,
					MinVersion:        "tls12",
					CertificateSource: "server",
					CipherSuites: []string{
						"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
						"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
						"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
						"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
						"TLS_RSA_WITH_AES_256_GCM_SHA384",
						"TLS_RSA_WITH_AES_256_CBC_SHA",
					},
				},
				Cors: CorsConfig{
					Enabled:              true,
					CustomAllowedHeaders: []string{},
					CustomExposedHeaders: []string{},
				},
				ContextPath:   "/dna",
				StaticPath:    "/www",
				StaticEnabled: true,
				TraceEnabled:  false,
				DebugEnabled:  false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewWebServerConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWebServerConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(testhelpers.Diff(tt.want, got))
			}
		})
	}
}
