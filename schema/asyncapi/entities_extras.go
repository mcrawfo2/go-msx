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
