package platform

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	platform "cto-github.cisco.com/NFV-BU/msx-platform-go-client"
	"net/http"
)

type WorkflowEventsApi interface {
	CreateWorkflowEvent(ctx context.Context, workflowEventCreate platform.WorkflowEventCreate) (platform.WorkflowEvent, *http.Response, error)
	DeleteWorkflowEvent(ctx context.Context, id string) (*http.Response, error)
	GetWorkflowEvent(ctx context.Context, id string) (platform.WorkflowEvent, *http.Response, error)
	GetWorkflowEventsList(ctx context.Context) ([]platform.WorkflowEvent, *http.Response, error)
	UpdateWorkflowEvent(ctx context.Context, id string, workflowEventUpdate platform.WorkflowEventUpdate) (platform.WorkflowEvent, *http.Response, error)
}

func NewWorkflowEventsApiService(ctx context.Context) *platform.WorkflowEventsApiService {
	cfg := newPlatformClientConfigFromContext(ctx, integration.ServiceNameWorkflow)
	return platform.NewAPIClient(cfg).WorkflowEventsApi
}
