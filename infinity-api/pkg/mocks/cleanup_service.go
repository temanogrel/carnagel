// Code generated by mockery v1.0.0. DO NOT EDIT.

package infinity_mocks

import context "context"

import mock "github.com/stretchr/testify/mock"

// CleanupService is an autogenerated mock type for the CleanupService type
type CleanupService struct {
	mock.Mock
}

// Run provides a mock function with given fields: ctx
func (_m *CleanupService) Run(ctx context.Context) {
	_m.Called(ctx)
}

// RunCleanup provides a mock function with given fields: ctx
func (_m *CleanupService) RunCleanup(ctx context.Context) {
	_m.Called(ctx)
}
