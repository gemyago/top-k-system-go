// Code generated by mockery. DO NOT EDIT.

//go:build !release

package services

import (
	context "context"

	kafka "github.com/segmentio/kafka-go"
	mock "github.com/stretchr/testify/mock"
)

// MockKafkaReader is an autogenerated mock type for the mockKafkaReader type
type MockKafkaReader struct {
	mock.Mock
}

type MockKafkaReader_Expecter struct {
	mock *mock.Mock
}

func (_m *MockKafkaReader) EXPECT() *MockKafkaReader_Expecter {
	return &MockKafkaReader_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with given fields:
func (_m *MockKafkaReader) Close() error {
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

// MockKafkaReader_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type MockKafkaReader_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *MockKafkaReader_Expecter) Close() *MockKafkaReader_Close_Call {
	return &MockKafkaReader_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *MockKafkaReader_Close_Call) Run(run func()) *MockKafkaReader_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockKafkaReader_Close_Call) Return(_a0 error) *MockKafkaReader_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockKafkaReader_Close_Call) RunAndReturn(run func() error) *MockKafkaReader_Close_Call {
	_c.Call.Return(run)
	return _c
}

// CommitMessages provides a mock function with given fields: ctx, msgs
func (_m *MockKafkaReader) CommitMessages(ctx context.Context, msgs ...kafka.Message) error {
	_va := make([]interface{}, len(msgs))
	for _i := range msgs {
		_va[_i] = msgs[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CommitMessages")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ...kafka.Message) error); ok {
		r0 = rf(ctx, msgs...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockKafkaReader_CommitMessages_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CommitMessages'
type MockKafkaReader_CommitMessages_Call struct {
	*mock.Call
}

// CommitMessages is a helper method to define mock.On call
//   - ctx context.Context
//   - msgs ...kafka.Message
func (_e *MockKafkaReader_Expecter) CommitMessages(ctx interface{}, msgs ...interface{}) *MockKafkaReader_CommitMessages_Call {
	return &MockKafkaReader_CommitMessages_Call{Call: _e.mock.On("CommitMessages",
		append([]interface{}{ctx}, msgs...)...)}
}

func (_c *MockKafkaReader_CommitMessages_Call) Run(run func(ctx context.Context, msgs ...kafka.Message)) *MockKafkaReader_CommitMessages_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]kafka.Message, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(kafka.Message)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *MockKafkaReader_CommitMessages_Call) Return(_a0 error) *MockKafkaReader_CommitMessages_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockKafkaReader_CommitMessages_Call) RunAndReturn(run func(context.Context, ...kafka.Message) error) *MockKafkaReader_CommitMessages_Call {
	_c.Call.Return(run)
	return _c
}

// FetchMessage provides a mock function with given fields: ctx
func (_m *MockKafkaReader) FetchMessage(ctx context.Context) (kafka.Message, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for FetchMessage")
	}

	var r0 kafka.Message
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (kafka.Message, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) kafka.Message); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(kafka.Message)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockKafkaReader_FetchMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FetchMessage'
type MockKafkaReader_FetchMessage_Call struct {
	*mock.Call
}

// FetchMessage is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockKafkaReader_Expecter) FetchMessage(ctx interface{}) *MockKafkaReader_FetchMessage_Call {
	return &MockKafkaReader_FetchMessage_Call{Call: _e.mock.On("FetchMessage", ctx)}
}

func (_c *MockKafkaReader_FetchMessage_Call) Run(run func(ctx context.Context)) *MockKafkaReader_FetchMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockKafkaReader_FetchMessage_Call) Return(_a0 kafka.Message, _a1 error) *MockKafkaReader_FetchMessage_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockKafkaReader_FetchMessage_Call) RunAndReturn(run func(context.Context) (kafka.Message, error)) *MockKafkaReader_FetchMessage_Call {
	_c.Call.Return(run)
	return _c
}

// ReadLastOffset provides a mock function with given fields: ctx
func (_m *MockKafkaReader) ReadLastOffset(ctx context.Context) (int64, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ReadLastOffset")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (int64, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) int64); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockKafkaReader_ReadLastOffset_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReadLastOffset'
type MockKafkaReader_ReadLastOffset_Call struct {
	*mock.Call
}

// ReadLastOffset is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockKafkaReader_Expecter) ReadLastOffset(ctx interface{}) *MockKafkaReader_ReadLastOffset_Call {
	return &MockKafkaReader_ReadLastOffset_Call{Call: _e.mock.On("ReadLastOffset", ctx)}
}

func (_c *MockKafkaReader_ReadLastOffset_Call) Run(run func(ctx context.Context)) *MockKafkaReader_ReadLastOffset_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockKafkaReader_ReadLastOffset_Call) Return(_a0 int64, _a1 error) *MockKafkaReader_ReadLastOffset_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockKafkaReader_ReadLastOffset_Call) RunAndReturn(run func(context.Context) (int64, error)) *MockKafkaReader_ReadLastOffset_Call {
	_c.Call.Return(run)
	return _c
}

// SetOffset provides a mock function with given fields: offset
func (_m *MockKafkaReader) SetOffset(offset int64) error {
	ret := _m.Called(offset)

	if len(ret) == 0 {
		panic("no return value specified for SetOffset")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(int64) error); ok {
		r0 = rf(offset)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockKafkaReader_SetOffset_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetOffset'
type MockKafkaReader_SetOffset_Call struct {
	*mock.Call
}

// SetOffset is a helper method to define mock.On call
//   - offset int64
func (_e *MockKafkaReader_Expecter) SetOffset(offset interface{}) *MockKafkaReader_SetOffset_Call {
	return &MockKafkaReader_SetOffset_Call{Call: _e.mock.On("SetOffset", offset)}
}

func (_c *MockKafkaReader_SetOffset_Call) Run(run func(offset int64)) *MockKafkaReader_SetOffset_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64))
	})
	return _c
}

func (_c *MockKafkaReader_SetOffset_Call) Return(_a0 error) *MockKafkaReader_SetOffset_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockKafkaReader_SetOffset_Call) RunAndReturn(run func(int64) error) *MockKafkaReader_SetOffset_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockKafkaReader creates a new instance of MockKafkaReader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockKafkaReader(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockKafkaReader {
	mock := &MockKafkaReader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
