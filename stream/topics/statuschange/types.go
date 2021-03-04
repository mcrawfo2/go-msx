package statuschange

const (
	StateCreated              = "created"
	StateUpdated              = "updated"
	StateDeleted              = "deleted"
	StateProvisioning         = "provisioning"
	StateProvisioned          = "provisioned"
	StateProvisioningFailed   = "provisioningFailed"
	StateDeprovisioning       = "deprovisioning"
	StateDeprovisioned        = "deprovisioned"
	StateDeprovisioningFailed = "deprovisioningFailed"

	TopicStatusChange = "STATUS_CHANGE_TOPIC"
	TopicName         = TopicStatusChange
)

type Message struct {
	EntityType string `json:"entityType"`
	EntityId   string `json:"entityId"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}
