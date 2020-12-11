package httpclient

import (
	"fmt"
	"net/http"
	_neturl "net/url"
	"testing"
	"time"
)

type ClientValidator struct {
	Description string
	Valid       func(c *http.Client) bool
}

func validateClientTimeout(timeout time.Duration) ClientValidator {
	return ClientValidator{
		Description: fmt.Sprintf("client.Timeout = %s", timeout.String()),
		Valid: func(c *http.Client) bool {
			return c.Timeout == timeout
		},
	}
}

func validateClientCookieJar(jar http.CookieJar) ClientValidator {
	return ClientValidator{
		Description: fmt.Sprintf("client.Jar = %v", jar),
		Valid: func(c *http.Client) bool {
			return c.Jar == jar
		},
	}
}

type DevNullJar struct{}

func (t DevNullJar) SetCookies(u *_neturl.URL, cookies []*http.Cookie) {
}

func (t DevNullJar) Cookies(u *_neturl.URL) []*http.Cookie {
	return nil
}

type TransportValidator struct {
	Description string
	Valid       func(c *http.Transport) bool
}

func validateTlsInsecure(insecure bool) TransportValidator {
	return TransportValidator{
		Description: fmt.Sprintf("transport.TLS.InsecureSkipVerify = %v", insecure),
		Valid: func(c *http.Transport) bool {
			return c.TLSClientConfig.InsecureSkipVerify == insecure
		},
	}
}

func validateNoProxy() TransportValidator {
	return TransportValidator{
		Description: fmt.Sprintf("transport.Proxy == nil"),
		Valid: func(c *http.Transport) bool {
			return c.Proxy == nil
		},
	}
}

func TestClientConfigurer_HttpClient(t *testing.T) {
	tests := []struct {
		name        string
		customizers []ClientConfigurationFunc
		validators  []ClientValidator
	}{
		{
			name: "ClientTimeout",
			customizers: []ClientConfigurationFunc{
				ClientTimeout(30 * time.Second),
			},
			validators: []ClientValidator{
				validateClientTimeout(30 * time.Second),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ClientConfigurer{
				ClientFuncs: tt.customizers,
			}

			client := new(http.Client)
			c.HttpClient(client)

			for _, validator := range tt.validators {
				if !validator.Valid(client) {
					t.Errorf("Client validation failed: %s", validator.Description)
				}
			}
		})
	}
}

func TestClientConfigurer_HttpTransport(t *testing.T) {
	tests := []struct {
		name        string
		customizers []TransportConfigurationFunc
		validators  []TransportValidator
	}{
		{
			name: "NoProxy",
			customizers: []TransportConfigurationFunc{
				NoProxy(),
			},
			validators: []TransportValidator{
				validateNoProxy(),
			},
		},
		{
			name: "TlsInsecure",
			customizers: []TransportConfigurationFunc{
				TlsInsecure(true),
			},
			validators: []TransportValidator{
				validateTlsInsecure(true),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ClientConfigurer{
				TransportFuncs: tt.customizers,
			}

			transport := new(http.Transport)
			transport.TLSClientConfig, _ = NewTlsConfig(&ClientConfig{})
			c.HttpTransport(transport)

			for _, validator := range tt.validators {
				if !validator.Valid(transport) {
					t.Errorf("transport validation failed: %s", validator.Description)
				}
			}
		})
	}
}

func TestCompositeConfigurer_HttpClient(t *testing.T) {
	jar := new(DevNullJar)

	type fields struct {
		Service  Configurer
		Endpoint Configurer
	}
	tests := []struct {
		name       string
		fields     fields
		validators []ClientValidator
	}{
		{
			name: "NoConflict",
			fields: fields{
				Service: ClientConfigurer{
					ClientFuncs: []ClientConfigurationFunc{
						ClientTimeout(3 * time.Millisecond),
					},
				},
				Endpoint: ClientConfigurer{
					ClientFuncs: []ClientConfigurationFunc{
						ClientCookieJar(jar),
					},
				},
			},
			validators: []ClientValidator{
				validateClientTimeout(3 * time.Millisecond),
				validateClientCookieJar(jar),
			},
		},
		{
			name: "Overwrite",
			fields: fields{
				Service: ClientConfigurer{
					ClientFuncs: []ClientConfigurationFunc{
						ClientTimeout(3 * time.Millisecond),
						ClientCookieJar(nil),
					},
				},
				Endpoint: ClientConfigurer{
					ClientFuncs: []ClientConfigurationFunc{
						ClientTimeout(10 * time.Millisecond),
					},
				},
			},
			validators: []ClientValidator{
				validateClientTimeout(10 * time.Millisecond),
				validateClientCookieJar(nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CompositeConfigurer{
				Service:  tt.fields.Service,
				Endpoint: tt.fields.Endpoint,
			}

			client := new(http.Client)
			c.HttpClient(client)

			for _, validator := range tt.validators {
				if !validator.Valid(client) {
					t.Errorf("Client validation failed: %s", validator.Description)
				}
			}
		})
	}
}

func TestCompositeConfigurer_HttpTransport(t *testing.T) {
	type fields struct {
		Service  Configurer
		Endpoint Configurer
	}
	tests := []struct {
		name       string
		fields     fields
		validators []TransportValidator
	}{
		{
			name: "NoConflict",
			fields: fields{
				Service: ClientConfigurer{
					TransportFuncs: []TransportConfigurationFunc{
						TlsInsecure(true),
					},
				},
				Endpoint: ClientConfigurer{
					TransportFuncs: []TransportConfigurationFunc{
						NoProxy(),
					},
				},
			},
			validators: []TransportValidator{
				validateTlsInsecure(true),
				validateNoProxy(),
			},
		},
		{
			name: "Overwrite",
			fields: fields{
				Service: ClientConfigurer{
					TransportFuncs: []TransportConfigurationFunc{
						TlsInsecure(true),
					},
				},
				Endpoint: ClientConfigurer{
					TransportFuncs: []TransportConfigurationFunc{
						TlsInsecure(false),
					},
				},
			},
			validators: []TransportValidator{
				validateTlsInsecure(false),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CompositeConfigurer{
				Service:  tt.fields.Service,
				Endpoint: tt.fields.Endpoint,
			}

			transport := new(http.Transport)
			transport.TLSClientConfig, _ = NewTlsConfig(&ClientConfig{})
			c.HttpTransport(transport)

			for _, validator := range tt.validators {
				if !validator.Valid(transport) {
					t.Errorf("transport validation failed: %s", validator.Description)
				}
			}
		})
	}
}
