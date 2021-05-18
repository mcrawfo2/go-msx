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
