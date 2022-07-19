// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package types

import (
	"reflect"
)

type Merge struct{}

// This is only necessary until go is upgrade to version 18, where "copy" functionality and generics become available
func (m Merge) RecursiveMerge(src map[string]interface{}, dest map[string]interface{}) map[string]interface{} {
	for key, value := range src {
		if dest[key] == nil {
			dest[key] = value
		} else if dest[key] != nil && reflect.TypeOf(value).Kind() == reflect.Map {
			dest[key] = m.RecursiveMerge(src[key].(map[string]interface{}), dest[key].(map[string]interface{}))
		} else if dest[key] != nil && reflect.TypeOf(value).Kind() == reflect.Slice && reflect.TypeOf(dest[key]).Kind() == reflect.Slice {
			// append the two arrays
			switch dest[key].(type) {
			case []interface{}:
				dest[key] = append(dest[key].([]interface{}), value.([]interface{})...)
			case []string:
				dest[key] = append(dest[key].([]string), value.([]string)...)
			case []int:
				dest[key] = append(dest[key].([]int), value.([]int)...)
			case []bool:
				dest[key] = append(dest[key].([]bool), value.([]bool)...)
			case []float32:
				dest[key] = append(dest[key].([]float32), value.([]float32)...)
			case []float64:
				dest[key] = append(dest[key].([]float64), value.([]float64)...)
			}

		} else {
			dest[key] = value
		}
	}
	return dest
}
