package fileprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProviderFactory_Name(t *testing.T) {
	assert.Equal(t, "file", new(ProviderFactory).Name())
}

func TestProviderFactory_New(t *testing.T) {
	ctx := configtest.ContextWithNewStaticConfig(
		context.Background(),
		map[string]string{"certificate.source.success.provider": "file"})

	tests := []struct {
		name         string
		configRoot   string
		wantProvider bool
		wantErr      bool
	}{
		{
			name:         "Success",
			configRoot:   "certificate.source.success",
			wantProvider: true,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ProviderFactory{}
			got, err := p.New(ctx, tt.configRoot)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if tt.wantProvider != (got != nil) {
				t.Errorf("New() got = %v, want nil = %v", got, tt.wantProvider)
			}
		})
	}
}

func TestProvider_GetCertificate(t *testing.T) {
	type fields struct {
		cfg ProviderConfig
	}
	tests := []struct {
		name     string
		fields   fields
		wantCert bool
		wantErr  bool
	}{
		{
			name: "Success",
			fields: fields{
				cfg: ProviderConfig{
					CertFile: "testdata/server.crt",
					KeyFile:  "testdata/server.key",
				},
			},
			wantCert: true,
		},
		{
			name: "NoSuchCertFile",
			fields: fields{
				cfg: ProviderConfig{
					CertFile: "testdata/missing.crt",
					KeyFile:  "testdata/server.key",
				},
			},
			wantCert: false,
			wantErr:  true,
		},
		{
			name: "NoSuchKeyFile",
			fields: fields{
				cfg: ProviderConfig{
					CertFile: "testdata/server.crt",
					KeyFile:  "testdata/missing.key",
				},
			},
			wantCert: false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Provider{
				cfg: tt.fields.cfg,
			}
			got, err := f.GetCertificate(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCertificate() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if tt.wantCert != (got != nil) {
				t.Errorf("GetCertificate() got = %v, wantCert = %v", got, tt.wantCert)
			}
		})
	}
}

func TestProvider_Renewable(t *testing.T) {
	assert.Equal(t, false, new(Provider).Renewable())
}

func TestRegisterFactory(t *testing.T) {
	err := RegisterFactory(context.Background())
	assert.NoError(t, err)
}
