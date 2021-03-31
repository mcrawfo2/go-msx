package manage

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/clienttest"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"net/http"
	"reflect"
	"strconv"
	"testing"
)

func TestNewIntegration(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	ctxWithConfig := configtest.ContextWithNewInMemoryConfig(
		context.Background(),
		map[string]string{
			"remoteservice.manageservice.service": "manageservice",
		})
	tests := []struct {
		name string
		args args
		want Api
	}{
		{
			name: "NonExisting",
			args: args{
				ctx: ctxWithConfig,
			},
			want: &Integration{
				MsxServiceExecutor: integration.NewMsxService(ctxWithConfig, serviceName, endpoints),
			},
		},
		{
			name: "Existing",
			args: args{
				ctx: ContextWithIntegration(ctxWithConfig, &Integration{}),
			},
			want: &Integration{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := NewIntegration(tt.args.ctx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIntegration() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type ManageIntegrationTest struct {
	*clienttest.EndpointTest
}

func NewManageIntegrationTest() *ManageIntegrationTest {
	return &ManageIntegrationTest{
		EndpointTest: new(clienttest.EndpointTest).WithEndpoints(endpoints),
	}
}

type ManageCall func(t *testing.T, api Api) (*integration.MsxResponse, error)

func (m *ManageIntegrationTest) WithCall(call ManageCall) *ManageIntegrationTest {
	m.EndpointTest.WithCall(func(t *testing.T, executor integration.MsxServiceExecutor) (*integration.MsxResponse, error) {
		return call(t, NewIntegrationWithExecutor(executor))
	})
	return m
}

func TestIntegration_GetAdminHealth(t *testing.T) {
	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetAdminHealth()
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(&integration.HealthDTO{
			Status: "Up",
		}).
		WithRequestPredicate(clienttest.EndpointRequestHasName(endpointNameGetAdminHealth)).
		WithRequestPredicate(clienttest.EndpointRequestHasToken(false)).
		WithEndpointPredicate(clienttest.ServiceEndpointHasMethod(http.MethodGet)).
		WithEndpointPredicate(clienttest.ServiceEndpointHasPath("/admin/health")).
		Test(t)
}

func TestIntegration_GetSubscription(t *testing.T) {
	const subscriptionId = "subscription-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetSubscription(subscriptionId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(Pojo)).
		WithRequestPredicate(clienttest.EndpointRequestHasName(endpointNameGetSubscription)).
		WithRequestPredicate(clienttest.EndpointRequestHasExpectEnvelope(true)).
		WithRequestPredicate(clienttest.EndpointRequestHasEndpointParameter("subscriptionId", subscriptionId)).
		WithEndpointPredicate(clienttest.ServiceEndpointHasMethod(http.MethodGet)).
		WithEndpointPredicate(clienttest.ServiceEndpointHasPath("/api/v2/subscriptions/{{.subscriptionId}}")).
		Test(t)
}

func TestIntegration_GetSubscriptionsV3(t *testing.T) {
	tests := []struct {
		name        string
		serviceType string
		page        int
		pageSize    int
	}{
		{
			name:        "Simple",
			serviceType: "service-type",
			page:        0,
			pageSize:    10,
		},
		{
			name:        "AnyServiceType",
			serviceType: "",
			page:        0,
			pageSize:    10,
		},
		{
			name:        "RandomPage",
			serviceType: "",
			page:        5,
			pageSize:    100,
		},
	}

	for _, tt := range tests {
		test := NewManageIntegrationTest().
			WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
				return api.GetSubscriptionsV3(tt.serviceType, tt.page, tt.pageSize)
			}).
			WithResponseStatus(http.StatusOK).
			WithResponseEnvelope().
			WithResponsePayload(new(Pojo)).
			WithRequestPredicates(
				clienttest.EndpointRequestHasName(endpointNameGetSubscriptionsV3),
				clienttest.EndpointRequestHasExpectEnvelope(true),
				clienttest.EndpointRequestHasQueryParam("page", strconv.Itoa(tt.page)),
				clienttest.EndpointRequestHasQueryParam("pageSize", strconv.Itoa(tt.pageSize))).
			WithEndpointPredicate(clienttest.ServiceEndpointHasMethod(http.MethodGet))

		if tt.serviceType != "" {
			test.WithRequestPredicate(clienttest.EndpointRequestHasQueryParam("serviceType", tt.serviceType))
		}

		t.Run(tt.name, test.Test)
	}
}

func TestIntegration_CreateSubscription(t *testing.T) {
	const tenantId = "tenant-id"
	const serviceType = "service-type"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.CreateSubscription(tenantId, serviceType, nil, nil, nil, nil, nil)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(Pojo)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameCreateSubscription),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("tenantId", tenantId),
			clienttest.EndpointRequestHasBodyJsonValue("serviceType", serviceType)).
		WithEndpointPredicate(clienttest.ServiceEndpointHasMethod(http.MethodPost)).
		Test(t)
}

func TestIntegration_UpdateSubscription(t *testing.T) {
	const subscriptionId = "subscription-id"
	const serviceType = "service-type"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.UpdateSubscription(subscriptionId, serviceType, nil, nil, nil, nil, nil)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(Pojo)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameUpdateSubscription),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("subscriptionId", subscriptionId),
			clienttest.EndpointRequestHasBodyJsonValue("serviceType", serviceType)).
		WithEndpointPredicate(clienttest.ServiceEndpointHasMethod(http.MethodPut)).
		Test(t)
}

func TestIntegration_DeleteSubscription(t *testing.T) {
	const subscriptionId = "subscription-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.DeleteSubscription(subscriptionId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameDeleteSubscription),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("subscriptionId", subscriptionId)).
		WithEndpointPredicate(clienttest.ServiceEndpointHasMethod(http.MethodDelete)).
		Test(t)
}

func TestIntegration_GetServiceInstance(t *testing.T) {
	const serviceInstanceId = "service-instance-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetServiceInstance(serviceInstanceId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(ServiceInstanceResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetServiceInstance),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("serviceInstanceId", serviceInstanceId),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v2/serviceinstances/{{.serviceInstanceId}}"),
		).
		Test(t)
}

func TestIntegration_GetSubscriptionServiceInstances(t *testing.T) {
	const subscriptionId = "subscription-id"
	const page = 0
	const pageSize = 100

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetSubscriptionServiceInstances(subscriptionId, page, pageSize)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(&paging.PaginatedResponse{
			Content: new(PojoArray),
		}).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetSubscriptionServiceInstances),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("subscriptionId", subscriptionId),
			clienttest.EndpointRequestHasQueryParam("page", strconv.Itoa(page)),
			clienttest.EndpointRequestHasQueryParam("pageSize", strconv.Itoa(pageSize)),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v2/serviceinstances/subscriptions/{{.subscriptionId}}"),
		).
		Test(t)
}

func TestIntegration_CreateServiceInstance(t *testing.T) {
	const subscriptionId = "subscription-id"
	const serviceInstanceId = "service-instance-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.CreateServiceInstance(
				subscriptionId,
				serviceInstanceId,
				map[string]string{
					"service-attribute-1": "service-value-1",
				},
				map[string]string{
					"service-def-attribute-1": "service-value-2",
				},
				map[string]string{
					"status-attribute-1": "service-value-3",
				})
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(ServiceInstanceResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameCreateServiceInstance),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("subscriptionId", subscriptionId),
			clienttest.EndpointRequestHasBodyJsonValue("serviceInstanceId", serviceInstanceId),
			clienttest.EndpointRequestHasBodyJsonValue("serviceAttribute.service-attribute-1", "service-value-1"),
			clienttest.EndpointRequestHasBodyJsonValue("serviceDefAttribute.service-def-attribute-1", "service-value-2"),
			clienttest.EndpointRequestHasBodyJsonValue("status.status-attribute-1", "service-value-3"),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPost),
			clienttest.ServiceEndpointHasPath("/api/v1/serviceinstances/subscriptions/{{.subscriptionId}}"),
		).
		Test(t)
}

func TestIntegration_UpdateServiceInstance(t *testing.T) {
	const serviceInstanceId = "service-instance-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.UpdateServiceInstance(
				serviceInstanceId,
				map[string]string{
					"service-attribute-1": "service-value-1",
				},
				map[string]string{
					"service-def-attribute-1": "service-value-2",
				},
				map[string]string{
					"status-attribute-1": "service-value-3",
				})
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(ServiceInstanceResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameUpdateServiceInstance),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("serviceInstanceId", serviceInstanceId),
			clienttest.EndpointRequestHasBodyJsonValue("serviceInstanceId", serviceInstanceId),
			clienttest.EndpointRequestHasBodyJsonValue("serviceAttribute.service-attribute-1", "service-value-1"),
			clienttest.EndpointRequestHasBodyJsonValue("serviceDefAttribute.service-def-attribute-1", "service-value-2"),
			clienttest.EndpointRequestHasBodyJsonValue("status.status-attribute-1", "service-value-3"),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPut),
			clienttest.ServiceEndpointHasPath("/api/v1/serviceinstances/{{.serviceInstanceId}}"),
		).
		Test(t)
}

func TestIntegration_DeleteServiceInstance(t *testing.T) {
	const serviceInstanceId = "service-instance-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.DeleteServiceInstance(serviceInstanceId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameDeleteServiceInstance),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("serviceInstanceId", serviceInstanceId)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodDelete),
			clienttest.ServiceEndpointHasPath("/api/v1/serviceinstances/{{.serviceInstanceId}}"),
		).
		Test(t)
}

func TestIntegration_GetSite(t *testing.T) {
	const siteId = "site-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetSite(siteId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(Pojo)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetSite),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("siteId", siteId),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v1/sites/{{.siteId}}"),
		).
		Test(t)
}

func TestIntegration_CreateSite(t *testing.T) {
	const subscriptionId = "subscription-id"
	const serviceInstanceId = "service-instance-id"
	const siteId = "site-id"
	const siteName = "site-name"
	const siteType = "site-type"
	const siteDisplayName = "site-display-name"
	var siteAttributes = map[string]string{
		"site-attribute-1": "site-value-1",
	}
	var siteDefAttributes = map[string]string{
		"site-def-attribute-1": "site-value-2",
	}
	var siteDevices = []string{
		"device-1",
		"device-2",
	}

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.CreateSite(
				subscriptionId,
				serviceInstanceId,
				types.NewOptionalStringFromString(siteId).Ptr(),
				types.NewOptionalStringFromString(siteName).Ptr(),
				types.NewOptionalStringFromString(siteType).Ptr(),
				types.NewOptionalStringFromString(siteDisplayName).Ptr(),
				siteAttributes,
				siteDefAttributes,
				siteDevices)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(Pojo)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameCreateSite),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("subscriptionId", subscriptionId),
			clienttest.EndpointRequestHasBodyJsonValue("serviceInstanceId", serviceInstanceId),
			clienttest.EndpointRequestHasBodyJsonValue("siteId", siteId),
			clienttest.EndpointRequestHasBodyJsonValue("siteName", siteName),
			clienttest.EndpointRequestHasBodyJsonValue("siteType", siteType),
			clienttest.EndpointRequestHasBodyJsonValue("displayName", siteDisplayName),
			clienttest.EndpointRequestHasBodyJsonValue("siteAttributes.site-attribute-1", siteAttributes["site-attribute-1"]),
			clienttest.EndpointRequestHasBodyJsonValue("siteDefAttributes.site-def-attribute-1", siteDefAttributes["site-def-attribute-1"]),
			clienttest.EndpointRequestHasBodyJsonValue("devices.#", float64(2)),
			clienttest.EndpointRequestHasBodyJsonValue("devices.0", siteDevices[0]),
			clienttest.EndpointRequestHasBodyJsonValue("devices.1", siteDevices[1]),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPost),
			clienttest.ServiceEndpointHasPath("/api/v1/sites/subscriptions/{{.subscriptionId}}"),
		).
		Test(t)
}

func TestIntegration_UpdateSite(t *testing.T) {
	const siteId = "site-id"
	const siteName = "site-name"
	const siteType = "site-type"
	const siteDisplayName = "site-display-name"
	var siteAttributes = map[string]string{
		"site-attribute-1": "site-value-1",
	}
	var siteDefAttributes = map[string]string{
		"site-def-attribute-1": "site-value-2",
	}
	var siteDevices = []string{
		"device-1",
		"device-2",
	}

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.UpdateSite(
				siteId,
				types.NewOptionalStringFromString(siteType).Ptr(),
				types.NewOptionalStringFromString(siteDisplayName).Ptr(),
				siteAttributes,
				siteDefAttributes,
				siteDevices)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(Pojo)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameUpdateSite),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("siteId", siteId),
			clienttest.EndpointRequestHasBodyJsonValue("siteId", siteId),
			clienttest.EndpointRequestHasBodyJsonValue("siteType", siteType),
			clienttest.EndpointRequestHasBodyJsonValue("displayName", siteDisplayName),
			clienttest.EndpointRequestHasBodyJsonValue("siteAttributes.site-attribute-1", siteAttributes["site-attribute-1"]),
			clienttest.EndpointRequestHasBodyJsonValue("siteDefAttributes.site-def-attribute-1", siteDefAttributes["site-def-attribute-1"]),
			clienttest.EndpointRequestHasBodyJsonValue("devices.#", float64(2)),
			clienttest.EndpointRequestHasBodyJsonValue("devices.0", siteDevices[0]),
			clienttest.EndpointRequestHasBodyJsonValue("devices.1", siteDevices[1]),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPut),
			clienttest.ServiceEndpointHasPath("/api/v1/sites/{{.siteId}}"),
		).
		Test(t)
}

func TestIntegration_DeleteSite(t *testing.T) {
	const siteId = "site-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.DeleteSite(siteId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameDeleteSite),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("siteId", siteId),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodDelete),
			clienttest.ServiceEndpointHasPath("/api/v1/site/{{.siteId}}"),
		).
		Test(t)
}

func TestIntegration_GetSitesV3(t *testing.T) {
	const page = 0
	const pageSize = 100
	const deviceInstanceId = "device-instance-id"
	const serviceInstanceId = "service-instance-id"
	const parentId = "parentId"
	const serviceType = "service-type"
	const tenantId = "tenant-id"
	const siteType = "site-type"
	const showImage = "true"
	var siteFilter = SiteQueryFilter{
		DeviceInstanceId:  types.NewOptionalStringFromString(deviceInstanceId).Ptr(),
		ParentId:          types.NewOptionalStringFromString(parentId).Ptr(),
		ServiceInstanceId: types.NewOptionalStringFromString(serviceInstanceId).Ptr(),
		ServiceType:       types.NewOptionalStringFromString(serviceType).Ptr(),
		ShowImage:         types.NewOptionalStringFromString(showImage).Ptr(),
		TenantId:          types.NewOptionalStringFromString(tenantId).Ptr(),
		Type:              types.NewOptionalStringFromString(siteType).Ptr(),
	}

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetSitesV3(siteFilter, page, pageSize)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(Pojo)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetSitesV3),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasQueryParam("deviceInstanceId", deviceInstanceId),
			clienttest.EndpointRequestHasQueryParam("parentId", parentId),
			clienttest.EndpointRequestHasQueryParam("serviceInstanceId", serviceInstanceId),
			clienttest.EndpointRequestHasQueryParam("serviceType", serviceType),
			clienttest.EndpointRequestHasQueryParam("showImage", showImage),
			clienttest.EndpointRequestHasQueryParam("tenantId", tenantId),
			clienttest.EndpointRequestHasQueryParam("type", siteType),
			clienttest.EndpointRequestHasQueryParam("page", strconv.Itoa(page)),
			clienttest.EndpointRequestHasQueryParam("pageSize", strconv.Itoa(pageSize)),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v3/sites"),
		).
		Test(t)
}

func TestIntegration_GetSiteV3(t *testing.T) {
	const siteId = "site-id"
	const showImage = "true"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetSiteV3(siteId, showImage)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(Pojo)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetSiteV3),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("siteId", siteId),
			clienttest.EndpointRequestHasQueryParam("showImage", showImage),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v3/sites/{{.siteId}}"),
		).
		Test(t)
}

func TestIntegration_CreateSiteV3(t *testing.T) {
	const tenantId = "tenant-id"
	const parentId = "parent-id"
	const siteName = "site-name"
	const siteDescription = "site-description"
	const siteImage = "site-image"
	const siteType = "site-type"
	const siteAddressName = "site-address-name"
	const siteContactName = "site-contact-name"

	var address = struct {
		Name     string `json:"name"`
		Company  string `json:"company"`
		Address1 string `json:"address1"`
		Address2 string `json:"address2"`
		City     string `json:"city"`
		State    string `json:"state"`
		Country  string `json:"country"`
		PostCode string `json:"postCode"`
	}{
		Name: siteAddressName,
	}
	var contact = struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}{
		Name: siteContactName,
	}
	var location = struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}{}
	var attributes = map[string]string{
		"site-attribute-1": "site-value-1",
	}
	var deviceInstanceIds = []string{
		"device-instance-1",
		"device-instance-2",
	}

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			siteCreateRequest := SiteCreateRequest{
				Address:           address,
				Attributes:        attributes,
				Contact:           contact,
				Description:       siteDescription,
				DeviceInstanceIds: deviceInstanceIds,
				Image:             siteImage,
				Location:          location,
				Name:              siteName,
				ParentId:          parentId,
				TenantId:          tenantId,
				Type:              siteType,
			}
			return api.CreateSiteV3(siteCreateRequest)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(Pojo)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameCreateSiteV3),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasBodyJsonValue("address.name", siteAddressName),
			clienttest.EndpointRequestHasBodyJsonValue("attributes.site-attribute-1", attributes["site-attribute-1"]),
			clienttest.EndpointRequestHasBodyJsonValue("contact.name", siteContactName),
			clienttest.EndpointRequestHasBodyJsonValue("description", siteDescription),
			clienttest.EndpointRequestHasBodyJsonValue("deviceInstanceIds.#", float64(2)),
			clienttest.EndpointRequestHasBodyJsonValue("deviceInstanceIds.0", deviceInstanceIds[0]),
			clienttest.EndpointRequestHasBodyJsonValue("deviceInstanceIds.1", deviceInstanceIds[1]),
			clienttest.EndpointRequestHasBodyJsonValue("image", siteImage),
			clienttest.EndpointRequestHasBodyJsonValue("name", siteName),
			clienttest.EndpointRequestHasBodyJsonValue("parentId", parentId),
			clienttest.EndpointRequestHasBodyJsonValue("tenantId", tenantId),
			clienttest.EndpointRequestHasBodyJsonValue("type", siteType),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPost),
			clienttest.ServiceEndpointHasPath("/api/v3/sites"),
		).
		Test(t)
}

func TestIntegration_UpdateSiteV3(t *testing.T) {
	const siteId = "site-id"
	const parentId = "parent-id"
	const siteName = "site-name"
	const siteDescription = "site-description"
	const siteImage = "site-image"
	const siteType = "site-type"
	const siteAddressName = "site-address-name"
	const siteContactName = "site-contact-name"

	var address = struct {
		Name     string `json:"name"`
		Company  string `json:"company"`
		Address1 string `json:"address1"`
		Address2 string `json:"address2"`
		City     string `json:"city"`
		State    string `json:"state"`
		Country  string `json:"country"`
		PostCode string `json:"postCode"`
	}{
		Name: siteAddressName,
	}
	var contact = struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}{
		Name: siteContactName,
	}
	var location = struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}{}
	var attributes = map[string]string{
		"site-attribute-1": "site-value-1",
	}
	const notification = "notification"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			siteCreateRequest := SiteUpdateRequest{
				Address:     address,
				Attributes:  attributes,
				Contact:     contact,
				Description: siteDescription,
				Image:       siteImage,
				Location:    location,
				Name:        siteName,
				ParentId:    parentId,
				Type:        siteType,
			}
			return api.UpdateSiteV3(siteCreateRequest, siteId, notification)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(Pojo)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameUpdateSiteV3),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("siteId", siteId),
			clienttest.EndpointRequestHasBodyJsonValue("address.name", siteAddressName),
			clienttest.EndpointRequestHasBodyJsonValue("attributes.site-attribute-1", attributes["site-attribute-1"]),
			clienttest.EndpointRequestHasBodyJsonValue("contact.name", siteContactName),
			clienttest.EndpointRequestHasBodyJsonValue("description", siteDescription),
			clienttest.EndpointRequestHasBodyJsonValue("image", siteImage),
			clienttest.EndpointRequestHasBodyJsonValue("name", siteName),
			clienttest.EndpointRequestHasBodyJsonValue("parentId", parentId),
			clienttest.EndpointRequestHasBodyJsonValue("type", siteType),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPut),
			clienttest.ServiceEndpointHasPath("/api/v3/sites/{{.siteId}}"),
		).
		Test(t)
}

func TestIntegration_DeleteSiteV3(t *testing.T) {
	const siteId = "site-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.DeleteSiteV3(siteId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameDeleteSiteV3),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("siteId", siteId),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodDelete),
			clienttest.ServiceEndpointHasPath("/api/v3/sites/{{.siteId}}"),
		).
		Test(t)
}

func TestIntegration_AddDeviceToSiteV3(t *testing.T) {
	const siteId = "site-id"
	const deviceId = "device-id"
	const notification = "notification"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.AddDeviceToSiteV3(deviceId, siteId, notification)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameAddDevicetoSiteV3),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("siteId", siteId),
			clienttest.EndpointRequestHasEndpointParameter("deviceId", deviceId),
			clienttest.EndpointRequestHasQueryParam("notification", notification),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPut),
			clienttest.ServiceEndpointHasPath("/api/v3/sites/{{.siteId}}/devices/{{.deviceId}}"),
		).
		Test(t)
}

func TestIntegration_DeleteDeviceFromSiteV3(t *testing.T) {
	const siteId = "site-id"
	const deviceId = "device-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.DeleteDeviceFromSiteV3(deviceId, siteId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameDeleteDeviceFromSiteV3),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("siteId", siteId),
			clienttest.EndpointRequestHasEndpointParameter("deviceId", deviceId),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodDelete),
			clienttest.ServiceEndpointHasPath("/api/v3/sites/{{.siteId}}/devices/{{.deviceId}}"),
		).
		Test(t)
}

func TestIntegration_UpdateSiteStatusV3(t *testing.T) {
	const siteId = "site-id"
	const value = "value"
	const severity = "severity"
	const lastUpdatedMessage = "lastUpdatedMessage"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			updateRequest := SiteStatusUpdateRequest{
				LastUpdatedMessage: lastUpdatedMessage,
				Severity:           severity,
				Value:              value,
			}
			return api.UpdateSiteStatusV3(updateRequest, siteId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameUpdateSiteStatusV3),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("siteId", siteId),
			clienttest.EndpointRequestHasBodyJsonValue("lastUpdatedMessage", lastUpdatedMessage),
			clienttest.EndpointRequestHasBodyJsonValue("severity", severity),
			clienttest.EndpointRequestHasBodyJsonValue("value", value),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPut),
			clienttest.ServiceEndpointHasPath("/api/v3/sites/{{.siteId}}/status"),
		).
		Test(t)
}

func TestIntegration_GetDevice(t *testing.T) {
	const deviceInstanceId = "device-instance-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetDevice(deviceInstanceId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(Pojo)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetDevice),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("deviceInstanceId", deviceInstanceId),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v1/devices/{{.deviceInstanceId}}"),
		).
		Test(t)
}

func TestIntegration_GetDevices(t *testing.T) {
	const deviceInstanceId = "device-instance-id"
	const subscriptionId = "subscription-id"
	const serialKey = "serial-key"
	const tenantId = "tenant-id"
	const page = 0
	const pageSize = 100

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetDevices(
				types.NewOptionalStringFromString(deviceInstanceId).Ptr(),
				types.NewOptionalStringFromString(subscriptionId).Ptr(),
				types.NewOptionalStringFromString(serialKey).Ptr(),
				types.NewOptionalStringFromString(tenantId).Ptr(),
				page,
				pageSize)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(Pojo)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetDevices),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasQueryParam("deviceID", deviceInstanceId),
			clienttest.EndpointRequestHasQueryParam("serialKey", serialKey),
			clienttest.EndpointRequestHasQueryParam("subscriptionID", subscriptionId),
			clienttest.EndpointRequestHasQueryParam("tenantId", tenantId),
			clienttest.EndpointRequestHasQueryParam("page", strconv.Itoa(page)),
			clienttest.EndpointRequestHasQueryParam("pageSize", strconv.Itoa(pageSize)),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v2/devices"),
		).
		Test(t)
}

func TestIntegration_CreateDevice(t *testing.T) {
	const deviceInstanceId = "device-instance-id"
	const subscriptionId = "subscription-id"
	const serialKey = "serial-key"
	const tenantId = "tenant-id"

	var attributes = map[string]string{
		"serialKey": serialKey,
		"tenantId":  tenantId,
	}
	var defAttributes = map[string]string{
		"device-def-attribute-1": "device-value-1",
	}
	var status = map[string]string{
		"device-status-1": "device-value-2",
	}

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.CreateDevice(
				subscriptionId,
				types.NewOptionalStringFromString(deviceInstanceId).Ptr(),
				attributes,
				defAttributes,
				status)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(Pojo)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameCreateDevice),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("subscriptionId", subscriptionId),
			clienttest.EndpointRequestHasBodyJsonValue("deviceInstanceId", deviceInstanceId),
			clienttest.EndpointRequestHasBodyJsonValue("deviceAttribute.serialKey", serialKey),
			clienttest.EndpointRequestHasBodyJsonValue("deviceAttribute.tenantId", tenantId),
			clienttest.EndpointRequestHasBodyJsonValue("deviceDefAttribute.device-def-attribute-1", defAttributes["device-def-attribute-1"]),
			clienttest.EndpointRequestHasBodyJsonValue("status.device-status-1", status["device-status-1"]),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPost),
			clienttest.ServiceEndpointHasPath("/api/v1/devices/subscriptions/{{.subscriptionId}}"),
		).
		Test(t)
}

func TestIntegration_UpdateDevice(t *testing.T) {
	const deviceInstanceId = "device-instance-id"
	const serialKey = "serial-key"
	const tenantId = "tenant-id"

	var attributes = map[string]string{
		"serialKey": serialKey,
		"tenantId":  tenantId,
	}
	var defAttributes = map[string]string{
		"device-def-attribute-1": "device-value-1",
	}
	var status = map[string]string{
		"device-status-1": "device-value-2",
	}

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.UpdateDevice(
				deviceInstanceId,
				attributes,
				defAttributes,
				status)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(Pojo)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameUpdateDevice),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("deviceInstanceId", deviceInstanceId),
			clienttest.EndpointRequestHasBodyJsonValue("deviceAttribute.serialKey", serialKey),
			clienttest.EndpointRequestHasBodyJsonValue("deviceAttribute.tenantId", tenantId),
			clienttest.EndpointRequestHasBodyJsonValue("deviceDefAttribute.device-def-attribute-1", defAttributes["device-def-attribute-1"]),
			clienttest.EndpointRequestHasBodyJsonValue("status.device-status-1", status["device-status-1"]),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPut),
			clienttest.ServiceEndpointHasPath("/api/v1/devices/{{.deviceInstanceId}}"),
		).
		Test(t)
}

func TestIntegration_DeleteDevice(t *testing.T) {
	const deviceInstanceId = "device-instance-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.DeleteDevice(deviceInstanceId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameDeleteDevice),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("deviceInstanceId", deviceInstanceId),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodDelete),
			clienttest.ServiceEndpointHasPath("/api/v1/devices/{{.deviceInstanceId}}"),
		).
		Test(t)
}

func TestIntegration_CreateManagedDevice(t *testing.T) {
	const tenantId = "tenant-id"
	const deviceModel = "device-model"
	const deviceOnboardType = "device-onboard-type"
	var deviceOnboardInfo = map[string]string{
		"device-onboard-info-1": "device-value-1",
	}

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.CreateManagedDevice(
				tenantId,
				deviceModel,
				deviceOnboardType,
				deviceOnboardInfo)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(Pojo)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameCreateManagedDevice),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasBodyJsonValue("tenantId", tenantId),
			clienttest.EndpointRequestHasBodyJsonValue("deviceModel", deviceModel),
			clienttest.EndpointRequestHasBodyJsonValue("deviceOnboardingType", deviceOnboardType),
			clienttest.EndpointRequestHasBodyJsonValue("deviceOnboardInfo.device-onboard-info-1", deviceOnboardInfo["device-onboard-info-1"]),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPost),
			clienttest.ServiceEndpointHasPath("/api/v3/devices"),
		).
		Test(t)
}

func TestIntegration_DeleteManagedDevice(t *testing.T) {
	const deviceInstanceId = "device-instance-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.DeleteManagedDevice(deviceInstanceId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameDeleteManagedDevice),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("deviceInstanceId", deviceInstanceId),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodDelete),
			clienttest.ServiceEndpointHasPath("/api/v3/devices/{{.deviceInstanceId}}"),
		).
		Test(t)
}

func TestIntegration_GetDeviceConfig(t *testing.T) {
	const deviceInstanceId = "device-instance-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetDeviceConfig(deviceInstanceId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetDeviceConfig),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("deviceInstanceId", deviceInstanceId),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v3/devices/{{.deviceInstanceId}}/config"),
		).
		Test(t)
}

func TestIntegration_GetDeviceV4(t *testing.T) {
	const deviceId = "device-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetDeviceV4(deviceId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetDeviceV4),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("deviceId", deviceId),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v4/devices/{{.deviceId}}"),
		).
		Test(t)
}

func TestIntegration_GetDevicesV4(t *testing.T) {
	const page = 0
	const pageSize = 100
	var query = map[string][]string{
		"query-param-1": {"query-arg-1"},
	}

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetDevicesV4(query, page, pageSize)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetDevicesV4),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasQueryParam("query-param-1", "query-arg-1"),
			clienttest.EndpointRequestHasQueryParam("page", strconv.Itoa(page)),
			clienttest.EndpointRequestHasQueryParam("pageSize", strconv.Itoa(pageSize)),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v4/devices"),
		).
		Test(t)
}

func TestIntegration_CreateDeviceV4(t *testing.T) {
	const deviceName = "device-name"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			createRequest := DeviceCreateRequest{
				Name: deviceName,
			}
			return api.CreateDeviceV4(createRequest)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameCreateDeviceV4),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasBodyJsonValue("name", deviceName),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPost),
			clienttest.ServiceEndpointHasPath("/api/v4/devices"),
		).
		Test(t)
}

func TestIntegration_UpdateDeviceV4(t *testing.T) {
	const deviceId = "device-id"
	const deviceName = "device-name"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			updateRequest := DeviceUpdateRequest{
				Name: deviceName,
			}
			return api.UpdateDeviceV4(updateRequest, deviceId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameUpdateDeviceV4),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("deviceId", deviceId),
			clienttest.EndpointRequestHasBodyJsonValue("name", deviceName),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPut),
			clienttest.ServiceEndpointHasPath("/api/v4/devices/{{.deviceId}}"),
		).
		Test(t)
}

func TestIntegration_UpdateDeviceStatusV4(t *testing.T) {
	const deviceId = "device-id"
	const status = "status"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			updateRequest := DeviceStatusUpdateRequest{Value: status}
			return api.UpdateDeviceStatusV4(updateRequest, deviceId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameUpdateDeviceStatusV4),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasBodyJsonValue("value", status),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPut),
			clienttest.ServiceEndpointHasPath("/api/v4/devices/{{.deviceId}}/status"),
		).
		Test(t)
}

func TestIntegration_DeleteDeviceV4(t *testing.T) {
	const deviceId = "device-id"
	const force = "force"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.DeleteDeviceV4(deviceId, force)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameDeleteDeviceV4),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("deviceId", deviceId),
			clienttest.EndpointRequestHasQueryParam("force", force),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodDelete),
			clienttest.ServiceEndpointHasPath("/api/v4/devices/{{.deviceId}}"),
		).
		Test(t)
}

func TestIntegration_CreateDeviceActions(t *testing.T) {
	const actionType = "action-type"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			deviceActions := DeviceActionCreateRequests{
				DeviceActionCreateRequest{
					ActionType: actionType,
				},
			}
			return api.CreateDeviceActions(deviceActions)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameCreateDeviceActions),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasBodyJsonValue("0.actionType", actionType),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPost),
			clienttest.ServiceEndpointHasPath("/api/v1/deviceActions"),
		).
		Test(t)
}

func TestIntegration_GetDeviceTemplateHistory(t *testing.T) {
	const deviceInstanceId = "device-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetDeviceTemplateHistory(deviceInstanceId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetDeviceTemplateHistory),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("deviceInstanceId", deviceInstanceId),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v3/devices/{{.deviceInstanceId}}/templates"),
		).
		Test(t)
}

func TestIntegration_AttachDeviceTemplates(t *testing.T) {
	const deviceInstanceId = "device-id"
	const templateId = "template-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			attachRequest := AttachTemplateRequest{
				TemplateDetails: []TemplateDetails{
					{
						TemplateID:     templateId,
						TemplateParams: []TemplateParams{},
					},
				},
			}
			return api.AttachDeviceTemplates(deviceInstanceId, attachRequest)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameAttachDeviceTemplates),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("deviceInstanceId", deviceInstanceId),
			clienttest.EndpointRequestHasBodyJsonValue("templateDetails.0.templateId", templateId),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPost),
			clienttest.ServiceEndpointHasPath("/api/v3/devices/{{.deviceInstanceId}}/templates"),
		).
		Test(t)
}

func TestIntegration_UpdateTemplateAccess(t *testing.T) {
	const templateId = "template-id"
	const tenantId = "tenant-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			accessRequest := DeviceTemplateAccess{
				Global:  false,
				Tenants: []string{
					tenantId,
				},
			}
			return api.UpdateTemplateAccess(templateId, accessRequest)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameUpdateTemplateAccess),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("templateId", templateId),
			clienttest.EndpointRequestHasBodyJsonValue("tenants.0", tenantId),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPut),
			clienttest.ServiceEndpointHasPath("/api/v1/devicetemplates/{{.templateId}}"),
		).
		Test(t)
}

func TestIntegration_AddDeviceTemplate(t *testing.T) {
	const configContent = "config-content"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			createRequest := DeviceTemplateCreateRequest{
				ConfigContent:        configContent,
			}
			return api.AddDeviceTemplate(createRequest)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameSetDeviceTemplate),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasBodyJsonValue("configContent", configContent),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPost),
			clienttest.ServiceEndpointHasPath("/api/v1/devicetemplates"),
		).
		Test(t)

}

func TestIntegration_GetAllControlPlanes(t *testing.T) {
	const tenantId = "tenant-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetAllControlPlanes(
				types.NewOptionalStringFromString(tenantId).Ptr())
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(PojoArray)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetAllControlPlanes),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasQueryParam("tenantId", tenantId),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v1/controlplanes"),
		).
		Test(t)
}

func TestIntegration_GetControlPlane(t *testing.T) {
	const controlPlaneId = "control-plane-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetControlPlane(controlPlaneId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetControlPlane),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("controlPlaneId", controlPlaneId),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v1/controlplanes/{{.controlPlaneId}}"),
		).
		Test(t)
}

func TestIntegration_CreateControlPlane(t *testing.T) {
	const tenantId = "tenant-id"
	const controlPlaneName = "control-plane-name"
	const controlPlaneUrl = "control-plane-url"
	const resourceProvider = "resource-provider"
	const authenticationType = "authentication-type"
	const tlsInsecure = true
	var attributes = map[string]string{
		"control-plane-attribute-1": "control-plane-value-1",
	}

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.CreateControlPlane(
				tenantId,
				controlPlaneName,
				controlPlaneUrl,
				resourceProvider,
				authenticationType,
				tlsInsecure,
				attributes)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameCreateControlPlane),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasBodyJsonValue("tenantId", tenantId),
			clienttest.EndpointRequestHasBodyJsonValue("name", controlPlaneName),
			clienttest.EndpointRequestHasBodyJsonValue("url", controlPlaneUrl),
			clienttest.EndpointRequestHasBodyJsonValue("resourceProvider", resourceProvider),
			clienttest.EndpointRequestHasBodyJsonValue("authenticationType", authenticationType),
			clienttest.EndpointRequestHasBodyJsonValue("tlsInsecure", tlsInsecure),
			clienttest.EndpointRequestHasBodyJsonValue("attributes.control-plane-attribute-1", attributes["control-plane-attribute-1"]),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPost),
			clienttest.ServiceEndpointHasPath("/api/v1/controlplanes"),
		).
		Test(t)
}

func TestIntegration_UpdateControlPlane(t *testing.T) {
	const controlPlaneId = "control-plane-id"
	const tenantId = "tenant-id"
	const controlPlaneName = "control-plane-name"
	const controlPlaneUrl = "control-plane-url"
	const resourceProvider = "resource-provider"
	const authenticationType = "authentication-type"
	const tlsInsecure = true
	var attributes = map[string]string{
		"control-plane-attribute-1": "control-plane-value-1",
	}

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.UpdateControlPlane(
				controlPlaneId,
				tenantId,
				controlPlaneName,
				controlPlaneUrl,
				resourceProvider,
				authenticationType,
				tlsInsecure,
				attributes)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameUpdateControlPlane),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("controlPlaneId", controlPlaneId),
			clienttest.EndpointRequestHasBodyJsonValue("tenantId", tenantId),
			clienttest.EndpointRequestHasBodyJsonValue("name", controlPlaneName),
			clienttest.EndpointRequestHasBodyJsonValue("url", controlPlaneUrl),
			clienttest.EndpointRequestHasBodyJsonValue("resourceProvider", resourceProvider),
			clienttest.EndpointRequestHasBodyJsonValue("authenticationType", authenticationType),
			clienttest.EndpointRequestHasBodyJsonValue("tlsInsecure", tlsInsecure),
			clienttest.EndpointRequestHasBodyJsonValue("attributes.control-plane-attribute-1", attributes["control-plane-attribute-1"]),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPut),
			clienttest.ServiceEndpointHasPath("/api/v1/controlplanes/{{.controlPlaneId}}"),
		).
		Test(t)
}

func TestIntegration_DeleteControlPlane(t *testing.T) {
	const controlPlaneId = "control-plane-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.DeleteControlPlane(controlPlaneId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameDeleteControlPlane),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("controlPlaneId", controlPlaneId),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodDelete),
			clienttest.ServiceEndpointHasPath("/api/v1/controlplanes/{{.controlPlaneId}}"),
		).
		Test(t)
}

func TestIntegration_ConnectControlPlane(t *testing.T) {
	const controlPlaneId = "control-plane-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.ConnectControlPlane(controlPlaneId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameConnectControlPlane),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("controlPlaneId", controlPlaneId),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPost),
			clienttest.ServiceEndpointHasPath("/api/v1/controlplanes/{{.controlPlaneId}}/connect"),
		).
		Test(t)
}

func TestIntegration_ConnectUnmanagedControlPlane(t *testing.T) {
	const username = "username"
	const password = "password"
	const controlPlaneUrl = "control-plane-url"
	const resourceProvider = "resource-provider"
	const tlsInsecure = true

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.ConnectUnmanagedControlPlane(
				username,
				password,
				controlPlaneUrl,
				resourceProvider,
				tlsInsecure)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameConnectUnmanagedControlPlane),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasBodyJsonValue("username", username),
			clienttest.EndpointRequestHasBodyJsonValue("password", password),
			clienttest.EndpointRequestHasBodyJsonValue("url", controlPlaneUrl),
			clienttest.EndpointRequestHasBodyJsonValue("resourceProvider", resourceProvider),
			clienttest.EndpointRequestHasBodyJsonValue("tlsInsecure", tlsInsecure),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPost),
			clienttest.ServiceEndpointHasPath("/api/v1/controlplanes/connect"),
		).
		Test(t)
}

func TestIntegration_GetEntityShard(t *testing.T) {
	const entityId = "entity-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetEntityShard(entityId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetEntityShard),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("entityId", entityId),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v2/shardmanagers/entity/{{.entityId}}"),
		).
		Test(t)
}

func TestIntegration_CreateDeviceConnection(t *testing.T) {
	const deviceInstanceId = "device-instance-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			createRequest := DeviceConnectionCreateRequest{
				DeviceInstanceId:  deviceInstanceId,
			}
			response, _, err := api.CreateDeviceConnection(createRequest)
			return response, err
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(DeviceConnectionResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameCreateDeviceConnection),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasBodyJsonValue("deviceInstanceId", deviceInstanceId),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPost),
			clienttest.ServiceEndpointHasPath("/api/v2/devices/connections"),
		).
		Test(t)
}


func TestIntegration_DeleteDeviceConnection(t *testing.T) {
	const deviceConnectionId = "device-connection-id"

	NewManageIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.DeleteDeviceConnection(deviceConnectionId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameDeleteDeviceConnection),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasEndpointParameter("deviceConnectionId", deviceConnectionId),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodDelete),
			clienttest.ServiceEndpointHasPath("/api/v2/devices/connections/{{.deviceConnectionId}}"),
		).
		Test(t)
}