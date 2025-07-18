/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright The KubeVirt Authors.
 *
 */

package virtconfig

/*
 This module is intended for exposing the virtualization configuration that is available at the cluster-level and its default settings.
*/

import (
	"kubevirt.io/client-go/log"

	k8sv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	v1 "kubevirt.io/api/core/v1"

	"kubevirt.io/kubevirt/pkg/virt-config/featuregate"
)

const (
	ParallelOutboundMigrationsPerNodeDefault uint32 = 2
	ParallelMigrationsPerClusterDefault      uint32 = 5
	BandwidthPerMigrationDefault                    = "0Mi"
	MigrationAllowAutoConverge               bool   = false
	MigrationAllowPostCopy                   bool   = false
	MigrationProgressTimeout                 int64  = 150
	MigrationCompletionTimeoutPerGiB         int64  = 150
	DefaultAMD64MachineType                         = "q35"
	DefaultPPC64LEMachineType                       = "pseries"
	DefaultAARCH64MachineType                       = "virt"
	DefaultS390XMachineType                         = "s390-ccw-virtio"
	DefaultCPURequest                               = "100m"
	DefaultMemoryOvercommit                         = 100
	DefaultAMD64EmulatedMachines                    = "q35*,pc-q35*"
	DefaultPPC64LEEmulatedMachines                  = "pseries*"
	DefaultAARCH64EmulatedMachines                  = "virt*"
	DefaultS390XEmulatedMachines                    = "s390-ccw-virtio*"
	DefaultLessPVCSpaceToleration                   = 10
	DefaultMinimumReservePVCBytes                   = 131072
	DefaultNodeSelectors                            = ""
	DefaultNetworkInterface                         = "bridge"
	DefaultImagePullPolicy                          = k8sv1.PullIfNotPresent
	DefaultAllowEmulation                           = false
	DefaultUnsafeMigrationOverride                  = false
	DefaultPermitSlirpInterface                     = false
	SmbiosConfigDefaultFamily                       = "KubeVirt"
	SmbiosConfigDefaultManufacturer                 = "KubeVirt"
	SmbiosConfigDefaultProduct                      = "None"
	DefaultPermitBridgeInterfaceOnPodNetwork        = true
	DefaultSELinuxLauncherType                      = ""
	SupportedGuestAgentVersions                     = "2.*,3.*,4.*,5.*"
	DefaultARCHOVMFPath                             = "/usr/share/OVMF"
	DefaultAARCH64OVMFPath                          = "/usr/share/AAVMF"
	DefaultS390xOVMFPath                            = ""
	DefaultMemBalloonStatsPeriod             uint32 = 10
	DefaultCPUAllocationRatio                       = 10
	DefaultDiskVerificationMemoryLimitBytes         = 2000 * 1024 * 1024
	DefaultVirtAPILogVerbosity                      = 2
	DefaultVirtControllerLogVerbosity               = 2
	DefaultVirtHandlerLogVerbosity                  = 2
	DefaultVirtLauncherLogVerbosity                 = 2
	DefaultVirtOperatorLogVerbosity                 = 2

	// Default REST configuration settings
	DefaultVirtHandlerQPS         float32 = 5
	DefaultVirtHandlerBurst               = 10
	DefaultVirtControllerQPS      float32 = 200
	DefaultVirtControllerBurst            = 400
	DefaultVirtAPIQPS             float32 = 5
	DefaultVirtAPIBurst                   = 10
	DefaultVirtWebhookClientQPS           = 200
	DefaultVirtWebhookClientBurst         = 400

	DefaultMaxHotplugRatio   = 4
	DefaultVMRolloutStrategy = v1.VMRolloutStrategyLiveUpdate
)

func IsARM64(arch string) bool {
	return arch == "arm64"
}

func (c *ClusterConfig) GetMemBalloonStatsPeriod() uint32 {
	return *c.GetConfig().MemBalloonStatsPeriod
}

func (c *ClusterConfig) AllowEmulation() bool {
	return c.GetConfig().DeveloperConfiguration.UseEmulation
}

func (c *ClusterConfig) GetMigrationConfiguration() *v1.MigrationConfiguration {
	migrationConfig := c.GetConfig().MigrationConfiguration
	// For backward compatibility, AllowWorkloadDisruption will follow the
	// value of AllowPostCopy, if not explicitly set
	if migrationConfig.AllowWorkloadDisruption == nil {
		allowPostCopy := false
		if migrationConfig.AllowPostCopy != nil {
			allowPostCopy = *migrationConfig.AllowPostCopy
		}
		migrationConfig.AllowWorkloadDisruption = &allowPostCopy
	}
	return migrationConfig
}

func (c *ClusterConfig) GetImagePullPolicy() (policy k8sv1.PullPolicy) {
	return c.GetConfig().ImagePullPolicy
}

func (c *ClusterConfig) GetResourceVersion() string {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.lastValidConfigResourceVersion
}

func (c *ClusterConfig) GetMachineType(arch string) string {
	if c.GetConfig().MachineType != "" {
		return c.GetConfig().MachineType
	}

	switch arch {
	case "arm64":
		return c.GetConfig().ArchitectureConfiguration.Arm64.MachineType
	case "ppc64le":
		return c.GetConfig().ArchitectureConfiguration.Ppc64le.MachineType
	case "s390x":
		return c.GetConfig().ArchitectureConfiguration.S390x.MachineType
	default:
		return c.GetConfig().ArchitectureConfiguration.Amd64.MachineType
	}
}

func (c *ClusterConfig) GetCPUModel() string {
	return c.GetConfig().CPUModel
}

func (c *ClusterConfig) GetCPURequest() *resource.Quantity {
	return c.GetConfig().CPURequest
}

func (c *ClusterConfig) GetDiskVerification() *v1.DiskVerification {
	return c.GetConfig().DeveloperConfiguration.DiskVerification
}

func (c *ClusterConfig) GetMemoryOvercommit() int {
	return c.GetConfig().DeveloperConfiguration.MemoryOvercommit
}

func (c *ClusterConfig) GetEmulatedMachines(arch string) []string {
	oldEmulatedMachines := c.GetConfig().EmulatedMachines
	if oldEmulatedMachines != nil {
		return oldEmulatedMachines
	}

	switch arch {
	case "arm64":
		return c.GetConfig().ArchitectureConfiguration.Arm64.EmulatedMachines
	case "ppc64le":
		return c.GetConfig().ArchitectureConfiguration.Ppc64le.EmulatedMachines
	case "s390x":
		return c.GetConfig().ArchitectureConfiguration.S390x.EmulatedMachines
	default:
		return c.GetConfig().ArchitectureConfiguration.Amd64.EmulatedMachines
	}
}

func (c *ClusterConfig) GetLessPVCSpaceToleration() int {
	return c.GetConfig().DeveloperConfiguration.LessPVCSpaceToleration
}

func (c *ClusterConfig) GetMinimumReservePVCBytes() uint64 {
	return c.GetConfig().DeveloperConfiguration.MinimumReservePVCBytes
}

func (c *ClusterConfig) GetNodeSelectors() map[string]string {
	return c.GetConfig().DeveloperConfiguration.NodeSelectors
}

func (c *ClusterConfig) GetDefaultNetworkInterface() string {
	return c.GetConfig().NetworkConfiguration.NetworkInterface
}

func (c *ClusterConfig) GetDefaultArchitecture() string {
	return c.GetConfig().ArchitectureConfiguration.DefaultArchitecture
}

func (c *ClusterConfig) GetSMBIOS() *v1.SMBiosConfiguration {
	return c.GetConfig().SMBIOSConfig
}

func (c *ClusterConfig) IsBridgeInterfaceOnPodNetworkEnabled() bool {
	return *c.GetConfig().NetworkConfiguration.PermitBridgeInterfaceOnPodNetwork
}

func (c *ClusterConfig) GetDefaultClusterConfig() *v1.KubeVirtConfiguration {
	return c.defaultConfig
}

func (c *ClusterConfig) GetSELinuxLauncherType() string {
	return c.GetConfig().SELinuxLauncherType
}

func (c *ClusterConfig) GetDefaultRuntimeClass() string {
	return c.GetConfig().DefaultRuntimeClass
}

func (c *ClusterConfig) GetSupportedAgentVersions() []string {
	return c.GetConfig().SupportedGuestAgentVersions
}

func (c *ClusterConfig) GetOVMFPath(arch string) string {
	oldOvmfPath := c.GetConfig().OVMFPath
	if oldOvmfPath != "" {
		return oldOvmfPath
	}

	switch arch {
	case "arm64":
		return c.GetConfig().ArchitectureConfiguration.Arm64.OVMFPath
	case "ppc64le":
		return c.GetConfig().ArchitectureConfiguration.Ppc64le.OVMFPath
	case "s390x":
		return c.GetConfig().ArchitectureConfiguration.S390x.OVMFPath
	default:
		return c.GetConfig().ArchitectureConfiguration.Amd64.OVMFPath
	}
}

func (c *ClusterConfig) GetCPUAllocationRatio() int {
	return c.GetConfig().DeveloperConfiguration.CPUAllocationRatio
}

func (c *ClusterConfig) GetMinimumClusterTSCFrequency() *int64 {
	return c.GetConfig().DeveloperConfiguration.MinimumClusterTSCFrequency
}

func (c *ClusterConfig) GetPermittedHostDevices() *v1.PermittedHostDevices {
	return c.GetConfig().PermittedHostDevices
}

func (c *ClusterConfig) GetSupportContainerRequest(typeName v1.SupportContainerType, resourceName k8sv1.ResourceName) *resource.Quantity {
	for _, containerResource := range c.GetConfig().SupportContainerResources {
		if containerResource.Type == typeName {
			quantity := containerResource.Resources.Requests[resourceName]
			if !quantity.IsZero() {
				return &quantity
			}
		}
	}
	return nil
}

func (c *ClusterConfig) GetSupportContainerLimit(typeName v1.SupportContainerType, resourceName k8sv1.ResourceName) *resource.Quantity {
	for _, containerResource := range c.GetConfig().SupportContainerResources {
		if containerResource.Type == typeName {
			quantity := containerResource.Resources.Limits[resourceName]
			if !quantity.IsZero() {
				return &quantity
			}
		}
	}
	return nil
}

func canSelectNode(nodeSelector map[string]string, node *k8sv1.Node) bool {
	for key, val := range nodeSelector {
		labelValue, exist := node.Labels[key]
		if !exist || val != labelValue {
			return false
		}
	}
	return true
}

func (c *ClusterConfig) GetDesiredMDEVTypes(node *k8sv1.Node) []string {
	mdevTypesConf := c.GetConfig().MediatedDevicesConfiguration
	if mdevTypesConf == nil {
		return []string{}
	}
	nodeMdevConf := mdevTypesConf.NodeMediatedDeviceTypes
	if nodeMdevConf != nil {
		mdevTypesMap := make(map[string]struct{})
		for _, nodeConfig := range nodeMdevConf {
			if canSelectNode(nodeConfig.NodeSelector, node) {
				types := nodeConfig.MediatedDeviceTypes
				// Handle deprecated spelling
				if len(types) == 0 {
					types = nodeConfig.MediatedDevicesTypes
				}
				for _, mdevType := range types {
					mdevTypesMap[mdevType] = struct{}{}
				}
			}
		}
		if len(mdevTypesMap) != 0 {
			mdevTypesList := []string{}
			for mdevType := range mdevTypesMap {
				mdevTypesList = append(mdevTypesList, mdevType)
			}
			return mdevTypesList
		}
	}
	// Handle deprecated spelling
	if len(mdevTypesConf.MediatedDeviceTypes) == 0 {
		return mdevTypesConf.MediatedDevicesTypes
	}
	return mdevTypesConf.MediatedDeviceTypes
}

type virtComponent int

const (
	virtHandler virtComponent = iota
	virtApi
	virtController
	virtOperator
	virtLauncher
	virtSynchronizationController
)

// Gets the component verbosity. nodeName can be empty, then it's ignored.
func (c *ClusterConfig) getComponentVerbosity(component virtComponent, nodeName string) uint {
	logConf := c.GetConfig().DeveloperConfiguration.LogVerbosity

	if nodeName != "" {
		if level := logConf.NodeVerbosity[nodeName]; level != 0 {
			return level
		}
	}

	switch component {
	case virtHandler:
		return logConf.VirtHandler
	case virtApi:
		return logConf.VirtAPI
	case virtController:
		return logConf.VirtController
	case virtOperator:
		return logConf.VirtOperator
	case virtLauncher:
		return logConf.VirtLauncher
	case virtSynchronizationController:
		return logConf.VirtSynchronizationController
	default:
		log.Log.Errorf("getComponentVerbosity called with an unknown virtComponent: %v", component)
		return 0
	}
}

func (c *ClusterConfig) GetVirtHandlerVerbosity(nodeName string) uint {
	return c.getComponentVerbosity(virtHandler, nodeName)
}

func (c *ClusterConfig) GetVirtAPIVerbosity(nodeName string) uint {
	return c.getComponentVerbosity(virtApi, nodeName)
}

func (c *ClusterConfig) GetVirtControllerVerbosity(nodeName string) uint {
	return c.getComponentVerbosity(virtController, nodeName)
}

func (c *ClusterConfig) GetVirtOperatorVerbosity(nodeName string) uint {
	return c.getComponentVerbosity(virtOperator, nodeName)
}

func (c *ClusterConfig) GetVirtLauncherVerbosity() uint {
	return c.getComponentVerbosity(virtLauncher, "")
}

func (c *ClusterConfig) GetVirtSynchronizationControllerVerbosity() uint {
	return c.getComponentVerbosity(virtSynchronizationController, "")
}

// GetMinCPUModel return minimal cpu which is used in node-labeller
func (c *ClusterConfig) GetMinCPUModel() string {
	return c.GetConfig().MinCPUModel
}

// GetObsoleteCPUModels return slice of obsolete cpus which are used in node-labeller
func (c *ClusterConfig) GetObsoleteCPUModels() map[string]bool {
	return c.GetConfig().ObsoleteCPUModels
}

// GetClusterCPUArch return the CPU architecture in ClusterConfig
func (c *ClusterConfig) GetClusterCPUArch() string {
	return c.cpuArch
}

// GetDeveloperConfigurationUseEmulation return the UseEmulation value in DeveloperConfiguration
func (c *ClusterConfig) GetDeveloperConfigurationUseEmulation() bool {
	config := c.GetConfig()

	if config.DeveloperConfiguration != nil {
		return config.DeveloperConfiguration.UseEmulation
	}

	return false
}

func (c *ClusterConfig) GetVMStateStorageClass() string {
	return c.GetConfig().VMStateStorageClass
}

func (c *ClusterConfig) IsFreePageReportingDisabled() bool {
	return c.GetConfig().VirtualMachineOptions != nil && c.GetConfig().VirtualMachineOptions.DisableFreePageReporting != nil
}

func (c *ClusterConfig) IsSerialConsoleLogDisabled() bool {
	return c.GetConfig().VirtualMachineOptions != nil && c.GetConfig().VirtualMachineOptions.DisableSerialConsoleLog != nil
}

func (c *ClusterConfig) GetKSMConfiguration() *v1.KSMConfiguration {
	return c.GetConfig().KSMConfiguration
}

func (c *ClusterConfig) GetMaximumCpuSockets() (numOfSockets uint32) {
	liveConfig := c.GetConfig().LiveUpdateConfiguration
	if liveConfig != nil && liveConfig.MaxCpuSockets != nil {
		numOfSockets = *liveConfig.MaxCpuSockets
	}

	return
}

func (c *ClusterConfig) GetMaximumGuestMemory() *resource.Quantity {
	liveConfig := c.GetConfig().LiveUpdateConfiguration
	if liveConfig != nil {
		return liveConfig.MaxGuest
	}
	return nil
}

func (c *ClusterConfig) GetMaxHotplugRatio() uint32 {
	liveConfig := c.GetConfig().LiveUpdateConfiguration
	if liveConfig == nil {
		return 1
	}

	return liveConfig.MaxHotplugRatio
}

func (c *ClusterConfig) IsVMRolloutStrategyLiveUpdate() bool {
	liveConfig := c.GetConfig().VMRolloutStrategy
	return liveConfig == nil || *liveConfig == v1.VMRolloutStrategyLiveUpdate
}

func (c *ClusterConfig) GetNetworkBindings() map[string]v1.InterfaceBindingPlugin {
	networkConfig := c.GetConfig().NetworkConfiguration
	if networkConfig != nil {
		return networkConfig.Binding
	}
	return nil
}

func (config *ClusterConfig) VGADisplayForEFIGuestsEnabled() bool {
	VGADisplayForEFIGuestsAnnotationExists := false
	kv := config.GetConfigFromKubeVirtCR()
	if kv != nil {
		_, VGADisplayForEFIGuestsAnnotationExists = kv.Annotations[v1.VGADisplayForEFIGuestsX86Annotation]
	}
	return VGADisplayForEFIGuestsAnnotationExists
}

func (c *ClusterConfig) GetInstancetypeReferencePolicy() v1.InstancetypeReferencePolicy {
	instancetypeConfig := c.GetConfig().Instancetype
	if instancetypeConfig != nil && instancetypeConfig.ReferencePolicy != nil {
		return *instancetypeConfig.ReferencePolicy
	}
	// Default to the Reference InstancetypeReferencePolicy
	return v1.Reference
}

func (c *ClusterConfig) ClusterProfilerEnabled() bool {
	return c.GetConfig().DeveloperConfiguration.ClusterProfiler ||
		c.isFeatureGateDefined(featuregate.ClusterProfiler)
}
