package skel

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/pkg/errors"
	"math/big"
	"os"
	"path"
	"time"
)

func init() {
	AddTarget("generate-certificate", "Generate an X.509 server certificate and private key", GenerateCertificate)
}

// https://github.com/Shyp/generate-tls-cert/blob/master/generate.go

func GenerateCertificate(args []string) error {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return errors.Wrap(err, "Failed to generate ecdsa key")
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:       []string{"Cisco Systems"},
			OrganizationalUnit: []string{"MSX"},
			CommonName:         skeletonConfig.AppName,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 180),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames: []string{
			skeletonConfig.AppName,
		},
	}

	/*
	   hosts := strings.Split(*host, ",")
	   for _, h := range hosts {
	   	if ip := net.ParseIP(h); ip != nil {
	   		template.IPAddresses = append(template.IPAddresses, ip)
	   	} else {
	   		template.DNSNames = append(template.DNSNames, h)
	   	}
	   }
	   if *isCA {
	   	template.IsCA = true
	   	template.KeyUsage |= x509.KeyUsageCertSign
	   }
	*/

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return errors.Wrap(err, "Failed to create certificate")
	}

	out := &bytes.Buffer{}
	_ = pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	targetFileName := path.Join(skeletonConfig.TargetDirectory(), "local", "server.crt")
	err = writeFileBytes(targetFileName, out.Bytes())
	if err != nil {
		return errors.Wrap(err, "Failed to write certificate file")
	}

	out.Reset()

	pemBlock, err := pemBlockForKey(priv)
	if err != nil {
		return errors.Wrap(err, "Failed to generate private key PEM block")
	}

	err = pem.Encode(out, pemBlock)
	if err != nil {
		return errors.Wrap(err, "Failed to encode private key")
	}
	targetFileName = path.Join(skeletonConfig.TargetDirectory(), "local", "server.key")
	err = writeFileBytes(targetFileName, out.Bytes())

	return nil
}

func writeFileBytes(targetFileName string, data []byte) error {
	logger.Infof("Writing %s", targetFileName)

	err := os.MkdirAll(path.Dir(targetFileName), 0755)
	if err != nil {
		return errors.Wrap(err, "Failed to create directory")
	}

	writer, err := os.Create(targetFileName)
	if err != nil {
		return errors.Wrap(err, "Failed to create file")
	}
	defer writer.Close()

	_, err = writer.Write(data)
	if err != nil {
		return errors.Wrap(err, "Failed to write file")
	}

	return nil
}

func pemBlockForKey(priv *ecdsa.PrivateKey) (*pem.Block, error) {
	b, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to marshal ECDSA private key: %v", err)
	}
	return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}, nil
}
