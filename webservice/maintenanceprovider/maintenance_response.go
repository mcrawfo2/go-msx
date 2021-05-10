package maintenanceprovider

type MaintenanceTask struct {
	Task    string `json:"task"`
	Mode    string `json:"mode"`
	Message string `json:"message"`
}

type MaintenanceResponse struct {
	Mode   string            `json:"mode"`
	Detail []MaintenanceTask `json:"detail"`
}
