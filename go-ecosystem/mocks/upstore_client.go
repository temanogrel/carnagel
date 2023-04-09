// Code generated by mockery v1.0.0. DO NOT EDIT.
package mocks

import mock "github.com/stretchr/testify/mock"
import os "os"

// UpstoreClient is an autogenerated mock type for the UpstoreClient type
type UpstoreClient struct {
	mock.Mock
}

// Upload provides a mock function with given fields: name, file
func (_m *UpstoreClient) Upload(name string, file *os.File) (string, error) {
	ret := _m.Called(name, file)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, *os.File) string); ok {
		r0 = rf(name, file)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *os.File) error); ok {
		r1 = rf(name, file)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
