// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package security

import "context"

type CertificateProvider interface {
	UserContextFromCertificate(ctx context.Context, certificate string) (*UserContext, error)
}

var certificateProvider CertificateProvider

func SetCertificateProvider(provider CertificateProvider) {
	if provider != nil {
		certificateProvider = provider
	}
}

func NewUserContextFromCertificate(ctx context.Context, certificate string) (userContext *UserContext, err error) {
	return certificateProvider.UserContextFromCertificate(ctx, certificate)
}
