// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package elasticsearch

import (
	"encoding/json"
	"time"
)

const timeFormat = "2006-01-02T15:04:05.000Z"

type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).UTC().Format(timeFormat))
}
