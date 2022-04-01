// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package certificate

import (
	"crypto/tls"
	"reflect"
	"testing"
)

func TestTlsConfig_Ciphers(t *testing.T) {
	type args struct {
		ciphers []string
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
				ciphers: []string{"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305"},
			},
			want:    []uint16{tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305},
			wantErr: false,
		},
		{
			name: "Invalid",
			args: args{
				ciphers: []string{"XYZ"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Ordered",
			args: args{
				ciphers: []string{
					"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
					"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
				},
			},
			want: []uint16{
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tlsConfig := TLSConfig{CipherSuites: tt.args.ciphers}
			got, err := tlsConfig.Ciphers()
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
