package security

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetTokenProvider(t *testing.T) {
	mockTokenProvider := new(MockTokenProvider)
	SetTokenProvider(mockTokenProvider)
	assert.Equal(t, mockTokenProvider, tokenProvider)
}

func TestNewUserContextFromToken(t *testing.T) {
	ctx := context.Background()

	mockUserContext := new(UserContext)

	mockTokenProvider := new(MockTokenProvider)
	mockTokenProvider.
		On("UserContextFromToken", ctx, "").
		Return(mockUserContext, nil)

	tokenProvider = mockTokenProvider

	actualUserContext, err := NewUserContextFromToken(ctx, "")
	assert.NoError(t, err)
	assert.Equal(t, mockUserContext, actualUserContext)
}
