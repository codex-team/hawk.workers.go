// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	worker "github.com/codex-team/hawk.workers.go/lib/worker"
	mock "github.com/stretchr/testify/mock"
)

// TaskHandler is an autogenerated mock type for the TaskHandler type
type TaskHandler struct {
	mock.Mock
}

// Execute provides a mock function with given fields: ctx
func (_m *TaskHandler) Execute(ctx worker.HandlerContext) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(worker.HandlerContext) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewTaskHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewTaskHandler creates a new instance of TaskHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTaskHandler(t mockConstructorTestingTNewTaskHandler) *TaskHandler {
	mock := &TaskHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}