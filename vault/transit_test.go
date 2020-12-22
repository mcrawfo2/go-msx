package vault

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCreateTransitKeyRequest(t *testing.T) {
	request := NewCreateTransitKeyRequest()
	assert.Equal(t, KeyTypeAes256Gcm96, request.Type)
	assert.Equal(t, false, *request.AllowPlaintextBackup)
	assert.Equal(t, false, *request.Exportable)
}
