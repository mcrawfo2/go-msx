package repository

import "github.com/pkg/errors"

var ErrNotFound = errors.New("Entity not found")
var ErrAlreadyExists = errors.New("Entity already exists")
