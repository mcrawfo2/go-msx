package tokensource
import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	tsconfig "cto-github.cisco.com/NFV-BU/go-msx/vault/tokensource/config"
)

//TokenSource interface represents a mechanism for retrival and management of Vault Tokens
type TokenSource interface {
	GetToken(cfg *config.Config) (token string, err error)
}

//GetTokenSource will return a TokenSource implementation based provided config
//Currently Config based source and Kubernetes Auth are implemented
func GetTokenSource(source string) TokenSource {
	switch source {
	default:
		return &tsconfig.ConfigSource{}
	}
}