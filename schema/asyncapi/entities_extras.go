// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package asyncapi

import (
	"encoding/json"
	"fmt"
)

func (r *Reference) UnmarshalJSON(data []byte) error {
	var rawMap map[string]json.RawMessage
	err := json.Unmarshal(data, &rawMap)
	if err != nil {
		return err
	}

	if ref, ok := rawMap["$ref"]; ok {
		return json.Unmarshal(ref, &r.Ref)
	} else {
		return fmt.Errorf("reference not found")
	}
}

func NewTag(name string) *Tag {
	return new(Tag).WithName(name)
}
