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
package export

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/uuid"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"go.uber.org/mock/gomock"

	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	k8sv1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"

	"kubevirt.io/kubevirt/pkg/pointer"

	vsv1 "github.com/kubernetes-csi/external-snapshotter/client/v4/apis/volumesnapshot/v1"
	framework "k8s.io/client-go/tools/cache/testing"
	virtv1 "kubevirt.io/api/core/v1"
	exportv1 "kubevirt.io/api/export/v1beta1"
	snapshotv1 "kubevirt.io/api/snapshot/v1beta1"
	"kubevirt.io/client-go/kubecli"
	kubevirtfake "kubevirt.io/client-go/kubevirt/fake"
	cdiv1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"

	apiinstancetype "kubevirt.io/api/instancetype"
	instancetypev1beta1 "kubevirt.io/api/instancetype/v1beta1"

	"kubevirt.io/kubevirt/pkg/certificates/triple"
	"kubevirt.io/kubevirt/pkg/certificates/triple/cert"
	certutil "kubevirt.io/kubevirt/pkg/certificates/triple/cert"
	virtcontroller "kubevirt.io/kubevirt/pkg/controller"
	backendstorage "kubevirt.io/kubevirt/pkg/storage/backend-storage"
	"kubevirt.io/kubevirt/pkg/testutils"
	"kubevirt.io/kubevirt/pkg/virt-controller/services"
	"kubevirt.io/kubevirt/pkg/virt-operator/resource/generate/components"
)

const (
	testNamespace   = "default"
	ingressSecret   = "ingress-secret"
	currentVersion  = "v1beta1"
	vmExportName    = "test"
	labelKey        = "label-key"
	labelValue      = "label-value"
	annotationKey   = "annotation-key"
	annotationValue = "annotation-value"
)

var (
	qemuGid            int64 = 107
	expectedPodEnvVars       = []k8sv1.EnvVar{
		{
			Name:  "EXPORT_VM_DEF_URI",
			Value: manifestsPath,
		}, {
			Name:  "CERT_FILE",
			Value: "/cert/tls.crt",
		}, {
			Name:  "KEY_FILE",
			Value: "/cert/tls.key",
		}, {
			Name:  "TOKEN_FILE",
			Value: "/token/token",
		}}
	tokenSecretName = "my-secret-token"
)

var _ = Describe("Export controller", func() {
	var (
		ctrl                        *gomock.Controller
		controller                  *VMExportController
		recorder                    *record.FakeRecorder
		pvcInformer                 cache.SharedIndexInformer
		podInformer                 cache.SharedIndexInformer
		cmInformer                  cache.SharedIndexInformer
		vmExportInformer            cache.SharedIndexInformer
		vmExportSource              *framework.FakeControllerSource
		serviceInformer             cache.SharedIndexInformer
		dvInformer                  cache.SharedIndexInformer
		vmSnapshotInformer          cache.SharedIndexInformer
		vmSnapshotContentInformer   cache.SharedIndexInformer
		secretInformer              cache.SharedIndexInformer
		vmInformer                  cache.SharedIndexInformer
		vmiInformer                 cache.SharedIndexInformer
		kvInformer                  cache.SharedIndexInformer
		crdInformer                 cache.SharedIndexInformer
		instancetypeInformer        cache.SharedIndexInformer
		clusterInstancetypeInformer cache.SharedIndexInformer
		preferenceInformer          cache.SharedIndexInformer
		clusterPreferenceInformer   cache.SharedIndexInformer
		controllerRevisionInformer  cache.SharedIndexInformer
		rqInformer                  cache.SharedIndexInformer
		nsInformer                  cache.SharedIndexInformer
		k8sClient                   *k8sfake.Clientset
		virtClient                  *kubecli.MockKubevirtClient
		vmExportClient              *kubevirtfake.Clientset
		fakeVolumeSnapshotProvider  *MockVolumeSnapshotProvider
		fakeCertManager             *MockCertManager
		mockVMExportQueue           *testutils.MockWorkQueue[string]
		routeCache                  cache.Store
		ingressCache                cache.Store
		stop                        chan struct{}
	)

	syncCaches := func(stop chan struct{}) {
		go vmExportInformer.Run(stop)
		go pvcInformer.Run(stop)
		go podInformer.Run(stop)
		go dvInformer.Run(stop)
		go cmInformer.Run(stop)
		go serviceInformer.Run(stop)
		go secretInformer.Run(stop)
		go vmSnapshotInformer.Run(stop)
		go vmSnapshotContentInformer.Run(stop)
		go vmInformer.Run(stop)
		go vmiInformer.Run(stop)
		go crdInformer.Run(stop)
		go kvInformer.Run(stop)
		go instancetypeInformer.Run(stop)
		go clusterInstancetypeInformer.Run(stop)
		go preferenceInformer.Run(stop)
		go clusterPreferenceInformer.Run(stop)
		go controllerRevisionInformer.Run(stop)
		go rqInformer.Run(stop)
		go nsInformer.Run(stop)
		Expect(cache.WaitForCacheSync(
			stop,
			vmExportInformer.HasSynced,
			pvcInformer.HasSynced,
			podInformer.HasSynced,
			dvInformer.HasSynced,
			cmInformer.HasSynced,
			serviceInformer.HasSynced,
			secretInformer.HasSynced,
			vmSnapshotInformer.HasSynced,
			vmSnapshotContentInformer.HasSynced,
			vmInformer.HasSynced,
			vmiInformer.HasSynced,
			crdInformer.HasSynced,
			kvInformer.HasSynced,
			instancetypeInformer.HasSynced,
			clusterInstancetypeInformer.HasSynced,
			preferenceInformer.HasSynced,
			clusterPreferenceInformer.HasSynced,
			controllerRevisionInformer.HasSynced,
			rqInformer.HasSynced,
			nsInformer.HasSynced,
		)).To(BeTrue())
	}

	BeforeEach(func() {
		stop = make(chan struct{})
		ctrl = gomock.NewController(GinkgoT())
		virtClient = kubecli.NewMockKubevirtClient(ctrl)
		pvcInformer, _ = testutils.NewFakeInformerFor(&k8sv1.PersistentVolumeClaim{})
		podInformer, _ = testutils.NewFakeInformerFor(&k8sv1.Pod{})
		cmInformer, _ = testutils.NewFakeInformerFor(&k8sv1.ConfigMap{})
		serviceInformer, _ = testutils.NewFakeInformerFor(&k8sv1.Service{})
		vmExportInformer, vmExportSource = testutils.NewFakeInformerWithIndexersFor(&exportv1.VirtualMachineExport{}, virtcontroller.GetVirtualMachineExportInformerIndexers())
		dvInformer, _ = testutils.NewFakeInformerFor(&cdiv1.DataVolume{})
		vmSnapshotInformer, _ = testutils.NewFakeInformerFor(&snapshotv1.VirtualMachineSnapshot{})
		vmSnapshotContentInformer, _ = testutils.NewFakeInformerFor(&snapshotv1.VirtualMachineSnapshotContent{})
		vmInformer, _ = testutils.NewFakeInformerFor(&virtv1.VirtualMachine{})
		vmiInformer, _ = testutils.NewFakeInformerFor(&virtv1.VirtualMachineInstance{})
		routeInformer, _ := testutils.NewFakeInformerFor(&routev1.Route{})
		routeCache = routeInformer.GetStore()
		ingressInformer, _ := testutils.NewFakeInformerFor(&networkingv1.Ingress{})
		ingressCache = ingressInformer.GetStore()
		secretInformer, _ = testutils.NewFakeInformerFor(&k8sv1.Secret{})
		kvInformer, _ = testutils.NewFakeInformerFor(&virtv1.KubeVirt{})
		crdInformer, _ = testutils.NewFakeInformerFor(&extv1.CustomResourceDefinition{})
		instancetypeInformer, _ = testutils.NewFakeInformerFor(&instancetypev1beta1.VirtualMachineInstancetype{})
		clusterInstancetypeInformer, _ = testutils.NewFakeInformerFor(&instancetypev1beta1.VirtualMachineClusterInstancetype{})
		preferenceInformer, _ = testutils.NewFakeInformerFor(&instancetypev1beta1.VirtualMachinePreference{})
		clusterPreferenceInformer, _ = testutils.NewFakeInformerFor(&instancetypev1beta1.VirtualMachineClusterPreference{})
		controllerRevisionInformer, _ = testutils.NewFakeInformerFor(&appsv1.ControllerRevision{})
		rqInformer, _ = testutils.NewFakeInformerFor(&k8sv1.ResourceQuota{})
		nsInformer, _ = testutils.NewFakeInformerFor(&k8sv1.Namespace{})
		fakeVolumeSnapshotProvider = &MockVolumeSnapshotProvider{
			volumeSnapshots: []*vsv1.VolumeSnapshot{},
		}
		fakeCertManager = &MockCertManager{}

		config, _, _ := testutils.NewFakeClusterConfigUsingKVConfig(&virtv1.KubeVirtConfiguration{})
		k8sClient = k8sfake.NewSimpleClientset()
		vmExportClient = kubevirtfake.NewSimpleClientset()
		recorder = record.NewFakeRecorder(100)

		virtClient.EXPECT().CoreV1().Return(k8sClient.CoreV1()).AnyTimes()
		virtClient.EXPECT().VirtualMachineExport(testNamespace).
			Return(vmExportClient.ExportV1beta1().VirtualMachineExports(testNamespace)).AnyTimes()

		controller = &VMExportController{
			Client:                      virtClient,
			Recorder:                    recorder,
			PVCInformer:                 pvcInformer,
			PodInformer:                 podInformer,
			ConfigMapInformer:           cmInformer,
			VMExportInformer:            vmExportInformer,
			ServiceInformer:             serviceInformer,
			DataVolumeInformer:          dvInformer,
			KubevirtNamespace:           "kubevirt",
			ManifestRenderer:            services.NewTemplateService("a", 240, "b", "c", "d", "e", "f", pvcInformer.GetStore(), virtClient, config, qemuGid, "g", rqInformer.GetStore(), nsInformer.GetStore()),
			caCertManager:               fakeCertManager,
			RouteCache:                  routeCache,
			IngressCache:                ingressCache,
			RouteConfigMapInformer:      cmInformer,
			SecretInformer:              secretInformer,
			VMSnapshotInformer:          vmSnapshotInformer,
			VMSnapshotContentInformer:   vmSnapshotContentInformer,
			VolumeSnapshotProvider:      fakeVolumeSnapshotProvider,
			VMInformer:                  vmInformer,
			VMIInformer:                 vmiInformer,
			CRDInformer:                 crdInformer,
			KubeVirtInformer:            kvInformer,
			InstancetypeInformer:        instancetypeInformer,
			ClusterInstancetypeInformer: clusterInstancetypeInformer,
			PreferenceInformer:          preferenceInformer,
			ClusterPreferenceInformer:   clusterPreferenceInformer,
			ControllerRevisionInformer:  controllerRevisionInformer,
		}
		initCert = func(ctrl *VMExportController) {
			ctrl.caCertManager.Start()
			Expect(ctrl.caCertManager.Current()).ToNot(BeNil())
		}

		controller.Init()
		mockVMExportQueue = testutils.NewMockWorkQueue(controller.vmExportQueue)
		controller.vmExportQueue = mockVMExportQueue

		Expect(
			cmInformer.GetStore().Add(&k8sv1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: controller.KubevirtNamespace,
					Name:      components.KubeVirtExportCASecretName,
				},
				Data: map[string]string{
					"ca-bundle": "replace me with ca cert",
				},
			}),
		).To(Succeed())

		Expect(
			kvInformer.GetStore().Add(&virtv1.KubeVirt{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: controller.KubevirtNamespace,
					Name:      "kv",
				},
				Spec: virtv1.KubeVirtSpec{
					CertificateRotationStrategy: virtv1.KubeVirtCertificateRotateStrategy{
						SelfSigned: &virtv1.KubeVirtSelfSignConfiguration{
							CA: &virtv1.CertConfig{
								Duration:    &metav1.Duration{Duration: 24 * time.Hour},
								RenewBefore: &metav1.Duration{Duration: 3 * time.Hour},
							},
							Server: &virtv1.CertConfig{
								Duration:    &metav1.Duration{Duration: 2 * time.Hour},
								RenewBefore: &metav1.Duration{Duration: 1 * time.Hour},
							},
						},
					},
				},
				Status: virtv1.KubeVirtStatus{
					Phase: virtv1.KubeVirtPhaseDeployed,
				},
			}),
		).To(Succeed())
	})

	generateCertFromTime := func(cn string, before, after *time.Time) string {
		defer GinkgoRecover()
		config := certutil.Config{
			CommonName: cn,
			NotBefore:  before,
			NotAfter:   after,
		}
		defer GinkgoRecover()
		caKeyPair, _ := triple.NewCA("kubevirt.io", time.Hour*24*7)

		intermediateKey, err := certutil.NewECDSAPrivateKey()
		Expect(err).ToNot(HaveOccurred())
		intermediateConfig := certutil.Config{
			CommonName: fmt.Sprintf("%s@%d", "intermediate", time.Now().Unix()),
			NotBefore:  before,
			NotAfter:   after,
		}
		intermediateConfig.Usages = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
		intermediateCert, err := certutil.NewSignedCert(intermediateConfig, intermediateKey, caKeyPair.Cert, caKeyPair.Key, time.Hour)
		Expect(err).ToNot(HaveOccurred())

		key, err := certutil.NewECDSAPrivateKey()
		Expect(err).ToNot(HaveOccurred())

		config.AltNames.DNSNames = []string{"hahaha.wwoo", "*.apps-crc.testing", "fgdgd.dfsgdf"}
		config.Usages = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
		cert, err := certutil.NewSignedCert(config, key, intermediateCert, intermediateKey, time.Hour)
		Expect(err).ToNot(HaveOccurred())
		pemOut := strings.Builder{}
		pem.Encode(&pemOut, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
		pem.Encode(&pemOut, &pem.Block{Type: "CERTIFICATE", Bytes: intermediateCert.Raw})
		pem.Encode(&pemOut, &pem.Block{Type: "CERTIFICATE", Bytes: caKeyPair.Cert.Raw})
		return strings.TrimSpace(pemOut.String())
	}

	generateFutureCert := func() string {
		before := time.Now().AddDate(0, 1, 0)
		after := before.AddDate(0, 0, 7)
		return generateCertFromTime("future cert", &before, &after)
	}

	generateExpiredCert := func() string {
		before := time.Now().AddDate(0, -1, 0)
		after := before.AddDate(0, 0, 7)
		return generateCertFromTime("expired cert", &before, &after)
	}

	generateExpectedCert := func() string {
		before := time.Now()
		after := before.AddDate(0, 0, 7)
		return generateCertFromTime("working cert", &before, &after)
	}

	var expectedFuturePemAll = generateFutureCert()
	var expectedExpiredPemAll = generateExpiredCert()

	generateExpectedPem := func(allPem string) string {
		now := time.Now()
		pemOut := strings.Builder{}
		certs, err := certutil.ParseCertsPEM([]byte(allPem))
		Expect(err).ToNot(HaveOccurred())
		for _, cert := range certs {
			if cert.NotAfter.After(now) && cert.NotBefore.Before(now) {
				pem.Encode(&pemOut, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
			}
		}
		return strings.TrimSpace(pemOut.String())
	}

	var expectedPem = generateExpectedCert()
	var expectedFuturePem = generateExpectedPem(expectedFuturePemAll)
	var expectedExpiredPem = generateExpectedPem(expectedExpiredPemAll)

	generateOverlappingCert := func() string {
		before := time.Now().AddDate(0, 0, -3)
		after := before.AddDate(0, 0, 7)
		return generateCertFromTime("overlapping cert", &before, &after) + "\n" + expectedPem
	}

	generateRouteCert := func() string {
		key, err := certutil.NewECDSAPrivateKey()
		Expect(err).ToNot(HaveOccurred())

		config := certutil.Config{
			CommonName: "something else",
		}

		cert, err := certutil.NewSelfSignedCACert(config, key, time.Hour)
		Expect(err).ToNot(HaveOccurred())
		pemOut := strings.Builder{}
		Expect(pem.Encode(&pemOut, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})).To(Succeed())
		return strings.TrimSpace(pemOut.String()) + "\n" + expectedPem
	}

	createRouteConfigMapFromString := func(ca string) *k8sv1.ConfigMap {
		return &k8sv1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      routeCAConfigMapName,
				Namespace: controller.KubevirtNamespace,
			},
			Data: map[string]string{
				routeCaKey: ca,
			},
		}
	}

	createRouteConfigMapFromFunc := func(certFunc func() string) *k8sv1.ConfigMap {
		return createRouteConfigMapFromString(certFunc())
	}

	createFutureRouteConfigMap := func() *k8sv1.ConfigMap {
		return createRouteConfigMapFromString(expectedFuturePemAll)
	}

	createExpiredRouteConfigMap := func() *k8sv1.ConfigMap {
		return createRouteConfigMapFromString(expectedExpiredPem)
	}

	createOverlappingRouteConfigMap := func() *k8sv1.ConfigMap {
		return createRouteConfigMapFromFunc(generateOverlappingCert)
	}

	createRouteConfigMap := func() *k8sv1.ConfigMap {
		return createRouteConfigMapFromFunc(generateRouteCert)
	}

	validIngressDefaultBackend := func(serviceName string) *networkingv1.Ingress {
		return &networkingv1.Ingress{
			Spec: networkingv1.IngressSpec{
				DefaultBackend: &networkingv1.IngressBackend{
					Service: &networkingv1.IngressServiceBackend{
						Name: serviceName,
					},
				},
				Rules: []networkingv1.IngressRule{
					{
						Host: "backend-host",
					},
				},
			},
		}
	}

	ingressDefaultBackendNoService := func() *networkingv1.Ingress {
		return &networkingv1.Ingress{
			Spec: networkingv1.IngressSpec{
				DefaultBackend: &networkingv1.IngressBackend{},
				Rules: []networkingv1.IngressRule{
					{
						Host: "backend-host",
					},
				},
			},
		}
	}

	validIngressRules := func(serviceName string) *networkingv1.Ingress {
		return &networkingv1.Ingress{
			Spec: networkingv1.IngressSpec{
				Rules: []networkingv1.IngressRule{
					{
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: serviceName,
											},
										},
									},
								},
							},
						},
						Host: "rule-host",
					},
				},
			},
		}
	}

	ingressRulesNoBackend := func() *networkingv1.Ingress {
		return &networkingv1.Ingress{
			Spec: networkingv1.IngressSpec{
				Rules: []networkingv1.IngressRule{
					{
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										Backend: networkingv1.IngressBackend{},
									},
								},
							},
						},
						Host: "rule-host",
					},
				},
			},
		}
	}

	It("should add vmexport to queue if matching PVC is added", func() {
		vmExport := createPVCVMExport()
		pvc := &k8sv1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testPVCName,
				Namespace: testNamespace,
			},
		}
		syncCaches(stop)
		mockVMExportQueue.ExpectAdds(1)
		vmExportInformer.GetStore().Add(vmExport)
		controller.handlePVC(pvc)
		mockVMExportQueue.Wait()
	})

	It("should add vmexport to queue if matching VMSnapshot is added", func() {
		vmExport := createSnapshotVMExport()
		vmSnapshot := &snapshotv1.VirtualMachineSnapshot{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testVmsnapshotName,
				Namespace: testNamespace,
			},
			Status: &snapshotv1.VirtualMachineSnapshotStatus{
				ReadyToUse: pointer.P(false),
			},
		}
		syncCaches(stop)
		mockVMExportQueue.ExpectAdds(1)
		vmExportInformer.GetStore().Add(vmExport)
		controller.handleVMSnapshot(vmSnapshot)
		mockVMExportQueue.Wait()
	})

	It("should add vmexport to queue if matching VM is added", func() {
		vmExport := createVMVMExport()
		vm := &virtv1.VirtualMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testVmName,
				Namespace: testNamespace,
			},
			Spec: virtv1.VirtualMachineSpec{
				Template: &virtv1.VirtualMachineInstanceTemplateSpec{
					Spec: virtv1.VirtualMachineInstanceSpec{
						Volumes: []virtv1.Volume{},
					},
				},
			},
		}
		syncCaches(stop)
		mockVMExportQueue.ExpectAdds(1)
		vmExportInformer.GetStore().Add(vmExport)
		controller.handleVM(vm)
		mockVMExportQueue.Wait()
	})

	It("should add vmexport to queue if matching VMI is added", func() {
		vmExport := createVMVMExport()
		vm := &virtv1.VirtualMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testVmName,
				Namespace: testNamespace,
			},
			Spec: virtv1.VirtualMachineSpec{
				Template: &virtv1.VirtualMachineInstanceTemplateSpec{
					Spec: virtv1.VirtualMachineInstanceSpec{
						Volumes: []virtv1.Volume{},
					},
				},
			},
		}
		vmi := &virtv1.VirtualMachineInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testVmName,
				Namespace: testNamespace,
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion: virtv1.GroupVersion.String(),
						Kind:       "VirtualMachine",
						Name:       testVmName,
						Controller: pointer.P(true),
					},
				},
			},
		}
		syncCaches(stop)
		mockVMExportQueue.ExpectAdds(1)
		vmExportInformer.GetStore().Add(vmExport)
		vmInformer.GetStore().Add(vm)
		controller.handleVMI(vmi)
		mockVMExportQueue.Wait()
	})

	It("should NOT add vmexport to queue if VMI is added without matching VM", func() {
		vmExport := createVMVMExport()
		vm := &virtv1.VirtualMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testVmName + "-other",
				Namespace: testNamespace,
			},
			Spec: virtv1.VirtualMachineSpec{
				Template: &virtv1.VirtualMachineInstanceTemplateSpec{
					Spec: virtv1.VirtualMachineInstanceSpec{
						Volumes: []virtv1.Volume{},
					},
				},
			},
		}
		vmi := &virtv1.VirtualMachineInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testVmName,
				Namespace: testNamespace,
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion: virtv1.GroupVersion.String(),
						Kind:       "VirtualMachine",
						Name:       testVmName,
						Controller: pointer.P(true),
					},
				},
			},
		}
		syncCaches(stop)
		mockVMExportQueue.ExpectAdds(1)
		vmExportSource.Add(vmExport)
		vmInformer.GetStore().Add(vm)
		controller.handleVMI(vmi)
		mockVMExportQueue.Wait()
	})

	It("should NOT add vmexport to queue if VMI is added without owner", func() {
		vmExport := createVMVMExport()
		vmi := &virtv1.VirtualMachineInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testVmName,
				Namespace: testNamespace,
			},
		}
		syncCaches(stop)
		mockVMExportQueue.ExpectAdds(1)
		vmExportSource.Add(vmExport)
		controller.handleVMI(vmi)
		mockVMExportQueue.Wait()
	})

	DescribeTable("should add vmexport to queue if VMI (pvc) is added that matches PVC export", func(source virtv1.VolumeSource) {
		vmExport := createPVCVMExport()
		vmi := &virtv1.VirtualMachineInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testVmName,
				Namespace: testNamespace,
			},
			Spec: virtv1.VirtualMachineInstanceSpec{
				Volumes: []virtv1.Volume{
					{
						Name:         "testVolume",
						VolumeSource: source,
					},
				},
			},
		}
		pvc := &k8sv1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testPVCName,
				Namespace: testNamespace,
			},
		}
		syncCaches(stop)
		mockVMExportQueue.ExpectAdds(2)
		vmExportSource.Add(vmExport)
		controller.processVMExportWorkItem()
		pvcInformer.GetStore().Add(pvc)
		vmiInformer.GetStore().Add(vmi)
		controller.handleVMI(vmi)
		mockVMExportQueue.Wait()
	},
		Entry("PVC", virtv1.VolumeSource{
			PersistentVolumeClaim: &virtv1.PersistentVolumeClaimVolumeSource{
				PersistentVolumeClaimVolumeSource: k8sv1.PersistentVolumeClaimVolumeSource{
					ClaimName: testPVCName,
				},
			},
		}),
		Entry("DV", virtv1.VolumeSource{
			DataVolume: &virtv1.DataVolumeSource{
				Name: testPVCName,
			},
		}),
	)

	It("should not add vmexport to queue if VMI (dv) is added that doesn't match a PVC export", func() {
		vmExport := createPVCVMExport()
		vmi := &virtv1.VirtualMachineInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testVmName,
				Namespace: testNamespace,
			},
			Spec: virtv1.VirtualMachineInstanceSpec{
				Volumes: []virtv1.Volume{
					{
						Name: "testVolume",
						VolumeSource: virtv1.VolumeSource{
							DataVolume: &virtv1.DataVolumeSource{
								Name: testPVCName,
							},
						},
					},
				},
			},
		}
		syncCaches(stop)
		mockVMExportQueue.ExpectAdds(1)
		vmExportSource.Add(vmExport)
		controller.processVMExportWorkItem()
		vmiInformer.GetStore().Add(vmi)
		controller.handleVMI(vmi)
		mockVMExportQueue.Wait()
	})

	It("Should create a service based on the name of the VMExport", func() {
		var service *k8sv1.Service
		testVMExport := createPVCVMExport()
		k8sClient.Fake.PrependReactor("create", "services", func(action testing.Action) (handled bool, obj runtime.Object, err error) {
			create, ok := action.(testing.CreateAction)
			Expect(ok).To(BeTrue())
			service, ok = create.GetObject().(*k8sv1.Service)
			Expect(ok).To(BeTrue())
			service.Status.Conditions = make([]metav1.Condition, 1)
			service.Status.Conditions[0].Type = "test"
			Expect(service.GetName()).To(Equal("virt-export-test"))
			Expect(service.GetNamespace()).To(Equal(testNamespace))
			Expect(service.Labels).To(And(
				HaveKeyWithValue(virtv1.AppLabel, "virt-exporter"),
				HaveKeyWithValue(labelKey, labelValue)))
			Expect(service.Annotations).To(HaveKeyWithValue(annotationKey, annotationValue))
			return true, service, nil
		})

		service, err := controller.getOrCreateExportService(testVMExport)
		Expect(err).ToNot(HaveOccurred())
		Expect(service).ToNot(BeNil())
		Expect(service.Status.Conditions[0].Type).To(Equal("test"))

		Expect(
			serviceInformer.GetStore().Add(&k8sv1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      controller.getExportServiceName(testVMExport),
					Namespace: testVMExport.Namespace,
				},
				Spec: k8sv1.ServiceSpec{},
				Status: k8sv1.ServiceStatus{
					Conditions: []metav1.Condition{
						{
							Type: "test2",
						},
					},
				},
			}),
		).To(Succeed())
		service, err = controller.getOrCreateExportService(testVMExport)
		Expect(err).ToNot(HaveOccurred())
		Expect(service).ToNot(BeNil())
		Expect(service.Status.Conditions[0].Type).To(Equal("test2"))
	})

	populateVmExportVM := func() *exportv1.VirtualMachineExport {
		testVMExport := createVMVMExport()
		vm := &virtv1.VirtualMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testVmName,
				Namespace: testNamespace,
			},
			Spec: virtv1.VirtualMachineSpec{
				Template: &virtv1.VirtualMachineInstanceTemplateSpec{
					Spec: virtv1.VirtualMachineInstanceSpec{
						Volumes: []virtv1.Volume{},
					},
				},
			},
		}
		vmInformer.GetStore().Add(vm)
		return testVMExport
	}

	populateVmExportVMSnapshot := func() *exportv1.VirtualMachineExport {
		testVMExport := createSnapshotVMExport()
		vm := &virtv1.VirtualMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testVmName,
				Namespace: testNamespace,
			},
			Spec: virtv1.VirtualMachineSpec{
				Template: &virtv1.VirtualMachineInstanceTemplateSpec{
					Spec: virtv1.VirtualMachineInstanceSpec{
						Volumes: []virtv1.Volume{},
					},
				},
			},
		}
		vmInformer.GetStore().Add(vm)
		snapshot := &snapshotv1.VirtualMachineSnapshot{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testVmsnapshotName,
				Namespace: testNamespace,
			},
			Spec: snapshotv1.VirtualMachineSnapshotSpec{
				Source: k8sv1.TypedLocalObjectReference{
					APIGroup: &virtv1.SchemeGroupVersion.Group,
					Kind:     "VirtualMachine",
					Name:     testVmName,
				},
			},
		}
		vmSnapshotInformer.GetStore().Add(snapshot)
		return testVMExport
	}

	DescribeTable("Should create a pod based on the name of the VMExport", func(populateExportFunc func() *exportv1.VirtualMachineExport, getPVCFromSource pvcFromSourceFunc, numberOfVolumes int) {
		testPVC := &k8sv1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testPVCName,
				Namespace: testNamespace,
			},
			Spec: k8sv1.PersistentVolumeClaimSpec{
				VolumeMode: (*k8sv1.PersistentVolumeMode)(pointer.P(string(k8sv1.PersistentVolumeBlock))),
			},
		}
		testVMExport := populateExportFunc()
		populateInitialVMExportStatus(testVMExport)
		err := controller.handleVMExportToken(testVMExport, getPVCFromSource)
		Expect(testVMExport.Status.TokenSecretRef).ToNot(BeNil())
		Expect(err).ToNot(HaveOccurred())
		k8sClient.Fake.PrependReactor("create", "pods", func(action testing.Action) (handled bool, obj runtime.Object, err error) {
			create, ok := action.(testing.CreateAction)
			Expect(ok).To(BeTrue())
			pod, ok := create.GetObject().(*k8sv1.Pod)
			Expect(ok).To(BeTrue())
			Expect(pod.GetName()).To(Equal(controller.getExportPodName(testVMExport)))
			Expect(len(pod.GetName())).To(BeNumerically("<=", validation.DNS1035LabelMaxLength))
			Expect(pod.GetNamespace()).To(Equal(testNamespace))
			return true, pod, nil
		})
		var service *k8sv1.Service
		k8sClient.Fake.PrependReactor("create", "services", func(action testing.Action) (handled bool, obj runtime.Object, err error) {
			create, ok := action.(testing.CreateAction)
			Expect(ok).To(BeTrue())
			service, ok = create.GetObject().(*k8sv1.Service)
			Expect(ok).To(BeTrue())
			service.Status.Conditions = make([]metav1.Condition, 1)
			service.Status.Conditions[0].Type = "test"
			Expect(service.GetName()).To(Equal(controller.getExportServiceName(testVMExport)))
			Expect(service.GetNamespace()).To(Equal(testNamespace))
			return true, service, nil
		})
		service, err = controller.getOrCreateExportService(testVMExport)
		Expect(err).ToNot(HaveOccurred())
		pod, err := controller.createExporterPod(testVMExport, service, []*k8sv1.PersistentVolumeClaim{testPVC})
		Expect(err).ToNot(HaveOccurred())
		Expect(pod).ToNot(BeNil())
		Expect(pod.Name).To(Equal(controller.getExportPodName(testVMExport)))
		Expect(pod.Spec.Volumes).To(HaveLen(numberOfVolumes), "There should be 3/4 volumes, one pvc, and two secrets (token and certs) (and vm def manifest if VM)")
		certSecretName := ""
		for _, volume := range pod.Spec.Volumes {
			if volume.Name == certificatesVolName {
				certSecretName = volume.Secret.SecretName
			}
		}
		Expect(pod.Spec.Volumes).To(ContainElement(k8sv1.Volume{
			Name: testPVCName,
			VolumeSource: k8sv1.VolumeSource{
				PersistentVolumeClaim: &k8sv1.PersistentVolumeClaimVolumeSource{
					ClaimName: testPVCName,
				},
			},
		}))
		Expect(pod.Spec.Volumes).To(ContainElement(k8sv1.Volume{
			Name: certificatesVolName,
			VolumeSource: k8sv1.VolumeSource{
				Secret: &k8sv1.SecretVolumeSource{
					SecretName: certSecretName,
				},
			},
		}))
		Expect(pod.Spec.Volumes).To(ContainElement(k8sv1.Volume{
			Name: tokenVolName,
			VolumeSource: k8sv1.VolumeSource{
				Secret: &k8sv1.SecretVolumeSource{
					SecretName: tokenSecretName,
				},
			},
		}))
		Expect(pod.Spec.Containers).To(HaveLen(1))
		Expect(pod.Spec.Containers[0].Name).To(Equal("exporter"))
		Expect(pod.Spec.Containers[0].VolumeMounts).To(HaveLen(numberOfVolumes - 1)) // The other is a block Volume
		Expect(pod.Spec.Containers[0].VolumeMounts).To(ContainElement(k8sv1.VolumeMount{
			Name:      certificatesVolName,
			MountPath: "/cert",
		}))
		Expect(pod.Spec.Containers[0].VolumeMounts).To(ContainElement(k8sv1.VolumeMount{
			Name:      tokenVolName,
			MountPath: "/token",
		}))
		Expect(pod.Spec.Containers[0].VolumeDevices).To(HaveLen(1))
		Expect(pod.Spec.Containers[0].VolumeDevices).To(ContainElement(k8sv1.VolumeDevice{
			Name:       testPVC.Name,
			DevicePath: fmt.Sprintf("%s/%s", blockVolumeMountPath, testPVC.Name),
		}))
		Expect(pod.Labels).To(And(
			HaveKeyWithValue(exportServiceLabel, controller.getExportLabelValue(testVMExport)),
			HaveKeyWithValue(labelKey, labelValue)))
		Expect(pod.Annotations).To(And(
			HaveKeyWithValue(annCertParams, fmt.Sprintf("{\"Duration\":%d,\"RenewBefore\":%d}",
				metav1.Duration{Duration: 2 * time.Hour}.Nanoseconds(),
				metav1.Duration{Duration: 1 * time.Hour}.Nanoseconds())),
			HaveKeyWithValue(annotationKey, annotationValue)))
		Expect(pod.Spec.Containers[0].Env).To(ContainElements(expectedPodEnvVars))
		Expect(pod.Spec.Containers[0].Resources.Requests.Cpu()).ToNot(BeNil())
		Expect(pod.Spec.Containers[0].Resources.Requests.Cpu().MilliValue()).To(Equal(int64(100)))
		Expect(pod.Spec.Containers[0].Resources.Requests.Memory()).ToNot(BeNil())
		Expect(pod.Spec.Containers[0].Resources.Requests.Memory().Value()).To(Equal(int64(209715200)))
		Expect(pod.Spec.Containers[0].Resources.Limits.Cpu()).ToNot(BeNil())
		Expect(pod.Spec.Containers[0].Resources.Limits.Cpu().MilliValue()).To(Equal(int64(1000)))
		Expect(pod.Spec.Containers[0].Resources.Limits.Memory()).ToNot(BeNil())
		Expect(pod.Spec.Containers[0].Resources.Limits.Memory().Value()).To(Equal(int64(1073741824)))
		Expect(pod.Spec.Containers[0].ReadinessProbe).ToNot(BeNil())
		Expect(pod.Spec.Containers[0].ReadinessProbe.ProbeHandler.HTTPGet.Path).To(Equal(ReadinessPath))
	},
		Entry("PVC", createPVCVMExport, controller.getPVCFromSourcePVC, 3),
		Entry("PVC, with long name export", createPVCVMExportLongName, controller.getPVCFromSourcePVC, 3),
		Entry("VM", populateVmExportVM, controller.getPVCFromSourceVM, 4),
		Entry("Snapshot", populateVmExportVMSnapshot, controller.getPVCFromSourceVMSnapshot, 4),
	)

	DescribeTable("Volumemount names should be trimmed depending on the PVC name", func(pvcName string) {
		testVMExport := createPVCVMExportWithName(pvcName)
		testPVC := &k8sv1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pvcName,
				Namespace: testNamespace,
			},
			Spec: k8sv1.PersistentVolumeClaimSpec{
				VolumeMode: (*k8sv1.PersistentVolumeMode)(pointer.P(string(k8sv1.PersistentVolumeBlock))),
			},
		}
		populateInitialVMExportStatus(testVMExport)
		err := controller.handleVMExportToken(testVMExport, controller.getPVCFromSourceVMSnapshot)
		Expect(testVMExport.Status.TokenSecretRef).ToNot(BeNil())
		Expect(err).ToNot(HaveOccurred())
		k8sClient.Fake.PrependReactor("create", "pods", func(action testing.Action) (handled bool, obj runtime.Object, err error) {
			create, ok := action.(testing.CreateAction)
			Expect(ok).To(BeTrue())
			pod, ok := create.GetObject().(*k8sv1.Pod)
			Expect(ok).To(BeTrue())
			Expect(pod.GetName()).To(Equal(controller.getExportPodName(testVMExport)))
			Expect(len(pod.GetName())).To(BeNumerically("<=", validation.DNS1035LabelMaxLength))
			Expect(pod.GetNamespace()).To(Equal(testNamespace))
			return true, pod, nil
		})
		var service *k8sv1.Service
		k8sClient.Fake.PrependReactor("create", "services", func(action testing.Action) (handled bool, obj runtime.Object, err error) {
			create, ok := action.(testing.CreateAction)
			Expect(ok).To(BeTrue())
			service, ok = create.GetObject().(*k8sv1.Service)
			Expect(ok).To(BeTrue())
			service.Status.Conditions = make([]metav1.Condition, 1)
			service.Status.Conditions[0].Type = "test"
			Expect(service.GetName()).To(Equal(controller.getExportServiceName(testVMExport)))
			Expect(service.GetNamespace()).To(Equal(testNamespace))
			return true, service, nil
		})
		service, err = controller.getOrCreateExportService(testVMExport)
		Expect(err).ToNot(HaveOccurred())
		pod, err := controller.createExporterPod(testVMExport, service, []*k8sv1.PersistentVolumeClaim{testPVC})
		Expect(err).ToNot(HaveOccurred())
		Expect(pod).ToNot(BeNil())
		Expect(pod.Spec.Containers).To(HaveLen(1))
		Expect(pod.Spec.Containers[0].VolumeDevices).To(HaveLen(1))
		Expect(pod.Spec.Containers[0].VolumeDevices).To(ContainElement(k8sv1.VolumeDevice{
			Name:       controller.getExportPodVolumeName(testPVC),
			DevicePath: fmt.Sprintf("%s/%s", blockVolumeMountPath, controller.getExportPodVolumeName(testPVC)),
		}))
		if len(pvcName) > validation.DNS1035LabelMaxLength {
			Expect(len(pod.Spec.Containers[0].VolumeDevices[0].Name)).To(BeNumerically("<", 63))
		} else {
			Expect(pod.Spec.Containers[0].VolumeDevices[0].Name).To(Equal(pvcName))
		}
	},
		Entry("PVC name within limit", "pvc-name-within-limit"),
		Entry("PVC name exceeding limit", strings.Repeat("a", validation.DNS1035LabelMaxLength+1)),
		Entry("PVC name with same length as limit", strings.Repeat("a", validation.DNS1035LabelMaxLength)),
	)

	DescribeTable("service name should be sanitized", func(exportName, expectedServiceName string) {
		var service *k8sv1.Service
		testVMExport := createPVCVMExport()
		testVMExport.Name = exportName
		k8sClient.Fake.PrependReactor("create", "services", func(action testing.Action) (handled bool, obj runtime.Object, err error) {
			create, ok := action.(testing.CreateAction)
			Expect(ok).To(BeTrue())
			service, ok = create.GetObject().(*k8sv1.Service)
			Expect(ok).To(BeTrue())
			service.Status.Conditions = make([]metav1.Condition, 1)
			service.Status.Conditions[0].Type = "test"
			Expect(service.GetName()).To(Equal(expectedServiceName))
			Expect(service.GetName()).ToNot(ContainSubstring("."))
			Expect(service.GetName()).ToNot(HavePrefix("-"))
			Expect(service.GetNamespace()).To(Equal(testNamespace))
			Expect(service.Labels).To(And(
				HaveKeyWithValue(virtv1.AppLabel, "virt-exporter"),
				HaveKeyWithValue(labelKey, labelValue)))
			Expect(service.Annotations).To(HaveKeyWithValue(annotationKey, annotationValue))
			return true, service, nil
		})

		service, err := controller.getOrCreateExportService(testVMExport)
		Expect(err).ToNot(HaveOccurred())
		Expect(service).ToNot(BeNil())
		Expect(service.Status.Conditions[0].Type).To(Equal("test"))
	},
		Entry("with dots", fmt.Sprintf("%s.com", vmExportName), "virt-export-test-com"),
		Entry("special case with - prefix", strings.Repeat("a", validation.DNS1035LabelMaxLength-10), "v-efa2f187-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
	)

	It("Should create a secret based on the vm export", func() {
		cp := &CertParams{Duration: 24 * time.Hour, RenewBefore: 2 * time.Hour}
		scp, err := serializeCertParams(cp)
		Expect(err).ToNot(HaveOccurred())
		testVMExport := createPVCVMExport()
		populateInitialVMExportStatus(testVMExport)
		err = controller.handleVMExportToken(testVMExport, controller.getPVCFromSourcePVC)
		Expect(err).ToNot(HaveOccurred())
		testExportPod := &k8sv1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-export-pod",
				Annotations: map[string]string{
					annCertParams: scp,
				},
			},
			Spec: k8sv1.PodSpec{
				Volumes: []k8sv1.Volume{
					{
						Name: certificatesVolName,
						VolumeSource: k8sv1.VolumeSource{
							Secret: &k8sv1.SecretVolumeSource{
								SecretName: "test-secret",
							},
						},
					},
				},
			},
		}
		k8sClient.Fake.PrependReactor("create", "secrets", func(action testing.Action) (handled bool, obj runtime.Object, err error) {
			create, ok := action.(testing.CreateAction)
			Expect(ok).To(BeTrue())
			secret, ok := create.GetObject().(*k8sv1.Secret)
			Expect(ok).To(BeTrue())
			Expect(secret.GetName()).To(Equal(controller.getExportSecretName(testExportPod)))
			Expect(secret.GetNamespace()).To(Equal(testNamespace))
			return true, secret, nil
		})
		err = controller.createCertSecret(testVMExport, testExportPod)
		Expect(err).ToNot(HaveOccurred())
		Expect(controller.Recorder.(*record.FakeRecorder).Events).To(HaveLen(1))
		testutils.ExpectEvent(recorder, secretCreatedEvent)
		By("Creating again, and returning exists")
		k8sClient.Fake.PrependReactor("create", "secrets", func(action testing.Action) (handled bool, obj runtime.Object, err error) {
			create, ok := action.(testing.CreateAction)
			Expect(ok).To(BeTrue())
			secret, ok := create.GetObject().(*k8sv1.Secret)
			Expect(ok).To(BeTrue())
			Expect(secret.GetName()).To(Equal(controller.getExportSecretName(testExportPod)))
			Expect(secret.GetNamespace()).To(Equal(testNamespace))
			return true, nil, errors.NewAlreadyExists(schema.GroupResource{}, secret.Name)
		})
		err = controller.createCertSecret(testVMExport, testExportPod)
		Expect(err).ToNot(HaveOccurred())
		Expect(controller.Recorder.(*record.FakeRecorder).Events).To(BeEmpty())
		k8sClient.Fake.PrependReactor("create", "secrets", func(action testing.Action) (handled bool, obj runtime.Object, err error) {
			create, ok := action.(testing.CreateAction)
			Expect(ok).To(BeTrue())
			secret, ok := create.GetObject().(*k8sv1.Secret)
			Expect(ok).To(BeTrue())
			Expect(secret.GetName()).To(Equal(controller.getExportSecretName(testExportPod)))
			Expect(secret.GetNamespace()).To(Equal(testNamespace))
			return true, nil, fmt.Errorf("failure")
		})
		err = controller.createCertSecret(testVMExport, testExportPod)
		Expect(err).To(HaveOccurred())
		Expect(controller.Recorder.(*record.FakeRecorder).Events).To(BeEmpty())
	})

	It("handleVMExportToken should create the export secret if no TokenSecretRef is specified, and the referenced PVC exists", func() {
		testVMExport := createPVCVMExportWithoutSecret()
		expectedName := getDefaultTokenSecretName(testVMExport)
		k8sClient.Fake.PrependReactor("create", "secrets", func(action testing.Action) (handled bool, obj runtime.Object, err error) {
			create, ok := action.(testing.CreateAction)
			Expect(ok).To(BeTrue())
			secret, ok := create.GetObject().(*k8sv1.Secret)
			Expect(ok).To(BeTrue())
			Expect(secret.GetName()).To(Equal(expectedName))
			Expect(secret.GetNamespace()).To(Equal(testNamespace))
			return true, secret, nil
		})
		Expect(controller.PVCInformer.GetStore().Add(&k8sv1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testPVCName,
				Namespace: testNamespace,
			},
		})).To(Succeed())
		populateInitialVMExportStatus(testVMExport)
		err := controller.handleVMExportToken(testVMExport, controller.getPVCFromSourcePVC)
		Expect(err).ToNot(HaveOccurred())
		Expect(testVMExport.Status.TokenSecretRef).ToNot(BeNil())
		Expect(*testVMExport.Status.TokenSecretRef).To(Equal(expectedName))
		testutils.ExpectEvent(recorder, secretCreatedEvent)
	})

	It("handleVMExportToken should not create the export secret if no TokenSecretRef is specified, and the referenced PVC doesn't exist", func() {
		testVMExport := createPVCVMExportWithoutSecret()
		expectedName := getDefaultTokenSecretName(testVMExport)
		k8sClient.Fake.PrependReactor("create", "secrets", func(action testing.Action) (handled bool, obj runtime.Object, err error) {
			create, ok := action.(testing.CreateAction)
			Expect(ok).To(BeTrue())
			secret, ok := create.GetObject().(*k8sv1.Secret)
			Expect(ok).To(BeTrue())
			Expect(secret.GetName()).To(Equal(expectedName))
			Expect(secret.GetNamespace()).To(Equal(testNamespace))
			return true, secret, nil
		})
		populateInitialVMExportStatus(testVMExport)
		err := controller.handleVMExportToken(testVMExport, controller.getPVCFromSourcePVC)
		Expect(err).ToNot(HaveOccurred())
		Expect(testVMExport.Status.TokenSecretRef).To(BeNil())
	})

	It("handleVMExportToken should use the already specified secret if the status is already populated", func() {
		testVMExport := createPVCVMExportWithoutSecret()
		oldSecretRef := "oldToken"
		newSecretRef := getDefaultTokenSecretName(testVMExport)
		testVMExport.Status = &exportv1.VirtualMachineExportStatus{
			TokenSecretRef: &oldSecretRef,
		}
		Expect(controller.PVCInformer.GetStore().Add(&k8sv1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testPVCName,
				Namespace: testNamespace,
			},
		})).To(Succeed())
		k8sClient.Fake.PrependReactor("create", "secrets", func(action testing.Action) (handled bool, obj runtime.Object, err error) {
			create, ok := action.(testing.CreateAction)
			Expect(ok).To(BeTrue())
			secret, ok := create.GetObject().(*k8sv1.Secret)
			Expect(ok).To(BeTrue())
			Expect(secret.GetName()).To(Equal(oldSecretRef))
			Expect(secret.GetNamespace()).To(Equal(testNamespace))
			return true, secret, nil
		})
		err := controller.handleVMExportToken(testVMExport, controller.getPVCFromSourcePVC)
		Expect(err).ToNot(HaveOccurred())
		Expect(testVMExport.Status.TokenSecretRef).ToNot(BeNil())
		Expect(*testVMExport.Status.TokenSecretRef).ToNot(Equal(newSecretRef))
		Expect(*testVMExport.Status.TokenSecretRef).To(Equal(oldSecretRef))
		testutils.ExpectEvent(recorder, secretCreatedEvent)
	})

	It("handleVMExportToken should use the user-specified secret if TokenSecretRef is specified", func() {
		testVMExport := createPVCVMExport()
		Expect(testVMExport.Spec.TokenSecretRef).ToNot(BeNil())
		expectedName := *testVMExport.Spec.TokenSecretRef
		populateInitialVMExportStatus(testVMExport)
		err := controller.handleVMExportToken(testVMExport, controller.getPVCFromSourcePVC)
		Expect(err).ToNot(HaveOccurred())
		Expect(testVMExport.Status.TokenSecretRef).ToNot(BeNil())
		Expect(*testVMExport.Status.TokenSecretRef).To(Equal(expectedName))
	})

	It("Should completely clean up VM export, when TTL is reached", func() {
		var deleted bool
		testVMExport := createPVCVMExport()
		ttl := &metav1.Duration{Duration: time.Minute}
		testVMExport.Spec.TTLDuration = ttl
		// Artificially reach TTL expiration time
		testVMExport.SetCreationTimestamp(metav1.NewTime(time.Now().Add(-1 * ttl.Duration)))
		pvc := &k8sv1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testPVCName,
				Namespace: testNamespace,
			},
			Status: k8sv1.PersistentVolumeClaimStatus{
				Phase: k8sv1.ClaimBound,
			},
		}
		Expect(controller.PVCInformer.GetStore().Add(pvc)).To(Succeed())

		vmExportClient.Fake.PrependReactor("delete", "virtualmachineexports", func(action testing.Action) (handled bool, obj runtime.Object, err error) {
			delete, ok := action.(testing.DeleteAction)
			Expect(ok).To(BeTrue())
			Expect(delete.GetName()).To(Equal(testVMExport.GetName()))
			deleted = true
			return true, nil, nil
		})
		retry, err := controller.updateVMExport(testVMExport)
		Expect(deleted).To(BeTrue())
		// Status update fails (call UPDATE on deleted VMExport), but its fine in real world
		// since requeue will back out of the reconcile loop if a deletion timestamp is set
		Expect(err).To(HaveOccurred())
		Expect(retry).To(BeEquivalentTo(0))
	})

	DescribeTable("Should ignore invalid VMExports kind/api combinations", func(kind, apigroup string) {
		testVMExport := createPVCVMExport()
		testVMExport.Spec.Source.Kind = kind
		testVMExport.Spec.Source.APIGroup = &apigroup
		retry, err := controller.updateVMExport(testVMExport)
		Expect(err).ToNot(HaveOccurred())
		Expect(retry).To(BeEquivalentTo(0))
	},
		Entry("VirtualMachineSnapshot kind blank apigroup", "VirtualMachineSnapshot", ""),
		Entry("VirtualMachineSnapshot kind invalid apigroup", "VirtualMachineSnapshot", "invalid"),
		Entry("PersistentVolumeClaim kind invalid apigroup", "PersistentVolumeClaim", "invalid"),
		Entry("PersistentVolumeClaim kind VMSnapshot apigroup", "PersistentVolumeClaim", snapshotv1.SchemeGroupVersion.Group),
	)

	DescribeTable("should find host when Ingress is defined", func(ingress *networkingv1.Ingress, hostname string) {
		Expect(controller.IngressCache.Add(ingress)).To(Succeed())
		host, _ := controller.getExternalLinkHostAndCert()
		Expect(hostname).To(Equal(host))
	},
		Entry("ingress with default backend host", validIngressDefaultBackend(components.VirtExportProxyServiceName), "backend-host"),
		Entry("ingress with default backend host different service", validIngressDefaultBackend("other-service"), ""),
		Entry("ingress with rules host", validIngressRules(components.VirtExportProxyServiceName), "rule-host"),
		Entry("ingress with rules host different service", validIngressRules("other-service"), ""),
		Entry("ingress with no default backend service", ingressDefaultBackendNoService(), ""),
		Entry("ingress with rules no backend service", ingressRulesNoBackend(), ""),
	)

	DescribeTable("should find host when route is defined", func(createCMFunc func() *k8sv1.ConfigMap, route *routev1.Route, hostname, expectedCert string) {
		controller.RouteCache.Add(route)
		controller.RouteConfigMapInformer.GetStore().Add(createCMFunc())
		host, cert := controller.getExternalLinkHostAndCert()
		Expect(host).To(Equal(hostname))
		Expect(cert).To(Equal(expectedCert))
	},
		Entry("route with service and host", createRouteConfigMap, routeToHostAndService(components.VirtExportProxyServiceName), "virt-exportproxy-kubevirt.apps-crc.testing", expectedPem),
		Entry("route with different service and host", createRouteConfigMap, routeToHostAndService("other-service"), "", ""),
		Entry("route with service and no ingress", createRouteConfigMap, routeToHostAndNoIngress(), "", ""),
		Entry("should not find route cert if in future", createFutureRouteConfigMap, routeToHostAndService(components.VirtExportProxyServiceName), "virt-exportproxy-kubevirt.apps-crc.testing", expectedFuturePem),
		Entry("should not find route cert if expired", createExpiredRouteConfigMap, routeToHostAndService(components.VirtExportProxyServiceName), "virt-exportproxy-kubevirt.apps-crc.testing", expectedExpiredPem),
		Entry("should find correct route cert if overlapping exists", createOverlappingRouteConfigMap, routeToHostAndService(components.VirtExportProxyServiceName), "virt-exportproxy-kubevirt.apps-crc.testing", expectedPem),
	)

	It("should pick ingress over route if both exist", func() {
		Expect(
			controller.IngressCache.Add(validIngressDefaultBackend(components.VirtExportProxyServiceName)),
		).To(Succeed())
		Expect(controller.RouteCache.Add(routeToHostAndService(components.VirtExportProxyServiceName))).To(Succeed())
		host, _ := controller.getExternalLinkHostAndCert()
		Expect("backend-host").To(Equal(host))
	})

	populateIngressSecret := func() {
		ingressCache.Add(ingressToHost())
		secretInformer.GetStore().Add(&k8sv1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      ingressSecret,
				Namespace: testNamespace,
			},
			Data: map[string][]byte{
				"tls.crt": []byte(expectedPem),
			},
		})
	}

	It("should create datamanifest and add it to the pod spec", func() {
		populateIngressSecret()
		testVMExport := createVMVMExportExternal()
		vm := &virtv1.VirtualMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testVmName,
				Namespace: testNamespace,
			},
			Spec: virtv1.VirtualMachineSpec{
				Template: &virtv1.VirtualMachineInstanceTemplateSpec{
					Spec: virtv1.VirtualMachineInstanceSpec{
						Volumes: []virtv1.Volume{},
					},
				},
			},
		}
		service := &k8sv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      controller.getExportServiceName(testVMExport),
				Namespace: testVMExport.Namespace,
			},
			Spec: k8sv1.ServiceSpec{},
			Status: k8sv1.ServiceStatus{
				Conditions: []metav1.Condition{
					{
						Type: "test2",
					},
				},
			},
		}
		testPod := &k8sv1.Pod{
			Spec: k8sv1.PodSpec{
				Containers: []k8sv1.Container{
					{
						VolumeMounts: []k8sv1.VolumeMount{},
					},
				},
				Volumes: []k8sv1.Volume{},
			},
		}
		cmName := controller.getVmManifestConfigMapName(testVMExport)
		vmBytes, err := controller.generateVMDefinitionFromVm(vm)
		Expect(err).ToNot(HaveOccurred())
		k8sClient.Fake.PrependReactor("create", "configmaps", func(action testing.Action) (handled bool, obj runtime.Object, err error) {
			create, ok := action.(testing.CreateAction)
			Expect(ok).To(BeTrue())
			cm, ok := create.GetObject().(*k8sv1.ConfigMap)
			Expect(ok).To(BeTrue())
			Expect(cm.GetName()).To(Equal(cmName))
			Expect(cm.GetNamespace()).To(Equal(testNamespace))
			Expect(cm.Data).ToNot(BeEmpty())
			Expect(cm.Data[internalHostKey]).To(Equal(fmt.Sprintf("%s.%s.svc", controller.getExportServiceName(testVMExport), service.Namespace)))
			Expect(cm.Data[vmManifest]).To(Equal(string(vmBytes)))
			return true, cm, nil
		})
		err = controller.createDataManifestAndAddToPod(testVMExport, vm, testPod, service)
		Expect(err).ToNot(HaveOccurred())
		Expect(testVMExport.Status).ToNot(BeNil())
	})

	createVM := func() *virtv1.VirtualMachine {
		return &virtv1.VirtualMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testVmName,
				Namespace: testNamespace,
			},
			Spec: virtv1.VirtualMachineSpec{
				Template: &virtv1.VirtualMachineInstanceTemplateSpec{
					Spec: virtv1.VirtualMachineInstanceSpec{
						Volumes: []virtv1.Volume{},
					},
				},
				Instancetype: &virtv1.InstancetypeMatcher{
					Name: "test-instance-type",
					Kind: apiinstancetype.SingularResourceName,
				},
			},
		}
	}
	It("Should properly expand instance types of VMs", func() {
		vm := createVM()
		testInstanceType := &instancetypev1beta1.VirtualMachineInstancetype{
			TypeMeta: metav1.TypeMeta{
				Kind:       apiinstancetype.SingularResourceName,
				APIVersion: instancetypev1beta1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-instance-type",
				Namespace: vm.Namespace,
			},
			Spec: instancetypev1beta1.VirtualMachineInstancetypeSpec{
				CPU: instancetypev1beta1.CPUInstancetype{
					Guest: uint32(2),
				},
			},
		}
		Expect(instancetypeInformer.GetStore().Add(testInstanceType)).To(Succeed())

		res, err := controller.generateVMDefinitionFromVm(vm)
		Expect(err).ToNot(HaveOccurred())
		Expect(res).ToNot(BeNil())

		resVM := &virtv1.VirtualMachine{}
		Expect(json.Unmarshal(res, resVM)).To(Succeed())
		Expect(resVM.Spec.Instancetype).To(BeNil())
		Expect(resVM.Spec.Template.Spec.Domain.CPU.Sockets).To(Equal(uint32(2)))
	})

	It("Should return error on conflict with instance types of VMs", func() {
		vm := createVM()
		vm.Spec.Template.Spec.Domain = virtv1.DomainSpec{
			CPU: &virtv1.CPU{
				Cores: uint32(1),
			},
		}
		testInstanceType := &instancetypev1beta1.VirtualMachineInstancetype{
			TypeMeta: metav1.TypeMeta{
				Kind:       apiinstancetype.SingularResourceName,
				APIVersion: instancetypev1beta1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-instance-type",
				Namespace: vm.Namespace,
			},
			Spec: instancetypev1beta1.VirtualMachineInstancetypeSpec{
				CPU: instancetypev1beta1.CPUInstancetype{
					Guest: uint32(2),
				},
			},
		}
		Expect(instancetypeInformer.GetStore().Add(testInstanceType)).To(Succeed())

		_, err := controller.generateVMDefinitionFromVm(vm)
		Expect(err).To(MatchError(ContainSubstring("VM field(s) spec.template.spec.domain.cpu.cores conflicts with selected instance type")))
	})

	createVMWithDVTemplateAndPVC := func() *virtv1.VirtualMachine {
		vm := createVM()
		vm.Spec.DataVolumeTemplates = []virtv1.DataVolumeTemplateSpec{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dv-template",
					Namespace: vm.Namespace,
				},
				Spec: cdiv1.DataVolumeSpec{
					Source: &cdiv1.DataVolumeSource{
						Blank: &cdiv1.DataVolumeBlankImage{},
					},
					SourceRef: &cdiv1.DataVolumeSourceRef{
						Kind:      "",
						Name:      "test",
						Namespace: pointer.P("default"),
					},
				},
			},
		}
		vm.Spec.Template.Spec.Volumes = append(vm.Spec.Template.Spec.Volumes, virtv1.Volume{
			Name: "template-for-dv",
			VolumeSource: virtv1.VolumeSource{
				DataVolume: &virtv1.DataVolumeSource{
					Name: "dv-template",
				},
			},
		})
		vm.Spec.Template.Spec.Volumes = append(vm.Spec.Template.Spec.Volumes, virtv1.Volume{
			Name: "non-dv-pvc",
			VolumeSource: virtv1.VolumeSource{
				PersistentVolumeClaim: &virtv1.PersistentVolumeClaimVolumeSource{
					PersistentVolumeClaimVolumeSource: k8sv1.PersistentVolumeClaimVolumeSource{
						ClaimName: "pvc",
					},
				},
			},
		})
		return vm
	}

	It("Should properly replace DVTemplates", func() {
		vm := createVMWithDVTemplateAndPVC()
		res, err := controller.updateHttpSourceDataVolumeTemplate(vm)
		Expect(err).ToNot(HaveOccurred())
		Expect(res).ToNot(BeNil())
		Expect(res.Spec.DataVolumeTemplates).To(HaveLen(1))
		Expect(res.Spec.DataVolumeTemplates[0].Spec.Source).ToNot(BeNil())
		Expect(res.Spec.DataVolumeTemplates[0].Spec.Source.HTTP).ToNot(BeNil())
		Expect(res.Spec.DataVolumeTemplates[0].Spec.Source.HTTP.URL).To(BeEmpty())
		Expect(res.Spec.DataVolumeTemplates[0].Spec.Source.Blank).To(BeNil())
		Expect(res.Spec.DataVolumeTemplates[0].Spec.SourceRef).To(BeNil())
	})

	It("Should generate DataVolumes from VM", func() {
		pvc := createPVC("pvc", string(cdiv1.DataVolumeKubeVirt))
		pvc.Spec.DataSource = &k8sv1.TypedLocalObjectReference{}
		pvc.Spec.DataSourceRef = &k8sv1.TypedObjectReference{}
		pvc.Spec.StorageClassName = pointer.P("somesc")
		pvcCopy := pvc.DeepCopy()
		pvcInformer.GetStore().Add(pvc)
		vm := createVMWithDVTemplateAndPVC()
		dvs, err := controller.generateDataVolumesFromVm(vm)
		Expect(err).ToNot(HaveOccurred())
		Expect(dvs).To(HaveLen(1))
		Expect(dvs[0]).ToNot(BeNil())
		Expect(dvs[0].Name).To(Equal("pvc"))
		Expect(dvs[0].Spec.Storage).To(
			gstruct.PointTo(
				gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
					"DataSource":       BeNil(),
					"DataSourceRef":    BeNil(),
					"StorageClassName": BeNil(),
					"VolumeMode":       Equal(pvcCopy.Spec.VolumeMode),
					"Resources":        Equal(pvcCopy.Spec.Resources),
					"AccessModes":      HaveExactElements(pvcCopy.Spec.AccessModes),
				}),
			),
		)
		Expect(dvs[0].Spec).To(
			gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"PVC":       BeNil(),
				"SourceRef": BeNil(),
			}),
		)
	})
})

func verifyLinksEmpty(vmExport *exportv1.VirtualMachineExport) {
	Expect(vmExport.Status).ToNot(BeNil())
	Expect(vmExport.Status.Links).ToNot(BeNil())
	Expect(vmExport.Status.Links.Internal).To(BeNil())
	Expect(vmExport.Status.Links.External).To(BeNil())
}

func verifyLinksInternal(vmExport *exportv1.VirtualMachineExport, expectedVolumeFormats ...exportv1.VirtualMachineExportVolumeFormat) {
	Expect(vmExport.Status).ToNot(BeNil())
	Expect(vmExport.Status.Links).ToNot(BeNil())
	Expect(vmExport.Status.Links.Internal).NotTo(BeNil())
	Expect(vmExport.Status.Links.Internal.Cert).NotTo(BeEmpty())
	Expect(vmExport.Status.Links.Internal.Volumes).To(HaveLen(len(expectedVolumeFormats) / 2))
	for _, volume := range vmExport.Status.Links.Internal.Volumes {
		Expect(volume.Formats).To(HaveLen(2))
		Expect(expectedVolumeFormats).To(ContainElements(volume.Formats))
	}
}

func verifyLinksExternal(vmExport *exportv1.VirtualMachineExport, link1Format exportv1.ExportVolumeFormat, link1Url string, link2Format exportv1.ExportVolumeFormat, link2Url string) {
	Expect(vmExport.Status.Links.External).ToNot(BeNil())
	Expect(vmExport.Status.Links.External.Cert).To(BeEmpty())
	Expect(vmExport.Status.Links.External.Volumes).To(HaveLen(1))
	Expect(vmExport.Status.Links.External.Volumes[0].Formats).To(HaveLen(2))
	Expect(vmExport.Status.Links.External.Volumes[0].Formats).To(ContainElements(exportv1.VirtualMachineExportVolumeFormat{
		Format: link1Format,
		Url:    link1Url,
	}, exportv1.VirtualMachineExportVolumeFormat{
		Format: link2Format,
		Url:    link2Url,
	}))
}

func verifyKubevirtInternal(vmExport *exportv1.VirtualMachineExport, exportName, namespace string, volumeNames ...string) {
	exportVolumeFormats := make([]exportv1.VirtualMachineExportVolumeFormat, 0)
	for _, volumeName := range volumeNames {
		exportVolumeFormats = append(exportVolumeFormats, exportv1.VirtualMachineExportVolumeFormat{
			Format: exportv1.KubeVirtRaw,
			Url:    fmt.Sprintf("https://%s.%s.svc/volumes/%s/disk.img", fmt.Sprintf("%s-%s", exportPrefix, exportName), namespace, volumeName),
		})
		exportVolumeFormats = append(exportVolumeFormats, exportv1.VirtualMachineExportVolumeFormat{
			Format: exportv1.KubeVirtGz,
			Url:    fmt.Sprintf("https://%s.%s.svc/volumes/%s/disk.img.gz", fmt.Sprintf("%s-%s", exportPrefix, exportName), namespace, volumeName),
		})
	}
	verifyLinksInternal(vmExport, exportVolumeFormats...)
}

func verifyKubevirtExternal(vmExport *exportv1.VirtualMachineExport, exportName, namespace, volumeName string) {
	verifyLinksExternal(vmExport,
		exportv1.KubeVirtRaw,
		fmt.Sprintf("https://virt-exportproxy-kubevirt.apps-crc.testing/api/export.kubevirt.io/%s/namespaces/%s/virtualmachineexports/%s/volumes/%s/disk.img", currentVersion, namespace, exportName, volumeName),
		exportv1.KubeVirtGz,
		fmt.Sprintf("https://virt-exportproxy-kubevirt.apps-crc.testing/api/export.kubevirt.io/%s/namespaces/%s/virtualmachineexports/%s/volumes/%s/disk.img.gz", currentVersion, namespace, exportName, volumeName))
}

func verifyArchiveInternal(vmExport *exportv1.VirtualMachineExport, exportName, namespace, volumeName string) {
	verifyLinksInternal(vmExport,
		exportv1.VirtualMachineExportVolumeFormat{
			Format: exportv1.Dir,
			Url:    fmt.Sprintf("https://%s.%s.svc/volumes/%s/dir", fmt.Sprintf("%s-%s", exportPrefix, exportName), namespace, volumeName),
		}, exportv1.VirtualMachineExportVolumeFormat{
			Format: exportv1.ArchiveGz,
			Url:    fmt.Sprintf("https://%s.%s.svc/volumes/%s/disk.tar.gz", fmt.Sprintf("%s-%s", exportPrefix, exportName), namespace, volumeName),
		})
}

func routeToHostAndService(serviceName string) *routev1.Route {
	return &routev1.Route{
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Name: serviceName,
			},
		},
		Status: routev1.RouteStatus{
			Ingress: []routev1.RouteIngress{
				{
					Host: fmt.Sprintf("%s-kubevirt.apps-crc.testing", components.VirtExportProxyServiceName),
				},
			},
		},
	}
}

func routeToHostAndNoIngress() *routev1.Route {
	return &routev1.Route{
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Name: components.VirtExportProxyServiceName,
			},
		},
		Status: routev1.RouteStatus{
			Ingress: []routev1.RouteIngress{},
		},
	}
}

func ingressToHost() *networkingv1.Ingress {
	return &networkingv1.Ingress{
		Spec: networkingv1.IngressSpec{
			TLS: []networkingv1.IngressTLS{
				{
					SecretName: ingressSecret,
				},
			},
			Rules: []networkingv1.IngressRule{
				{
					Host: "test-host",
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: components.VirtExportProxyServiceName,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func verifyArchiveExternal(vmExport *exportv1.VirtualMachineExport, exportName, namespace, volumeName string) {
	verifyLinksExternal(vmExport,
		exportv1.Dir,
		fmt.Sprintf("https://virt-exportproxy-kubevirt.apps-crc.testing/api/export.kubevirt.io/%s/namespaces/%s/virtualmachineexports/%s/volumes/%s/dir", currentVersion, namespace, exportName, volumeName),
		exportv1.ArchiveGz,
		fmt.Sprintf("https://virt-exportproxy-kubevirt.apps-crc.testing/api/export.kubevirt.io/%s/namespaces/%s/virtualmachineexports/%s/volumes/%s/disk.tar.gz", currentVersion, namespace, exportName, volumeName))
}

func createVMExportMeta(name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:              name,
		Namespace:         testNamespace,
		UID:               uuid.NewUUID(),
		CreationTimestamp: metav1.Now(),
		Labels:            map[string]string{labelKey: labelValue},
		Annotations:       map[string]string{annotationKey: annotationValue},
	}
}

func createPVCVMExport() *exportv1.VirtualMachineExport {
	return createPVCVMExportWithName(vmExportName)
}

func createPVCVMExportLongName() *exportv1.VirtualMachineExport {
	return createPVCVMExportWithName(vmExportName + strings.Repeat("a", 63))
}

func createPVCVMExportWithName(name string) *exportv1.VirtualMachineExport {
	return &exportv1.VirtualMachineExport{
		ObjectMeta: createVMExportMeta(name),
		Spec: exportv1.VirtualMachineExportSpec{
			Source: k8sv1.TypedLocalObjectReference{
				APIGroup: &k8sv1.SchemeGroupVersion.Group,
				Kind:     "PersistentVolumeClaim",
				Name:     testPVCName,
			},
			TokenSecretRef: &tokenSecretName,
		},
	}
}

func createPVCVMExportWithoutSecret() *exportv1.VirtualMachineExport {
	return &exportv1.VirtualMachineExport{
		ObjectMeta: createVMExportMeta(vmExportName),
		Spec: exportv1.VirtualMachineExportSpec{
			Source: k8sv1.TypedLocalObjectReference{
				APIGroup: &k8sv1.SchemeGroupVersion.Group,
				Kind:     "PersistentVolumeClaim",
				Name:     testPVCName,
			},
		},
	}
}

func createSnapshotVMExport() *exportv1.VirtualMachineExport {
	return &exportv1.VirtualMachineExport{
		ObjectMeta: createVMExportMeta(vmExportName),
		Spec: exportv1.VirtualMachineExportSpec{
			Source: k8sv1.TypedLocalObjectReference{
				APIGroup: &snapshotv1.SchemeGroupVersion.Group,
				Kind:     "VirtualMachineSnapshot",
				Name:     testVmsnapshotName,
			},
			TokenSecretRef: &tokenSecretName,
		},
	}
}

func createVMVMExport() *exportv1.VirtualMachineExport {
	return &exportv1.VirtualMachineExport{
		ObjectMeta: createVMExportMeta(vmExportName),
		Spec: exportv1.VirtualMachineExportSpec{
			Source: k8sv1.TypedLocalObjectReference{
				APIGroup: &virtv1.SchemeGroupVersion.Group,
				Kind:     "VirtualMachine",
				Name:     testVmName,
			},
			TokenSecretRef: &tokenSecretName,
		},
	}
}

func createVMVMExportExternal() *exportv1.VirtualMachineExport {
	res := createVMVMExport()
	res.Status = &exportv1.VirtualMachineExportStatus{
		Links: &exportv1.VirtualMachineExportLinks{
			External: &exportv1.VirtualMachineExportLink{
				Cert:    "test-cert",
				Volumes: []exportv1.VirtualMachineExportVolume{},
			},
		},
	}
	return res
}

func createPVC(name, contentType string) *k8sv1.PersistentVolumeClaim {
	return &k8sv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
			Annotations: map[string]string{
				annContentType: contentType,
			},
		},
		Spec: k8sv1.PersistentVolumeClaimSpec{
			VolumeMode:  pointer.P(k8sv1.PersistentVolumeMode("testvolumemode")),
			AccessModes: []k8sv1.PersistentVolumeAccessMode{k8sv1.PersistentVolumeAccessMode("testaccessmode")},
			Resources: k8sv1.VolumeResourceRequirements{
				Requests: k8sv1.ResourceList{
					k8sv1.ResourceStorage: resource.MustParse("99Gi"),
				},
			},
		},
		Status: k8sv1.PersistentVolumeClaimStatus{
			Phase: k8sv1.ClaimBound,
		},
	}
}

func createBackendPVC(vmName string) *k8sv1.PersistentVolumeClaim {
	return &k8sv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s-", backendstorage.PVCPrefix, vmName),
			Namespace: testNamespace,
			Annotations: map[string]string{
				annContentType: "archive",
			},
			Labels: map[string]string{
				backendstorage.PVCPrefix: vmName,
			},
		},
		Status: k8sv1.PersistentVolumeClaimStatus{
			Phase: k8sv1.ClaimBound,
		},
	}
}

func expectExporterCreate(k8sClient *k8sfake.Clientset, phase k8sv1.PodPhase) {
	k8sClient.Fake.PrependReactor("create", "pods", func(action testing.Action) (handled bool, obj runtime.Object, err error) {
		create, ok := action.(testing.CreateAction)
		Expect(ok).To(BeTrue())
		exportPod, ok := create.GetObject().(*k8sv1.Pod)
		Expect(ok).To(BeTrue())
		exportPod.Status = k8sv1.PodStatus{
			Phase: phase,
		}
		return true, exportPod, nil
	})
}

func expectExporterDelete(k8sClient *k8sfake.Clientset, expectedName string) {
	k8sClient.Fake.PrependReactor("delete", "pods", func(action testing.Action) (handled bool, obj runtime.Object, err error) {
		delete, ok := action.(testing.DeleteAction)
		Expect(ok).To(BeTrue())
		Expect(delete.GetName()).To(Equal(expectedName))
		return true, nil, nil
	})
}

// A mock to implement volumeSnapshotProvider interface
type MockVolumeSnapshotProvider struct {
	volumeSnapshots []*vsv1.VolumeSnapshot
}

func (v *MockVolumeSnapshotProvider) GetVolumeSnapshot(namespace, name string) (*vsv1.VolumeSnapshot, error) {
	if len(v.volumeSnapshots) == 0 {
		return nil, nil
	}
	vs := v.volumeSnapshots[0]
	v.volumeSnapshots = v.volumeSnapshots[1:]
	return vs, nil
}

func (v *MockVolumeSnapshotProvider) Add(s *vsv1.VolumeSnapshot) {
	v.volumeSnapshots = append(v.volumeSnapshots, s)
}

// A mock to implement the certificate.Manager interface for the export controller
type MockCertManager struct {
	crt *tls.Certificate
}

func (f *MockCertManager) Start() {
	caKeyPair, _ := triple.NewCA("test.kubevirt.io", time.Hour)

	encodedCert := cert.EncodeCertPEM(caKeyPair.Cert)
	encodedKey := cert.EncodePrivateKeyPEM(caKeyPair.Key)

	crt, err := tls.X509KeyPair(encodedCert, encodedKey)
	Expect(err).ToNot(HaveOccurred())
	leaf, err := cert.ParseCertsPEM(encodedCert)
	Expect(err).ToNot(HaveOccurred())
	crt.Leaf = leaf[0]
	f.crt = &crt
}

func (f *MockCertManager) Stop() {
}

func (f *MockCertManager) Current() *tls.Certificate {
	return f.crt
}

func (f *MockCertManager) ServerHealthy() bool {
	return true
}
