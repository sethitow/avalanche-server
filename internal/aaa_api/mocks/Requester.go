// Code generated by mockery v2.34.2. DO NOT EDIT.

package mocks

import (
	aaa_api "avalancheserver/internal/aaa_api"

	mock "github.com/stretchr/testify/mock"
)

// Requester is an autogenerated mock type for the Requester type
type Requester struct {
	mock.Mock
}

// GetForecastByCenter provides a mock function with given fields: _a0
func (_m *Requester) GetForecastByCenter(_a0 string) (aaa_api.Root, error) {
	ret := _m.Called(_a0)

	var r0 aaa_api.Root
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (aaa_api.Root, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) aaa_api.Root); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(aaa_api.Root)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRequester creates a new instance of Requester. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRequester(t interface {
	mock.TestingT
	Cleanup(func())
}) *Requester {
	mock := &Requester{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
