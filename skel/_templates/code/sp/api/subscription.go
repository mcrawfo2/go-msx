package api

type SubscriptionCreateRequest struct {
	OfferId   string `json:"offerId"`
	TenantId  string `json:"tenantId"`
	ServiceId string `json:"serviceId"`
}

type SubscriptionCreateResponse struct {
	SubscriptonId     string `json:"subscriptionId"`
	ServiceInstanceId string `json:"serviceInstanceId"`
}
