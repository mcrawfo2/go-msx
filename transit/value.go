package transit

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"github.com/pkg/errors"
	"strings"
)

const (
	valueFieldSeparator = ":"
	valueFieldCount     = 4
	valueVersion1       = "1"
	valueTypeEncrypted  = "e"
	valueTypePlain      = "p"
	valueTypeEmpty      = ""

	payloadJsonPrefix = `["java.util.HashMap",`
	payloadJsonSuffix = `]`
)

var (
	ErrValueSerializationInvalid = errors.New("Invalid serialization of value")
	ErrValueSerializationUnknown = errors.New("Unknown serialization type")
	ErrValueEncrypted            = errors.New("Value is encrypted")
)

type Value struct {
	version   string
	keyId     types.UUID
	encrypted bool
	payload   string
}

func (v Value) String() string {
	valueType := valueTypeEncrypted
	if 0 == len(v.payload) {
		valueType = valueTypeEmpty
	} else if !v.encrypted {
		valueType = valueTypePlain
	}

	return strings.Join([]string{
		v.version,
		strings.ToLower(v.keyId.String()),
		valueType,
		v.payload,
	}, valueFieldSeparator)
}

func (v Value) KeyName() string {
	return strings.ToLower(v.keyId.String())
}

func (v Value) KeyId() types.UUID {
	return v.keyId
}

func (v Value) UsesKey(key types.UUID) bool {
	return key.String() == v.keyId.String()
}

func (v Value) IsEmpty() bool {
	return len(v.payload) == 0
}

func (v Value) Payload() (map[string]*string, error) {
	if v.encrypted {
		return nil, ErrValueEncrypted
	}

	if v.IsEmpty() {
		return nil, nil
	}

	return deserializePayload(v.payload)
}

func (v Value) RawPayload() string {
	return v.payload
}

func (v Value) WithEncryptedPayload(payload string) Value {
	return Value{
		version:   v.version,
		keyId:     v.keyId,
		encrypted: true,
		payload:   payload,
	}
}

func (v Value) WithDecryptedPayload(payload string) Value {
	return Value{
		version:   v.version,
		keyId:     v.keyId,
		encrypted: false,
		payload:   payload,
	}
}

func deserializePayload(payload string) (map[string]*string, error) {
	if !strings.HasPrefix(payload, payloadJsonPrefix) || !strings.HasSuffix(payload, payloadJsonSuffix) {
		return nil, ErrValueSerializationUnknown
	}

	// golang interface{} slices are unhelpful here
	payload = strings.TrimPrefix(payload, payloadJsonPrefix)
	payload = strings.TrimSuffix(payload, payloadJsonSuffix)

	var result map[string]*string
	if err := json.Unmarshal([]byte(payload), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func serializePayload(payload map[string]*string) (string, error) {
	if len(payload) == 0 {
		return "", nil
	}

	payloadSlice := []interface{}{
		"java.util.HashMap",
		payload,
	}
	payloadBytes, err := json.Marshal(payloadSlice)
	if err != nil {
		return "", err
	}
	return string(payloadBytes), nil
}

func NewValue(keyName types.UUID, payload map[string]*string) (Value, error) {
	payloadJson, err := serializePayload(payload)
	if err != nil {
		return Value{}, err
	}

	return Value{
		version:   valueVersion1,
		keyId:     keyName,
		encrypted: false,
		payload:   payloadJson,
	}, nil
}

func NewSecureValue(keyName types.UUID, securePayload string) Value {
	return Value{
		version:   valueVersion1,
		keyId:     keyName,
		encrypted: true,
		payload:   securePayload,
	}
}

func ParseValue(value string) (Value, error) {
	parts := strings.SplitN(value, valueFieldSeparator, valueFieldCount)
	if len(parts) != valueFieldCount {
		return Value{}, ErrValueSerializationInvalid
	}

	if parts[0] != valueVersion1 {
		return Value{}, ErrValueSerializationInvalid
	}

	keyName, err := types.ParseUUID(parts[1])
	if err != nil {
		return Value{}, ErrValueSerializationInvalid
	}

	if len(parts[2]) == 0 && len(parts[3]) == 0 {
		return Value{
			version:   parts[0],
			keyId:     keyName,
			encrypted: false,
			payload:   "",
		}, nil
	}

	if (parts[2] != valueTypeEncrypted) &&
		(parts[2] != valueTypePlain) {
		return Value{}, ErrValueSerializationInvalid
	}

	if len(parts[3]) == 0 {
		return Value{}, ErrValueSerializationInvalid
	}

	return Value{
		version:   parts[0],
		keyId:     keyName,
		encrypted: parts[2] == valueTypeEncrypted,
		payload:   parts[3],
	}, nil
}
