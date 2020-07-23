package subscription

import "cto-github.cisco.com/NFV-BU/go-msx/skel/templates/code/sp/api"

type subscriptionConverter struct{}

func (c *subscriptionConverter) ToCreateResponse(subscription subscription) api.SubscriptionCreateResponse {
	return api.SubscriptionCreateResponse{
		SubscriptonId:     subscription.SubscriptionId,
		ServiceInstanceId: subscription.ServiceInstanceId,
	}
}
