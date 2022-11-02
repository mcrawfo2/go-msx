// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package types

import (
	"bytes"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"sort"
)

type Filterable interface {
	Filter() error
}

type CompositeError interface {
	Errors() interface{}
}

type ErrorList []error

func (l ErrorList) Errors() interface{} {
	var result []interface{}
	for _, err := range l {
		switch typedErr := err.(type) {
		case CompositeError:
			result = append(result, typedErr.Errors())
		case error:
			result = append(result, typedErr.Error())
		}
	}
	return result
}

func (l ErrorList) Error() string {
	var buffer bytes.Buffer
	for i, err := range l {
		if i > 0 {
			buffer.WriteString("; ")
		}
		buffer.WriteString(err.Error())
	}
	return buffer.String()
}

func (l ErrorList) Filter() error {
	return FilterList(l)
}

func (l ErrorList) Strings() []string {
	var errs []string
	for _, err := range l {
		errs = append(errs, err.Error())
	}
	return errs
}

func FilterList(source ErrorList) error {
	var result ErrorList
	for _, v := range source {
		if filteredItem := FilterItem(v); filteredItem != (error)(nil) {
			result = append(result, filteredItem)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

type ErrorMap map[string]error

func (m ErrorMap) Errors() interface{} {
	var result map[string]interface{}
	for k, err := range m {
		switch typedErr := err.(type) {
		case CompositeError:
			result[k] = typedErr.Errors()
		case error:
			result[k] = typedErr.Error()
		}
	}
	return result
}

// Error returns the error string of Errors.
func (m ErrorMap) Error() string {
	if len(m) == 0 {
		return ""
	}

	keys := []string{}
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	s := ""
	for i, key := range keys {
		if i > 0 {
			s += "; "
		}
		s += fmt.Sprintf("%v: %v", key, m[key].Error())
	}
	return s + "."
}

func (m ErrorMap) Filter() error {
	return FilterMap(m)
}

func FilterMap(source map[string]error) error {
	var result = make(ErrorMap)
	for k, v := range source {
		if filteredItem := FilterItem(v); filteredItem != (error)(nil) {
			result[k] = filteredItem
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func FilterItem(v error) error {

	if v == (error)(nil) {
		return nil
	}

	switch val := v.(type) {
	case ErrorMap:
		return FilterMap(val)
	case validation.Errors:
		return FilterMap(val)
	case ErrorList:
		return FilterList(val)
	default:
		return v
	}
}

// May is a helper function that throws away its error argument
// intended for calls where we don't care about the error but want to inline the fn call
// note: this form is possible because of the special case of passing back matching
// return values from a function call
func May[vtype any](v vtype, err error) vtype {
	_ = err
	return v
}
