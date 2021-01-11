// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	context "context"

	paging "cto-github.cisco.com/NFV-BU/go-msx/paging"
	mock "github.com/stretchr/testify/mock"
)

// CrudRepositoryApi is an autogenerated mock type for the CrudRepositoryApi type
type CrudRepositoryApi struct {
	mock.Mock
}

// DeleteBy provides a mock function with given fields: ctx, where
func (_m *CrudRepositoryApi) DeleteBy(ctx context.Context, where map[string]interface{}) error {
	ret := _m.Called(ctx, where)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}) error); ok {
		r0 = rf(ctx, where)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAll provides a mock function with given fields: ctx, dest
func (_m *CrudRepositoryApi) FindAll(ctx context.Context, dest interface{}) error {
	ret := _m.Called(ctx, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = rf(ctx, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAllBy provides a mock function with given fields: ctx, where, dest
func (_m *CrudRepositoryApi) FindAllBy(ctx context.Context, where map[string]interface{}, dest interface{}) error {
	ret := _m.Called(ctx, where, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, interface{}) error); ok {
		r0 = rf(ctx, where, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAllPagedBy provides a mock function with given fields: ctx, where, preq, dest
func (_m *CrudRepositoryApi) FindAllPagedBy(ctx context.Context, where map[string]interface{}, preq paging.Request, dest interface{}) (paging.Response, error) {
	ret := _m.Called(ctx, where, preq, dest)

	var r0 paging.Response
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, paging.Request, interface{}) paging.Response); ok {
		r0 = rf(ctx, where, preq, dest)
	} else {
		r0 = ret.Get(0).(paging.Response)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, map[string]interface{}, paging.Request, interface{}) error); ok {
		r1 = rf(ctx, where, preq, dest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAllSortedBy provides a mock function with given fields: ctx, where, sortOrder, dest
func (_m *CrudRepositoryApi) FindAllSortedBy(ctx context.Context, where map[string]interface{}, sortOrder paging.SortOrder, dest interface{}) error {
	ret := _m.Called(ctx, where, sortOrder, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, paging.SortOrder, interface{}) error); ok {
		r0 = rf(ctx, where, sortOrder, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindOneBy provides a mock function with given fields: ctx, where, dest
func (_m *CrudRepositoryApi) FindOneBy(ctx context.Context, where map[string]interface{}, dest interface{}) error {
	ret := _m.Called(ctx, where, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, interface{}) error); ok {
		r0 = rf(ctx, where, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindOneSortedBy provides a mock function with given fields: ctx, where, sortOrder, dest
func (_m *CrudRepositoryApi) FindOneSortedBy(ctx context.Context, where map[string]interface{}, sortOrder paging.SortOrder, dest interface{}) error {
	ret := _m.Called(ctx, where, sortOrder, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, paging.SortOrder, interface{}) error); ok {
		r0 = rf(ctx, where, sortOrder, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Insert provides a mock function with given fields: ctx, value
func (_m *CrudRepositoryApi) Insert(ctx context.Context, value interface{}) error {
	ret := _m.Called(ctx, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = rf(ctx, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Save provides a mock function with given fields: ctx, value
func (_m *CrudRepositoryApi) Save(ctx context.Context, value interface{}) error {
	ret := _m.Called(ctx, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = rf(ctx, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveAll provides a mock function with given fields: ctx, values
func (_m *CrudRepositoryApi) SaveAll(ctx context.Context, values []interface{}) error {
	ret := _m.Called(ctx, values)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []interface{}) error); ok {
		r0 = rf(ctx, values)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Truncate provides a mock function with given fields: ctx
func (_m *CrudRepositoryApi) Truncate(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, where, value
func (_m *CrudRepositoryApi) Update(ctx context.Context, where map[string]interface{}, value interface{}) error {
	ret := _m.Called(ctx, where, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, interface{}) error); ok {
		r0 = rf(ctx, where, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
