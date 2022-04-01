// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package prepared

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/pkg/errors"
)

// JsonMapStringInterface represents a JSON object in the database mapping string keys to arbitrary values.
// Expected column type: JSON/JSONB
type JsonMapStringInterface map[string]interface{}

func (a JsonMapStringInterface) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *JsonMapStringInterface) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.Errorf("Cannot convert %T to JsonMapStringInterface", value)
	}

	return json.Unmarshal(b, &a)
}
