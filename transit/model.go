// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package transit

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"database/sql/driver"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
)

var ErrKeyNotSet = errors.New("Key Id not set for encryption")
var ErrDataInvalid = errors.New("Data not valid type")

type SecureData struct {
	ctx     context.Context
	dirty   bool
	keyId   types.UUID
	secure  Value
	payload map[string]*string
}

func (s *SecureData) cleanValue() (string, error) {
	if s == nil || s.keyId == nil {
		return "", ErrKeyNotSet
	}

	// Lazy encrypt on dirty write
	if s.dirty {
		securePayload, encrypted, err := encrypterFactory.
			Create(s.ctx, s.keyId).
			Encrypt(s.payload)
		if err != nil {
			return "", err
		}
		if encrypted {
			s.secure = NewSecureValue(s.keyId, securePayload)
		} else {
			s.secure, err = NewValue(s.keyId, s.payload)
			if err != nil {
				return "", err
			}
		}
		s.dirty = false
	}

	return s.secure.String(), nil
}

func (s *SecureData) UnmarshalCQL(_ gocql.TypeInfo, data []byte) (err error) {
	s.dirty = false
	s.secure, err = ParseValue(string(data))
	if err != nil {
		return err
	}
	s.keyId = s.secure.KeyId()
	return nil
}

func (s *SecureData) MarshalCQL(_ gocql.TypeInfo) ([]byte, error) {
	result, err := s.cleanValue()
	if err != nil {
		return nil, err
	}
	return []byte(result), nil
}

func (s *SecureData) Value() (driver.Value, error) {
	return s.cleanValue()
}

func (s *SecureData) Scan(src interface{}) (err error) {
	data, ok := src.(string)
	if !ok {
		return ErrDataInvalid
	}

	s.dirty = false
	s.secure, err = ParseValue(data)
	if err != nil {
		return err
	}
	s.keyId = s.secure.KeyId()
	return nil
}

func (s *SecureData) Field(ctx context.Context, name string) (value *string, err error) {
	if s == nil {
		return nil, nil
	}
	
	// Lazy decrypt on first read
	s.ctx = ctx
	if s.payload == nil {
		s.payload, err = s.secure.Payload()
		if err == ErrValueEncrypted {
			s.payload, err = encrypterFactory.
				Create(s.ctx, s.keyId).
				Decrypt(s.secure.payload)
		}
		if err != nil {
			return nil, err
		}
	}
	value, ok := s.payload[name]
	if !ok {
		return nil, nil
	}
	return value, nil
}

func (s *SecureData) SetField(ctx context.Context, name string, value *string) (*SecureData, error) {
	s.ctx = ctx

	if s.payload == nil {
		// Initialize payload
		switch {
		case s.secure.encrypted:
			payload, err := encrypterFactory.
				Create(ctx, s.secure.keyId).
				Decrypt(s.secure.RawPayload())
			if err != nil {
				return nil, err
			}

			decryptedValue, err := NewValue(s.secure.keyId, payload)
			if err != nil {
				return nil, err
			}

			s.payload = payload
			s.secure = decryptedValue

		case !s.secure.IsEmpty():
			payload, err := s.secure.Payload()
			if err != nil {
				return nil, err
			}

			s.payload = payload

		default:
			s.payload = make(map[string]*string)
		}
	}

	if cur, ok := s.payload[name]; ok {
		// Matching pointers
		if cur == value {
			return s, nil
		}
		// Mismatched nils
		if cur == nil || value == nil {
			s.payload[name] = value
			s.dirty = true
			// Mismatched values
		} else if *cur != *value {
			s.payload[name] = value
			s.dirty = true
		}
	} else {
		s.payload[name] = value
		s.dirty = true
	}
	return s, nil
}

func (s *SecureData) SetKeyId(ctx context.Context, keyId types.UUID) *SecureData {
	s.ctx = ctx
	s.keyId = keyId
	s.dirty = true
	return s
}

func (s *SecureData) KeyId() types.UUID {
	return s.keyId
}

func (s *SecureData) withDecryptedValue(insecureValue Value) (err error) {
	s.payload, err = insecureValue.Payload()
	if err != nil {
		return err
	}
	s.dirty = false
	s.secure = insecureValue
	return nil
}

type WithSecureData struct {
	SecureData *SecureData `db:"secure_data"`
}

func (g *WithSecureData) SecureValue(ctx context.Context, fieldName string) (string, error) {
	if g == nil {
		return "", nil
	}
	optionalValue, err := g.SecureOptionalValue(ctx, fieldName)
	if err != nil {
		return "", err
	}
	return optionalValue.OrElse(""), nil
}

func (g *WithSecureData) SecureOptionalValue(ctx context.Context, fieldName string) (types.OptionalString, error) {
	if g == nil {
		return types.OptionalString{}, nil
	}

	value, err := g.SecureData.Field(ctx, fieldName)
	if err != nil {
		return types.NewOptionalString(nil), err
	}
	return types.NewOptionalString(value), nil
}

func (g *WithSecureData) SetSecureValue(ctx context.Context, keyId types.UUID, fieldName string, value *string) error {
	if keyId == nil || keyId.IsEmpty() {
		return ErrKeyNotSet
	}
	if err := keyId.Validate(); err != nil {
		return err
	}

	if g.SecureData == nil {
		g.SecureData = new(SecureData)
	}
	_, err := g.SecureData.
		SetKeyId(ctx, keyId).
		SetField(ctx, fieldName, value)

	return err
}
