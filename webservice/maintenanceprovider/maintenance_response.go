// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package maintenanceprovider

type MaintenanceTask struct {
	Task    string `json:"task"`
	Mode    string `json:"mode"`
	Message string `json:"message"`
}

type MaintenanceResponse struct {
	Mode   string            `json:"mode"`
	Detail []MaintenanceTask `json:"detail"`
}
