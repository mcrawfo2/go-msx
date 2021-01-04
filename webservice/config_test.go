package webservice

import (
	"crypto/tls"
	"reflect"
	"testing"
)

func TestParseCiphers(t *testing.T) {
	type args struct {
		cipherStr string
	}
	tests := []struct {
		name    string
		args    args
		want    []uint16
		wantErr bool
	}{
		{
			name: "Single",
			args: args{
				cipherStr: "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
			},
			want:    []uint16{tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305},
			wantErr: false,
		},
		{
			name: "Invalid",
			args: args{
				cipherStr: "XYZ",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Ordered",
			args: args{
				cipherStr: "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
			},
			want: []uint16{
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCiphers(tt.args.cipherStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCiphers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCiphers() got = %v, want %v", got, tt.want)
			}
		})
	}
}

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
				Tls: TLSConfig{
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
