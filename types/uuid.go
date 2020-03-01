package types

import (
	"bytes"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gocql/gocql"
	"github.com/hashicorp/go-uuid"
)

type UUID []byte

func (u UUID) MarshalJSON() ([]byte, error) {
	if u == nil {
		return json.Marshal(nil)
	}
	str, err := uuid.FormatUUID(u)
	if err != nil {
		return nil, err
	}
	return json.Marshal(str)
}

func (u *UUID) UnmarshalJSON(data []byte) error {
	var uuidString string
	if err := json.Unmarshal(data, &uuidString); err != nil {
		return err
	}
	if uuidBytes, err := uuid.ParseUUID(uuidString); err != nil {
		return err
	} else {
		*u = uuidBytes[:]
	}
	return nil
}

// DEPRECATED
func (u UUID) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return u[:], nil
}

// DEPRECATED
func (u *UUID) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	*u = data
	return nil
}

func (u UUID) MarshalText() (string, error) {
	return uuid.FormatUUID(u)
}

func (u UUID) MustMarshalText() string {
	text, err := u.MarshalText()
	if err != nil {
		panic(err)
	}
	return text
}

func (u *UUID) UnmarshalText(data string) error {
	if uuidBytes, err := uuid.ParseUUID(data); err != nil {
		return err
	} else {
		*u = uuidBytes[:]
	}
	return nil
}

func (u UUID) Equals(other UUID) bool {
	return bytes.Compare(u, other) == 0
}

func (u UUID) IsEmpty() bool {
	for _, v := range u {
		if v != 0 {
			return false
		}
	}
	return true
}

func (u UUID) String() string {
	return u.MustMarshalText()
}

func ParseUUID(value string) (UUID, error) {
	return uuid.ParseUUID(value)
}

func MustParseUUID(value string) UUID {
	result, err := uuid.ParseUUID(value)
	if err != nil {
		panic(err)
	}
	return result
}

func EmptyUUID() UUID {
	return UUID([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}

func NewUUID() (UUID, error) {
	return uuid.GenerateRandomBytes(16)
}

func (u UUID) Validate() error {
	return validation.Validate([]byte(u), validation.Length(16, 16))
}

func (u UUID) ToByteArray() [16]byte {
	result := [16]byte{}
	copy(result[:], u[0:16])
	return result
}
