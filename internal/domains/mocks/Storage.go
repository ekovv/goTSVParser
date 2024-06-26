// Code generated by mockery v3.0.0-alpha.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	shema "goTSVParser/internal/shema"
)

// Storage is an autogenerated mock type for the Storage type
type Storage struct {
	mock.Mock
}

// GetAllGuids provides a mock function with given fields: ctx, unitGuid
func (_m *Storage) GetAllGuids(ctx context.Context, unitGuid string) ([]shema.Tsv, error) {
	ret := _m.Called(ctx, unitGuid)

	var r0 []shema.Tsv
	if rf, ok := ret.Get(0).(func(context.Context, string) []shema.Tsv); ok {
		r0 = rf(ctx, unitGuid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]shema.Tsv)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, unitGuid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCheckedFiles provides a mock function with given fields:
func (_m *Storage) GetCheckedFiles() ([]shema.ParsedFiles, error) {
	ret := _m.Called()

	var r0 []shema.ParsedFiles
	if rf, ok := ret.Get(0).(func() []shema.ParsedFiles); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]shema.ParsedFiles)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: sh
func (_m *Storage) Save(sh shema.Tsv) error {
	ret := _m.Called(sh)

	var r0 error
	if rf, ok := ret.Get(0).(func(shema.Tsv) error); ok {
		r0 = rf(sh)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveFiles provides a mock function with given fields: fileName
func (_m *Storage) SaveFiles(fileName string) error {
	ret := _m.Called(fileName)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(fileName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveFilesWithErr provides a mock function with given fields: sh
func (_m *Storage) SaveFilesWithErr(sh shema.Files) error {
	ret := _m.Called(sh)

	var r0 error
	if rf, ok := ret.Get(0).(func(shema.Files) error); ok {
		r0 = rf(sh)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ShutDown provides a mock function with given fields:
func (_m *Storage) ShutDown() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewStorage interface {
	mock.TestingT
	Cleanup(func())
}

// NewStorage creates a new instance of Storage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStorage(t mockConstructorTestingTNewStorage) *Storage {
	mock := &Storage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
