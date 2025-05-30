// Code generated by MockGen. DO NOT EDIT.
// Source: manager.go
//
// Generated by this command:
//
//	mockgen -source manager.go -package=virtwrap -destination=generated_mock_manager.go
//

// Package virtwrap is a generated GoMock package.
package virtwrap

import (
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
	v1 "kubevirt.io/api/core/v1"

	v10 "kubevirt.io/kubevirt/pkg/handler-launcher-com/cmd/v1"
	cmdclient "kubevirt.io/kubevirt/pkg/virt-handler/cmd-client"
	api "kubevirt.io/kubevirt/pkg/virt-launcher/virtwrap/api"
	stats "kubevirt.io/kubevirt/pkg/virt-launcher/virtwrap/stats"
)

// MockDomainManager is a mock of DomainManager interface.
type MockDomainManager struct {
	ctrl     *gomock.Controller
	recorder *MockDomainManagerMockRecorder
	isgomock struct{}
}

// MockDomainManagerMockRecorder is the mock recorder for MockDomainManager.
type MockDomainManagerMockRecorder struct {
	mock *MockDomainManager
}

// NewMockDomainManager creates a new mock instance.
func NewMockDomainManager(ctrl *gomock.Controller) *MockDomainManager {
	mock := &MockDomainManager{ctrl: ctrl}
	mock.recorder = &MockDomainManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDomainManager) EXPECT() *MockDomainManagerMockRecorder {
	return m.recorder
}

// CancelVMIMigration mocks base method.
func (m *MockDomainManager) CancelVMIMigration(arg0 *v1.VirtualMachineInstance) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CancelVMIMigration", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CancelVMIMigration indicates an expected call of CancelVMIMigration.
func (mr *MockDomainManagerMockRecorder) CancelVMIMigration(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelVMIMigration", reflect.TypeOf((*MockDomainManager)(nil).CancelVMIMigration), arg0)
}

// DeleteVMI mocks base method.
func (m *MockDomainManager) DeleteVMI(arg0 *v1.VirtualMachineInstance) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteVMI", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteVMI indicates an expected call of DeleteVMI.
func (mr *MockDomainManagerMockRecorder) DeleteVMI(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteVMI", reflect.TypeOf((*MockDomainManager)(nil).DeleteVMI), arg0)
}

// Exec mocks base method.
func (m *MockDomainManager) Exec(arg0, arg1 string, arg2 []string, arg3 int32) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exec", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *MockDomainManagerMockRecorder) Exec(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockDomainManager)(nil).Exec), arg0, arg1, arg2, arg3)
}

// FinalizeVirtualMachineMigration mocks base method.
func (m *MockDomainManager) FinalizeVirtualMachineMigration(arg0 *v1.VirtualMachineInstance, arg1 *v10.VirtualMachineOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FinalizeVirtualMachineMigration", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// FinalizeVirtualMachineMigration indicates an expected call of FinalizeVirtualMachineMigration.
func (mr *MockDomainManagerMockRecorder) FinalizeVirtualMachineMigration(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FinalizeVirtualMachineMigration", reflect.TypeOf((*MockDomainManager)(nil).FinalizeVirtualMachineMigration), arg0, arg1)
}

// FreezeVMI mocks base method.
func (m *MockDomainManager) FreezeVMI(arg0 *v1.VirtualMachineInstance, arg1 int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FreezeVMI", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// FreezeVMI indicates an expected call of FreezeVMI.
func (mr *MockDomainManagerMockRecorder) FreezeVMI(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FreezeVMI", reflect.TypeOf((*MockDomainManager)(nil).FreezeVMI), arg0, arg1)
}

// GetDomainDirtyRateStats mocks base method.
func (m *MockDomainManager) GetDomainDirtyRateStats(calculationDuration time.Duration) (*stats.DomainStatsDirtyRate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDomainDirtyRateStats", calculationDuration)
	ret0, _ := ret[0].(*stats.DomainStatsDirtyRate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDomainDirtyRateStats indicates an expected call of GetDomainDirtyRateStats.
func (mr *MockDomainManagerMockRecorder) GetDomainDirtyRateStats(calculationDuration any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDomainDirtyRateStats", reflect.TypeOf((*MockDomainManager)(nil).GetDomainDirtyRateStats), calculationDuration)
}

// GetDomainStats mocks base method.
func (m *MockDomainManager) GetDomainStats() (*stats.DomainStats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDomainStats")
	ret0, _ := ret[0].(*stats.DomainStats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDomainStats indicates an expected call of GetDomainStats.
func (mr *MockDomainManagerMockRecorder) GetDomainStats() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDomainStats", reflect.TypeOf((*MockDomainManager)(nil).GetDomainStats))
}

// GetFilesystems mocks base method.
func (m *MockDomainManager) GetFilesystems() []v1.VirtualMachineInstanceFileSystem {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilesystems")
	ret0, _ := ret[0].([]v1.VirtualMachineInstanceFileSystem)
	return ret0
}

// GetFilesystems indicates an expected call of GetFilesystems.
func (mr *MockDomainManagerMockRecorder) GetFilesystems() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilesystems", reflect.TypeOf((*MockDomainManager)(nil).GetFilesystems))
}

// GetGuestInfo mocks base method.
func (m *MockDomainManager) GetGuestInfo() v1.VirtualMachineInstanceGuestAgentInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGuestInfo")
	ret0, _ := ret[0].(v1.VirtualMachineInstanceGuestAgentInfo)
	return ret0
}

// GetGuestInfo indicates an expected call of GetGuestInfo.
func (mr *MockDomainManagerMockRecorder) GetGuestInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGuestInfo", reflect.TypeOf((*MockDomainManager)(nil).GetGuestInfo))
}

// GetGuestOSInfo mocks base method.
func (m *MockDomainManager) GetGuestOSInfo() *api.GuestOSInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGuestOSInfo")
	ret0, _ := ret[0].(*api.GuestOSInfo)
	return ret0
}

// GetGuestOSInfo indicates an expected call of GetGuestOSInfo.
func (mr *MockDomainManagerMockRecorder) GetGuestOSInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGuestOSInfo", reflect.TypeOf((*MockDomainManager)(nil).GetGuestOSInfo))
}

// GetLaunchMeasurement mocks base method.
func (m *MockDomainManager) GetLaunchMeasurement(arg0 *v1.VirtualMachineInstance) (*v1.SEVMeasurementInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLaunchMeasurement", arg0)
	ret0, _ := ret[0].(*v1.SEVMeasurementInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLaunchMeasurement indicates an expected call of GetLaunchMeasurement.
func (mr *MockDomainManagerMockRecorder) GetLaunchMeasurement(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLaunchMeasurement", reflect.TypeOf((*MockDomainManager)(nil).GetLaunchMeasurement), arg0)
}

// GetQemuVersion mocks base method.
func (m *MockDomainManager) GetQemuVersion() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQemuVersion")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQemuVersion indicates an expected call of GetQemuVersion.
func (mr *MockDomainManagerMockRecorder) GetQemuVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQemuVersion", reflect.TypeOf((*MockDomainManager)(nil).GetQemuVersion))
}

// GetSEVInfo mocks base method.
func (m *MockDomainManager) GetSEVInfo() (*v1.SEVPlatformInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSEVInfo")
	ret0, _ := ret[0].(*v1.SEVPlatformInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSEVInfo indicates an expected call of GetSEVInfo.
func (mr *MockDomainManagerMockRecorder) GetSEVInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSEVInfo", reflect.TypeOf((*MockDomainManager)(nil).GetSEVInfo))
}

// GetUsers mocks base method.
func (m *MockDomainManager) GetUsers() []v1.VirtualMachineInstanceGuestOSUser {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers")
	ret0, _ := ret[0].([]v1.VirtualMachineInstanceGuestOSUser)
	return ret0
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockDomainManagerMockRecorder) GetUsers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockDomainManager)(nil).GetUsers))
}

// GuestPing mocks base method.
func (m *MockDomainManager) GuestPing(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GuestPing", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// GuestPing indicates an expected call of GuestPing.
func (mr *MockDomainManagerMockRecorder) GuestPing(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GuestPing", reflect.TypeOf((*MockDomainManager)(nil).GuestPing), arg0)
}

// HotplugHostDevices mocks base method.
func (m *MockDomainManager) HotplugHostDevices(vmi *v1.VirtualMachineInstance) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HotplugHostDevices", vmi)
	ret0, _ := ret[0].(error)
	return ret0
}

// HotplugHostDevices indicates an expected call of HotplugHostDevices.
func (mr *MockDomainManagerMockRecorder) HotplugHostDevices(vmi any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HotplugHostDevices", reflect.TypeOf((*MockDomainManager)(nil).HotplugHostDevices), vmi)
}

// InjectLaunchSecret mocks base method.
func (m *MockDomainManager) InjectLaunchSecret(arg0 *v1.VirtualMachineInstance, arg1 *v1.SEVSecretOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InjectLaunchSecret", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InjectLaunchSecret indicates an expected call of InjectLaunchSecret.
func (mr *MockDomainManagerMockRecorder) InjectLaunchSecret(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InjectLaunchSecret", reflect.TypeOf((*MockDomainManager)(nil).InjectLaunchSecret), arg0, arg1)
}

// InterfacesStatus mocks base method.
func (m *MockDomainManager) InterfacesStatus() []api.InterfaceStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InterfacesStatus")
	ret0, _ := ret[0].([]api.InterfaceStatus)
	return ret0
}

// InterfacesStatus indicates an expected call of InterfacesStatus.
func (mr *MockDomainManagerMockRecorder) InterfacesStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InterfacesStatus", reflect.TypeOf((*MockDomainManager)(nil).InterfacesStatus))
}

// KillVMI mocks base method.
func (m *MockDomainManager) KillVMI(arg0 *v1.VirtualMachineInstance) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KillVMI", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// KillVMI indicates an expected call of KillVMI.
func (mr *MockDomainManagerMockRecorder) KillVMI(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KillVMI", reflect.TypeOf((*MockDomainManager)(nil).KillVMI), arg0)
}

// ListAllDomains mocks base method.
func (m *MockDomainManager) ListAllDomains() ([]*api.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAllDomains")
	ret0, _ := ret[0].([]*api.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAllDomains indicates an expected call of ListAllDomains.
func (mr *MockDomainManagerMockRecorder) ListAllDomains() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAllDomains", reflect.TypeOf((*MockDomainManager)(nil).ListAllDomains))
}

// MarkGracefulShutdownVMI mocks base method.
func (m *MockDomainManager) MarkGracefulShutdownVMI() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MarkGracefulShutdownVMI")
}

// MarkGracefulShutdownVMI indicates an expected call of MarkGracefulShutdownVMI.
func (mr *MockDomainManagerMockRecorder) MarkGracefulShutdownVMI() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkGracefulShutdownVMI", reflect.TypeOf((*MockDomainManager)(nil).MarkGracefulShutdownVMI))
}

// MemoryDump mocks base method.
func (m *MockDomainManager) MemoryDump(vmi *v1.VirtualMachineInstance, dumpPath string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MemoryDump", vmi, dumpPath)
	ret0, _ := ret[0].(error)
	return ret0
}

// MemoryDump indicates an expected call of MemoryDump.
func (mr *MockDomainManagerMockRecorder) MemoryDump(vmi, dumpPath any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MemoryDump", reflect.TypeOf((*MockDomainManager)(nil).MemoryDump), vmi, dumpPath)
}

// MigrateVMI mocks base method.
func (m *MockDomainManager) MigrateVMI(arg0 *v1.VirtualMachineInstance, arg1 *cmdclient.MigrationOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MigrateVMI", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// MigrateVMI indicates an expected call of MigrateVMI.
func (mr *MockDomainManagerMockRecorder) MigrateVMI(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MigrateVMI", reflect.TypeOf((*MockDomainManager)(nil).MigrateVMI), arg0, arg1)
}

// PauseVMI mocks base method.
func (m *MockDomainManager) PauseVMI(arg0 *v1.VirtualMachineInstance) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PauseVMI", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// PauseVMI indicates an expected call of PauseVMI.
func (mr *MockDomainManagerMockRecorder) PauseVMI(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PauseVMI", reflect.TypeOf((*MockDomainManager)(nil).PauseVMI), arg0)
}

// PrepareMigrationTarget mocks base method.
func (m *MockDomainManager) PrepareMigrationTarget(arg0 *v1.VirtualMachineInstance, arg1 bool, arg2 *v10.VirtualMachineOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrepareMigrationTarget", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// PrepareMigrationTarget indicates an expected call of PrepareMigrationTarget.
func (mr *MockDomainManagerMockRecorder) PrepareMigrationTarget(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrepareMigrationTarget", reflect.TypeOf((*MockDomainManager)(nil).PrepareMigrationTarget), arg0, arg1, arg2)
}

// ResetVMI mocks base method.
func (m *MockDomainManager) ResetVMI(arg0 *v1.VirtualMachineInstance) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetVMI", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetVMI indicates an expected call of ResetVMI.
func (mr *MockDomainManagerMockRecorder) ResetVMI(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetVMI", reflect.TypeOf((*MockDomainManager)(nil).ResetVMI), arg0)
}

// SignalShutdownVMI mocks base method.
func (m *MockDomainManager) SignalShutdownVMI(arg0 *v1.VirtualMachineInstance) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignalShutdownVMI", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SignalShutdownVMI indicates an expected call of SignalShutdownVMI.
func (mr *MockDomainManagerMockRecorder) SignalShutdownVMI(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignalShutdownVMI", reflect.TypeOf((*MockDomainManager)(nil).SignalShutdownVMI), arg0)
}

// SoftRebootVMI mocks base method.
func (m *MockDomainManager) SoftRebootVMI(arg0 *v1.VirtualMachineInstance) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SoftRebootVMI", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SoftRebootVMI indicates an expected call of SoftRebootVMI.
func (mr *MockDomainManagerMockRecorder) SoftRebootVMI(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SoftRebootVMI", reflect.TypeOf((*MockDomainManager)(nil).SoftRebootVMI), arg0)
}

// SyncVMI mocks base method.
func (m *MockDomainManager) SyncVMI(arg0 *v1.VirtualMachineInstance, arg1 bool, arg2 *v10.VirtualMachineOptions) (*api.DomainSpec, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SyncVMI", arg0, arg1, arg2)
	ret0, _ := ret[0].(*api.DomainSpec)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SyncVMI indicates an expected call of SyncVMI.
func (mr *MockDomainManagerMockRecorder) SyncVMI(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncVMI", reflect.TypeOf((*MockDomainManager)(nil).SyncVMI), arg0, arg1, arg2)
}

// UnfreezeVMI mocks base method.
func (m *MockDomainManager) UnfreezeVMI(arg0 *v1.VirtualMachineInstance) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnfreezeVMI", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnfreezeVMI indicates an expected call of UnfreezeVMI.
func (mr *MockDomainManagerMockRecorder) UnfreezeVMI(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnfreezeVMI", reflect.TypeOf((*MockDomainManager)(nil).UnfreezeVMI), arg0)
}

// UnpauseVMI mocks base method.
func (m *MockDomainManager) UnpauseVMI(arg0 *v1.VirtualMachineInstance) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnpauseVMI", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnpauseVMI indicates an expected call of UnpauseVMI.
func (mr *MockDomainManagerMockRecorder) UnpauseVMI(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnpauseVMI", reflect.TypeOf((*MockDomainManager)(nil).UnpauseVMI), arg0)
}

// UpdateGuestMemory mocks base method.
func (m *MockDomainManager) UpdateGuestMemory(vmi *v1.VirtualMachineInstance) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGuestMemory", vmi)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateGuestMemory indicates an expected call of UpdateGuestMemory.
func (mr *MockDomainManagerMockRecorder) UpdateGuestMemory(vmi any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGuestMemory", reflect.TypeOf((*MockDomainManager)(nil).UpdateGuestMemory), vmi)
}

// UpdateVCPUs mocks base method.
func (m *MockDomainManager) UpdateVCPUs(vmi *v1.VirtualMachineInstance, options *v10.VirtualMachineOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateVCPUs", vmi, options)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateVCPUs indicates an expected call of UpdateVCPUs.
func (mr *MockDomainManagerMockRecorder) UpdateVCPUs(vmi, options any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateVCPUs", reflect.TypeOf((*MockDomainManager)(nil).UpdateVCPUs), vmi, options)
}
