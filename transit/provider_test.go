package transit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterProvider(t *testing.T) {
	t.Run("NoRegister", func(t *testing.T) {
		encryptionProvider = nil
		err := RegisterProvider(nil)
		assert.Error(t, err)
		assert.Equal(t, ErrNotRegistered, err)
		assert.Nil(t, encryptionProvider)
	})

	t.Run("Register", func(t *testing.T) {
		mockProvider := new(MockProvider)
		err := RegisterProvider(mockProvider)
		assert.NoError(t, err)
		assert.Equal(t, mockProvider, encryptionProvider)
	})
}

func Test_provider(t *testing.T) {
	t.Run("NotRegistered", func(t *testing.T) {
		encryptionProvider = nil
		actualEncryptionProvider, err := provider()
		assert.Error(t, err)
		assert.Equal(t, ErrNotRegistered, err)
		assert.Nil(t, actualEncryptionProvider)
	})

	t.Run("Registered", func(t *testing.T) {
		encryptionProvider = new(MockProvider)
		actualEncryptionProvider, err := provider()
		assert.NoError(t, err)
		assert.Equal(t, encryptionProvider, actualEncryptionProvider)
	})
}
