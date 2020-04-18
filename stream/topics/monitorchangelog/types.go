package monitorchangelog

const (
	ChangeTypeServiceProvisioned   = "serviceProvisioned"
	ChangeTypeServiceDeprovisioned = "serviceDeprovisioned"
	ChangeTypeDeviceProvisioned    = "deviceProvisioned"
	ChangeTypeDeviceDeprovisioned  = "deviceDeprovisioned"
)

type Message struct {
	ChangeType        string `json:"changeType"`
	DeviceInstanceId  string `json:"deviceInstanceId"`
	ServiceInstanceId string `json:"serviceInstanceId"`
	TenantId          string `json:"tenantId"`
	Name              string `json:"name"`
	Profile           string `json:"profile"`
	Type              string `json:"type"`
	SpecificType      string `json:"specificType"`
	Category          string `json:"category"`
	HostName          string `json:"hostName"`
	IpAddress         string `json:"ipAddress"`
	SerialKey         string `json:"serialKey"`
	CreatedOn         string `json:"createdOn"`
	CreatedBy         string `json:"createdBy"`
	ModifiedOn        string `json:"modifiedOn"`
	ModifiedBy        string `json:"modifiedBy"`
}
