package certificate

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"cto-github.cisco.com/NFV-BU/go-msx/background"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/pkg/errors"
	"math/rand"
	"sync"
	"time"
)

var logger = log.NewLogger("msx.certificate")
var random = rand.New(rand.NewSource(time.Now().UnixNano()))
var sources = make(map[string]*Source)
var sourceLock sync.Mutex

type Source struct {
	sync.Mutex
	certificate *tls.Certificate
	provider    Provider
}

// TlsCertificate fetches the certificate for tls.Config
func (c *Source) TlsCertificate(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return c.certificate, nil
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

// renew continuously refreshes the certificate after approximately half of its validity period.
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

		t := time.NewTimer(d)

		select {
		case <-ctx.Done():
			return

		case <-t.C:
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

	validity := parsedCert.NotAfter.Sub(time.Now())
	windowMin, windowMax := int64(0), int64(0)
	if validity >= 10*time.Minute {
		windowMin, windowMax = int64(2*time.Minute), int64(5*time.Minute)
	} else {
		windowMin, windowMax = int64(validity)*1/4, int64(validity)*3/8
	}
	offset := time.Duration(random.Int63n(windowMax-windowMin) + windowMin)
	return (validity / 2) + offset, nil
}

func (c *Source) renewOnce(ctx context.Context) error {
	cert, err := c.provider.GetCertificate(ctx)
	if err != nil {
		logger.Errorf("Problem retrieving certificate from provider: %s", err.Error())
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
