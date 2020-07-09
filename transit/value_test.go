package transit

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestParseValue(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    Value
		wantErr bool
	}{
		{
			name:    "InvalidNotEnoughFields",
			value:   "::",
			wantErr: true,
		},
		{
			name:  "Empty",
			value: "1:22a342bf-3278-4126-9a02-f1ac0c9cf05f::",
			want: Value{
				version: "1",
				keyId:   types.MustParseUUID("22a342bf-3278-4126-9a02-f1ac0c9cf05f"),
			},
			wantErr: false,
		},
		{
			name:  "Plain",
			value: `1:22a342bf-3278-4126-9a02-f1ac0c9cf05f:p:["java.util.HashMap",{}]`,
			want: Value{
				version:   "1",
				keyId:     types.MustParseUUID("22a342bf-3278-4126-9a02-f1ac0c9cf05f"),
				encrypted: false,
				payload:   `["java.util.HashMap",{}]`,
			},
		},
		{
			name:  "Encrypted",
			value: `1:22a342bf-3278-4126-9a02-f1ac0c9cf05f:e:ABCDEF`,
			want: Value{
				version:   "1",
				keyId:     types.MustParseUUID("22a342bf-3278-4126-9a02-f1ac0c9cf05f"),
				encrypted: true,
				payload:   `ABCDEF`,
			},
		},
		{
			name:    "InvalidEncrypted",
			value:   `1:22a342bf-3278-4126-9a02-f1ac0c9cf05f:s:ABCDEF`,
			want:    Value{},
			wantErr: true,
		},
		{
			name:    "InvalidEmptyPlain",
			value:   `1:22a342bf-3278-4126-9a02-f1ac0c9cf05f:p:`,
			want:    Value{},
			wantErr: true,
		},
		{
			name:    "InvalidEmptyEncrypted",
			value:   `1:22a342bf-3278-4126-9a02-f1ac0c9cf05f:e:`,
			want:    Value{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseValue(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewValue(t *testing.T) {
	payload := map[string]*string{
		"key1": types.NewStringPtr("value1"),
		"key2": nil,
	}

	keyId := types.MustParseUUID("22a342bf-3278-4126-9a02-f1ac0c9cf05f")

	t.Run("Empty", func(t *testing.T) {
		value, err := NewValue(keyId, nil)
		assert.NoError(t, err)
		assert.Equal(t, "1", value.version)
		assert.Equal(t, keyId, value.KeyId())
		assert.Equal(t, false, value.encrypted)
		assert.Equal(t, `1:22a342bf-3278-4126-9a02-f1ac0c9cf05f::`, value.String())
	})

	t.Run("Payload", func(t *testing.T) {
		value, err := NewValue(keyId, payload)
		assert.NoError(t, err)
		assert.Equal(t, `["java.util.HashMap",{"key1":"value1","key2":null}]`, value.RawPayload())
		assert.Equal(t, "1", value.version)
		assert.Equal(t, keyId, value.KeyId())
		assert.Equal(t, false, value.encrypted)
		assert.Equal(t, `1:22a342bf-3278-4126-9a02-f1ac0c9cf05f:p:["java.util.HashMap",{"key1":"value1","key2":null}]`, value.String())
	})
}

func TestNewSecureValue(t *testing.T) {
	payload := `ABCD`
	keyId := types.MustParseUUID("22a342bf-3278-4126-9a02-f1ac0c9cf05f")

	value := NewSecureValue(keyId, payload)
	assert.Equal(t, payload, value.RawPayload())
	assert.Equal(t, "1", value.version)
	assert.Equal(t, true, value.encrypted)
	assert.Equal(t, `1:22a342bf-3278-4126-9a02-f1ac0c9cf05f:e:ABCD`, value.String())
}

func TestValue_String(t *testing.T) {
	keyId := types.MustParseUUID("22a342bf-3278-4126-9a02-f1ac0c9cf05f")
	tests := []struct {
		name  string
		value Value
		want  string
	}{
		{
			name: "Empty",
			value: Value{
				version: "1",
				keyId:   keyId,
			},
			want: "1:22a342bf-3278-4126-9a02-f1ac0c9cf05f::",
		},
		{
			name: "Plain",
			value: Value{
				version:   "1",
				keyId:     keyId,
				encrypted: false,
				payload:   "[]",
			},
			want: "1:22a342bf-3278-4126-9a02-f1ac0c9cf05f:p:[]",
		},
		{
			name: "Encrypted",
			value: Value{
				version:   "1",
				keyId:     keyId,
				encrypted: true,
				payload:   "ABCD",
			},
			want: "1:22a342bf-3278-4126-9a02-f1ac0c9cf05f:e:ABCD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.value.String()
			if got != tt.want {
				t.Errorf("String() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValue_KeyName(t *testing.T) {
	tests := []struct {
		name  string
		value Value
		want  string
	}{
		{
			name: "Arbitrary",
			value: Value{
				version: "1",
				keyId:   types.MustParseUUID("22a342bf-3278-4126-9a02-f1ac0c9cf05f"),
			},
			want: "22a342bf-3278-4126-9a02-f1ac0c9cf05f",
		},
		{
			name: "Empty",
			value: Value{
				version: "1",
				keyId:   types.EmptyUUID(),
			},
			want: "00000000-0000-0000-0000-000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.value.KeyName()
			if got != tt.want {
				t.Errorf("KeyName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValue_UsesKey(t *testing.T) {
	keyId := types.MustParseUUID("22a342bf-3278-4126-9a02-f1ac0c9cf05f")
	tests := []struct {
		name  string
		value Value
		key   types.UUID
		want  bool
	}{
		{
			name: "Match",
			value: Value{
				version: "1",
				keyId:   keyId,
			},
			key:  keyId,
			want: true,
		},
		{
			name: "MatchEmpty",
			value: Value{
				version: "1",
				keyId:   types.EmptyUUID(),
			},
			key:  types.EmptyUUID(),
			want: true,
		},
		{
			name: "NoMatch",
			value: Value{
				version: "1",
				keyId:   keyId,
			},
			key:  types.EmptyUUID(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.value.UsesKey(tt.key)
			if got != tt.want {
				t.Errorf("UsesKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValue_IsEmpty(t *testing.T) {
	tests := []struct {
		name  string
		value Value
		want  bool
	}{
		{
			name: "Empty",
			value: Value{
				version: "1",
				keyId:   types.EmptyUUID(),
			},
			want: true,
		},
		{
			name: "NotEmpty",
			value: Value{
				version: "1",
				keyId:   types.EmptyUUID(),
				payload: "[]",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.value.IsEmpty()
			if got != tt.want {
				t.Errorf("IsEmpty() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValue_Payload(t *testing.T) {
	tests := []struct {
		name    string
		value   Value
		want    map[string]*string
		wantErr bool
	}{
		{
			name: "Unencrypted",
			value: Value{
				version:   "1",
				keyId:     types.EmptyUUID(),
				encrypted: false,
				payload:   `["java.util.HashMap",{"key1":"value1"}]`,
			},
			want: map[string]*string{
				"key1": types.NewStringPtr("value1"),
			},
		},
		{
			name: "Empty",
			value: Value{
				version:   "1",
				keyId:     types.EmptyUUID(),
				encrypted: false,
				payload:   ``,
			},
			want: nil,
		},
		{
			name: "Encrypted",
			value: Value{
				version:   "1",
				keyId:     types.EmptyUUID(),
				encrypted: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.value.Payload()
			if (err != nil) != tt.wantErr {
				t.Errorf("Payload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Payload() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValue_RawPayload(t *testing.T) {
	tests := []struct {
		name  string
		value Value
		want  string
	}{
		{
			name: "Unencrypted",
			value: Value{
				version:   "1",
				keyId:     types.EmptyUUID(),
				encrypted: false,
				payload:   `["java.util.HashMap",{"key1":"value1"}]`,
			},
			want: `["java.util.HashMap",{"key1":"value1"}]`,
		},
		{
			name: "Empty",
			value: Value{
				version:   "1",
				keyId:     types.EmptyUUID(),
				encrypted: false,
				payload:   ``,
			},
			want: ``,
		},
		{
			name: "Encrypted",
			value: Value{
				version:   "1",
				keyId:     types.EmptyUUID(),
				encrypted: true,
				payload:   `ABC`,
			},
			want: `ABC`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.value.RawPayload()
			if got != tt.want {
				t.Errorf("RawPayload() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValue_WithEncryptedPayload(t *testing.T) {
	keyId := types.MustParseUUID("22a342bf-3278-4126-9a02-f1ac0c9cf05f")
	tests := []struct {
		name    string
		value   Value
		payload string
		want    Value
	}{
		{
			name: "PreviouslyEmpty",
			value: Value{
				version: "1",
				keyId:   keyId,
			},
			payload: `ABC`,
			want: Value{
				version:   "1",
				keyId:     keyId,
				encrypted: true,
				payload:   `ABC`,
			},
		},
		{
			name: "PreviouslyFull",
			value: Value{
				version:   "1",
				keyId:     keyId,
				encrypted: false,
				payload:   `[]`,
			},
			payload: `ABC`,
			want: Value{
				version:   "1",
				keyId:     keyId,
				encrypted: true,
				payload:   `ABC`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.value.WithEncryptedPayload(tt.payload)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithEncryptedPayload() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValue_WithDecryptedPayload(t *testing.T) {
	keyId := types.MustParseUUID("22a342bf-3278-4126-9a02-f1ac0c9cf05f")
	tests := []struct {
		name    string
		value   Value
		payload string
		want    Value
	}{
		{
			name: "PreviouslyEmpty",
			value: Value{
				version: "1",
				keyId:   keyId,
			},
			payload: `["java.util.HashMap",{"key1":"value1"}]`,
			want: Value{
				version:   "1",
				keyId:     keyId,
				encrypted: false,
				payload:   `["java.util.HashMap",{"key1":"value1"}]`,
			},
		},
		{
			name: "PreviouslyFull",
			value: Value{
				version:   "1",
				keyId:     keyId,
				encrypted: false,
				payload:   `["java.util.HashMap",{"key2":"value2"}]`,
			},
			payload: `["java.util.HashMap",{"key1":"value1"}]`,
			want: Value{
				version:   "1",
				keyId:     keyId,
				encrypted: false,
				payload:   `["java.util.HashMap",{"key1":"value1"}]`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.value.WithDecryptedPayload(tt.payload)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDecryptedPayload() got = %v, want %v", got, tt.want)
			}
		})
	}
}
