package platform

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	platform "cto-github.cisco.com/NFV-BU/msx-platform-go-client"
	"net/http"
)

type WorkflowTargetsApi interface {
	CreateWorkflowTarget(ctx context.Context, workflowTargetCreate platform.WorkflowTargetCreate) (platform.WorkflowTarget, *http.Response, error)
	DeleteWorkflowTarget(ctx context.Context, id string) (*http.Response, error)
	GetWorkflowTarget(ctx context.Context, id string) (platform.WorkflowTarget, *http.Response, error)
	GetWorkflowTargetsList(ctx context.Context, localVarOptionals *platform.GetWorkflowTargetsListOpts) ([]platform.WorkflowTarget, *http.Response, error)
	UpdateWorkflowTarget(ctx context.Context, id string, workflowTargetUpdate platform.WorkflowTargetUpdate) (platform.WorkflowTarget, *http.Response, error)
}

func NewWorkflowTargetsApiService(ctx context.Context) *platform.WorkflowTargetsApiService {
	cfg := newPlatformClientConfigFromContext(ctx, integration.ServiceNameWorkflow)
	return platform.NewAPIClient(cfg).WorkflowTargetsApi
}
