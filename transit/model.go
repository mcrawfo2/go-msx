package transit

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
)

var ErrKeyNotSet = errors.New("Key Id not set for encryption")

type SecureData struct {
	ctx     context.Context
	dirty   bool
	keyId   types.UUID
	secure  Value
	payload map[string]*string
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

func (s *SecureData) MarshalCQL(_ gocql.TypeInfo) (data []byte, err error) {
	if s == nil || s.keyId == nil {
		return nil, ErrKeyNotSet
	}

	// Lazy encrypt on dirty write
	if s.dirty {
		payload, encrypted, err := NewEncrypter(s.ctx, s.keyId).Encrypt(s.payload)
		if err != nil {
			return nil, err
		}
		if encrypted {
			s.secure = NewSecureValue(s.keyId, payload)
		} else {
			s.secure, err = NewValue(s.keyId, s.payload)
			if err != nil {
				return nil, err
			}
		}
		s.dirty = false
	}

	return []byte(s.secure.String()), nil
}

func (s *SecureData) Field(ctx context.Context, name string) (value *string, err error) {
	// Lazy decrypt on first read
	s.ctx = ctx
	if s.payload == nil {
		s.payload, err = s.secure.Payload()
		if err == ErrValueEncrypted {
			s.payload, err = NewEncrypter(s.ctx, s.keyId).Decrypt(s.secure.payload)
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

func (s *SecureData) SetField(ctx context.Context, name string, value *string) *SecureData {
	s.ctx = ctx
	if s.payload == nil {
		s.payload = make(map[string]*string)
	}
	if cur, ok := s.payload[name]; ok {
		if cur == value {
			return s
		}
		if cur == nil || value == nil {
			s.payload[name] = value
			s.dirty = true
		} else if *cur != *value {
			s.payload[name] = value
			s.dirty = true
		}
	} else {
		s.payload[name] = value
		s.dirty = true
	}
	return s
}

func (s *SecureData) SetKeyId(ctx context.Context, keyId types.UUID) *SecureData {
	s.ctx = ctx
	s.keyId = keyId
	s.dirty = true
	return s
}