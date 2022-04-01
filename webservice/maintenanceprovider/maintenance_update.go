// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package maintenanceprovider

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	validation "github.com/go-ozzo/ozzo-validation"
)

type MaintenanceUpdate struct {
	Mode string `json:"mode" enum:"NORMAL|MAINTENANCE"`
}

func (v *MaintenanceUpdate) Validate() error {
	return types.ErrorMap{
		"mode": validation.Validate(&v.Mode, validation.Required,
			validation.In(&v.Mode, "NORMAL", "MAINTENANCE")),
	}
}
