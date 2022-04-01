// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package vaultprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"cto-github.cisco.com/NFV-BU/go-msx/transit"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/vault"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestProvider_CreateKey(t *testing.T) {
	type fields struct {
		cfg *Config
	}
	type args struct {
		ctx     context.Context
		keyName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Disabled",
			fields: fields{
				cfg: func() *Config {
					cfg, _ := NewEncryptionConfig(configtest.NewInMemoryConfig(nil))
					return cfg
				}(),
			},
			args: args{
				ctx: func() context.Context {
					mockConnection := new(vault.MockConnection)
					return vault.ContextWithConnection(context.Background(), mockConnection)
				}(),
				keyName: "my-key",
			},
		},
		{
			name: "Enabled",
			fields: fields{
				cfg: func() *Config {
					cfg, _ := NewEncryptionConfig(configtest.NewInMemoryConfig(map[string]string{
						"per-tenant-encryption.enabled": "true",
					}))
					return cfg
				}(),
			},
			args: args{
				ctx: func() context.Context {
					mockConnection := new(vault.MockConnection)
					mockConnection.
						On("CreateTransitKey",
							mock.AnythingOfType("*context.valueCtx"),
							"my-key",
							mock.AnythingOfType("vault.CreateTransitKeyRequest")).
						Return(nil)
					return vault.ContextWithConnection(context.Background(), mockConnection)
				}(),
				keyName: "my-key",
			},
		},
		{
			name: "Error",
			fields: fields{
				cfg: func() *Config {
					cfg, _ := NewEncryptionConfig(configtest.NewInMemoryConfig(map[string]string{
						"per-tenant-encryption.enabled": "true",
					}))
					return cfg
				}(),
			},
			args: args{
				ctx: func() context.Context {
					mockConnection := new(vault.MockConnection)
					mockConnection.
						On("CreateTransitKey",
							mock.AnythingOfType("*context.valueCtx"),
							"my-key",
							mock.AnythingOfType("vault.CreateTransitKeyRequest")).
						Return(errors.New("error"))
					return vault.ContextWithConnection(context.Background(), mockConnection)
				}(),
				keyName: "my-key",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Provider{
				cfg: tt.fields.cfg,
			}
			if err := p.CreateKey(tt.args.ctx, tt.args.keyName); (err != nil) != tt.wantErr {
				t.Errorf("CreateKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProvider_Decrypt(t *testing.T) {
	keyName := "a4b1d114-22fb-4135-9a57-a5b60d06d175"
	keyId := types.MustParseUUID(keyName)

	emptyValue, _ := transit.NewValue(keyId, map[string]*string{})
	nonEmptyCiphertext := "ABC"
	nonEmptyPlaintext := `["java.util.HashMap",{"key":"value"}]`
	nonEmptyDecryptedValue, _ := transit.NewValue(keyId, map[string]*string{
		"key": types.NewOptionalStringFromString("value").Ptr(),
	})
	nonEmptyEncryptedValue := nonEmptyDecryptedValue.WithEncryptedPayload(nonEmptyCiphertext)

	type fields struct {
		cfg *Config
	}
	type args struct {
		ctx         context.Context
		secureValue transit.Value
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantValue transit.Value
		wantErr   bool
	}{
		{
			name: "EmptyValue",
			args: args{
				ctx:         vault.ContextWithConnection(context.Background(), new(vault.MockConnection)),
				secureValue: emptyValue,
			},
			wantValue: emptyValue,
		},
		{
			name: "Disabled",
			fields: fields{
				cfg: func() *Config {
					cfg, _ := NewEncryptionConfig(configtest.NewInMemoryConfig(nil))
					return cfg
				}(),
			},
			args: args{
				ctx:         vault.ContextWithConnection(context.Background(), new(vault.MockConnection)),
				secureValue: nonEmptyDecryptedValue,
			},
			wantValue: nonEmptyDecryptedValue,
		},
		{
			name: "Decrypted",
			fields: fields{
				cfg: func() *Config {
					cfg, _ := NewEncryptionConfig(configtest.NewInMemoryConfig(map[string]string{
						"per-tenant-encryption.enabled": "true",
					}))
					return cfg
				}(),
			},
			args: args{
				ctx: func() context.Context {
					mockConnection := new(vault.MockConnection)
					mockConnection.
						On("TransitDecrypt",
							mock.AnythingOfType("*context.valueCtx"),
							keyName,
							nonEmptyCiphertext).
						Return(nonEmptyPlaintext, nil)

					return vault.ContextWithConnection(context.Background(), mockConnection)
				}(),
				secureValue: nonEmptyEncryptedValue,
			},
			wantValue: nonEmptyDecryptedValue,
		},
		{
			name: "DecryptError",
			fields: fields{
				cfg: func() *Config {
					cfg, _ := NewEncryptionConfig(configtest.NewInMemoryConfig(map[string]string{
						"per-tenant-encryption.enabled": "true",
					}))
					return cfg
				}(),
			},
			args: args{
				ctx: func() context.Context {
					mockConnection := new(vault.MockConnection)
					mockConnection.
						On("TransitDecrypt",
							mock.AnythingOfType("*context.valueCtx"),
							keyName,
							nonEmptyCiphertext).
						Return("", errors.New("error"))

					return vault.ContextWithConnection(context.Background(), mockConnection)
				}(),
				secureValue: nonEmptyEncryptedValue,
			},
			wantValue: transit.Value{},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Provider{
				cfg: tt.fields.cfg,
			}
			gotValue, err := p.Decrypt(tt.args.ctx, tt.args.secureValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("Decrypt() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

func TestProvider_Encrypt(t *testing.T) {
	keyName := "a4b1d114-22fb-4135-9a57-a5b60d06d175"
	keyId := types.MustParseUUID(keyName)

	emptyValue, _ := transit.NewValue(keyId, map[string]*string{})
	nonEmptyCiphertext := "ABC"
	nonEmptyPlainText := `["java.util.HashMap",{"key":"value"}]`
	nonEmptyValue, _ := transit.NewValue(keyId, map[string]*string{
		"key": types.NewOptionalStringFromString("value").Ptr(),
	})
	nonEmptyEncryptedValue := nonEmptyValue.WithEncryptedPayload(nonEmptyCiphertext)

	type fields struct {
		cfg *Config
	}
	type args struct {
		ctx   context.Context
		value transit.Value
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantSecureValue transit.Value
		wantErr         bool
	}{
		{
			name: "EmptyValue",
			args: args{
				ctx:   vault.ContextWithConnection(context.Background(), new(vault.MockConnection)),
				value: emptyValue,
			},
			wantSecureValue: emptyValue,
		},
		{
			name: "Disabled",
			fields: fields{
				cfg: func() *Config {
					cfg, _ := NewEncryptionConfig(configtest.NewInMemoryConfig(nil))
					return cfg
				}(),
			},
			args: args{
				ctx:   vault.ContextWithConnection(context.Background(), new(vault.MockConnection)),
				value: nonEmptyValue,
			},
			wantSecureValue: nonEmptyValue,
		},
		{
			name: "Encrypted",
			fields: fields{
				cfg: func() *Config {
					cfg, _ := NewEncryptionConfig(configtest.NewInMemoryConfig(map[string]string{
						"per-tenant-encryption.enabled": "true",
					}))
					return cfg
				}(),
			},
			args: args{
				ctx: func() context.Context {
					mockConnection := new(vault.MockConnection)
					mockConnection.
						On("TransitEncrypt",
							mock.AnythingOfType("*context.valueCtx"),
							keyName,
							nonEmptyPlainText).
						Return(nonEmptyCiphertext, nil)

					return vault.ContextWithConnection(context.Background(), mockConnection)
				}(),
				value: nonEmptyValue,
			},
			wantSecureValue: nonEmptyEncryptedValue,
		},
		{
			name: "EncryptError",
			fields: fields{
				cfg: func() *Config {
					cfg, _ := NewEncryptionConfig(configtest.NewInMemoryConfig(map[string]string{
						"per-tenant-encryption.enabled": "true",
					}))
					return cfg
				}(),
			},
			args: args{
				ctx: func() context.Context {
					mockConnection := new(vault.MockConnection)
					mockConnection.
						On("TransitEncrypt",
							mock.AnythingOfType("*context.valueCtx"),
							keyName,
							nonEmptyPlainText).
						Return("", errors.New("error"))

					return vault.ContextWithConnection(context.Background(), mockConnection)
				}(),
				value: nonEmptyValue,
			},
			wantSecureValue: transit.Value{},
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Provider{
				cfg: tt.fields.cfg,
			}
			gotSecureValue, err := p.Encrypt(tt.args.ctx, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSecureValue, tt.wantSecureValue) {
				t.Errorf("Encrypt() gotSecureValue = %v, want %v", gotSecureValue, tt.wantSecureValue)
			}
		})
	}
}

func TestRegisterVaultTransitProvider(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ConfigError",
			args: args{
				ctx: configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
					"per-tenant-encryption.enabled": "falsy",
				}),
			},
			wantErr: true,
		},
		{
			name: "Success",
			args: args{
				ctx: configtest.ContextWithNewInMemoryConfig(context.Background(), nil),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RegisterVaultTransitProvider(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("RegisterVaultTransitProvider() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProvider_DecryptBulk(t *testing.T) {
	keyName := "a4b1d114-22fb-4135-9a57-a5b60d06d175"
	keyId := types.MustParseUUID(keyName)

	cfg := &Config{Enabled: true}

	emptyValue, _ := transit.NewValue(keyId, map[string]*string{})
	nonEmptyCiphertext := "ABC"
	nonEmptyPlaintext := `["java.util.HashMap",{"key":"value"}]`
	nonEmptyDecryptedValue, _ := transit.NewValue(keyId, map[string]*string{
		"key": types.NewOptionalStringFromString("value").Ptr(),
	})
	nonEmptyEncryptedValue := nonEmptyDecryptedValue.WithEncryptedPayload(nonEmptyCiphertext)

	tests := []struct {
		name         string
		cfg          *Config
		ctx          context.Context
		secureValues []transit.Value
		wantValues   []transit.Value
		wantErr      bool
	}{
		{
			name:         "EmptyValue",
			cfg:          cfg,
			ctx:          vault.ContextWithConnection(context.Background(), new(vault.MockConnection)),
			secureValues: []transit.Value{emptyValue},
			wantValues:   []transit.Value{emptyValue},
			wantErr:      false,
		},
		{
			name:         "DecryptedValue",
			cfg:          cfg,
			ctx:          vault.ContextWithConnection(context.Background(), new(vault.MockConnection)),
			secureValues: []transit.Value{nonEmptyDecryptedValue},
			wantValues:   []transit.Value{nonEmptyDecryptedValue},
			wantErr:      false,
		},
		{
			name: "Success",
			cfg:  cfg,
			ctx: func() context.Context {
				mockConnection := new(vault.MockConnection)
				mockConnection.
					On("TransitBulkDecrypt",
						mock.AnythingOfType("*context.valueCtx"),
						keyName,
					nonEmptyCiphertext).
					Return([]string{
						nonEmptyPlaintext,
					}, nil)
				return vault.ContextWithConnection(context.Background(), mockConnection)
			}(),
			secureValues: []transit.Value{nonEmptyEncryptedValue},
			wantValues:   []transit.Value{nonEmptyDecryptedValue},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Provider{
				cfg: tt.cfg,
			}
			gotValues, err := p.DecryptBulk(tt.ctx, tt.secureValues)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptBulk() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotValues, tt.wantValues) {
				t.Errorf("DecryptBulk() gotValues = %v, want %v", gotValues, tt.wantValues)
			}
		})
	}
}
