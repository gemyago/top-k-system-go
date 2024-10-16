// Code generated by mockery. DO NOT EDIT.

//go:build !release

package aggregation

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// mockItemEventsAggregator is an autogenerated mock type for the itemEventsAggregator type
type mockItemEventsAggregator struct {
	mock.Mock
}

type mockItemEventsAggregator_Expecter struct {
	mock *mock.Mock
}

func (_m *mockItemEventsAggregator) EXPECT() *mockItemEventsAggregator_Expecter {
	return &mockItemEventsAggregator_Expecter{mock: &_m.Mock}
}

// beginAggregating provides a mock function with given fields: _a0, state, opts
func (_m *mockItemEventsAggregator) beginAggregating(_a0 context.Context, state aggregationState, opts beginAggregatingOpts) error {
	ret := _m.Called(_a0, state, opts)

	if len(ret) == 0 {
		panic("no return value specified for beginAggregating")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, aggregationState, beginAggregatingOpts) error); ok {
		r0 = rf(_a0, state, opts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockItemEventsAggregator_beginAggregating_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'beginAggregating'
type mockItemEventsAggregator_beginAggregating_Call struct {
	*mock.Call
}

// beginAggregating is a helper method to define mock.On call
//   - _a0 context.Context
//   - state aggregationState
//   - opts beginAggregatingOpts
func (_e *mockItemEventsAggregator_Expecter) beginAggregating(_a0 interface{}, state interface{}, opts interface{}) *mockItemEventsAggregator_beginAggregating_Call {
	return &mockItemEventsAggregator_beginAggregating_Call{Call: _e.mock.On("beginAggregating", _a0, state, opts)}
}

func (_c *mockItemEventsAggregator_beginAggregating_Call) Run(run func(_a0 context.Context, state aggregationState, opts beginAggregatingOpts)) *mockItemEventsAggregator_beginAggregating_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(aggregationState), args[2].(beginAggregatingOpts))
	})
	return _c
}

func (_c *mockItemEventsAggregator_beginAggregating_Call) Return(_a0 error) *mockItemEventsAggregator_beginAggregating_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockItemEventsAggregator_beginAggregating_Call) RunAndReturn(run func(context.Context, aggregationState, beginAggregatingOpts) error) *mockItemEventsAggregator_beginAggregating_Call {
	_c.Call.Return(run)
	return _c
}

// newMockItemEventsAggregator creates a new instance of mockItemEventsAggregator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockItemEventsAggregator(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockItemEventsAggregator {
	mock := &mockItemEventsAggregator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
