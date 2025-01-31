// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// KafkaClient is an autogenerated mock type for the KafkaClient type
type KafkaClient struct {
	mock.Mock
}

// Close provides a mock function with no fields
func (_m *KafkaClient) Close() {
	_m.Called()
}

// Produce provides a mock function with given fields: topic, key, value
func (_m *KafkaClient) Produce(topic string, key []byte, value []byte) error {
	ret := _m.Called(topic, key, value)

	if len(ret) == 0 {
		panic("no return value specified for Produce")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []byte, []byte) error); ok {
		r0 = rf(topic, key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewKafkaClient creates a new instance of KafkaClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewKafkaClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *KafkaClient {
	mock := &KafkaClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
