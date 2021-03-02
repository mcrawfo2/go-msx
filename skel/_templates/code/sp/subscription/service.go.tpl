package subscription

import (
	"context"
	consume "cto-github.cisco.com/NFV-BU/catalogservice/pkg/api"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/manage"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/_templates/code/sp/api"
	"github.com/pkg/errors"
)

const (
	ServiceLifecycleStateOrdering     = "ORDERING"
	ServiceLifecycleStateProvisioning = "PROVISIONING"
	ServiceLifecycleStateProvisioned  = "PROVISIONED"
	ServiceLifecycleStateDeleting     = "DELETING"
	ServiceLifecycleStateOrderFailed  = "ORDER_FAILED"
)

var (
	errPayloadConversion = errors.New("error converting payload")
)

type subscriptionServiceApi interface {
	CreateSubscription(ctx context.Context, req api.SubscriptionCreateRequest) (subscription, error)
	DeleteSubscription(ctx context.Context, serviceInstanceId string) error
}

type subscriptionService struct {
	subscriptionConverter subscriptionConverter
}

func (s *subscriptionService) CreateSubscription(ctx context.Context, req api.SubscriptionCreateRequest) (subscription, error) {
	consumeApi := consume.NewIntegration(ctx)

	offerResponse, err := consumeApi.GetOffer(req.OfferId)
	if err != nil {
		return subscription{}, err
	}

	offerPayload, ok := offerResponse.Payload.(*consume.ServiceOffering)
	if !ok {
		return subscription{}, errPayloadConversion
	}

	manageApi, err := manage.NewIntegration(ctx)
	if err != nil {
		return subscription{}, err
	}

	response, err := manageApi.CreateSubscription(
		req.TenantId,
		"${service.type}",
		nil,
		map[string]string{},
		map[string]string{
			"id": req.OfferId,
		},
		map[string]string{},
		map[string]string{})
	if err != nil {
		return subscription{}, err
	}

	subscriptionPayload, ok := response.Payload.(*manage.CreateSubscriptionResponse)

	if !ok {
		return subscription{}, errPayloadConversion
	}

	response, err = manageApi.CreateServiceInstance(
		subscriptionPayload.SubscriptionID,
		"",
		map[string]string{}, map[string]string{
			"type":      "${service.type}",
			"offerName": offerPayload.Name,
			"id":        req.ServiceId,
		},
		map[string]string{
			"lifeCycleStatus": "Ordering",
			"txStatus":        "ORDERING",
		})
	if err != nil {
		return subscription{}, errors.Wrap(err, "Failed to submit create subscription request")
	}

	serviceInstanceResponse, ok := response.Payload.(*manage.ServiceInstanceResponse)
	if !ok {
		return subscription{}, errPayloadConversion
	}

	// Set service to provisioned
	err = s.UpdateServiceInstanceStatus(ctx, serviceInstanceResponse.ServiceInstanceID, ServiceLifecycleStateProvisioned)
	if err != nil {
		return subscription{}, errPayloadConversion
	}

	return subscription{
		SubscriptionId:    subscriptionPayload.SubscriptionID,
		ServiceInstanceId: serviceInstanceResponse.ServiceInstanceID,
	}, err
}

func (s *subscriptionService) DeleteSubscription(ctx context.Context, serviceInstanceId string) error {
	manageApi, err := manage.NewIntegration(ctx)

	if err != nil {
		return err
	}

	serviceInstanceResp, err := manageApi.GetServiceInstance(serviceInstanceId)

	if err != nil {
		return err
	}

	serviceInstance, _ := serviceInstanceResp.Payload.(*manage.ServiceInstanceResponse)

	_, err = manageApi.DeleteSubscription(serviceInstance.ServiceInstanceID)

	return nil
}

func (s *subscriptionService) UpdateServiceInstanceStatus(ctx context.Context, serviceInstanceId string, status string) error {
	manageApi, err := manage.NewIntegration(ctx)
	if err != nil {
		return err
	}

	lifeCycleStatus, err := getLifeCycleStatusFromTxStatus(status)
	if err != nil {
		return err
	}

	// mark service instance as deleting
	_, err = manageApi.UpdateServiceInstance(serviceInstanceId, nil, nil, map[string]string{
		"lifeCycleStatus": lifeCycleStatus,
		"txStatus":        status,
	})

	return err
}

func (s *subscriptionService) ServiceInstanceExists(ctx context.Context, serviceInstanceId string) (bool, error) {
	manageApi, err := manage.NewIntegration(ctx)

	if err != nil {
		return false, err
	}

	if _, err := manageApi.GetServiceInstance(serviceInstanceId); err != nil {
		return false, err
	}

	return true, nil
}

func newSubscriptionService(ctx context.Context) subscriptionServiceApi {
	service := serviceFromContext(ctx)
	if service == nil {
		service = &subscriptionService{
			subscriptionConverter: subscriptionConverter{},
		}
	}
	return service
}

func getLifeCycleStatusFromTxStatus(txStatus string) (string, error) {
	switch txStatus {
	case ServiceLifecycleStateDeleting:
		return "Deleting", nil
	case ServiceLifecycleStateOrdering:
		return "Ordering", nil
	case ServiceLifecycleStateProvisioned:
		return "Provisioned", nil
	case ServiceLifecycleStateProvisioning:
		return "Provisioning", nil
	case ServiceLifecycleStateOrderFailed:
		return "Order Failed", nil
	}

	return "", errors.Errorf("Unknown status %q", txStatus)
}
