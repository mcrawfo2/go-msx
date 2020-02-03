package manage

import "cto-github.cisco.com/NFV-BU/go-msx/integration"

type Pojo integration.Pojo
type PojoArray integration.PojoArray
type HealthResult integration.HealthResult
type ErrorDTO integration.ErrorDTO
type ErrorDTO2 integration.ErrorDTO2

type EntityShard struct {
	Name       string      `json:"name"`
	ShardID    string      `json:"shardId"`
	PnpURL     string      `json:"pnpUrl"`
	Host       string      `json:"host"`
	Port       int         `json:"port"`
	Capability string      `json:"capability"`
	EntityID   string      `json:"entityId"`
	EntityType interface{} `json:"entityType"`
	CreatedOn  string      `json:"createdOn"`
	CreatedBy  string      `json:"createdBy"`
	ModifiedOn string      `json:"modifiedOn"`
	ModifiedBy string      `json:"modifiedBy"`
}
