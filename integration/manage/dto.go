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

type CreateSubscriptionResponse struct {
	SubscriptionID   string `json:"subscriptionId"`
	SubscriptionName string `json:"subscriptionName"`
	UserID           string `json:"userId"`
	ProviderID       string `json:"providerId"`
	TenantID         string `json:"tenantId"`
	ServiceType      string `json:"serviceType"`
	CostAttribute    struct {
		CustomAttribute string `json:"customAttribute"`
	} `json:"costAttribute"`
	OfferDefAttribute struct {
		ID string `json:"id"`
	} `json:"offerDefAttribute"`
	OfferSelectionDetail struct {
		PriceDetail           string `json:"priceDetail"`
		ServiceInstanceDetail string `json:"serviceInstanceDetail"`
	} `json:"offerSelectionDetail"`
	SubscriptionAttribute struct {
		Configuration    string `json:"configuration"`
		NsoResponseTypes string `json:"nsoResponseTypes"`
	} `json:"subscriptionAttribute"`
	CreatedOn       string      `json:"createdOn"`
	ModifiedOn      string      `json:"modifiedOn"`
	ServiceList     interface{} `json:"serviceList"`
	RemoteUserCount interface{} `json:"remoteUserCount"`
}

type ServiceInstanceResponse struct {
	ProviderID        string      `json:"providerId"`
	TenantID          string      `json:"tenantId"`
	UserID            string      `json:"userId"`
	ServiceInstanceID string      `json:"serviceInstanceId"`
	SubscriptionID    string      `json:"subscriptionId"`
	TenantName        string      `json:"tenantName"`
	CreatedOn         string      `json:"createdOn"`
	ModifiedOn        interface{} `json:"modifiedOn"`
	ProvisionedOn     interface{} `json:"provisionedOn"`
	Status            struct {
		TxStatus        string `json:"txStatus"`
		LifeCycleStatus string `json:"lifeCycleStatus"`
	} `json:"status"`
	ServiceDefAttribute interface{} `json:"serviceDefAttribute"`
	ServiceAttribute    interface{} `json:"serviceAttribute"`
}

type DeviceCreateRequest struct {
	Name          				string            `json:"name"`
	TenantId          			string            `json:"tenantId"`
	SubscriptionId				string			  `json:"subscriptionId"`
	ServiceInstanceId			string			  `json:"serviceInstanceId"`
	ServiceType          		string            `json:"serviceType"`
	Model          				string            `json:"model"`
	Type         				string            `json:"type"`
	SubType       				string            `json:"subType"`
	SerialKey        			string            `json:"serialKey"`
	DeviceOnboardingType 		string            `json:"deviceOnboardingType"`
	Managed          			bool         	  `json:"managed"`
	Version          			string            `json:"version"`
	Tags    					map[string]string `json:"tags"`
	Attributes    				map[string]string `json:"attributes"`
	DeviceOnboardInformation    map[string]string `json:"deviceOnboardInformation"`
}

type DeviceResponse struct {
	Id          				string            `json:"id"`
	Name          				string            `json:"name"`
	Model          				string            `json:"model"`
	Type         				string            `json:"type"`
	SubType       				string            `json:"subType"`
	Tags    					map[string]string `json:"tags"`
	SerialKey        			string            `json:"serialKey"`
	Version          			string            `json:"version"`
	ServiceInstanceId			string			  `json:"serviceInstanceId"`
	SubscriptionId				string 			  `json:"subscriptionId"`
	ServiceType          		string            `json:"serviceType"`
	Managed          			bool         	  `json:"managed"`
	DeviceOnboardingType 		string            `json:"deviceOnboardingType"`
	Attributes    				map[string]string `json:"attributes"`
	DeviceOnboardInformation    map[string]string `json:"deviceOnboardInformation"`
	StatusDetails    			struct{
		HealthStatus			DeviceStatusDetail
		PnpStatus				DeviceStatusDetail
		LifeCycleStatus			DeviceStatusDetail
		SyncStatus				DeviceStatusDetail
		TunnelStatus			DeviceStatusDetail
	} `json:"statusDetails"`
	TenantId          			string            `json:"tenantId"`
	ProviderId          		string            `json:"providerId"`
	UserId          			string            `json:"userId"`
	CreatedOn					string			  `json:"createdOn"`
	ModifiedOn					string			  `json:"modifiedOn"`
}

type DevicesResponse struct {
	Content []DeviceResponse `json:"content"`
}

type DeviceStatusDetail struct {
	Type   				string	`json:"type"`
	Name					string	`json:"name"`
	Value					string	`json:"value"`
	Severity    			string	`json:"severity"`
	LastUpdated			string	`json:"lastUpdated"`
	LastUpdatedMessage		string	`json:"lastUpdatedMessage"`
}

type DeviceStatusUpdateRequest struct {
	Message   				string	`json:"message"`
	Type   					string	`json:"type"`
	Value   				string	`json:"value"`
}