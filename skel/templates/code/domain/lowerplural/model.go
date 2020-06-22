package lowerplural

//#if TENANT_DOMAIN
//#if REPOSITORY_COCKROACH
import "github.com/google/uuid"

//#else REPOSITORY_COCKROACH
import "github.com/gocql/gocql"

//#endif REPOSITORY_COCKROACH
//#endif TENANT_DOMAIN

type lowerCamelSingular struct {
	Name string `db:"name"`
	//#if TENANT_DOMAIN
	//#if REPOSITORY_COCKROACH
	TenantId uuid.UUID `db:"tenant_id"`
	//#else REPOSITORY_COCKROACH
	TenantId gocql.UUID `db:"tenant_id"`
	//#endif REPOSITORY_COCKROACH
	//#endif TENANT_DOMAIN
	Data string `db:"data"`
}
