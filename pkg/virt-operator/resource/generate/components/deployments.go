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

package components

import (
	"fmt"
	"path"
	"strings"

	networkv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	virtv1 "kubevirt.io/api/core/v1"

	"kubevirt.io/kubevirt/pkg/pointer"
	"kubevirt.io/kubevirt/pkg/virt-operator/resource/placement"
	operatorutil "kubevirt.io/kubevirt/pkg/virt-operator/util"
)

const (
	nodeLabellerVolumePath = "/var/lib/kubevirt-node-labeller"

	VirtAPIName                       = "virt-api"
	VirtControllerName                = "virt-controller"
	VirtOperatorName                  = "virt-operator"
	VirtExportProxyName               = "virt-exportproxy"
	VirtSynchronizationControllerName = "virt-synchronization-controller"

	kubevirtLabelKey = "kubevirt.io"

	portName = "--port"
)

func NewPrometheusService(namespace string) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      "kubevirt-prometheus-metrics",
			Labels: map[string]string{
				virtv1.AppLabel:    "",
				prometheusLabelKey: prometheusLabelValue,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				prometheusLabelKey: prometheusLabelValue,
			},
			Ports: []corev1.ServicePort{
				{
					Name: "metrics",
					Port: 443,
					TargetPort: intstr.IntOrString{
						Type:   intstr.String,
						StrVal: "metrics",
					},
					Protocol: corev1.ProtocolTCP,
				},
			},
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: corev1.ClusterIPNone,
		},
	}
}

func NewApiServerService(namespace string) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      VirtAPIName,
			Labels: map[string]string{
				virtv1.AppLabel: VirtAPIName,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				virtv1.AppLabel: VirtAPIName,
			},
			Ports: []corev1.ServicePort{
				{
					Port: 443,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 8443,
					},
					Protocol: corev1.ProtocolTCP,
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
}

func NewExportProxyService(namespace string) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      VirtExportProxyName,
			Labels: map[string]string{
				virtv1.AppLabel: VirtExportProxyName,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				virtv1.AppLabel: VirtExportProxyName,
			},
			Ports: []corev1.ServicePort{
				{
					Port: 443,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 8443,
					},
					Protocol: corev1.ProtocolTCP,
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
}

func newPodTemplateSpec(podName, imageName, repository, version, productName, productVersion, productComponent, image string, pullPolicy corev1.PullPolicy, imagePullSecrets []corev1.LocalObjectReference, podAffinity *corev1.Affinity, envVars *[]corev1.EnvVar) *corev1.PodTemplateSpec {

	if image == "" {
		image = fmt.Sprintf("%s/%s%s", repository, imageName, AddVersionSeparatorPrefix(version))
	}

	podTemplateSpec := &corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				virtv1.AppLabel:    podName,
				prometheusLabelKey: prometheusLabelValue,
			},
			Name: podName,
		},
		Spec: corev1.PodSpec{
			PriorityClassName: "kubevirt-cluster-critical",
			Affinity:          podAffinity,
			Tolerations:       criticalAddonsToleration(),
			Containers: []corev1.Container{
				{
					Name:                     podName,
					Image:                    image,
					ImagePullPolicy:          pullPolicy,
					TerminationMessagePolicy: corev1.TerminationMessageFallbackToLogsOnError,
				},
			},
		},
	}

	if len(imagePullSecrets) > 0 {
		podTemplateSpec.Spec.ImagePullSecrets = imagePullSecrets
	}

	if productVersion != "" {
		podTemplateSpec.ObjectMeta.Labels[virtv1.AppVersionLabel] = productVersion
	}

	if productName != "" {
		podTemplateSpec.ObjectMeta.Labels[virtv1.AppPartOfLabel] = productName
	}

	if productComponent != "" {
		podTemplateSpec.ObjectMeta.Labels[virtv1.AppComponentLabel] = productComponent
	}

	if envVars != nil && len(*envVars) != 0 {
		podTemplateSpec.Spec.Containers[0].Env = *envVars
	}

	return podTemplateSpec
}

func attachProfileVolume(spec *corev1.PodSpec) {

	volume := corev1.Volume{
		Name: "profile-data",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}
	volumeMount := corev1.VolumeMount{
		Name:      "profile-data",
		MountPath: "/profile-data",
	}
	spec.Volumes = append(spec.Volumes, volume)
	spec.Containers[0].VolumeMounts = append(spec.Containers[0].VolumeMounts, volumeMount)

}

func attachCertificateSecret(spec *corev1.PodSpec, secretName string, mountPath string) {
	True := true
	secretVolume := corev1.Volume{
		Name: secretName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: secretName,
				Optional:   &True,
			},
		},
	}

	secretVolumeMount := corev1.VolumeMount{
		Name:      secretName,
		ReadOnly:  true,
		MountPath: mountPath,
	}
	spec.Volumes = append(spec.Volumes, secretVolume)
	spec.Containers[0].VolumeMounts = append(spec.Containers[0].VolumeMounts, secretVolumeMount)
}

func newBaseDeployment(deploymentName, imageName, namespace, repository, version, productName, productVersion, productComponent, image string, pullPolicy corev1.PullPolicy, imagePullSecrets []corev1.LocalObjectReference, podAffinity *corev1.Affinity, envVars *[]corev1.EnvVar) *appsv1.Deployment {

	podTemplateSpec := newPodTemplateSpec(deploymentName, imageName, repository, version, productName, productVersion, productComponent, image, pullPolicy, imagePullSecrets, podAffinity, envVars)

	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      deploymentName,
			Labels: map[string]string{
				virtv1.AppLabel:     deploymentName,
				virtv1.AppNameLabel: deploymentName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.P(int32(2)),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					kubevirtLabelKey: deploymentName,
				},
			},
			Template: *podTemplateSpec,
		},
	}

	if productVersion != "" {
		deployment.ObjectMeta.Labels[virtv1.AppVersionLabel] = productVersion
	}

	if productName != "" {
		deployment.ObjectMeta.Labels[virtv1.AppPartOfLabel] = productName
	}

	if productComponent != "" {
		deployment.ObjectMeta.Labels[virtv1.AppComponentLabel] = productComponent
	}

	return deployment
}

func newPodAntiAffinity(key, topologyKey string, operator metav1.LabelSelectorOperator, values []string) *corev1.Affinity {
	return &corev1.Affinity{
		PodAntiAffinity: &corev1.PodAntiAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
				{
					Weight: 1,
					PodAffinityTerm: corev1.PodAffinityTerm{
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									Key:      key,
									Operator: operator,
									Values:   values,
								},
							},
						},
						TopologyKey: topologyKey,
					},
				},
			},
		},
	}
}

func NewApiServerDeployment(namespace, repository, imagePrefix, version, productName, productVersion, productComponent, image string, pullPolicy corev1.PullPolicy, imagePullSecrets []corev1.LocalObjectReference, verbosity string, extraEnv map[string]string) *appsv1.Deployment {
	podAntiAffinity := newPodAntiAffinity(kubevirtLabelKey, corev1.LabelHostname, metav1.LabelSelectorOpIn, []string{VirtAPIName})
	deploymentName := VirtAPIName
	imageName := fmt.Sprintf("%s%s", imagePrefix, deploymentName)
	env := operatorutil.NewEnvVarMap(extraEnv)
	deployment := newBaseDeployment(deploymentName, imageName, namespace, repository, version, productName, productVersion, productComponent, image, pullPolicy, imagePullSecrets, podAntiAffinity, env)

	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = make(map[string]string)
	}
	deployment.Spec.Template.Annotations["openshift.io/required-scc"] = "restricted-v2"

	attachCertificateSecret(&deployment.Spec.Template.Spec, VirtApiCertSecretName, "/etc/virt-api/certificates")
	attachCertificateSecret(&deployment.Spec.Template.Spec, VirtHandlerCertSecretName, "/etc/virt-handler/clientcertificates")
	attachProfileVolume(&deployment.Spec.Template.Spec)

	pod := &deployment.Spec.Template.Spec
	pod.ServiceAccountName = ApiServiceAccountName
	pod.SecurityContext = &corev1.PodSecurityContext{
		RunAsNonRoot:   pointer.P(true),
		SeccompProfile: &corev1.SeccompProfile{Type: corev1.SeccompProfileTypeRuntimeDefault},
	}

	container := &deployment.Spec.Template.Spec.Containers[0]
	container.Command = []string{
		VirtAPIName,
	}
	container.Args = []string{
		portName,
		"8443",
		"--console-server-port",
		"8186",
		"--subresources-only",
		"-v",
		verbosity,
	}
	container.Ports = []corev1.ContainerPort{
		{
			Name:          VirtAPIName,
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: 8443,
		},
		{
			Name:          "metrics",
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: 8443,
		},
	}
	container.ReadinessProbe = &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Scheme: corev1.URISchemeHTTPS,
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 8443,
				},
				Path: path.Join("/apis/subresources.kubevirt.io", virtv1.SubresourceGroupVersions[0].Version, "healthz"),
			},
		},
		InitialDelaySeconds: 15,
		PeriodSeconds:       10,
	}

	container.Resources = corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("5m"),
			corev1.ResourceMemory: resource.MustParse("500Mi"),
		},
	}

	container.SecurityContext = &corev1.SecurityContext{
		AllowPrivilegeEscalation: pointer.P(false),
		Capabilities: &corev1.Capabilities{
			Drop: []corev1.Capability{"ALL"},
		},
		SeccompProfile: &corev1.SeccompProfile{Type: corev1.SeccompProfileTypeRuntimeDefault},
	}
	return deployment
}

func NewControllerDeployment(namespace, repository, imagePrefix, controllerVersion, launcherVersion, exportServerVersion, sidecarVersion, productName, productVersion, productComponent, image, launcherImage, exporterImage, sidecarImage string, pullPolicy corev1.PullPolicy, imagePullSecrets []corev1.LocalObjectReference, verbosity string, extraEnv map[string]string) *appsv1.Deployment {
	podAntiAffinity := newPodAntiAffinity(kubevirtLabelKey, corev1.LabelHostname, metav1.LabelSelectorOpIn, []string{VirtControllerName})
	deploymentName := VirtControllerName
	imageName := fmt.Sprintf("%s%s", imagePrefix, deploymentName)
	env := operatorutil.NewEnvVarMap(extraEnv)
	deployment := newBaseDeployment(deploymentName, imageName, namespace, repository, controllerVersion, productName, productVersion, productComponent, image, pullPolicy, imagePullSecrets, podAntiAffinity, env)

	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = make(map[string]string)
	}
	deployment.Spec.Template.Annotations["openshift.io/required-scc"] = "restricted-v2"

	if launcherImage == "" {
		launcherImage = fmt.Sprintf("%s/%s%s%s", repository, imagePrefix, "virt-launcher", AddVersionSeparatorPrefix(launcherVersion))
	}
	if exporterImage == "" {
		exporterImage = fmt.Sprintf("%s/%s%s%s", repository, imagePrefix, "virt-exportserver", AddVersionSeparatorPrefix(exportServerVersion))
	}

	pod := &deployment.Spec.Template.Spec
	pod.ServiceAccountName = ControllerServiceAccountName
	pod.SecurityContext = &corev1.PodSecurityContext{
		RunAsNonRoot:   pointer.P(true),
		SeccompProfile: &corev1.SeccompProfile{Type: corev1.SeccompProfileTypeRuntimeDefault},
	}

	container := &deployment.Spec.Template.Spec.Containers[0]
	container.Command = []string{
		VirtControllerName,
	}
	container.Args = []string{
		"--launcher-image",
		launcherImage,
		"--exporter-image",
		exporterImage,
		portName,
		"8443",
		"-v",
		verbosity,
	}

	container.Ports = []corev1.ContainerPort{
		{
			Name:          "metrics",
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: 8443,
		},
	}
	container.LivenessProbe = &corev1.Probe{
		FailureThreshold: 8,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Scheme: corev1.URISchemeHTTPS,
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 8443,
				},
				Path: "/healthz",
			},
		},
		InitialDelaySeconds: 15,
		TimeoutSeconds:      10,
	}
	container.ReadinessProbe = &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Scheme: corev1.URISchemeHTTPS,
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 8443,
				},
				Path: "/leader",
			},
		},
		InitialDelaySeconds: 15,
		TimeoutSeconds:      10,
	}

	attachCertificateSecret(pod, VirtControllerCertSecretName, "/etc/virt-controller/certificates")
	attachCertificateSecret(pod, KubeVirtExportCASecretName, "/etc/virt-controller/exportca")
	attachProfileVolume(pod)

	container.Resources = corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("10m"),
			corev1.ResourceMemory: resource.MustParse("275Mi"),
		},
	}

	container.SecurityContext = &corev1.SecurityContext{
		AllowPrivilegeEscalation: pointer.P(false),
		Capabilities: &corev1.Capabilities{
			Drop: []corev1.Capability{"ALL"},
		},
		SeccompProfile: &corev1.SeccompProfile{Type: corev1.SeccompProfileTypeRuntimeDefault},
	}

	if sidecarImage == "" {
		sidecarImage = fmt.Sprintf("%s/%s%s%s", repository, imagePrefix, "sidecar-shim", AddVersionSeparatorPrefix(sidecarVersion))
	}
	container.Env = append(container.Env, corev1.EnvVar{Name: operatorutil.SidecarShimImageEnvName, Value: sidecarImage})

	return deployment
}

// Used for manifest generation only
func NewOperatorDeployment(namespace, repository, imagePrefix, version, verbosity, kubeVirtVersionEnv, runbookURLTemplate, virtApiImageEnv, virtControllerImageEnv, virtHandlerImageEnv, virtLauncherImageEnv, virtExportProxyImageEnv, virtExportServerImageEnv, virtSynchronizationControllerImageEnv, gsImage, prHelperImage, sidecarShimImage,
	image string, pullPolicy corev1.PullPolicy) *appsv1.Deployment {

	const kubernetesOSLinux = "linux"
	podAntiAffinity := newPodAntiAffinity(kubevirtLabelKey, corev1.LabelHostname, metav1.LabelSelectorOpIn, []string{VirtOperatorName})
	version = AddVersionSeparatorPrefix(version)
	if image == "" {
		image = fmt.Sprintf("%s/%s%s%s", repository, imagePrefix, VirtOperatorName, version)
	}

	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      VirtOperatorName,
			Labels: map[string]string{
				virtv1.AppLabel: VirtOperatorName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.P(int32(2)),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					virtv1.AppLabel: VirtOperatorName,
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						virtv1.AppLabel:    VirtOperatorName,
						virtv1.AppName:     VirtOperatorName,
						prometheusLabelKey: prometheusLabelValue,
					},
					Name: VirtOperatorName,
				},
				Spec: corev1.PodSpec{
					PriorityClassName:  "kubevirt-cluster-critical",
					Tolerations:        criticalAddonsToleration(),
					Affinity:           podAntiAffinity,
					ServiceAccountName: "kubevirt-operator",
					NodeSelector: map[string]string{
						corev1.LabelOSStable: kubernetesOSLinux,
					},
					Containers: []corev1.Container{
						{
							Name:            VirtOperatorName,
							Image:           image,
							ImagePullPolicy: pullPolicy,
							Command: []string{
								VirtOperatorName,
							},
							Args: []string{
								portName,
								"8443",
								"-v",
								verbosity,
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "metrics",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 8443,
								},
								{
									Name:          "webhooks",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 8444,
								},
							},
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Scheme: corev1.URISchemeHTTPS,
										Port: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 8443,
										},
										Path: "/metrics",
									},
								},
								InitialDelaySeconds: 5,
								TimeoutSeconds:      10,
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Scheme: corev1.URISchemeHTTPS,
										Port: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 8443,
										},
										Path: "/metrics",
									},
								},
								InitialDelaySeconds: 5,
								TimeoutSeconds:      10,
							},
							Env: []corev1.EnvVar{
								{
									Name:  operatorutil.VirtOperatorImageEnvName,
									Value: image,
								},
								{
									Name: "WATCH_NAMESPACE", // not used yet
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.annotations['olm.targetNamespaces']", // filled by OLM
										},
									},
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("10m"),
									corev1.ResourceMemory: resource.MustParse("450Mi"),
								},
							},
							SecurityContext: &corev1.SecurityContext{
								AllowPrivilegeEscalation: pointer.P(false),
								Capabilities: &corev1.Capabilities{
									Drop: []corev1.Capability{"ALL"},
								},
								SeccompProfile: &corev1.SeccompProfile{Type: corev1.SeccompProfileTypeRuntimeDefault},
							},
							TerminationMessagePolicy: corev1.TerminationMessageFallbackToLogsOnError,
						},
					},
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot:   pointer.P(true),
						SeccompProfile: &corev1.SeccompProfile{Type: corev1.SeccompProfileTypeRuntimeDefault},
					},
				},
			},
		},
	}

	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = make(map[string]string)
	}
	deployment.Spec.Template.Annotations["openshift.io/required-scc"] = "restricted-v2"

	envVars := generateVirtOperatorEnvVars(
		runbookURLTemplate, virtApiImageEnv, virtControllerImageEnv, virtHandlerImageEnv, virtLauncherImageEnv, virtExportProxyImageEnv,
		virtExportServerImageEnv, virtSynchronizationControllerImageEnv, gsImage, prHelperImage, sidecarShimImage, kubeVirtVersionEnv,
	)

	if envVars != nil {
		deployment.Spec.Template.Spec.Containers[0].Env = append(deployment.Spec.Template.Spec.Containers[0].Env, envVars...)
	}

	attachCertificateSecret(&deployment.Spec.Template.Spec, VirtOperatorCertSecretName, "/etc/virt-operator/certificates")
	attachProfileVolume(&deployment.Spec.Template.Spec)
	placement.InjectPlacementMetadata(nil, &deployment.Spec.Template.Spec, placement.RequireControlPlanePreferNonWorker)

	return deployment
}

func NewExportProxyDeployment(namespace, repository, imagePrefix, version, productName, productVersion, productComponent, image string, pullPolicy corev1.PullPolicy, imagePullSecrets []corev1.LocalObjectReference, verbosity string, extraEnv map[string]string) *appsv1.Deployment {
	podAntiAffinity := newPodAntiAffinity(kubevirtLabelKey, corev1.LabelHostname, metav1.LabelSelectorOpIn, []string{VirtAPIName})
	deploymentName := VirtExportProxyName
	imageName := fmt.Sprintf("%s%s", imagePrefix, deploymentName)
	env := operatorutil.NewEnvVarMap(extraEnv)
	deployment := newBaseDeployment(deploymentName, imageName, namespace, repository, version, productName, productVersion, productComponent, image, pullPolicy, imagePullSecrets, podAntiAffinity, env)

	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = make(map[string]string)
	}
	deployment.Spec.Template.Annotations["openshift.io/required-scc"] = "restricted-v2"

	attachCertificateSecret(&deployment.Spec.Template.Spec, VirtExportProxyCertSecretName, "/etc/virt-exportproxy/certificates")
	attachProfileVolume(&deployment.Spec.Template.Spec)

	pod := &deployment.Spec.Template.Spec
	pod.ServiceAccountName = ExportProxyServiceAccountName
	pod.SecurityContext = &corev1.PodSecurityContext{
		RunAsNonRoot: pointer.P(true),
	}

	const shortName = "exportproxy"
	container := &deployment.Spec.Template.Spec.Containers[0]
	// virt-exportproxy too long
	container.Name = shortName
	container.Command = []string{
		VirtExportProxyName,
		portName,
		"8443",
		"-v",
		verbosity,
	}
	container.Ports = []corev1.ContainerPort{
		{
			Name:          shortName,
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: 8443,
		},
		{
			Name:          "metrics",
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: 8443,
		},
	}

	container.ReadinessProbe = &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Scheme: corev1.URISchemeHTTPS,
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 8443,
				},
				Path: "/healthz",
			},
		},
		InitialDelaySeconds: 15,
		PeriodSeconds:       10,
	}

	container.Resources = corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("5m"),
			corev1.ResourceMemory: resource.MustParse("150Mi"),
		},
	}

	return deployment
}

func NewSynchronizationControllerDeployment(
	namespace,
	repository,
	imagePrefix,
	version,
	productName,
	productVersion,
	productComponent,
	image string,
	pullPolicy corev1.PullPolicy,
	imagePullSecrets []corev1.LocalObjectReference,
	migrationNetwork *string,
	syncPort int32,
	verbosity string,
	extraEnv map[string]string) *appsv1.Deployment {

	podAntiAffinity := newPodAntiAffinity(kubevirtLabelKey, corev1.LabelHostname, metav1.LabelSelectorOpIn, []string{VirtSynchronizationControllerName})
	deploymentName := VirtSynchronizationControllerName
	imageName := fmt.Sprintf("%s%s", imagePrefix, deploymentName)

	env := operatorutil.NewEnvVarMap(extraEnv)
	*env = append(*env, corev1.EnvVar{
		Name: "MY_POD_IP",
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: "status.podIP",
			},
		},
	})

	deployment := newBaseDeployment(deploymentName, imageName, namespace, repository, version, productName, productVersion, productComponent, image, pullPolicy, imagePullSecrets, podAntiAffinity, env)

	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = make(map[string]string)
	}
	// remove the prometheus label key, so prometheus doesn't try to scrape anything of the synchronization controller.
	delete(deployment.Spec.Template.Labels, prometheusLabelKey)
	deployment.Spec.Template.Annotations["openshift.io/required-scc"] = "restricted-v2"
	if migrationNetwork != nil {
		// Join the pod to the migration network and name the corresponding interface "migration0"
		deployment.Spec.Template.ObjectMeta.Annotations[networkv1.NetworkAttachmentAnnot] = *migrationNetwork + "@" + virtv1.MigrationInterfaceName
	}

	attachCertificateSecret(&deployment.Spec.Template.Spec, VirtSynchronizationControllerCertSecretName, "/etc/virt-sync-controller/clientcertificates")
	attachCertificateSecret(&deployment.Spec.Template.Spec, VirtSynchronizationControllerServerCertSecretName, "/etc/virt-sync-controller/servercertificates")
	attachProfileVolume(&deployment.Spec.Template.Spec)

	pod := &deployment.Spec.Template.Spec
	pod.ServiceAccountName = SynchronizationControllerServiceAccountName
	pod.SecurityContext = &corev1.PodSecurityContext{
		RunAsNonRoot: pointer.P(true),
	}

	const shortName = "sync"
	container := &deployment.Spec.Template.Spec.Containers[0]
	// synchronization-controller too long
	container.Name = shortName
	container.Command = []string{
		VirtSynchronizationControllerName,
		"--v",
		verbosity,
		"--port",
		fmt.Sprintf("%d", syncPort),
	}
	container.Ports = []corev1.ContainerPort{
		{
			Name:          "metrics",
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: 8443,
		},
		{
			Name:          shortName,
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: syncPort,
		},
	}

	container.ReadinessProbe = &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Scheme: corev1.URISchemeHTTPS,
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 8443,
				},
				Path: "/healthz",
			},
		},
		InitialDelaySeconds: 15,
		PeriodSeconds:       10,
	}

	container.Resources = corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("5m"),
			corev1.ResourceMemory: resource.MustParse("150Mi"),
		},
	}

	return deployment
}

func criticalAddonsToleration() []corev1.Toleration {
	return []corev1.Toleration{
		{
			Key:      "CriticalAddonsOnly",
			Operator: corev1.TolerationOpExists,
		},
	}
}

func AddVersionSeparatorPrefix(version string) string {
	// version can be a template or a tag
	// prefix tags with ":"
	// templates have to deal with the correct image/version separator themselves
	if !strings.HasPrefix(version, "{{if") {
		version = fmt.Sprintf(":%s", version)
	}
	return version
}

func NewPodDisruptionBudgetForDeployment(deployment *appsv1.Deployment) *policyv1.PodDisruptionBudget {
	pdbName := deployment.Name + "-pdb"
	minAvailable := intstr.FromInt(1)
	if deployment.Spec.Replicas != nil {
		minAvailable = intstr.FromInt(int(*deployment.Spec.Replicas - 1))
	}
	selector := deployment.Spec.Selector.DeepCopy()
	podDisruptionBudget := &policyv1.PodDisruptionBudget{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: deployment.Namespace,
			Name:      pdbName,
			Labels: map[string]string{
				virtv1.AppLabel: pdbName,
			},
		},
		Spec: policyv1.PodDisruptionBudgetSpec{
			MinAvailable: &minAvailable,
			Selector:     selector,
		},
	}
	return podDisruptionBudget
}

func generateVirtOperatorEnvVars(runbookURLTemplate, virtApiImageEnv, virtControllerImageEnv, virtHandlerImageEnv, virtLauncherImageEnv, virtExportProxyImageEnv,
	virtExportServerImageEnv, virtSynchronizationControllerImageEnv, gsImage, prHelperImage, sidecarShimImage, kubeVirtVersionEnv string) (envVars []corev1.EnvVar) {

	addEnvVar := func(envVarName, envVarValue string) {
		envVars = append(envVars, corev1.EnvVar{
			Name:  envVarName,
			Value: envVarValue,
		})
	}

	if virtApiImageEnv != "" {
		addEnvVar(operatorutil.VirtApiImageEnvName, virtApiImageEnv)
	}

	if virtControllerImageEnv != "" {
		addEnvVar(operatorutil.VirtControllerImageEnvName, virtControllerImageEnv)
	}

	if virtHandlerImageEnv != "" {
		addEnvVar(operatorutil.VirtHandlerImageEnvName, virtHandlerImageEnv)
	}

	if virtLauncherImageEnv != "" {
		addEnvVar(operatorutil.VirtLauncherImageEnvName, virtLauncherImageEnv)
	}

	if virtExportProxyImageEnv != "" {
		addEnvVar(operatorutil.VirtExportProxyImageEnvName, virtExportProxyImageEnv)
	}

	if virtExportServerImageEnv != "" {
		addEnvVar(operatorutil.VirtExportServerImageEnvName, virtExportServerImageEnv)
	}

	if virtSynchronizationControllerImageEnv != "" {
		addEnvVar(operatorutil.VirtSynchronizationControllerImageEnvName, virtSynchronizationControllerImageEnv)
	}

	if gsImage != "" {
		addEnvVar(operatorutil.GsImageEnvName, gsImage)
	}

	if runbookURLTemplate != "" {
		addEnvVar(operatorutil.RunbookURLTemplate, runbookURLTemplate)
	}
	if prHelperImage != "" {
		addEnvVar(operatorutil.PrHelperImageEnvName, prHelperImage)
	}

	if sidecarShimImage != "" {
		addEnvVar(operatorutil.SidecarShimImageEnvName, sidecarShimImage)
	}

	if kubeVirtVersionEnv != "" {
		addEnvVar(operatorutil.KubeVirtVersionEnvName, kubeVirtVersionEnv)
	}

	return envVars
}
