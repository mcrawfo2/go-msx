package tokensource
import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/hashicorp/vault/api"
)

var logger  = log.NewLogger("msx.vault.tokensource")

//TokenSource interface represents a mechanism for retrival and management of Vault Tokens
type TokenSource interface {
	GetToken(cfg *config.Config) (token string, err error)
	StartRenewer(client *api.Client)
}

//GetTokenSource will return a TokenSource implementation based provided config
//Currently Config based source and Kubernetes Auth are implemented
func GetTokenSource(source string) TokenSource {
	switch source {
	default:
		return &ConfigSource{}
	}
}