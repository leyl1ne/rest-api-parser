// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import (
	models "github.com/leyl1ne/rest-api-parser/pkg/models"
	mock "github.com/stretchr/testify/mock"
)

// SongGetterAll is an autogenerated mock type for the SongGetterAll type
type SongGetterAll struct {
	mock.Mock
}

// GetAllSong provides a mock function with given fields:
func (_m *SongGetterAll) GetAllSong() ([]models.Song, error) {
	ret := _m.Called()

	var r0 []models.Song
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]models.Song, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []models.Song); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Song)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewSongGetterAll interface {
	mock.TestingT
	Cleanup(func())
}

// NewSongGetterAll creates a new instance of SongGetterAll. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSongGetterAll(t mockConstructorTestingTNewSongGetterAll) *SongGetterAll {
	mock := &SongGetterAll{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}