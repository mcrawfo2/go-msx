// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package certdetailsprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"github.com/pkg/errors"
)

var ErrCertificateMapperNotRegistered = errors.New("Certificate mapper not registered.")

type TokenDetailsProvider struct {
}

func (t TokenDetailsProvider) TokenDetails(ctx context.Context) (*security.UserContextDetails, error) {
	if mapper == nil {
		return nil, ErrCertificateMapperNotRegistered
	}
	userContext := security.UserContextFromContext(ctx)
	return mapper.CertificateDetails(ctx, userContext)
}

func (t TokenDetailsProvider) IsTokenActive(ctx context.Context) (bool, error) {
	tokenDetails, err := t.TokenDetails(ctx)
	if err != nil {
		return false, err
	}
	return tokenDetails.Active, nil
}
