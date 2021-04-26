package prepared

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/pkg/errors"
)

// JsonMapStringString represents a JSON object in the database mapping string keys to string values.
// Expected column type: JSON/JSONB
type JsonMapStringString map[string]string

func (a JsonMapStringString) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *JsonMapStringString) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.Errorf("Cannot convert %T to JsonMapStringString", value)
	}

	return json.Unmarshal(b, &a)
}
