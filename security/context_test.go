package security

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContextWithUserContext(t *testing.T) {
	userContext := new(UserContext)
	ctx := context.Background()

	newCtx := ContextWithUserContext(ctx, userContext)
	injectedUserContext := newCtx.Value(contextKeyUserContext)
	assert.Equal(t, userContext, injectedUserContext)
}

func TestUserContextFromContext(t *testing.T) {
	userContext := new(UserContext)
	ctx := context.Background()
	newCtx := context.WithValue(ctx, contextKeyUserContext, userContext)

	injectedUserContext := UserContextFromContext(newCtx)
	assert.Equal(t, userContext, injectedUserContext)
}

func TestUserContextFromContext_DefaultUserContext(t *testing.T) {
	ctx := context.Background()
	userContext := UserContextFromContext(ctx)
	assert.NotNil(t, userContext)
	assert.Equal(t, contextDefaultUserName, userContext.UserName)
}

func TestUserNameFromContext(t *testing.T) {
	userName := "cisco123"
	userContext := &UserContext{
		UserName: userName,
	}
	ctx := context.Background()

	newCtx := ContextWithUserContext(ctx, userContext)
	retrievedUserName := UserNameFromContext(newCtx)
	assert.Equal(t, userName, retrievedUserName)
}
