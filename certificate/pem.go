package certificate

import (
	"crypto/x509"
	"encoding/pem"
	"github.com/pkg/errors"
)

func ParsePemCertificate(data []byte) (*x509.Certificate, error) {
	pemBlock, _ := pem.Decode(data)
	if pemBlock == nil || pemBlock.Type != "CERTIFICATE" {
		return nil, errors.New("PEM file does not contain valid certificate")
	}

	return x509.ParseCertificate(pemBlock.Bytes)
}
