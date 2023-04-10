// Code generated by mockery v1.0.0. DO NOT EDIT.

package infinity_mocks

import context "context"

import mock "github.com/stretchr/testify/mock"

// CryptoExchangeRateService is an autogenerated mock type for the CryptoExchangeRateService type
type CryptoExchangeRateService struct {
	mock.Mock
}

// ConvertUSDToBtc provides a mock function with given fields: amount
func (_m *CryptoExchangeRateService) ConvertUSDToBtc(amount float64) (float64, float64, error) {
	ret := _m.Called(amount)

	var r0 float64
	if rf, ok := ret.Get(0).(func(float64) float64); ok {
		r0 = rf(amount)
	} else {
		r0 = ret.Get(0).(float64)
	}

	var r1 float64
	if rf, ok := ret.Get(1).(func(float64) float64); ok {
		r1 = rf(amount)
	} else {
		r1 = ret.Get(1).(float64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(float64) error); ok {
		r2 = rf(amount)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Run provides a mock function with given fields: ctx
func (_m *CryptoExchangeRateService) Run(ctx context.Context) {
	_m.Called(ctx)
}
