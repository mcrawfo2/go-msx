package certificate

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/pkg/errors"
)

const (
	pemKeyHeaderRsa = "RSA PRIVATE KEY"
	pemKeyHeaderEc  = "EC PRIVATE KEY"
)

func ParsePemCertificate(data []byte) (*x509.Certificate, error) {
	pemBlock, _ := pem.Decode(data)
	if pemBlock == nil || pemBlock.Type != "CERTIFICATE" {
		return nil, errors.New("PEM file does not contain valid certificate")
	}

	return x509.ParseCertificate(pemBlock.Bytes)
}

func ParsePemPrivateKey(data []byte) (interface{}, error) {
	pemBlock, data := pem.Decode(data)
	for pemBlock != nil {
		switch pemBlock.Type {
		case pemKeyHeaderRsa:
			return x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
		case pemKeyHeaderEc:
			return x509.ParseECPrivateKey(pemBlock.Bytes)
		}
	}

	return nil, errors.New("No EC/RSA private key found")
}

func RenderPemCertificate(cert *x509.Certificate) ([]byte, error) {
	out := &bytes.Buffer{}
	err := pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func RenderPemPrivateKey(key crypto.PrivateKey) ([]byte, error) {
	if ecdsaPrivateKey, ok := key.(*ecdsa.PrivateKey); ok {
		privateKeyBytes, err := x509.MarshalECPrivateKey(ecdsaPrivateKey)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to marshal ecdsa private key")
		}

		pemBytes := pem.EncodeToMemory(&pem.Block{
			Type:  pemKeyHeaderEc,
			Bytes: privateKeyBytes,
		})

		return pemBytes, nil

	} else if rsaPrivateKey, ok := key.(*rsa.PrivateKey); ok {
		privateKeyBytes := x509.MarshalPKCS1PrivateKey(rsaPrivateKey)

		pemBytes := pem.EncodeToMemory(&pem.Block{
			Type:  pemKeyHeaderRsa,
			Bytes: privateKeyBytes,
		})

		return pemBytes, nil
	} else {
		return nil, errors.Errorf("Could not PEM encode private key from %T", key)
	}
}

type Issued struct {
	Authority  *x509.Certificate
	Identity   *x509.Certificate
	PrivateKey interface{}
}

func (i Issued) AuthorityPem() ([]byte, error) {
	return RenderPemCertificate(i.Authority)
}

func (i Issued) IdentityPem() ([]byte, error) {
	return RenderPemCertificate(i.Identity)
}

func (i Issued) PrivateKeyPem() ([]byte, error) {
	return RenderPemPrivateKey(i.PrivateKey)
}
