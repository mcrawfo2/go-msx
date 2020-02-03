package types

import (
	"bytes"
	"fmt"
	"sort"
)

type Filterable interface {
	Filter() error
}

type ErrorList []error

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
	} else if errorMap, ok := v.(ErrorMap); ok {
		return FilterMap(errorMap)
		// TODO: ozzo-validation.Errors
		//} else if validationErrors, ok := v.(validation.Errors); ok {
		//	return FilterMap(validationErrors)
	} else if errorList, ok := v.(ErrorList); ok {
		return FilterList(errorList)
	} else {
		return v
	}
}
