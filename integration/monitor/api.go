// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

//go:generate mockery --inpackage --name=Api --structname=MockMonitor
package monitor

import "cto-github.cisco.com/NFV-BU/go-msx/integration"

type Api interface {
	GetDeviceHealth(deviceIds string) (*integration.MsxResponse, error)
}

// Ensure MockMonitor implementation is up-to-date
var _ Api = new(MockMonitor)
