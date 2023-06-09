// Code generated by mockery v1.0.0. DO NOT EDIT.
package mocks

import common "git.misc.vee.bz/carnagel/go-ecosystem/protobuf/common"
import context "context"
import grpc "google.golang.org/grpc"
import mock "github.com/stretchr/testify/mock"

// ContentClient is an autogenerated mock type for the ContentClient type
type ContentClient struct {
	mock.Mock
}

// DeleteRecording provides a mock function with given fields: ctx, in, opts
func (_m *ContentClient) DeleteRecording(ctx context.Context, in *common.RecordingIdentifier, opts ...grpc.CallOption) (*common.DeleteRecordingResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *common.DeleteRecordingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *common.RecordingIdentifier, ...grpc.CallOption) *common.DeleteRecordingResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*common.DeleteRecordingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *common.RecordingIdentifier, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertRecording provides a mock function with given fields: ctx, in, opts
func (_m *ContentClient) UpsertRecording(ctx context.Context, in *common.RecordingIdentifier, opts ...grpc.CallOption) (*common.UpsertRecordingResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *common.UpsertRecordingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *common.RecordingIdentifier, ...grpc.CallOption) *common.UpsertRecordingResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*common.UpsertRecordingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *common.RecordingIdentifier, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
