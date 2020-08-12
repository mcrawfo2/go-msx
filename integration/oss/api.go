//go:generate mockery --inpackage --name=Api --structname=MockOss

package oss

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

const (
	endpointServiceCancellationCharge = "serviceCancellationCharge"
	endpointNotificationUrl           = "notificationUrl"
	endpointServiceAccess             = "serviceAccess"
	endpointPricingOptions            = "pricingoptions"
	endpointAllowedValues             = "allowedValues"
)

type Api interface {
	GetPricePlanOptions(serviceId, offerId types.UUID, options PricingOptionsRequest) (PricingOptionsResponse, error)
	GetAccessibleServices() (ServicesResponse, error)
	GetAllowedValues(serviceId string, propertyName string) (AllowedValuesResponse, error)
}

type PricingOptionsRequest struct {
	Currency              string                 `json:"currency"`
	Language              string                 `json:"language"`
	SubscriptionId        string                 `json:"subscriptionId"`
	ServiceInstanceDetail map[string]interface{} `json:"serviceInstanceDetail"`
}

type PricingOptionsResponse struct {
	PricePlans []PricingOptionResponse `json:"pricePlans"`
}

type PricingOptionResponse struct {
	Id               types.UUID                       `json:"id"`
	Name             string                           `json:"name"`
	PricePlanDetails map[string]PricePlanOptionDetail `json:"pricePlanDetails"`
}

type PricePlanOptionDetail struct {
	Value                   string `json:"value"`
	OneTimePrice            int    `json:"oneTimePrice"`
	PeriodicPrice           int    `json:"periodicPrice"`
	TimePeriod              string `json:"timePeriod"`
	IncludeQuantity         int    `json:"includeQuantity"`
	AdditionalOneTimePrice  int    `json:"additionalOneTimePrice"`
	AdditionalPeriodicPrice int    `json:"additionalOneTimePrice"`
	AdditionalQuantity      int    `json:"additionalOneTimePrice"`
}

type ServicesResponse struct {
	Services []Service `json:"services"`
}

type Service struct {
	Id string `json:"id"`
}

type AllowedValuesResponse struct {
	AllowedValues []map[string]interface{} `json:"allowedValues"`
}
