package lowerplural

//#if TENANT_DOMAIN
import "github.com/gocql/gocql"

//#endif TENANT_DOMAIN

type lowerCamelSingular struct {
	Name string `db:"name"`
	//#if TENANT_DOMAIN
	TenantId gocql.UUID `db:"tenant_id"`
	//#endif TENANT_DOMAIN
	Data string `db:"data"`
}
