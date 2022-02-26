package certificate

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"cto-github.cisco.com/NFV-BU/go-msx/background"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"github.com/thejerf/abtime"
	"math/rand"
	"sync"
	"time"
)

var logger = log.NewLogger("msx.certificate")
var random = rand.New(rand.NewSource(time.Now().UnixNano()))
var sources = make(map[string]*Source)
var sourceLock sync.Mutex

const renewTimerId = iota

type Source struct {
	sync.Mutex
	certificate *tls.Certificate
	provider    Provider
	clock       abtime.AbstractTime
}

func (c *Source) Provider() Provider {
	return c.provider
}

// TlsCertificate fetches the server certificate for tls.Config
func (c *Source) TlsCertificate(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return c.Certificate(), nil
}

// TlsClientCertificate fetches the client certificate for tls.Config
func (c *Source) TlsClientCertificate(info *tls.CertificateRequestInfo) (*tls.Certificate, error) {
	return c.Certificate(), nil
}

// Certificate fetches the most recently cached certificate
func (c *Source) Certificate() *tls.Certificate {
	c.Lock()
	defer c.Unlock()
	return c.certificate
}

func (c *Source) setCertificate(certificate *tls.Certificate) {
	c.Lock()
	defer c.Unlock()
	c.certificate = certificate
}

// renew continuously refreshes the certificate after approximately half of its remaining validity period.
func (c *Source) renew(ctx context.Context) {
	defer logger.Info("Exiting certificate renewal")

	for {
		d, err := c.period()
		if err != nil {
			err = errors.Wrap(err, "Failed to calculate certificate renewal period")
			background.ErrorReporterFromContext(ctx).Fatal(err)
			return
		}

		logger.Infof("Renewing Certificate in %f minutes", d.Minutes())

		t := c.clock.NewTimer(d, renewTimerId)

		select {
		case <-ctx.Done():
			t.Stop()
			return

		case <-t.Channel():
			err = c.renewOnce(ctx)
			if err != nil {
				err = errors.Wrap(err, "Failed to renew certificate")
				background.ErrorReporterFromContext(ctx).Fatal(err)
				return
			}
		}
	}
}

func (c *Source) period() (time.Duration, error) {
	cert := c.Certificate()
	parsedCert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return 0, errors.Wrap(err, "Problem parsing certificate")
	}

	validFrom := parsedCert.NotBefore
	validTo := parsedCert.NotAfter
	now := c.clock.Now()

	totalValidity := validTo.Sub(validFrom)
	halfPeriod := 0.5 * float64(totalValidity)
	validHalf := validFrom.Add(time.Duration(halfPeriod))

	if validHalf.Before(now) {
		// Retry after 60-75 seconds each time
		return c.jitter(60, 15), nil
	}

	// Sleep until the certificate is half expired
	return validHalf.Sub(now) + c.jitter(0, 15), nil
}

func (c *Source) jitter(min, fuzzy int) time.Duration {
	return (time.Duration(min) + time.Duration(rand.Int63n(int64(fuzzy)))) * time.Second
}

func (c *Source) renewOnce(ctx context.Context) error {
	cert, err := c.provider.GetCertificate(ctx)
	if err != nil {
		logger.Errorf("Problem retrieving certificate from provider: %s", err.Error())
		return err
	}

	c.setCertificate(cert)
	return nil
}

func NewSource(ctx context.Context, bindingName string) (*Source, error) {
	sourceLock.Lock()
	defer sourceLock.Unlock()

	src, ok := sources[bindingName]
	if ok && src != nil {
		return src, nil
	}

	provider, err := NewProvider(ctx, bindingName)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create certificate provider")
	}

	logger.Infof("Creating certificate source %q", bindingName)

	cert, err := provider.GetCertificate(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve certificate from provider")
	}

	source := &Source{
		certificate: cert,
		provider:    provider,
		clock:       types.NewClock(ctx),
	}

	if provider.Renewable() {
		logger.Infof("Starting periodic renewal for %q certificate provider", bindingName)
		go source.renew(ctx)
	} else {
		logger.Infof("Certificates from %q binding are not renewable.", bindingName)
	}

	sources[bindingName] = source

	return source, nil
}
