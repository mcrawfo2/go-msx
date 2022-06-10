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
			dest[key] = append(dest[key].([]interface{}), src[key].([]interface{}))
		} else {
			dest[key] = value
		}
	}
	return dest
}
