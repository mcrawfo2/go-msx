package platform

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	platform "cto-github.cisco.com/NFV-BU/msx-platform-go-client"
	"net/http"
)

type WorkflowsApi interface {
	DeleteWorkflow(ctx context.Context, id string) (*http.Response, error)
	ExportWorkflow(ctx context.Context, id string) (string, *http.Response, error)
	GetWorkflow(ctx context.Context, id string) (platform.Workflow, *http.Response, error)
	GetWorkflowStartConfig(ctx context.Context, id string) (platform.WorkflowStartConfig, *http.Response, error)
	GetWorkflowsList(ctx context.Context, localVarOptionals *platform.GetWorkflowsListOpts) ([]platform.Workflow, *http.Response, error)
	ImportWorkflow(ctx context.Context, requestBody map[string]map[string]interface{}, localVarOptionals *platform.ImportWorkflowOpts) (platform.WorkflowMapping, *http.Response, error)
	StartWorkflow(ctx context.Context, id string, workflowStartConfig platform.WorkflowStartConfig, localVarOptionals *platform.StartWorkflowOpts) ([]platform.StartWorkflowResponse, *http.Response, error)
	UpdateWorkflow(ctx context.Context, id string, requestBody map[string]map[string]interface{}, localVarOptionals *platform.UpdateWorkflowOpts) (platform.WorkflowMapping, *http.Response, error)
	ValidateWorkflow(ctx context.Context, id string) (platform.ValidateWorkflowResponse, *http.Response, error)
}

func NewWorkflowsApiService(ctx context.Context) *platform.WorkflowsApiService {
	cfg := newPlatformClientConfigFromContext(ctx, integration.ServiceNameWorkflow)
	return platform.NewAPIClient(cfg).WorkflowsApi
}
