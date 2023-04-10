// Code generated by mockery v1.0.0. DO NOT EDIT.

package infinity_mocks

import context "context"
import gobcy "github.com/blockcypher/gobcy"
import infinity "git.misc.vee.bz/carnagel/infinity-api/pkg"
import mock "github.com/stretchr/testify/mock"

// PaymentService is an autogenerated mock type for the PaymentService type
type PaymentService struct {
	mock.Mock
}

// AddressForwardingCallback provides a mock function with given fields: transaction, addrForward
func (_m *PaymentService) AddressForwardingCallback(transaction *infinity.PaymentTransaction, addrForward *gobcy.Payback) error {
	ret := _m.Called(transaction, addrForward)

	var r0 error
	if rf, ok := ret.Get(0).(func(*infinity.PaymentTransaction, *gobcy.Payback) error); ok {
		r0 = rf(transaction, addrForward)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InitiatePurchase provides a mock function with given fields: user, plan
func (_m *PaymentService) InitiatePurchase(user *infinity.User, plan *infinity.PaymentPlan) (*infinity.PaymentTransaction, error) {
	ret := _m.Called(user, plan)

	var r0 *infinity.PaymentTransaction
	if rf, ok := ret.Get(0).(func(*infinity.User, *infinity.PaymentPlan) *infinity.PaymentTransaction); ok {
		r0 = rf(user, plan)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*infinity.PaymentTransaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*infinity.User, *infinity.PaymentPlan) error); ok {
		r1 = rf(user, plan)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Run provides a mock function with given fields: ctx
func (_m *PaymentService) Run(ctx context.Context) {
	_m.Called(ctx)
}

// UpdatePurchase provides a mock function with given fields: transaction
func (_m *PaymentService) UpdatePurchase(transaction *infinity.PaymentTransaction) error {
	ret := _m.Called(transaction)

	var r0 error
	if rf, ok := ret.Get(0).(func(*infinity.PaymentTransaction) error); ok {
		r0 = rf(transaction)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}