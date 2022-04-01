// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

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
