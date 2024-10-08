// Code generated by mockery. DO NOT EDIT.

//go:build !release

package aggregation

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// mockCheckPointer is an autogenerated mock type for the checkPointer type
type mockCheckPointer struct {
	mock.Mock
}

type mockCheckPointer_Expecter struct {
	mock *mock.Mock
}

func (_m *mockCheckPointer) EXPECT() *mockCheckPointer_Expecter {
	return &mockCheckPointer_Expecter{mock: &_m.Mock}
}

// dumpState provides a mock function with given fields: ctx, counters1
func (_m *mockCheckPointer) dumpState(ctx context.Context, counters1 counters) error {
	ret := _m.Called(ctx, counters1)

	if len(ret) == 0 {
		panic("no return value specified for dumpState")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, counters) error); ok {
		r0 = rf(ctx, counters1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockCheckPointer_dumpState_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'dumpState'
type mockCheckPointer_dumpState_Call struct {
	*mock.Call
}

// dumpState is a helper method to define mock.On call
//   - ctx context.Context
//   - counters1 counters
func (_e *mockCheckPointer_Expecter) dumpState(ctx interface{}, counters1 interface{}) *mockCheckPointer_dumpState_Call {
	return &mockCheckPointer_dumpState_Call{Call: _e.mock.On("dumpState", ctx, counters1)}
}

func (_c *mockCheckPointer_dumpState_Call) Run(run func(ctx context.Context, counters1 counters)) *mockCheckPointer_dumpState_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(counters))
	})
	return _c
}

func (_c *mockCheckPointer_dumpState_Call) Return(_a0 error) *mockCheckPointer_dumpState_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCheckPointer_dumpState_Call) RunAndReturn(run func(context.Context, counters) error) *mockCheckPointer_dumpState_Call {
	_c.Call.Return(run)
	return _c
}

// restoreState provides a mock function with given fields: ctx, counters1
func (_m *mockCheckPointer) restoreState(ctx context.Context, counters1 counters) error {
	ret := _m.Called(ctx, counters1)

	if len(ret) == 0 {
		panic("no return value specified for restoreState")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, counters) error); ok {
		r0 = rf(ctx, counters1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockCheckPointer_restoreState_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'restoreState'
type mockCheckPointer_restoreState_Call struct {
	*mock.Call
}

// restoreState is a helper method to define mock.On call
//   - ctx context.Context
//   - counters1 counters
func (_e *mockCheckPointer_Expecter) restoreState(ctx interface{}, counters1 interface{}) *mockCheckPointer_restoreState_Call {
	return &mockCheckPointer_restoreState_Call{Call: _e.mock.On("restoreState", ctx, counters1)}
}

func (_c *mockCheckPointer_restoreState_Call) Run(run func(ctx context.Context, counters1 counters)) *mockCheckPointer_restoreState_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(counters))
	})
	return _c
}

func (_c *mockCheckPointer_restoreState_Call) Return(_a0 error) *mockCheckPointer_restoreState_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCheckPointer_restoreState_Call) RunAndReturn(run func(context.Context, counters) error) *mockCheckPointer_restoreState_Call {
	_c.Call.Return(run)
	return _c
}

// newMockCheckPointer creates a new instance of mockCheckPointer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockCheckPointer(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockCheckPointer {
	mock := &mockCheckPointer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
