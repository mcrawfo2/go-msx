package security

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testUserName = "test"

func TestSetTokenDetailsProvider(t *testing.T) {
	ctx := context.Background()

	userContextDetails := new(UserContextDetails)
	userContextDetails.Username = types.NewStringPtr(testUserName)

	tokenDetailsProvider := new(MockTokenDetailsProvider)
	tokenDetailsProvider.
		On("TokenDetails", ctx).
		Return(userContextDetails, nil)

	SetTokenDetailsProvider(tokenDetailsProvider)

	actualUserContextDetails, err := NewUserContextDetails(ctx)
	assert.NoError(t, err)
	assert.Equal(t, userContextDetails, actualUserContextDetails)
}

func TestNewUserContextDetails(t *testing.T) {
	ctx := context.Background()

	userContextDetails := new(UserContextDetails)
	userContextDetails.Username = types.NewStringPtr(testUserName)

	mockTokenDetailsProvider := new(MockTokenDetailsProvider)
	mockTokenDetailsProvider.
		On("TokenDetails", ctx).
		Return(userContextDetails, nil)

	tokenDetailsProvider = mockTokenDetailsProvider

	actualUserContextDetails, err := NewUserContextDetails(ctx)
	assert.NoError(t, err)
	assert.Equal(t, userContextDetails, actualUserContextDetails)
}

func TestIsTokenActive(t *testing.T) {
	ctx := context.Background()

	userContextDetails := new(UserContextDetails)
	userContextDetails.Username = types.NewStringPtr(testUserName)

	mockTokenDetailsProvider := new(MockTokenDetailsProvider)
	mockTokenDetailsProvider.
		On("IsTokenActive", ctx).
		Return(true, nil)

	tokenDetailsProvider = mockTokenDetailsProvider

	isActive, err := IsTokenActive(ctx)
	assert.NoError(t, err)
	assert.True(t, isActive)
}
