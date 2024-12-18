// Code generated by mockery. DO NOT EDIT.

//go:build !release

package aggregation

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// mockCheckPointerModel is an autogenerated mock type for the checkPointerModel type
type mockCheckPointerModel struct {
	mock.Mock
}

type mockCheckPointerModel_Expecter struct {
	mock *mock.Mock
}

func (_m *mockCheckPointerModel) EXPECT() *mockCheckPointerModel_Expecter {
	return &mockCheckPointerModel_Expecter{mock: &_m.Mock}
}

// readCounters provides a mock function with given fields: ctx, blobFileName
func (_m *mockCheckPointerModel) readCounters(ctx context.Context, blobFileName string) (map[string]int64, error) {
	ret := _m.Called(ctx, blobFileName)

	if len(ret) == 0 {
		panic("no return value specified for readCounters")
	}

	var r0 map[string]int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (map[string]int64, error)); ok {
		return rf(ctx, blobFileName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) map[string]int64); ok {
		r0 = rf(ctx, blobFileName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]int64)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, blobFileName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockCheckPointerModel_readCounters_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'readCounters'
type mockCheckPointerModel_readCounters_Call struct {
	*mock.Call
}

// readCounters is a helper method to define mock.On call
//   - ctx context.Context
//   - blobFileName string
func (_e *mockCheckPointerModel_Expecter) readCounters(ctx interface{}, blobFileName interface{}) *mockCheckPointerModel_readCounters_Call {
	return &mockCheckPointerModel_readCounters_Call{Call: _e.mock.On("readCounters", ctx, blobFileName)}
}

func (_c *mockCheckPointerModel_readCounters_Call) Run(run func(ctx context.Context, blobFileName string)) *mockCheckPointerModel_readCounters_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *mockCheckPointerModel_readCounters_Call) Return(_a0 map[string]int64, _a1 error) *mockCheckPointerModel_readCounters_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockCheckPointerModel_readCounters_Call) RunAndReturn(run func(context.Context, string) (map[string]int64, error)) *mockCheckPointerModel_readCounters_Call {
	_c.Call.Return(run)
	return _c
}

// readItems provides a mock function with given fields: ctx, blobFileName
func (_m *mockCheckPointerModel) readItems(ctx context.Context, blobFileName string) ([]*topKItem, error) {
	ret := _m.Called(ctx, blobFileName)

	if len(ret) == 0 {
		panic("no return value specified for readItems")
	}

	var r0 []*topKItem
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]*topKItem, error)); ok {
		return rf(ctx, blobFileName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []*topKItem); ok {
		r0 = rf(ctx, blobFileName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*topKItem)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, blobFileName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockCheckPointerModel_readItems_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'readItems'
type mockCheckPointerModel_readItems_Call struct {
	*mock.Call
}

// readItems is a helper method to define mock.On call
//   - ctx context.Context
//   - blobFileName string
func (_e *mockCheckPointerModel_Expecter) readItems(ctx interface{}, blobFileName interface{}) *mockCheckPointerModel_readItems_Call {
	return &mockCheckPointerModel_readItems_Call{Call: _e.mock.On("readItems", ctx, blobFileName)}
}

func (_c *mockCheckPointerModel_readItems_Call) Run(run func(ctx context.Context, blobFileName string)) *mockCheckPointerModel_readItems_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *mockCheckPointerModel_readItems_Call) Return(_a0 []*topKItem, _a1 error) *mockCheckPointerModel_readItems_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockCheckPointerModel_readItems_Call) RunAndReturn(run func(context.Context, string) ([]*topKItem, error)) *mockCheckPointerModel_readItems_Call {
	_c.Call.Return(run)
	return _c
}

// readManifest provides a mock function with given fields: ctx
func (_m *mockCheckPointerModel) readManifest(ctx context.Context) (checkPointManifest, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for readManifest")
	}

	var r0 checkPointManifest
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (checkPointManifest, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) checkPointManifest); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(checkPointManifest)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockCheckPointerModel_readManifest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'readManifest'
type mockCheckPointerModel_readManifest_Call struct {
	*mock.Call
}

// readManifest is a helper method to define mock.On call
//   - ctx context.Context
func (_e *mockCheckPointerModel_Expecter) readManifest(ctx interface{}) *mockCheckPointerModel_readManifest_Call {
	return &mockCheckPointerModel_readManifest_Call{Call: _e.mock.On("readManifest", ctx)}
}

func (_c *mockCheckPointerModel_readManifest_Call) Run(run func(ctx context.Context)) *mockCheckPointerModel_readManifest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *mockCheckPointerModel_readManifest_Call) Return(_a0 checkPointManifest, _a1 error) *mockCheckPointerModel_readManifest_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockCheckPointerModel_readManifest_Call) RunAndReturn(run func(context.Context) (checkPointManifest, error)) *mockCheckPointerModel_readManifest_Call {
	_c.Call.Return(run)
	return _c
}

// writeCounters provides a mock function with given fields: ctx, blobFileName, val
func (_m *mockCheckPointerModel) writeCounters(ctx context.Context, blobFileName string, val map[string]int64) error {
	ret := _m.Called(ctx, blobFileName, val)

	if len(ret) == 0 {
		panic("no return value specified for writeCounters")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]int64) error); ok {
		r0 = rf(ctx, blobFileName, val)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockCheckPointerModel_writeCounters_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'writeCounters'
type mockCheckPointerModel_writeCounters_Call struct {
	*mock.Call
}

// writeCounters is a helper method to define mock.On call
//   - ctx context.Context
//   - blobFileName string
//   - val map[string]int64
func (_e *mockCheckPointerModel_Expecter) writeCounters(ctx interface{}, blobFileName interface{}, val interface{}) *mockCheckPointerModel_writeCounters_Call {
	return &mockCheckPointerModel_writeCounters_Call{Call: _e.mock.On("writeCounters", ctx, blobFileName, val)}
}

func (_c *mockCheckPointerModel_writeCounters_Call) Run(run func(ctx context.Context, blobFileName string, val map[string]int64)) *mockCheckPointerModel_writeCounters_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(map[string]int64))
	})
	return _c
}

func (_c *mockCheckPointerModel_writeCounters_Call) Return(_a0 error) *mockCheckPointerModel_writeCounters_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCheckPointerModel_writeCounters_Call) RunAndReturn(run func(context.Context, string, map[string]int64) error) *mockCheckPointerModel_writeCounters_Call {
	_c.Call.Return(run)
	return _c
}

// writeItems provides a mock function with given fields: ctx, blobFileName, val
func (_m *mockCheckPointerModel) writeItems(ctx context.Context, blobFileName string, val []*topKItem) error {
	ret := _m.Called(ctx, blobFileName, val)

	if len(ret) == 0 {
		panic("no return value specified for writeItems")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []*topKItem) error); ok {
		r0 = rf(ctx, blobFileName, val)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockCheckPointerModel_writeItems_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'writeItems'
type mockCheckPointerModel_writeItems_Call struct {
	*mock.Call
}

// writeItems is a helper method to define mock.On call
//   - ctx context.Context
//   - blobFileName string
//   - val []*topKItem
func (_e *mockCheckPointerModel_Expecter) writeItems(ctx interface{}, blobFileName interface{}, val interface{}) *mockCheckPointerModel_writeItems_Call {
	return &mockCheckPointerModel_writeItems_Call{Call: _e.mock.On("writeItems", ctx, blobFileName, val)}
}

func (_c *mockCheckPointerModel_writeItems_Call) Run(run func(ctx context.Context, blobFileName string, val []*topKItem)) *mockCheckPointerModel_writeItems_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].([]*topKItem))
	})
	return _c
}

func (_c *mockCheckPointerModel_writeItems_Call) Return(_a0 error) *mockCheckPointerModel_writeItems_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCheckPointerModel_writeItems_Call) RunAndReturn(run func(context.Context, string, []*topKItem) error) *mockCheckPointerModel_writeItems_Call {
	_c.Call.Return(run)
	return _c
}

// writeManifest provides a mock function with given fields: ctx, manifest
func (_m *mockCheckPointerModel) writeManifest(ctx context.Context, manifest checkPointManifest) error {
	ret := _m.Called(ctx, manifest)

	if len(ret) == 0 {
		panic("no return value specified for writeManifest")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, checkPointManifest) error); ok {
		r0 = rf(ctx, manifest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockCheckPointerModel_writeManifest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'writeManifest'
type mockCheckPointerModel_writeManifest_Call struct {
	*mock.Call
}

// writeManifest is a helper method to define mock.On call
//   - ctx context.Context
//   - manifest checkPointManifest
func (_e *mockCheckPointerModel_Expecter) writeManifest(ctx interface{}, manifest interface{}) *mockCheckPointerModel_writeManifest_Call {
	return &mockCheckPointerModel_writeManifest_Call{Call: _e.mock.On("writeManifest", ctx, manifest)}
}

func (_c *mockCheckPointerModel_writeManifest_Call) Run(run func(ctx context.Context, manifest checkPointManifest)) *mockCheckPointerModel_writeManifest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(checkPointManifest))
	})
	return _c
}

func (_c *mockCheckPointerModel_writeManifest_Call) Return(_a0 error) *mockCheckPointerModel_writeManifest_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCheckPointerModel_writeManifest_Call) RunAndReturn(run func(context.Context, checkPointManifest) error) *mockCheckPointerModel_writeManifest_Call {
	_c.Call.Return(run)
	return _c
}

// newMockCheckPointerModel creates a new instance of mockCheckPointerModel. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockCheckPointerModel(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockCheckPointerModel {
	mock := &mockCheckPointerModel{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
