// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	entities "github.com/abx123/library/entities"
)

// IdbRepo is an autogenerated mock type for the IdbRepo type
type IdbRepo struct {
	mock.Mock
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *IdbRepo) Get(_a0 context.Context, _a1 *entities.Book) (*entities.Book, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *entities.Book
	if rf, ok := ret.Get(0).(func(context.Context, *entities.Book) *entities.Book); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Book)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *entities.Book) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *IdbRepo) List(_a0 context.Context, _a1 int64, _a2 int64, _a3 string) ([]*entities.Book, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 []*entities.Book
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, string) []*entities.Book); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entities.Book)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, int64, string) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Upsert provides a mock function with given fields: _a0, _a1
func (_m *IdbRepo) Upsert(_a0 context.Context, _a1 *entities.Book) (*entities.Book, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *entities.Book
	if rf, ok := ret.Get(0).(func(context.Context, *entities.Book) *entities.Book); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Book)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *entities.Book) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
