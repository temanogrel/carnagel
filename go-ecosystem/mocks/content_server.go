// Code generated by mockery v1.0.0. DO NOT EDIT.
package mocks

import common "git.misc.vee.bz/carnagel/go-ecosystem/protobuf/common"
import context "context"
import mock "github.com/stretchr/testify/mock"

// ContentServer is an autogenerated mock type for the ContentServer type
type ContentServer struct {
	mock.Mock
}

// DeleteRecording provides a mock function with given fields: _a0, _a1
func (_m *ContentServer) DeleteRecording(_a0 context.Context, _a1 *common.RecordingIdentifier) (*common.DeleteRecordingResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *common.DeleteRecordingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *common.RecordingIdentifier) *common.DeleteRecordingResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*common.DeleteRecordingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *common.RecordingIdentifier) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertRecording provides a mock function with given fields: _a0, _a1
func (_m *ContentServer) UpsertRecording(_a0 context.Context, _a1 *common.RecordingIdentifier) (*common.UpsertRecordingResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *common.UpsertRecordingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *common.RecordingIdentifier) *common.UpsertRecordingResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*common.UpsertRecordingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *common.RecordingIdentifier) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
