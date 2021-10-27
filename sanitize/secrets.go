package sanitize

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"regexp"
	"strings"
)

const (
	configRootSanitizeSecret = "sanitize.secrets"
	keyPlaceHolder           = "__KEY__"
	secretPlaceHolder        = "*****"
)

type Replacer struct {
	From     *regexp.Regexp
	ToString string
	ToBytes  []byte
}

func (r Replacer) ReplaceAllString(source string) (result string) {
	return r.From.ReplaceAllString(source, r.ToString)
}

func (r Replacer) ReplaceAllBytes(source []byte) (result []byte) {
	return r.From.ReplaceAll(source, r.ToBytes)
}

type SecretPattern struct {
	From string `config:"from"`
	To   string `config:"to,default="`
}

func (p SecretPattern) Replacer(key string) (Replacer, error) {
	from := p.From
	if key != "" {
		from = strings.ReplaceAll(from, keyPlaceHolder, regexp.QuoteMeta(key))
	}
	reFrom, err := regexp.Compile(from)
	if err != nil {
		return Replacer{}, err
	}

	to := p.To
	if to == "" {
		to = "${prefix}" + secretPlaceHolder + "${postfix}"
	}

	return Replacer{
		From:     reFrom,
		ToString: to,
		ToBytes:  []byte(to),
	}, nil
}

type SecretPatterns struct {
	Enabled  bool            `config:"enabled,default=true"`
	Patterns []SecretPattern `config:"patterns"`
}

func (p SecretPatterns) Merge(o SecretPatterns) SecretPatterns {
	var result = SecretPatterns{
		Enabled: p.Enabled,
	}

	result.Patterns = append(result.Patterns, o.Patterns...)
	result.Patterns = append(result.Patterns, p.Patterns...)

	return result
}

type SecretConfig struct {
	Enabled        bool           `config:"enabled,default=true"`
	TrustedLoggers []string       `config:"trusted-loggers,default=msx"`
	Keys           []string       `config:"keys"`
	Custom         SecretPatterns `config:"custom"`
	ToString       SecretPatterns `config:"to-string"`
	Json           SecretPatterns `config:"json"`
	Xml            SecretPatterns `config:"xml"`
}

func (c SecretConfig) MergeDefaults() SecretConfig {
	var result = SecretConfig{
		Enabled: c.Enabled,
	}

	o := defaultSecretConfig

	result.TrustedLoggers = append(result.TrustedLoggers, o.TrustedLoggers...)
	result.TrustedLoggers = append(result.TrustedLoggers, c.TrustedLoggers...)

	result.Keys = append(result.Keys, o.Keys...)
	result.Keys = append(result.Keys, c.Keys...)

	result.Custom = c.Custom
	result.ToString = c.ToString.Merge(o.ToString)
	result.Json = c.Json.Merge(o.Json)
	result.Xml = c.Xml.Merge(o.Xml)

	return result
}

func (c SecretConfig) Replacers() ([]Replacer, error) {
	var replacers []Replacer
	var patterns []SecretPattern
	if c.ToString.Enabled {
		patterns = append(patterns, c.ToString.Patterns...)
	}
	if c.Json.Enabled {
		patterns = append(patterns, c.Json.Patterns...)
	}
	if c.Xml.Enabled {
		patterns = append(patterns, c.Xml.Patterns...)
	}

	for _, pattern := range patterns {
		for _, key := range c.Keys {
			replacer, err := pattern.Replacer(key)
			if err != nil {
				return nil, err
			}
			replacers = append(replacers, replacer)
		}
	}

	if c.Custom.Enabled {
		for _, pattern := range c.Custom.Patterns {
			replacer, err := pattern.Replacer("")
			if err != nil {
				return nil, err
			}
			replacers = append(replacers, replacer)
		}
	}

	// Always add our default custom replacers
	for _, pattern := range defaultSecretConfig.Custom.Patterns {
		replacer, err := pattern.Replacer("")
		if err != nil {
			return nil, err
		}
		replacers = append(replacers, replacer)
	}

	return replacers, nil
}

var defaultSecretConfig = SecretConfig{
	TrustedLoggers: nil,
	Keys: []string{
		// keys involved in authentication
		"password",
		"access_token",
		"refresh_token",
		"id_token",
		// Security clients
		"clientSecret",
		// Used when registering NSO via API
		"nsoEncryptionKey",
	},
	Custom: SecretPatterns{
		// configure default patterns to sanitize - don't want this to change easily for security reasons
		Patterns: []SecretPattern{
			{
				From: `access_token":"[a-zA-Z0-9-_.]+"`,
				To:   `access_token":"` + secretPlaceHolder + `"`,
			},
			{
				From: "x-client-secret,value:[^;]+;",
				To:   "x-client-secret,value: " + secretPlaceHolder + " ;",
			},
			// Passwords in typical toString() implementations (assumes the terminating character is not part of the value itself)
			{
				From: `password=[^,)}]+`,
				To:   "password=" + secretPlaceHolder,
			},
			{
				From: `"client_token":"[a-zA-Z0-9-_.]+"`,
				To:   `"client_token":"` + secretPlaceHolder + `"`,
			},
			{
				From: `"accessor":"[a-zA-Z0-9-_.]+"`,
				To:   `"accessor":"` + secretPlaceHolder + `"`,
			},
		},
	},
}

func NewSecretConfig(ctx context.Context) (*SecretConfig, error) {
	var cfg SecretConfig
	if err := config.FromContext(ctx).Populate(&cfg, configRootSanitizeSecret); err != nil {
		return nil, err
	}

	cfg = cfg.MergeDefaults()
	return &cfg, nil
}

type SecretSanitizer struct {
	Replacers []Replacer
}

var secretSanitizer = &SecretSanitizer{
	Replacers: func() []Replacer {
		replacers, _ := defaultSecretConfig.Replacers()
		return replacers
	}(),
}

func (l *SecretSanitizer) SetConfig(cfg *SecretConfig) (err error) {
	l.Replacers, err = cfg.Replacers()
	return err
}

func (l *SecretSanitizer) Secrets(value string) string {
	for _, replacer := range l.Replacers {
		value = replacer.ReplaceAllString(value)
	}
	return value
}

func ConfigureSecretSanitizer(ctx context.Context) error {
	cfg, err := NewSecretConfig(ctx)
	if err != nil {
		return err
	}

	return secretSanitizer.SetConfig(cfg)
}
