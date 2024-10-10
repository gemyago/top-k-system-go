// Code generated by mockery. DO NOT EDIT.

//go:build !release

package services

import mock "github.com/stretchr/testify/mock"

// mockKafkaConn is an autogenerated mock type for the kafkaConn type
type mockKafkaConn struct {
	mock.Mock
}

type mockKafkaConn_Expecter struct {
	mock *mock.Mock
}

func (_m *mockKafkaConn) EXPECT() *mockKafkaConn_Expecter {
	return &mockKafkaConn_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with given fields:
func (_m *mockKafkaConn) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockKafkaConn_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type mockKafkaConn_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *mockKafkaConn_Expecter) Close() *mockKafkaConn_Close_Call {
	return &mockKafkaConn_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *mockKafkaConn_Close_Call) Run(run func()) *mockKafkaConn_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockKafkaConn_Close_Call) Return(_a0 error) *mockKafkaConn_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockKafkaConn_Close_Call) RunAndReturn(run func() error) *mockKafkaConn_Close_Call {
	_c.Call.Return(run)
	return _c
}

// ReadLastOffset provides a mock function with given fields:
func (_m *mockKafkaConn) ReadLastOffset() (int64, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ReadLastOffset")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func() (int64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockKafkaConn_ReadLastOffset_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReadLastOffset'
type mockKafkaConn_ReadLastOffset_Call struct {
	*mock.Call
}

// ReadLastOffset is a helper method to define mock.On call
func (_e *mockKafkaConn_Expecter) ReadLastOffset() *mockKafkaConn_ReadLastOffset_Call {
	return &mockKafkaConn_ReadLastOffset_Call{Call: _e.mock.On("ReadLastOffset")}
}

func (_c *mockKafkaConn_ReadLastOffset_Call) Run(run func()) *mockKafkaConn_ReadLastOffset_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockKafkaConn_ReadLastOffset_Call) Return(_a0 int64, _a1 error) *mockKafkaConn_ReadLastOffset_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockKafkaConn_ReadLastOffset_Call) RunAndReturn(run func() (int64, error)) *mockKafkaConn_ReadLastOffset_Call {
	_c.Call.Return(run)
	return _c
}

// newMockKafkaConn creates a new instance of mockKafkaConn. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockKafkaConn(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockKafkaConn {
	mock := &mockKafkaConn{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
