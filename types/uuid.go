package types

import (
	"bytes"
	"github.com/hashicorp/go-uuid"
)

type UUID []byte

func ParseUUID(value string) (UUID, error) {
	return uuid.ParseUUID(value)
}

func (u UUID) Equals(other UUID) bool {
	return bytes.Compare(u, other) == 0
}
