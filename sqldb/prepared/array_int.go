// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

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
