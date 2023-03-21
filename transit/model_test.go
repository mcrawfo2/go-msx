// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package transit

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func testModelEncrypterFactory() EncrypterFactory {
	return func(ctx context.Context, keyName types.UUID) Encrypter {
		encrypter := new(MockEncrypter)

		encrypter.
			On("Encrypt", map[string]*string{
				"key1": types.NewStringPtr("value1"),
			}).
			Return(`ABCD`, true, nil)

		encrypter.
			On("Encrypt", map[string]*string{
				"key1": types.NewStringPtr("value2"),
			}).
			Return(`EFGH`, true, nil)

		encrypter.
			On("Encrypt", map[string]*string{
				"key1": types.NewStringPtr("value1"),
				"key2": types.NewStringPtr("value2"),
			}).
			Return(`IJKL`, true, nil)

		encrypter.
			On("Encrypt", map[string]*string{
				"key2": types.NewStringPtr("value1"),
			}).
			Return(`1:22a342bf-3278-4126-9a02-f1ac0c9cf05f:p:["java.util.HashMap",{"key1":"value1"}]`, false, nil)

		encrypter.
			On("Encrypt", map[string]*string{}).
			Return(``, false, errors.New("encryption error"))

		encrypter.
			On("Encrypt", map[string]*string{
				"key3": types.NewStringPtr("value3"),
			}).
			Return(``, false, errors.New("encryption error"))

		encrypter.
			On("Decrypt", `ABCD`).
			Return(map[string]*string{
				"key1": types.NewStringPtr("value1"),
			}, nil)

		encrypter.
			On("Decrypt", `EFGH`).
			Return(map[string]*string{
				"key1": types.NewStringPtr("value2"),
			}, nil)

		encrypter.
			On("Decrypt", `MNOP`).
			Return(nil, errors.New("decryption error"))

		return encrypter
	}
}

func TestSecureData_Field(t *testing.T) {
	ctx := context.Background()
	keyId := types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`)
	encrypterFactory = testModelEncrypterFactory()

	tests := []struct {
		name      string
		value     SecureData
		fieldName string
		want      *string
		wantErr   bool
	}{
		{
			name: "Plain",
			value: SecureData{
				ctx:   ctx,
				dirty: false,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: false,
					payload:   `["java.util.HashMap",{"key1":"value1"}]`,
				},
				payload: nil,
			},
			fieldName: "key1",
			want:      types.NewStringPtr("value1"),
			wantErr:   false,
		},
		{
			name: "Encrypted",
			value: SecureData{
				ctx:   ctx,
				dirty: false,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: true,
					payload:   `ABCD`,
				},
				payload: nil,
			},
			fieldName: "key1",
			want:      types.NewStringPtr("value1"),
			wantErr:   false,
		},
		{
			name: "Missing",
			value: SecureData{
				ctx:   ctx,
				dirty: false,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: true,
					payload:   `ABCD`,
				},
				payload: nil,
			},
			fieldName: "key2",
			want:      nil,
			wantErr:   false,
		},
		{
			name: "InvalidEncrypted",
			value: SecureData{
				ctx:   ctx,
				dirty: false,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: true,
					payload:   `MNOP`,
				},
				payload: nil,
			},
			fieldName: "key1",
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.value
			sp := &s
			got, err := sp.Field(ctx, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Field() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Field() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecureData_SetField(t *testing.T) {
	ctx := context.Background()
	keyId := types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`)
	encrypterFactory = testModelEncrypterFactory()

	tests := []struct {
		name       string
		value      SecureData
		fieldName  string
		fieldValue *string
		want       SecureData
		wantErr    bool
	}{
		{
			name:       "Empty",
			value:      SecureData{ctx: ctx, keyId: keyId},
			fieldName:  "key1",
			fieldValue: types.NewStringPtr("value1"),
			want: SecureData{
				ctx:   ctx,
				dirty: true,
				keyId: keyId,
				payload: map[string]*string{
					"key1": types.NewStringPtr("value1"),
				},
			},
		},
		{
			name: "Plain",
			value: SecureData{
				ctx:   ctx,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: false,
					payload:   `["java.util.HashMap",{"key1":"value1"}]`,
				},
			},
			fieldName:  "key1",
			fieldValue: types.NewStringPtr("value2"),
			want: SecureData{
				ctx:   ctx,
				dirty: true,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: false,
					payload:   `["java.util.HashMap",{"key1":"value1"}]`,
				},
				payload: map[string]*string{
					"key1": types.NewStringPtr("value2"),
				},
			},
		},
		{
			name: "NoChange",
			value: SecureData{
				ctx:   ctx,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: false,
					payload:   `["java.util.HashMap",{"key1":"value1"}]`,
				},
			},
			fieldName:  "key1",
			fieldValue: types.NewStringPtr("value1"),
			want: SecureData{
				ctx:   ctx,
				dirty: false,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: false,
					payload:   `["java.util.HashMap",{"key1":"value1"}]`,
				},
				payload: map[string]*string{
					"key1": types.NewStringPtr("value1"),
				},
			},
		},
		{
			name: "Encrypted",
			value: SecureData{
				ctx:   ctx,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: true,
					payload:   `ABCD`,
				},
			},
			fieldName:  "key1",
			fieldValue: types.NewStringPtr("value2"),
			want: SecureData{
				ctx:   ctx,
				dirty: true,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: false,
					payload:   `["java.util.HashMap",{"key1":"value1"}]`,
				},
				payload: map[string]*string{
					"key1": types.NewStringPtr("value2"),
				},
			},
		},
		{
			name: "DecryptError",
			value: SecureData{
				ctx:   ctx,
				dirty: false,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: true,
					payload:   `MNOP`,
				},
			},
			fieldName:  "key3",
			fieldValue: types.NewStringPtr("value3"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s = tt.value
			var sp = &s
			var err error

			sp, err = sp.SetField(ctx, tt.fieldName, tt.fieldValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetField() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.want, s) {
				t.Errorf("SetField() got = %+v want %+v", s, tt.want)
			}
		})
	}
}

func TestSecureData_SetKeyId(t *testing.T) {
	ctx := context.Background()
	keyId := types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`)
	value := SecureData{}

	value.SetKeyId(ctx, keyId)
	assert.Equal(t, value.ctx, ctx)
	assert.Equal(t, value.keyId, keyId)
}

func TestSecureData_KeyId(t *testing.T) {
	ctx := context.Background()
	keyId := types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`)
	value := SecureData{
		ctx:   ctx,
		keyId: keyId,
	}

	assert.Equal(t, keyId, value.KeyId())
}

func TestWithSecureData_SetSecureValue(t *testing.T) {
	ctx := context.Background()
	keyId := types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`)

	tests := []struct {
		name       string
		value      WithSecureData
		keyId      types.UUID
		fieldName  string
		fieldValue *string
		want       WithSecureData
		wantErr    bool
	}{
		{
			name:    "EmptyKey",
			keyId:   types.EmptyUUID(),
			wantErr: true,
		},
		{
			name:    "InvalidKey",
			keyId:   types.UUID([]byte{1, 2, 3}),
			wantErr: true,
		},
		{
			name:       "CreateSecureData",
			keyId:      keyId,
			fieldName:  "key1",
			fieldValue: types.NewStringPtr("value1"),
			want: WithSecureData{
				SecureData: &SecureData{
					ctx:    ctx,
					dirty:  true,
					keyId:  keyId,
					secure: Value{},
					payload: map[string]*string{
						"key1": types.NewStringPtr("value1"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ws = tt.value
			var wsp = &ws
			var err error

			err = wsp.SetSecureValue(ctx, tt.keyId, tt.fieldName, tt.fieldValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetSecureValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.want, ws) {
				t.Errorf("SetSecureValue() got = %+v want %+v", ws, tt.want)
			}
		})
	}
}

func TestWithSecureData_SecureValue(t *testing.T) {
	ctx := context.Background()
	keyId := types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`)

	tests := []struct {
		name      string
		value     WithSecureData
		keyId     types.UUID
		fieldName string
		want      string
		wantErr   bool
	}{
		{
			name:      "ValueExists",
			keyId:     keyId,
			fieldName: "key1",
			value: WithSecureData{
				SecureData: &SecureData{
					ctx:    ctx,
					dirty:  true,
					keyId:  keyId,
					secure: Value{},
					payload: map[string]*string{
						"key1": types.NewStringPtr("value1"),
					},
				},
			},
			want: "value1",
		},
		{
			name:      "ValueNotExists",
			keyId:     keyId,
			fieldName: "key1",
			value: WithSecureData{
				SecureData: &SecureData{
					ctx:     ctx,
					dirty:   true,
					keyId:   keyId,
					secure:  Value{},
					payload: map[string]*string{},
				},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ws = tt.value
			var wsp = &ws
			var err error
			var got string

			got, err = wsp.SecureValue(ctx, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("SecureValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("SecureValue() got = %+v want %+v", got, tt.want)
			}
		})
	}
}

func TestWithSecureData_SecureOptionalValue(t *testing.T) {
	ctx := context.Background()
	keyId := types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`)

	tests := []struct {
		name      string
		value     WithSecureData
		keyId     types.UUID
		fieldName string
		want      types.OptionalString
		wantErr   bool
	}{
		{
			name:      "ValueExists",
			keyId:     keyId,
			fieldName: "key1",
			value: WithSecureData{
				SecureData: &SecureData{
					ctx:    ctx,
					dirty:  true,
					keyId:  keyId,
					secure: Value{},
					payload: map[string]*string{
						"key1": types.NewStringPtr("value1"),
					},
				},
			},
			want: types.NewOptionalStringFromString("value1"),
		},
		{
			name:      "ValueNotExists",
			keyId:     keyId,
			fieldName: "key1",
			value: WithSecureData{
				SecureData: &SecureData{
					ctx:     ctx,
					dirty:   true,
					keyId:   keyId,
					secure:  Value{},
					payload: map[string]*string{},
				},
			},
			want: types.NewOptionalString(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ws = tt.value
			var wsp = &ws
			var err error
			var got types.OptionalString

			got, err = wsp.SecureOptionalValue(ctx, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("SecureOptionalValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SecureOptionalValue() got = %+v want %+v", got, tt.want)
			}
		})
	}
}
