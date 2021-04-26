package prepared

import (
	"database/sql/driver"
	"github.com/lib/pq"
)

// StringArray is a sequence of strings.
// Expected column type: ARRAY
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	return pq.StringArray(a).Value()
}

func (a *StringArray) Scan(value interface{}) error {
	v := &pq.StringArray{}
	err := v.Scan(value)
	if err != nil {
		return err
	}
	*a = []string(*v)
	return nil
}
