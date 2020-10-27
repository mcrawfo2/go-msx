// Code generated by mockery v2.3.0. DO NOT EDIT.

package platform

import (
	context "context"
	http "net/http"

	mock "github.com/stretchr/testify/mock"

	msx_platform_client "cto-github.cisco.com/NFV-BU/msx-platform-go-client"
)

// MockWorkflowTargetsApi is an autogenerated mock type for the WorkflowTargetsApi type
type MockWorkflowTargetsApi struct {
	mock.Mock
}

// CreateWorkflowTarget provides a mock function with given fields: ctx, workflowTargetCreate
func (_m *MockWorkflowTargetsApi) CreateWorkflowTarget(ctx context.Context, workflowTargetCreate msx_platform_client.WorkflowTargetCreate) (msx_platform_client.WorkflowTarget, *http.Response, error) {
	ret := _m.Called(ctx, workflowTargetCreate)

	var r0 msx_platform_client.WorkflowTarget
	if rf, ok := ret.Get(0).(func(context.Context, msx_platform_client.WorkflowTargetCreate) msx_platform_client.WorkflowTarget); ok {
		r0 = rf(ctx, workflowTargetCreate)
	} else {
		r0 = ret.Get(0).(msx_platform_client.WorkflowTarget)
	}

	var r1 *http.Response
	if rf, ok := ret.Get(1).(func(context.Context, msx_platform_client.WorkflowTargetCreate) *http.Response); ok {
		r1 = rf(ctx, workflowTargetCreate)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*http.Response)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, msx_platform_client.WorkflowTargetCreate) error); ok {
		r2 = rf(ctx, workflowTargetCreate)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// DeleteWorkflowTarget provides a mock function with given fields: ctx, id
func (_m *MockWorkflowTargetsApi) DeleteWorkflowTarget(ctx context.Context, id string) (*http.Response, error) {
	ret := _m.Called(ctx, id)

	var r0 *http.Response
	if rf, ok := ret.Get(0).(func(context.Context, string) *http.Response); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWorkflowTarget provides a mock function with given fields: ctx, id
func (_m *MockWorkflowTargetsApi) GetWorkflowTarget(ctx context.Context, id string) (msx_platform_client.WorkflowTarget, *http.Response, error) {
	ret := _m.Called(ctx, id)

	var r0 msx_platform_client.WorkflowTarget
	if rf, ok := ret.Get(0).(func(context.Context, string) msx_platform_client.WorkflowTarget); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(msx_platform_client.WorkflowTarget)
	}

	var r1 *http.Response
	if rf, ok := ret.Get(1).(func(context.Context, string) *http.Response); ok {
		r1 = rf(ctx, id)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*http.Response)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, id)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetWorkflowTargetsList provides a mock function with given fields: ctx, localVarOptionals
func (_m *MockWorkflowTargetsApi) GetWorkflowTargetsList(ctx context.Context, localVarOptionals *msx_platform_client.GetWorkflowTargetsListOpts) ([]msx_platform_client.WorkflowTarget, *http.Response, error) {
	ret := _m.Called(ctx, localVarOptionals)

	var r0 []msx_platform_client.WorkflowTarget
	if rf, ok := ret.Get(0).(func(context.Context, *msx_platform_client.GetWorkflowTargetsListOpts) []msx_platform_client.WorkflowTarget); ok {
		r0 = rf(ctx, localVarOptionals)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]msx_platform_client.WorkflowTarget)
		}
	}

	var r1 *http.Response
	if rf, ok := ret.Get(1).(func(context.Context, *msx_platform_client.GetWorkflowTargetsListOpts) *http.Response); ok {
		r1 = rf(ctx, localVarOptionals)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*http.Response)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, *msx_platform_client.GetWorkflowTargetsListOpts) error); ok {
		r2 = rf(ctx, localVarOptionals)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// UpdateWorkflowTarget provides a mock function with given fields: ctx, id, workflowTargetUpdate
func (_m *MockWorkflowTargetsApi) UpdateWorkflowTarget(ctx context.Context, id string, workflowTargetUpdate msx_platform_client.WorkflowTargetUpdate) (msx_platform_client.WorkflowTarget, *http.Response, error) {
	ret := _m.Called(ctx, id, workflowTargetUpdate)

	var r0 msx_platform_client.WorkflowTarget
	if rf, ok := ret.Get(0).(func(context.Context, string, msx_platform_client.WorkflowTargetUpdate) msx_platform_client.WorkflowTarget); ok {
		r0 = rf(ctx, id, workflowTargetUpdate)
	} else {
		r0 = ret.Get(0).(msx_platform_client.WorkflowTarget)
	}

	var r1 *http.Response
	if rf, ok := ret.Get(1).(func(context.Context, string, msx_platform_client.WorkflowTargetUpdate) *http.Response); ok {
		r1 = rf(ctx, id, workflowTargetUpdate)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*http.Response)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, msx_platform_client.WorkflowTargetUpdate) error); ok {
		r2 = rf(ctx, id, workflowTargetUpdate)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}