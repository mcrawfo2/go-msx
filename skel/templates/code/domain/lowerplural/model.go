package lowerplural

//#if TENANT_DOMAIN
import "github.com/gocql/gocql"

//#endif TENANT_DOMAIN

type lowerCamelSingular struct {
	Name string
	//#if TENANT_DOMAIN
	TenantId gocql.UUID
	//#endif TENANT_DOMAIN
	Data string
}
