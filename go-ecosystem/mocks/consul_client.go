// Code generated by mockery v1.0.0. DO NOT EDIT.
package mocks

import api "github.com/hashicorp/consul/api"
import context "context"

import mock "github.com/stretchr/testify/mock"

// ConsulClient is an autogenerated mock type for the ConsulClient type
type ConsulClient struct {
	mock.Mock
}

// API provides a mock function with given fields:
func (_m *ConsulClient) API() *api.Client {
	ret := _m.Called()

	var r0 *api.Client
	if rf, ok := ret.Get(0).(func() *api.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Client)
		}
	}

	return r0
}

// GetBool provides a mock function with given fields: key, defaultValue
func (_m *ConsulClient) GetBool(key string, defaultValue bool) bool {
	ret := _m.Called(key, defaultValue)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, bool) bool); ok {
		r0 = rf(key, defaultValue)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// GetInt16 provides a mock function with given fields: key, defaultValue
func (_m *ConsulClient) GetInt16(key string, defaultValue int16) int16 {
	ret := _m.Called(key, defaultValue)

	var r0 int16
	if rf, ok := ret.Get(0).(func(string, int16) int16); ok {
		r0 = rf(key, defaultValue)
	} else {
		r0 = ret.Get(0).(int16)
	}

	return r0
}

// GetInt32 provides a mock function with given fields: key, defaultValue
func (_m *ConsulClient) GetInt32(key string, defaultValue int32) int32 {
	ret := _m.Called(key, defaultValue)

	var r0 int32
	if rf, ok := ret.Get(0).(func(string, int32) int32); ok {
		r0 = rf(key, defaultValue)
	} else {
		r0 = ret.Get(0).(int32)
	}

	return r0
}

// GetInt64 provides a mock function with given fields: key, defaultValue
func (_m *ConsulClient) GetInt64(key string, defaultValue int64) int64 {
	ret := _m.Called(key, defaultValue)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string, int64) int64); ok {
		r0 = rf(key, defaultValue)
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// GetInt8 provides a mock function with given fields: key, defaultValue
func (_m *ConsulClient) GetInt8(key string, defaultValue int8) int8 {
	ret := _m.Called(key, defaultValue)

	var r0 int8
	if rf, ok := ret.Get(0).(func(string, int8) int8); ok {
		r0 = rf(key, defaultValue)
	} else {
		r0 = ret.Get(0).(int8)
	}

	return r0
}

// GetString provides a mock function with given fields: key, defaultValue
func (_m *ConsulClient) GetString(key string, defaultValue string) string {
	ret := _m.Called(key, defaultValue)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(key, defaultValue)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetUint16 provides a mock function with given fields: key, defaultValue
func (_m *ConsulClient) GetUint16(key string, defaultValue uint16) uint16 {
	ret := _m.Called(key, defaultValue)

	var r0 uint16
	if rf, ok := ret.Get(0).(func(string, uint16) uint16); ok {
		r0 = rf(key, defaultValue)
	} else {
		r0 = ret.Get(0).(uint16)
	}

	return r0
}

// GetUint32 provides a mock function with given fields: key, defaultValue
func (_m *ConsulClient) GetUint32(key string, defaultValue uint32) uint32 {
	ret := _m.Called(key, defaultValue)

	var r0 uint32
	if rf, ok := ret.Get(0).(func(string, uint32) uint32); ok {
		r0 = rf(key, defaultValue)
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// GetUint64 provides a mock function with given fields: key, defaultValue
func (_m *ConsulClient) GetUint64(key string, defaultValue uint64) uint64 {
	ret := _m.Called(key, defaultValue)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(string, uint64) uint64); ok {
		r0 = rf(key, defaultValue)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// GetUint8 provides a mock function with given fields: key, defaultValue
func (_m *ConsulClient) GetUint8(key string, defaultValue uint8) uint8 {
	ret := _m.Called(key, defaultValue)

	var r0 uint8
	if rf, ok := ret.Get(0).(func(string, uint8) uint8); ok {
		r0 = rf(key, defaultValue)
	} else {
		r0 = ret.Get(0).(uint8)
	}

	return r0
}

// PollBool provides a mock function with given fields: ctx, key
func (_m *ConsulClient) PollBool(ctx context.Context, key string) (chan bool, chan error) {
	ret := _m.Called(ctx, key)

	var r0 chan bool
	if rf, ok := ret.Get(0).(func(context.Context, string) chan bool); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan bool)
		}
	}

	var r1 chan error
	if rf, ok := ret.Get(1).(func(context.Context, string) chan error); ok {
		r1 = rf(ctx, key)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(chan error)
		}
	}

	return r0, r1
}

// PollString provides a mock function with given fields: ctx, key
func (_m *ConsulClient) PollString(ctx context.Context, key string) (chan string, chan error) {
	ret := _m.Called(ctx, key)

	var r0 chan string
	if rf, ok := ret.Get(0).(func(context.Context, string) chan string); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan string)
		}
	}

	var r1 chan error
	if rf, ok := ret.Get(1).(func(context.Context, string) chan error); ok {
		r1 = rf(ctx, key)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(chan error)
		}
	}

	return r0, r1
}
