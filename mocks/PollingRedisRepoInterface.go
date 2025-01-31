// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// PollingRedisRepoInterface is an autogenerated mock type for the PollingRedisRepoInterface type
type PollingRedisRepoInterface struct {
	mock.Mock
}

// SetNX provides a mock function with given fields: hashKey, actionID, expiration
func (_m *PollingRedisRepoInterface) SetNX(hashKey string, actionID string, expiration time.Duration) (bool, error) {
	ret := _m.Called(hashKey, actionID, expiration)

	if len(ret) == 0 {
		panic("no return value specified for SetNX")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, time.Duration) (bool, error)); ok {
		return rf(hashKey, actionID, expiration)
	}
	if rf, ok := ret.Get(0).(func(string, string, time.Duration) bool); ok {
		r0 = rf(hashKey, actionID, expiration)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string, string, time.Duration) error); ok {
		r1 = rf(hashKey, actionID, expiration)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidateActionGlobalUUID provides a mock function with given fields: field
func (_m *PollingRedisRepoInterface) ValidateActionGlobalUUID(field *string) (bool, error) {
	ret := _m.Called(field)

	if len(ret) == 0 {
		panic("no return value specified for ValidateActionGlobalUUID")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(*string) (bool, error)); ok {
		return rf(field)
	}
	if rf, ok := ret.Get(0).(func(*string) bool); ok {
		r0 = rf(field)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(*string) error); ok {
		r1 = rf(field)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewPollingRedisRepoInterface creates a new instance of PollingRedisRepoInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPollingRedisRepoInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *PollingRedisRepoInterface {
	mock := &PollingRedisRepoInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
