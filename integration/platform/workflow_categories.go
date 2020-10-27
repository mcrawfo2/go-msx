package platform

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	platform "cto-github.cisco.com/NFV-BU/msx-platform-go-client"
	"net/http"
)

type WorkflowCategoriesApi interface {
	CreateWorkflowCategory(ctx context.Context, workflowCategoryCreate platform.WorkflowCategoryCreate, localVarOptionals *platform.CreateWorkflowCategoryOpts) (platform.WorkflowCategory, *http.Response, error)
	DeleteWorkflowCategory(ctx context.Context, id string) (*http.Response, error)
	GetWorkflowCategoriesList(ctx context.Context, localVarOptionals *platform.GetWorkflowCategoriesListOpts) ([]platform.WorkflowCategory, *http.Response, error)
	GetWorkflowCategory(ctx context.Context, id string) (platform.WorkflowCategory, *http.Response, error)
	UpdateWorkflowCategory(ctx context.Context, id string, workflowCategoryUpdate platform.WorkflowCategoryUpdate) (platform.WorkflowCategory, *http.Response, error)
}

func NewWorkflowCategoriesApiService(ctx context.Context) *platform.WorkflowCategoriesApiService {
	cfg := newPlatformClientConfigFromContext(ctx, integration.ServiceNameWorkflow)
	return platform.NewAPIClient(cfg).WorkflowCategoriesApi
}
