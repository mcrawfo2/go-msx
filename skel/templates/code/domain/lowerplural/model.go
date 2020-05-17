package lowerplural

import "github.com/gocql/gocql"

type lowerCamelSingular struct {
	Name string
	//#if TENANT_DOMAIN
	TenantId gocql.UUID
	//#endif TENANT_DOMAIN
	Data string
}
