package prepared

import (
	"database/sql/driver"

	"github.com/lib/pq"
)

// Expected column type: ARRAY
type IntArray []int64

func (a IntArray) Value() (driver.Value, error) {
	return pq.Int64Array(a).Value()
}

func (a *IntArray) Scan(value interface{}) error {
	v := &pq.Int64Array{}
	err := v.Scan(value)
	if err != nil {
		return err
	}
	*a = []int64(*v)
	return nil
}
