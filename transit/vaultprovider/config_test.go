package vaultprovider

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestNewEncryptionConfig(t *testing.T) {
	var falsehood = false

	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "Embedded",
			args: args{
				cfg: configtest.NewStaticConfig(map[string]string{
					"per-tenant-encryption.enabled":                               "false",
					"per-tenant-encryption.always-create-keys":                    "false",
					"per-tenant-encryption.key-properties.type":                   "aes256-gcm96",
					"per-tenant-encryption.key-properties.exportable":             "false",
					"per-tenant-encryption.key-properties.allow-plaintext-backup": "false",
				}),
			},
			want: &Config{
				Enabled:          false,
				AlwaysCreateKeys: false,
				KeyProperties: KeyPropertiesConfig{
					Type:                 "aes256-gcm96",
					Exportable:           &falsehood,
					AllowPlaintextBackup: &falsehood,
				},
			},
		},
		{
			name: "Custom",
			args: args{
				cfg: configtest.NewStaticConfig(map[string]string{
					"per-tenant-encryption.enabled":                               "false",
					"per-tenant-encryption.always-create-keys":                    "false",
					"per-tenant-encryption.key-properties.type":                   "aes256-gcm96",
					"per-tenant-encryption.key-properties.exportable":             "false",
					"per-tenant-encryption.key-properties.allow-plaintext-backup": "false",
				}),
			},
			want: &Config{
				Enabled:          false,
				AlwaysCreateKeys: false,
				KeyProperties: KeyPropertiesConfig{
					Type:                 "aes256-gcm96",
					Exportable:           &falsehood,
					AllowPlaintextBackup: &falsehood,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEncryptionConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEncryptionConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEncryptionConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
