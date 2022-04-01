// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

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
