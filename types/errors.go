package types

import "bytes"

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

func FilterList(source ErrorList) error {
	var result ErrorList
	for _, v := range source {
		if v != (error)(nil) {
			result = append(result, v)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}
