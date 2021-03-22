package sqldb

import (
	"errors"
)

var ErrDataInvalid = errors.New("invalid data, expecting: []byte")