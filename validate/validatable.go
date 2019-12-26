package validate

import "cto-github.cisco.com/NFV-BU/go-msx/types"

type Validatable interface {
	Validate() error
}

func Validate(validatable Validatable) error {
	err := validatable.Validate()
	if err != nil {
		if filterable, ok := err.(types.Filterable); ok {
			err = filterable.Filter()
		}
	}
	return err
}
