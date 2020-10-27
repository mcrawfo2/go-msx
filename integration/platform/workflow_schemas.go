package platform

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	platform "cto-github.cisco.com/NFV-BU/msx-platform-go-client"
	"net/http"
)

type WorkflowSchemasApi interface {
	GetWorkflowSchema(ctx context.Context, id string, localVarOptionals *platform.GetWorkflowSchemaOpts) (platform.WorkflowSchemaByTypeResponse, *http.Response, error)
	GetWorkflowSchemasList(ctx context.Context, baseType string, localVarOptionals *platform.GetWorkflowSchemasListOpts) ([]platform.WorkflowSchema, *http.Response, error)
}

func NewWorkflowSchemasApiService(ctx context.Context) *platform.WorkflowSchemasApiService {
	cfg := newPlatformClientConfigFromContext(ctx, integration.ServiceNameWorkflow)
	return platform.NewAPIClient(cfg).WorkflowSchemasApi
}
