package manage

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"time"
)

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

func (r ServiceInstanceResponse) ServiceAttributes() types.Pojo {
	if r.ServiceAttribute == nil {
		return nil
	}
	if pojo, ok := r.ServiceAttribute.(map[string]interface{}); ok {
		return pojo
	}
	return nil
}

func (r ServiceInstanceResponse) ServiceDefAttributes() types.Pojo {
	if r.ServiceAttribute == nil {
		return nil
	}
	if pojo, ok := r.ServiceAttribute.(map[string]interface{}); ok {
		return pojo
	}
	return nil
}

type DeviceCreateRequest struct {
	Name               string                 `json:"name"`
	TenantId           string                 `json:"tenantId"`
	SubscriptionId     *string                `json:"subscriptionId"`
	ServiceInstanceId  *string                `json:"serviceInstanceId"`
	ServiceType        string                 `json:"serviceType"`
	Model              string                 `json:"model"`
	Type               string                 `json:"type"`
	SubType            string                 `json:"subType"`
	SerialKey          string                 `json:"serialKey"`
	OnboardType        string                 `json:"onboardType"`
	Managed            bool                   `json:"managed"`
	Version            string                 `json:"version"`
	Tags               map[string]string      `json:"tags"`
	Attributes         map[string]string      `json:"attributes"`
	OnboardInformation map[string]interface{} `json:"onboardInformation"`
}

type DeviceUpdateRequest struct {
	Name               string                 `json:"name"`
	ServiceType        string                 `json:"serviceType"`
	Model              string                 `json:"model"`
	Type               string                 `json:"type"`
	SubType            string                 `json:"subType"`
	SerialKey          string                 `json:"serialKey"`
	OnboardType        string                 `json:"onboardType"`
	Managed            bool                   `json:"managed"`
	Version            string                 `json:"version"`
	Tags               map[string]string      `json:"tags"`
	Attributes         map[string]string      `json:"attributes"`
	OnboardInformation map[string]interface{} `json:"onboardInformation"`
}

type DeviceListResponse []DeviceResponse

type DeviceResponse struct {
	Id                 string                 `json:"id"`
	Name               string                 `json:"name"`
	Model              string                 `json:"model"`
	Type               string                 `json:"type"`
	SubType            string                 `json:"subType"`
	Tags               map[string]string      `json:"tags"`
	SerialKey          string                 `json:"serialKey"`
	Version            string                 `json:"version"`
	ServiceInstanceId  string                 `json:"serviceInstanceId"`
	SubscriptionId     string                 `json:"subscriptionId"`
	ServiceType        string                 `json:"serviceType"`
	Managed            bool                   `json:"managed"`
	OnboardType        string                 `json:"onboardType"`
	Attributes         map[string]string      `json:"attributes"`
	OnboardInformation map[string]interface{} `json:"onboardInformation"`
	StatusDetails      struct {
		HealthStatus    DeviceStatusDetail
		PnpStatus       DeviceStatusDetail
		LifeCycleStatus DeviceStatusDetail
		SyncStatus      DeviceStatusDetail
		TunnelStatus    DeviceStatusDetail
	} `json:"statusDetails"`
	TenantId   string `json:"tenantId"`
	ProviderId string `json:"providerId"`
	UserId     string `json:"userId"`
	CreatedOn  string `json:"createdOn"`
	ModifiedOn string `json:"modifiedOn"`
}

type DeviceStatusDetail struct {
	Type               string `json:"type"`
	Name               string `json:"name"`
	Value              string `json:"value"`
	Severity           string `json:"severity"`
	LastUpdated        string `json:"lastUpdated"`
	LastUpdatedMessage string `json:"lastUpdatedMessage"`
}

type DeviceStatusUpdateRequest struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Value   string `json:"value"`
}

type SiteQueryFilter struct {
	DeviceInstanceId  *string `json:"deviceInstanceId"`
	ParentId          *string `json:"parentId"`
	ServiceInstanceId *string `json:"serviceInstanceId"`
	ServiceType       *string `json:"serviceType"`
	ShowImage         *string `json:"showImage"`
	TenantId          *string `json:"tenantId"`
	Type              *string `json:"type"`
}

type SiteCreateRequest struct {
	Address struct {
		Name     string `json:"name"`
		Company  string `json:"company"`
		Address1 string `json:"address1"`
		Address2 string `json:"address2"`
		City     string `json:"city"`
		State    string `json:"state"`
		Country  string `json:"country"`
		PostCode string `json:"postCode"`
	} `json:"address"`
	Attributes map[string]string `json:"attributes"`
	Contact    struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	} `json:"contact"`
	Description       string   `json:"description"`
	DeviceInstanceIds []string `json:"deviceInstanceIds,omitempty"`
	Image             string   `json:"image"`
	Location          struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Name     string `json:"name"`
	ParentId string `json:"parentId"`
	TenantId string `json:"tenantId"`
	Type     string `json:"type"`
}

type SiteUpdateRequest struct {
	Address struct {
		Name     string `json:"name"`
		Company  string `json:"company"`
		Address1 string `json:"address1"`
		Address2 string `json:"address2"`
		City     string `json:"city"`
		State    string `json:"state"`
		Country  string `json:"country"`
		PostCode string `json:"postCode"`
	} `json:"address"`
	Attributes map[string]string `json:"attributes"`
	Contact    struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	} `json:"contact"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Location    struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Name     string `json:"name"`
	ParentId string `json:"parentId"`
	Type     string `json:"type"`
}

type SiteStatusUpdateRequest struct {
	LastUpdatedMessage string `json:"lastUpdatedMessage"`
	Severity           string `json:"severity"`
	Value              string `json:"value"`
}

type ControlPlaneResponse struct {
	ControlPlaneId     types.UUID        `json:"controlPlaneId"`
	ResourceProvider   string            `json:"resourceProvider"`
	TenantId           types.UUID        `json:"tenantId"`
	AuthenticationType string            `json:"authenticationType"`
	Url                string            `json:"url"`
	Name               string            `json:"name"`
	TlsInsecure        bool              `json:"tlsInsecure"`
	Attributes         map[string]string `json:"attributes"`
}

type AttachTemplateRequest struct {
	TemplateDetails []TemplateDetails `json:"templateDetails"`
}

type AttachTemplateResponse []struct {
	ID             string           `json:"id"`
	DeviceID       string           `json:"deviceId"`
	TemplateID     string           `json:"templateId"`
	TemplateParams []TemplateParams `json:"templateParams"`
	Status         string           `json:"status"`
	LastUpdated    string           `json:"lastUpdated"`
	UserID         string           `json:"userId"`
}

type TemplateDetails struct {
	TemplateID     string           `json:"templateId"`
	TemplateParams []TemplateParams `json:"templateParams"`
}

type TemplateParams struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DeviceTemplateAccess struct {
	Global  bool     `json:"global,omitempty"`
	Tenants []string `json:"tenants,omitempty"`
}

type DeviceTemplateAccessResponse struct {
	Global               bool     `json:"global"`
	SuccessListOfTenants []string `json:successListOfTenants`
	FailureListOfTenants []string `json:failureListOfTenants`
}

type DeviceTemplateCreateRequest struct {
	DeviceTemplateAccess DeviceTemplateAccess       `json:"access"`
	ConfigContent        string                     `json:"configContent"`
	Description          string                     `json:"description"`
	DeviceModels         []string                   `json:"deviceModels"`
	Name                 string                     `json:"name"`
	ResourceProvider     string                     `json:"resourceProvider"`
	ServiceType          string                     `json:"serviceType"`
	TemplateStandard     string                     `json:"templateStandard"`
	Validators           []DeviceTemplateValidators `json:"validators"`
	Version              string                     `json:"version"`
}
type DeviceActionCreateRequest struct {
	ActionConfig          map[string]string `json:"actionConfig"`
	ActionType            string            `json:"actionType"`
	Description           string            `json:"description"`
	DeviceTypes           []string          `json:"deviceTypes"`
	Name                  string            `json:"name"`
	Owner                 string            `json:"owner"`
	SupportsBulkOperation bool              `json:"supportsBulkOperation"`
}

type DeviceActionCreateRequests []DeviceActionCreateRequest

type DeviceTemplateCreateResponse struct {
	Id                   types.UUID                 `json:"id"`
	DeviceTemplateAccess DeviceTemplateAccess       `json:"access"`
	ConfigContent        string                     `json:"configContent"`
	Description          string                     `json:"description"`
	DeviceModels         []string                   `json:"deviceModels"`
	Name                 string                     `json:"name"`
	ResourceProvider     string                     `json:"resourceProvider"`
	ServiceType          string                     `json:"serviceType"`
	TemplateStandard     string                     `json:"templateStandard"`
	Validators           []DeviceTemplateValidators `json:"validators"`
	Version              string                     `json:"version"`
}

type DeviceTemplateListItemResponse struct {
	Id               types.UUID                 `json:"id"`
	Name             string                     `json:"name"`
	Version          string                     `json:"version"`
	ServiceType      string                     `json:"serviceType"`
	DeviceModels     []string                   `json:"deviceModels"`
	ResourceProvider string                     `json:"resourceProvider"`
	TemplateStandard string                     `json:"templateStandard"`
	CreatedOn        time.Time                  `json:"createdOn"`
	UserId           *types.UUID                `json:"userId"`
	Description      string                     `json:"description"`
	Validators       []DeviceTemplateValidators `json:"validators"`
}

type DeviceTemplateResponse struct {
	Id                   types.UUID                 `json:"id"`
	Name                 string                     `json:"name"`
	Version              string                     `json:"version"`
	ServiceType          string                     `json:"serviceType"`
	DeviceModels         []string                   `json:"deviceModels"`
	ConfigContent        string                     `json:"configContent"`
	ResourceProvider     string                     `json:"resourceProvider"`
	TemplateStandard     string                     `json:"templateStandard"`
	CreatedOn            time.Time                  `json:"createdOn"`
	UserId               *types.UUID                `json:"userId"`
	Description          string                     `json:"description"`
	Validators           []DeviceTemplateValidators `json:"validators"`
	DeviceTemplateAccess DeviceTemplateAccess       `json:"access"`
}

type DeviceTemplateValidators struct {
	AllowedValues []string `json:"allowedValues"`
	DisplayType   string   `json:"displayType"`
	HintText      string   `json:"hintText"`
	Label         string   `json:"label"`
	Name          string   `json:"name"`
	ToolTipText   string   `json:"toolTipText"`
	Type          string   `json:"type"`
}

type DeviceConnectionCreateRequest struct {
	Category          string `json:"category"`
	DeviceInstanceId  string `json:"deviceInstanceId"`
	HostName          string `json:"hostName"`
	IpAddress         string `json:"ipAddress"`
	Name              string `json:"name"`
	Profile           string `json:"profile"`
	SerialKey         string `json:"serialKey"`
	ServiceInstanceId string `json:"serviceInstanceId"`
	SpecificType      string `json:"specificType"`
	Type              string `json:"type"`
}

type DeviceConnectionResponse struct {
	Category          string `json:"category"`
	CreatedBy         string `json:"createdBy"`
	CreatedOn         string `json:"createdOn"`
	DeviceInstanceId  string `json:"deviceInstanceId"`
	HostName          string `json:"hostName"`
	IpAddress         string `json:"ipAddress"`
	ModifiedBy        string `json:"modifiedBy"`
	ModifiedOn        string `json:"modifiedOn"`
	Name              string `json:"name"`
	Profile           string `json:"profile"`
	SerialKey         string `json:"serialKey"`
	ServiceInstanceId string `json:"serviceInstanceId"`
	SpecificType      string `json:"specificType"`
	TenantId          string `json:"tenantId"`
	Type              string `json:"type"`
}

const (
	ServiceLifecycleStateOrdering           = "ORDERING"
	ServiceLifecycleStateOrderFailed        = "ORDER_FAILED"
	ServiceLifecycleStateProvisioning       = "PROVISIONING"
	ServiceLifecycleStateProvisioned        = "PROVISIONED"
	ServiceLifecycleStateProvisioningFailed = "PROVISIONING_FAILED"
	ServiceLifecycleStateDeleting           = "DELETING"
	ServiceLifecycleStateDeleteFailed       = "DELETE_FAILED"
)
