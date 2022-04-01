// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package certificate

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func createCertificate(t *testing.T) *x509.Certificate {
	_, cert := createCertificatePair(t)
	return cert
}

func createKey(t *testing.T) *ecdsa.PrivateKey {
	key, _ := createCertificatePair(t)
	return key
}

func createCertificatePair(t *testing.T) (*ecdsa.PrivateKey, *x509.Certificate) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.SkipNow()
	}

	appName := uuid.New().String()

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:       []string{"Cisco Systems"},
			OrganizationalUnit: []string{"MSX"},
			CommonName:         appName,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 180),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames: []string{
			appName,
		},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		t.SkipNow()
	}

	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		t.SkipNow()
	}

	return priv, cert
}

func TestRenderPemCertificate(t *testing.T) {
	type args struct {
		cert *x509.Certificate
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				cert: createCertificate(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RenderPemCertificate(tt.args.cert)
			if (err != nil) != (tt.wantErr) {
				t.Error(testhelpers.Diff(nil, err))
			}
			cert, err := ParsePemCertificate(got)
			assert.NoError(t, err)
			assert.Equal(t, tt.args.cert.Raw, cert.Raw)
		})
	}
}

func TestRenderPemPrivateKey(t *testing.T) {
	type args struct {
		key *ecdsa.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				key: createKey(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RenderPemPrivateKey(tt.args.key)
			if (err != nil) != (tt.wantErr) {
				t.Error(testhelpers.Diff(nil, err))
			}

			key, err := ParsePemPrivateKey(got)
			assert.NoError(t, err)

			ecdsaKey, ok := key.(*ecdsa.PrivateKey)
			assert.True(t, ok)
			assert.True(t, tt.args.key.Equal(ecdsaKey))
		})
	}
}
