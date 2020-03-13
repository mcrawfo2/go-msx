package security

import "github.com/pkg/errors"

var (
	ErrTokenNotFound = errors.New("Token missing from context")
)
