package platform

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	platform "cto-github.cisco.com/NFV-BU/msx-platform-go-client"
	"net/http"
)

type WorkflowInstancesApi interface {
	CancelWorkflowInstance(ctx context.Context, id string) (platform.WorkflowInstance, *http.Response, error)
	DeleteWorkflowInstance(ctx context.Context, id string) (platform.WorkflowInstanceDeleteResponse, *http.Response, error)
	GetWorkflowInstance(ctx context.Context, id string) (platform.WorkflowInstance, *http.Response, error)
	GetWorkflowInstanceAction(ctx context.Context, id string, actionId string) (platform.WorkflowAction, *http.Response, error)
	GetWorkflowInstancesList(ctx context.Context, id string, page int32, pageSize int32, localVarOptionals *platform.GetWorkflowInstancesListOpts) ([]platform.WorkflowInstance, *http.Response, error)
}

func NewWorkflowInstancesApiService(ctx context.Context) *platform.WorkflowInstancesApiService {
	cfg := newPlatformClientConfigFromContext(ctx, integration.ServiceNameWorkflow)
	return platform.NewAPIClient(cfg).WorkflowInstancesApi
}
