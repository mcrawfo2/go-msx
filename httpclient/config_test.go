package httpclient

import (
	"context"
	"crypto/tls"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestNewClientConfig(t *testing.T) {
	cfg := config.NewConfig(
		config.NewStatic("static", map[string]string{
			"http.client.timeout":       "90s",
			"http.client.idle-timeout":  "120s",
			"http.client.tls-insecure":  "false",
			"http.client.local-ca-file": "ca.crt",
			"http.client.cert-file":     "client.crt",
			"http.client.key-file":      "client.key",
		}))
	err := cfg.Load(context.Background())
	assert.NoError(t, err)

	got, err := NewClientConfig(cfg)
	assert.NoError(t, err)
	assert.Equal(t, 90*time.Second, got.Timeout)
	assert.Equal(t, 120*time.Second, got.IdleTimeout)
	assert.Equal(t, false, got.TlsInsecure)
	assert.Equal(t, "ca.crt", got.LocalCaFile)
	assert.Equal(t, "client.crt", got.CertFile)
	assert.Equal(t, "client.key", got.KeyFile)
}

func TestNewTlsConfig(t *testing.T) {
	rootCAs, _ := getRootCAs(&ClientConfig{})

	tests := []struct {
		name         string
		clientConfig ClientConfig
		want         *tls.Config
		wantConfig   bool
		wantErr      bool
	}{
		{
			name: "Success",
			clientConfig: ClientConfig{
				Timeout:     90 * time.Second,
				IdleTimeout: 120 * time.Second,
				TlsInsecure: false,
			},
			wantConfig: true,
			want: &tls.Config{
				InsecureSkipVerify: false,
				RootCAs:            rootCAs,
			},
			wantErr: false,
		},
		{
			name: "BadLocalCa",
			clientConfig: ClientConfig{
				LocalCaFile: "testdata/none.crt",
			},
			wantConfig: false,
			wantErr: true,
		},
		{
			name: "BadClientCert",
			clientConfig: ClientConfig{
				CertFile: "testdata/none.crt",
			},
			wantConfig: false,
			wantErr: true,
		},
		{
			name: "ValidClientCert",
			clientConfig: ClientConfig{
				CertFile: "testdata/server.crt",
				KeyFile: "testdata/server.key",
			},
			wantConfig: true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTlsConfig(&tt.clientConfig)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewTlsConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got != nil) != tt.wantConfig {
				t.Errorf("NewTlsConfig() got = %v, wantConfig %v", err, tt.wantConfig)
				return
			}

			if tt.want != nil {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewTlsConfig() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_getClientCerts(t *testing.T) {
	tests := []struct {
		name     string
		cfg      ClientConfig
		wantCert bool
		wantErr  bool
		wantLog  log.Check
	}{
		{
			name:     "NotSpecified",
			cfg:      ClientConfig{},
			wantCert: false,
			wantErr:  false,
			wantLog: log.Check{
				Validators: []log.EntryPredicate{
					log.HasLevel(logrus.WarnLevel),
				},
			},
		},
		{
			name: "InvalidSpecified",
			cfg: ClientConfig{
				CertFile: "testdata/server.crt",
			},
			wantCert: false,
			wantErr:  true,
		},
		{
			name: "BadSpecified",
			cfg: ClientConfig{
				CertFile: "testdata/server.crt",
				KeyFile:  "testdata/server2.key",
			},
			wantCert: false,
			wantErr:  true,
		},
		{
			name: "Specified",
			cfg: ClientConfig{
				CertFile: "testdata/server.crt",
				KeyFile:  "testdata/server.key",
			},
			wantCert: true,
			wantErr:  false,
			wantLog: log.Check{
				Validators: []log.EntryPredicate{
					log.HasLevel(logrus.InfoLevel),
					log.HasMessage(`Loaded client certificate from "testdata/server.crt"`),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recording := log.RecordLogging()

			got, err := getClientCerts(&tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("getClientCerts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.wantCert {
				t.Errorf("getClientCerts() got = %v, want %v", got, tt.wantCert)
			}

			errors := tt.wantLog.Check(recording)
			assert.Len(t, errors, 0)
		})
	}
}

func Test_getRootCAs(t *testing.T) {
	tests := []struct {
		name     string
		cfg      ClientConfig
		wantPool bool
		wantErr  bool
		wantLog  log.Check
	}{
		{
			name:     "NoLocalCA",
			cfg:      ClientConfig{},
			wantPool: true,
			wantErr:  false,
		},
		{
			name: "LocalCA",
			cfg: ClientConfig{
				LocalCaFile: "testdata/server.crt",
			},
			wantPool: true,
			wantErr:  false,
			wantLog: log.Check{
				Validators: []log.EntryPredicate{
					log.HasLevel(logrus.InfoLevel),
					log.HasMessage(`Added certificates from "testdata/server.crt" to RootCAs`),
				},
			},
		},
		{
			name: "BadLocalCA",
			cfg: ClientConfig{
				LocalCaFile: "testdata/none.crt",
			},
			wantPool: false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recording := log.RecordLogging()

			got, err := getRootCAs(&tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("getClientCerts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.wantPool {
				t.Errorf("getClientCerts() got = %v, want %v", got, tt.wantPool)
			}

			errors := tt.wantLog.Check(recording)
			assert.Len(t, errors, 0)
		})
	}
}
