package vault

import (
	"strings"
	"time"
)

type IssueCertificateRequest struct {
	CommonName string
	Ttl time.Duration
	AltNames []string
	IpSans []string
}

func (r IssueCertificateRequest) Data() map[string]interface{} {
	return map[string]interface{}{
		"common_name": r.CommonName,
		"ttl":         r.Ttl.String(),
		"alt_names":   strings.Join(r.AltNames, ","),
		"ip_sans":     strings.Join(r.IpSans, ","),
	}
}
