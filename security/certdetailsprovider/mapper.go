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
