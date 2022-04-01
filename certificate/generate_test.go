// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package certificate

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/thejerf/abtime"
	"math/big"
	"testing"
	"time"
)

func generateCertificate(t *testing.T, clock abtime.AbstractTime, duration time.Duration) *tls.Certificate {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Error(err.Error())
		t.SkipNow()
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Dummy Co"},
		},
		NotBefore:             clock.Now(),
		NotAfter:              clock.Now().Add(duration),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		t.Error(err.Error())
		t.SkipNow()
	}

	return &tls.Certificate{
		Certificate: [][]byte{certBytes},
	}
}
