// Code generated by mockery v2.22.1. DO NOT EDIT.

package sqldb

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockSqlRepositoryApi is an autogenerated mock type for the SqlRepositoryApi type
type MockSqlRepositoryApi struct {
	mock.Mock
}

type MockSqlRepositoryApi_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSqlRepositoryApi) EXPECT() *MockSqlRepositoryApi_Expecter {
	return &MockSqlRepositoryApi_Expecter{mock: &_m.Mock}
}

// SqlExecute provides a mock function with given fields: ctx, stmt, args
func (_m *MockSqlRepositoryApi) SqlExecute(ctx context.Context, stmt string, args []interface{}) error {
	ret := _m.Called(ctx, stmt, args)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []interface{}) error); ok {
		r0 = rf(ctx, stmt, args)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSqlRepositoryApi_SqlExecute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SqlExecute'
type MockSqlRepositoryApi_SqlExecute_Call struct {
	*mock.Call
}

// SqlExecute is a helper method to define mock.On call
//   - ctx context.Context
//   - stmt string
//   - args []interface{}
func (_e *MockSqlRepositoryApi_Expecter) SqlExecute(ctx interface{}, stmt interface{}, args interface{}) *MockSqlRepositoryApi_SqlExecute_Call {
	return &MockSqlRepositoryApi_SqlExecute_Call{Call: _e.mock.On("SqlExecute", ctx, stmt, args)}
}

func (_c *MockSqlRepositoryApi_SqlExecute_Call) Run(run func(ctx context.Context, stmt string, args []interface{})) *MockSqlRepositoryApi_SqlExecute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].([]interface{}))
	})
	return _c
}

func (_c *MockSqlRepositoryApi_SqlExecute_Call) Return(_a0 error) *MockSqlRepositoryApi_SqlExecute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSqlRepositoryApi_SqlExecute_Call) RunAndReturn(run func(context.Context, string, []interface{}) error) *MockSqlRepositoryApi_SqlExecute_Call {
	_c.Call.Return(run)
	return _c
}

// SqlGet provides a mock function with given fields: ctx, stmt, args, dest
func (_m *MockSqlRepositoryApi) SqlGet(ctx context.Context, stmt string, args []interface{}, dest interface{}) error {
	ret := _m.Called(ctx, stmt, args, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []interface{}, interface{}) error); ok {
		r0 = rf(ctx, stmt, args, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSqlRepositoryApi_SqlGet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SqlGet'
type MockSqlRepositoryApi_SqlGet_Call struct {
	*mock.Call
}

// SqlGet is a helper method to define mock.On call
//   - ctx context.Context
//   - stmt string
//   - args []interface{}
//   - dest interface{}
func (_e *MockSqlRepositoryApi_Expecter) SqlGet(ctx interface{}, stmt interface{}, args interface{}, dest interface{}) *MockSqlRepositoryApi_SqlGet_Call {
	return &MockSqlRepositoryApi_SqlGet_Call{Call: _e.mock.On("SqlGet", ctx, stmt, args, dest)}
}

func (_c *MockSqlRepositoryApi_SqlGet_Call) Run(run func(ctx context.Context, stmt string, args []interface{}, dest interface{})) *MockSqlRepositoryApi_SqlGet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].([]interface{}), args[3].(interface{}))
	})
	return _c
}

func (_c *MockSqlRepositoryApi_SqlGet_Call) Return(_a0 error) *MockSqlRepositoryApi_SqlGet_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSqlRepositoryApi_SqlGet_Call) RunAndReturn(run func(context.Context, string, []interface{}, interface{}) error) *MockSqlRepositoryApi_SqlGet_Call {
	_c.Call.Return(run)
	return _c
}

// SqlSelect provides a mock function with given fields: ctx, stmt, args, dest
func (_m *MockSqlRepositoryApi) SqlSelect(ctx context.Context, stmt string, args []interface{}, dest interface{}) error {
	ret := _m.Called(ctx, stmt, args, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []interface{}, interface{}) error); ok {
		r0 = rf(ctx, stmt, args, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSqlRepositoryApi_SqlSelect_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SqlSelect'
type MockSqlRepositoryApi_SqlSelect_Call struct {
	*mock.Call
}

// SqlSelect is a helper method to define mock.On call
//   - ctx context.Context
//   - stmt string
//   - args []interface{}
//   - dest interface{}
func (_e *MockSqlRepositoryApi_Expecter) SqlSelect(ctx interface{}, stmt interface{}, args interface{}, dest interface{}) *MockSqlRepositoryApi_SqlSelect_Call {
	return &MockSqlRepositoryApi_SqlSelect_Call{Call: _e.mock.On("SqlSelect", ctx, stmt, args, dest)}
}

func (_c *MockSqlRepositoryApi_SqlSelect_Call) Run(run func(ctx context.Context, stmt string, args []interface{}, dest interface{})) *MockSqlRepositoryApi_SqlSelect_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].([]interface{}), args[3].(interface{}))
	})
	return _c
}

func (_c *MockSqlRepositoryApi_SqlSelect_Call) Return(_a0 error) *MockSqlRepositoryApi_SqlSelect_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSqlRepositoryApi_SqlSelect_Call) RunAndReturn(run func(context.Context, string, []interface{}, interface{}) error) *MockSqlRepositoryApi_SqlSelect_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockSqlRepositoryApi interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockSqlRepositoryApi creates a new instance of MockSqlRepositoryApi. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockSqlRepositoryApi(t mockConstructorTestingTNewMockSqlRepositoryApi) *MockSqlRepositoryApi {
	mock := &MockSqlRepositoryApi{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
