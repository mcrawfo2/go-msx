// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package certdetailsprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
)

type CertificateMapper interface {
	CertificateDetails(ctx context.Context, userContext *security.UserContext) (*security.UserContextDetails, error)
}

var mapper CertificateMapper

func RegisterCertificateMapper(c CertificateMapper) {
	if c != nil {
		mapper = c
	}
}
