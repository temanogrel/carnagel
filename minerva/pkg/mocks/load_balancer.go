// Code generated by mockery v1.0.0. DO NOT EDIT.
package infinity_mocks

import minerva "git.misc.vee.bz/carnagel/minerva/pkg"
import mock "github.com/stretchr/testify/mock"

// LoadBalancer is an autogenerated mock type for the LoadBalancer type
type LoadBalancer struct {
	mock.Mock
}

// RecommendDownload provides a mock function with given fields: source, fileUUID
func (_m *LoadBalancer) RecommendDownload(source minerva.Hostname, fileUUID string) (minerva.DownloadPath, error) {
	ret := _m.Called(source, fileUUID)

	var r0 minerva.DownloadPath
	if rf, ok := ret.Get(0).(func(minerva.Hostname, string) minerva.DownloadPath); ok {
		r0 = rf(source, fileUUID)
	} else {
		r0 = ret.Get(0).(minerva.DownloadPath)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(minerva.Hostname, string) error); ok {
		r1 = rf(source, fileUUID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RecommendStorage provides a mock function with given fields: source, size
func (_m *LoadBalancer) RecommendStorage(source minerva.Hostname, size uint64) (*minerva.Server, string, error) {
	ret := _m.Called(source, size)

	var r0 *minerva.Server
	if rf, ok := ret.Get(0).(func(minerva.Hostname, uint64) *minerva.Server); ok {
		r0 = rf(source, size)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*minerva.Server)
		}
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(minerva.Hostname, uint64) string); ok {
		r1 = rf(source, size)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(minerva.Hostname, uint64) error); ok {
		r2 = rf(source, size)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// RedistributeData provides a mock function with given fields: top, bottom, amountPerServer
func (_m *LoadBalancer) RedistributeData(top int, bottom int, amountPerServer uint64) minerva.RedistributionReport {
	ret := _m.Called(top, bottom, amountPerServer)

	var r0 minerva.RedistributionReport
	if rf, ok := ret.Get(0).(func(int, int, uint64) minerva.RedistributionReport); ok {
		r0 = rf(top, bottom, amountPerServer)
	} else {
		r0 = ret.Get(0).(minerva.RedistributionReport)
	}

	return r0
}
