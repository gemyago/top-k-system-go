// Code generated by mockery. DO NOT EDIT.

//go:build !release

package aggregation

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockItemEventsAggregator is an autogenerated mock type for the ItemEventsAggregator type
type MockItemEventsAggregator struct {
	mock.Mock
}

type MockItemEventsAggregator_Expecter struct {
	mock *mock.Mock
}

func (_m *MockItemEventsAggregator) EXPECT() *MockItemEventsAggregator_Expecter {
	return &MockItemEventsAggregator_Expecter{mock: &_m.Mock}
}

// BeginAggregating provides a mock function with given fields: _a0, counters, opts
func (_m *MockItemEventsAggregator) BeginAggregating(_a0 context.Context, counters Counters, opts BeginAggregatingOpts) error {
	ret := _m.Called(_a0, counters, opts)

	if len(ret) == 0 {
		panic("no return value specified for BeginAggregating")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, Counters, BeginAggregatingOpts) error); ok {
		r0 = rf(_a0, counters, opts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockItemEventsAggregator_BeginAggregating_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BeginAggregating'
type MockItemEventsAggregator_BeginAggregating_Call struct {
	*mock.Call
}

// BeginAggregating is a helper method to define mock.On call
//   - _a0 context.Context
//   - counters Counters
//   - opts BeginAggregatingOpts
func (_e *MockItemEventsAggregator_Expecter) BeginAggregating(_a0 interface{}, counters interface{}, opts interface{}) *MockItemEventsAggregator_BeginAggregating_Call {
	return &MockItemEventsAggregator_BeginAggregating_Call{Call: _e.mock.On("BeginAggregating", _a0, counters, opts)}
}

func (_c *MockItemEventsAggregator_BeginAggregating_Call) Run(run func(_a0 context.Context, counters Counters, opts BeginAggregatingOpts)) *MockItemEventsAggregator_BeginAggregating_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(Counters), args[2].(BeginAggregatingOpts))
	})
	return _c
}

func (_c *MockItemEventsAggregator_BeginAggregating_Call) Return(_a0 error) *MockItemEventsAggregator_BeginAggregating_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockItemEventsAggregator_BeginAggregating_Call) RunAndReturn(run func(context.Context, Counters, BeginAggregatingOpts) error) *MockItemEventsAggregator_BeginAggregating_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockItemEventsAggregator creates a new instance of MockItemEventsAggregator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockItemEventsAggregator(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockItemEventsAggregator {
	mock := &MockItemEventsAggregator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
