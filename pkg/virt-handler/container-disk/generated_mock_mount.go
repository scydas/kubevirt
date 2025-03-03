// Automatically generated by MockGen. DO NOT EDIT!
// Source: mount.go

package container_disk

import (
	time "time"

	gomock "github.com/golang/mock/gomock"
	v1 "kubevirt.io/api/core/v1"
)

// Mock of Mounter interface
type MockMounter struct {
	ctrl     *gomock.Controller
	recorder *_MockMounterRecorder
}

// Recorder for MockMounter (not exported)
type _MockMounterRecorder struct {
	mock *MockMounter
}

func NewMockMounter(ctrl *gomock.Controller) *MockMounter {
	mock := &MockMounter{ctrl: ctrl}
	mock.recorder = &_MockMounterRecorder{mock}
	return mock
}

func (_m *MockMounter) EXPECT() *_MockMounterRecorder {
	return _m.recorder
}

func (_m *MockMounter) ContainerDisksReady(vmi *v1.VirtualMachineInstance, notInitializedSince time.Time) (bool, error) {
	ret := _m.ctrl.Call(_m, "ContainerDisksReady", vmi, notInitializedSince)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockMounterRecorder) ContainerDisksReady(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ContainerDisksReady", arg0, arg1)
}

func (_m *MockMounter) MountAndVerify(vmi *v1.VirtualMachineInstance) error {
	ret := _m.ctrl.Call(_m, "MountAndVerify", vmi)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockMounterRecorder) MountAndVerify(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "MountAndVerify", arg0)
}

func (_m *MockMounter) Unmount(vmi *v1.VirtualMachineInstance) error {
	ret := _m.ctrl.Call(_m, "Unmount", vmi)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockMounterRecorder) Unmount(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Unmount", arg0)
}

func (_m *MockMounter) ComputeChecksums(vmi *v1.VirtualMachineInstance) (*DiskChecksums, error) {
	ret := _m.ctrl.Call(_m, "ComputeChecksums", vmi)
	ret0, _ := ret[0].(*DiskChecksums)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockMounterRecorder) ComputeChecksums(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ComputeChecksums", arg0)
}
