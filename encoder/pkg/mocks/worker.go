// Code generated by mockery v1.0.0. DO NOT EDIT.
package encoder_mocks

import context "context"

import mock "github.com/stretchr/testify/mock"

// Worker is an autogenerated mock type for the Worker type
type Worker struct {
	mock.Mock
}

// Run provides a mock function with given fields: ctx, routines
func (_m *Worker) Run(ctx context.Context, routines int) {
	_m.Called(ctx, routines)
}