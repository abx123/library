// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"
	entities "library/entities"

	mock "github.com/stretchr/testify/mock"
)

// Igoodread is an autogenerated mock type for the Igoodread type
type Igoodread struct {
	mock.Mock
}

// GetFromGoodread provides a mock function with given fields: _a0, _a1
func (_m *Igoodread) GetFromGoodread(_a0 context.Context, _a1 string) (*entities.Book, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *entities.Book
	if rf, ok := ret.Get(0).(func(context.Context, string) *entities.Book); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Book)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}