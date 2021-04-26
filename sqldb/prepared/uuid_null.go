package prepared

import (
	"database/sql/driver"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// NullUUID represents a nullable-UUID in the database.
// Expected column type: NULLABLE UUID
type NullUUID struct {
	valid bool
	value uuid.UUID
}

func (a NullUUID) Empty() bool {
	return !a.valid
}

func (a NullUUID) Valid() bool {
	return a.valid
}

func (a NullUUID) String() string {
	if !a.valid {
		return "[Missing UUID]"
	}
	return a.value.String()
}

func (a NullUUID) UUID() uuid.UUID {
	return a.value
}

func (a NullUUID) Value() (driver.Value, error) {
	if !a.valid {
		return nil, nil
	}
	return a.value.String(), nil
}

func (a *NullUUID) Scan(value interface{}) error {
	var err error
	switch vt := value.(type) {
	case []byte:
		a.value, err = uuid.ParseBytes(vt)
		a.valid = err == nil
	case string:
		a.value, err = uuid.Parse(vt)
		a.valid = err == nil
	case nil:
		a.valid = false
	default:
		return errors.Errorf("Cannot convert %T to NullUUID", value)
	}

	return nil
}

func NewNullUUID(value uuid.UUID) NullUUID {
	return NullUUID{
		valid: true,
		value: value,
	}
}
