// Code generated by mockery v1.0.0. DO NOT EDIT.

package infinity_mocks

import context "context"
import infinity "git.misc.vee.bz/carnagel/infinity-api/pkg"
import mock "github.com/stretchr/testify/mock"
import uuid "github.com/satori/go.uuid"

// PerformerSearchService is an autogenerated mock type for the PerformerSearchService type
type PerformerSearchService struct {
	mock.Mock
}

// AddPerformerAlias provides a mock function with given fields: performerId, stageName
func (_m *PerformerSearchService) AddPerformerAlias(performerId uuid.UUID, stageName string) error {
	ret := _m.Called(performerId, stageName)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, string) error); ok {
		r0 = rf(performerId, stageName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Matching provides a mock function with given fields: criteria
func (_m *PerformerSearchService) Matching(criteria *infinity.PerformerRepositoryCriteria) ([]uuid.UUID, int, error) {
	ret := _m.Called(criteria)

	var r0 []uuid.UUID
	if rf, ok := ret.Get(0).(func(*infinity.PerformerRepositoryCriteria) []uuid.UUID); ok {
		r0 = rf(criteria)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]uuid.UUID)
		}
	}

	var r1 int
	if rf, ok := ret.Get(1).(func(*infinity.PerformerRepositoryCriteria) int); ok {
		r1 = rf(criteria)
	} else {
		r1 = ret.Get(1).(int)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*infinity.PerformerRepositoryCriteria) error); ok {
		r2 = rf(criteria)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// PerformerAliasExists provides a mock function with given fields: performerId, stageName
func (_m *PerformerSearchService) PerformerAliasExists(performerId uuid.UUID, stageName string) (bool, error) {
	ret := _m.Called(performerId, stageName)

	var r0 bool
	if rf, ok := ret.Get(0).(func(uuid.UUID, string) bool); ok {
		r0 = rf(performerId, stageName)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, string) error); ok {
		r1 = rf(performerId, stageName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Run provides a mock function with given fields: ctx
func (_m *PerformerSearchService) Run(ctx context.Context) {
	_m.Called(ctx)
}
