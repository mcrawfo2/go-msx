// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package certprovider

import (
	"context"
	"crypto/x509"
	"cto-github.cisco.com/NFV-BU/go-msx/cache/lru"
	"cto-github.cisco.com/NFV-BU/go-msx/certificate"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"github.com/benbjohnson/clock"
	"time"
)

const authorityProvider = "msx-ca"
const authorityCacheKey = "msx-ca"

var logger = log.NewLogger("msx.security.certprovider")

type Provider struct {
	provider certificate.Provider
	cache    lru.Cache
}

func (p *Provider) UserContextFromCertificate(ctx context.Context, pem string) (*security.UserContext, error) {
	cert, err := certificate.ParsePemCertificate([]byte(pem))
	if err != nil {
		return nil, err
	}
	err = p.validateCertificate(ctx, cert)
	if err != nil {
		return nil, err
	}

	// TODO: Revocation check

	userContext := security.UserContextFromContext(ctx)
	userContext = userContext.Clone()
	userContext.UserName = cert.Subject.CommonName
	userContext.Subject = cert.Subject.String()
	userContext.Issuer = cert.Issuer.String()
	userContext.Certificate = cert
	userContext.IssuedAt = int(cert.NotBefore.Unix())
	userContext.Exp = int(cert.NotAfter.Unix())
	userContext.Authorities = []string{
		"ROLE_CLIENT",
	}
	userContext.Scopes = []string{
		"read", "write",
	}

	return userContext, nil
}

func (p *Provider) validateCertificate(ctx context.Context, cert *x509.Certificate) error {
	authorityCert, err := p.getAuthorityCertificate(ctx)
	if err != nil {
		return err
	}

	authorityPool := x509.NewCertPool()
	authorityPool.AddCert(authorityCert)

	verifyOptions := x509.VerifyOptions{
		Roots:       authorityPool,
		CurrentTime: time.Now(),
	}

	if _, err = cert.Verify(verifyOptions); err != nil {
		return err
	}

	// TODO: Revocation checking

	return nil
}

func (p *Provider) getAuthorityCertificate(ctx context.Context) (*x509.Certificate, error) {
	certObject, ok := p.cache.Get(authorityCacheKey)
	if ok {
		return certObject.(*x509.Certificate), nil
	}

	cert, err := p.provider.GetCaCertificate(ctx)
	if err != nil {
		return nil, err
	}

	p.cache.Set(authorityCacheKey, cert)

	return cert, nil
}

func newCertificateProvider(ctx context.Context) (*Provider, error) {
	certProvider, err := certificate.NewProvider(ctx, authorityProvider)
	if err != nil {
		return nil, err
	}

	return &Provider{
		provider: certProvider,
		cache:    lru.NewCache(1*time.Minute, 10, 1*time.Minute, clock.New()),
	}, nil
}

func RegisterCertificateProvider(ctx context.Context) error {
	logger.Info("Registering certificate provider")

	certProvider, err := newCertificateProvider(ctx)
	if err != nil {
		return err
	}

	security.SetCertificateProvider(certProvider)

	return nil
}
