package sqldb

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
)

var ErrDataInvalid = errors.New("invalid data")
var errNotByteArray = errors.Wrap(ErrDataInvalid, "expecting []byte")

// Deprecated
type MapStrStr map[string]string

func (a MapStrStr) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *MapStrStr) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errNotByteArray
	}

	return json.Unmarshal(b, &a)
}

// Deprecated
type MapStrIface map[string]interface{}

func (a MapStrIface) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *MapStrIface) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

// Deprecated
type Bytes []byte

func (b Bytes) Value() (driver.Value, error) {
	data := make([]byte, len(b)*2+3)
	data[0] = 'x'
	data[1] = '\''
	hex.Encode(data[2:], b)
	data[len(b)*2+2] = '\''
	return data, nil
}

func (b *Bytes) Scan(value interface{}) error {
	data, ok := value.([]byte)
	if !ok {
		return errNotByteArray
	}

	if data[0] != 'x' || data[1] != '\'' || data[len(data)-1] != '\'' {
		return errors.New("incorrect format")
	}

	*b = make([]byte, (len(data)-3)/2)
	_, err := hex.Decode(*b, data[2:len(data)-1])
	return err
}
