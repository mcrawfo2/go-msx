package types

import (
	"bytes"
	"encoding/json"
	"github.com/hashicorp/go-uuid"
)

type UUID []byte

func (u UUID) MarshalJSON() ([]byte, error) {
	str, err := uuid.FormatUUID(u)
	return []byte(str), err
}

func (u UUID) UnmarshalJSON(data []byte) error {
	var uuidString string
	if err := json.Unmarshal(data, &uuidString); err != nil {
		return err
	}
	if uuidBytes, err := uuid.ParseUUID(uuidString); err != nil {
		return err
	} else {
		copy(u, uuidBytes)
	}
	return nil
}

func ParseUUID(value string) (UUID, error) {
	return uuid.ParseUUID(value)
}

func (u UUID) Equals(other UUID) bool {
	return bytes.Compare(u, other) == 0
}
