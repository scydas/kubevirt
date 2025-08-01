module kubevirt.io/kubevirt

go 1.23.0

require (
	github.com/cheggaaa/pb/v3 v3.1.0
	github.com/containernetworking/cni v1.2.3
	github.com/containernetworking/plugins v1.5.1
	github.com/coreos/go-semver v0.3.1
	github.com/emicklei/go-restful/v3 v3.11.0
	github.com/evanphx/json-patch v5.6.0+incompatible
	github.com/fsnotify/fsnotify v1.7.0
	github.com/go-kit/log v0.2.1
	github.com/go-openapi/spec v0.21.0
	github.com/go-openapi/strfmt v0.23.0
	github.com/go-openapi/validate v0.24.0
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.4
	github.com/google/go-cmp v0.6.0
	github.com/google/go-github/v32 v32.0.0
	github.com/google/goexpect v0.0.0-20190425035906-112704a48083
	github.com/google/gofuzz v1.2.0
	github.com/google/uuid v1.6.0
	github.com/gordonklaus/ineffassign v0.0.0-20210209182638-d0e41b2fc8ed
	github.com/gorilla/websocket v1.5.0
	github.com/imdario/mergo v0.3.16
	github.com/insomniacslk/dhcp v0.0.0-20230908212754-65c27093e38a
	github.com/k8snetworkplumbingwg/network-attachment-definition-client v1.3.0
	github.com/kisielk/errcheck v1.6.2
	github.com/klauspost/pgzip v1.2.6
	github.com/krolaw/dhcp4 v0.0.0-20180925202202-7cead472c414
	github.com/kubernetes-csi/external-snapshotter/client/v4 v4.2.0
	github.com/kubevirt/monitoring/pkg/metrics/parser v0.0.0-20230627123556-81a891d4462a
	github.com/mdlayher/vsock v1.2.1
	github.com/mitchellh/go-ps v0.0.0-20190716172923-621e5597135b
	github.com/mitchellh/go-vnc v0.0.0-20150629162542-723ed9867aed
	github.com/moby/sys/mountinfo v0.7.2
	github.com/nunnatsa/ginkgolinter v0.14.0
	github.com/nxadm/tail v1.4.8
	github.com/onsi/ginkgo/v2 v2.22.1
	github.com/onsi/gomega v1.36.2
	github.com/opencontainers/runc v1.1.14
	github.com/opencontainers/selinux v1.11.0
	github.com/openshift/api v0.0.0
	github.com/openshift/client-go v0.0.0
	github.com/openshift/library-go v0.0.0-20211220195323-eca2c467c492
	github.com/operator-framework/operator-lifecycle-manager v0.0.0-20190725173916-b56e63a643cc
	github.com/pkg/errors v0.9.1
	github.com/povsister/scp v0.0.0-20210427074412-33febfd9f13e
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.68.0
	github.com/prometheus/client_golang v1.19.1
	github.com/prometheus/client_model v0.6.1
	github.com/prometheus/common v0.55.0
	github.com/prometheus/procfs v0.15.1
	github.com/rhobs/operator-observability-toolkit v0.0.29
	github.com/spf13/cobra v1.8.1
	github.com/spf13/pflag v1.0.5
	github.com/vishvananda/netlink v1.3.0
	github.com/wadey/gocovmerge v0.0.0-20160331181800-b5bfa59ec0ad
	go.uber.org/mock v0.5.1
	golang.org/x/crypto v0.36.0
	golang.org/x/net v0.38.0
	golang.org/x/sync v0.12.0
	golang.org/x/sys v0.31.0
	golang.org/x/term v0.30.0
	golang.org/x/time v0.7.0
	golang.org/x/tools v0.28.0
	google.golang.org/grpc v1.65.0
	k8s.io/api v0.32.5
	k8s.io/apiextensions-apiserver v0.32.5
	k8s.io/apimachinery v0.32.5
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/kube-aggregator v0.26.4
	k8s.io/kube-openapi v0.31.0
	k8s.io/kubectl v0.0.0-00010101000000-000000000000
	k8s.io/utils v0.0.0-20241104100929-3ea5e8cea738
	kubevirt.io/api v0.0.0-00010101000000-000000000000
	kubevirt.io/client-go v0.0.0-00010101000000-000000000000
	kubevirt.io/containerized-data-importer-api v1.60.3-0.20241105012228-50fbed985de9
	kubevirt.io/qe-tools v0.1.8
	libvirt.org/go/libvirt v1.11005.0
	libvirt.org/go/libvirtxml v1.11000.1
	mvdan.cc/sh/v3 v3.8.0
	sigs.k8s.io/controller-runtime v0.20.2
	sigs.k8s.io/json v0.0.0-20241010143419-9aa6b5e7a4b3
	sigs.k8s.io/yaml v1.4.0
)

require (
	github.com/VividCortex/ewma v1.2.0 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cilium/ebpf v0.11.0 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/cyphar/filepath-securejoin v0.3.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-openapi/analysis v0.23.0 // indirect
	github.com/go-openapi/errors v0.22.0 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/loads v0.22.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/go-toolsmith/astcopy v1.1.0 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/golang/glog v1.2.4 // indirect
	github.com/google/btree v1.1.3 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/google/pprof v0.0.0-20241210010833-40e02aabc2ad // indirect
	github.com/google/renameio/v2 v2.0.0 // indirect
	github.com/grafana/regexp v0.0.0-20221122212121-6b5c0a4cb7fd // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/josharian/native v1.1.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mdlayher/socket v0.4.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/moby/spdystream v0.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mxk/go-flowrate v0.0.0-20140419014527-cca7078d478f // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/opencontainers/runtime-spec v1.2.0 // indirect
	github.com/openshift/custom-resource-status v1.1.2 // indirect
	github.com/pierrec/lz4/v4 v4.1.18 // indirect
	github.com/pkg/diff v0.0.0-20210226163009-20ebb0f2a09e // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/seccomp/libseccomp-golang v0.10.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/u-root/uio v0.0.0-20230220225925-ffce2a382923 // indirect
	github.com/vishvananda/netns v0.0.4 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	go.mongodb.org/mongo-driver v1.14.0 // indirect
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56 // indirect
	golang.org/x/exp/typeparams v0.0.0-20230203172020-98cc5a0785f9 // indirect
	golang.org/x/mod v0.22.0 // indirect
	golang.org/x/oauth2 v0.27.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240826202546-f6391c0de4c7 // indirect
	google.golang.org/protobuf v1.36.1 // indirect
	gopkg.in/evanphx/json-patch.v4 v4.12.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apiserver v0.32.5 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	kubevirt.io/controller-lifecycle-operator-sdk/api v0.0.0-20220329064328-f3cc58c6ed90 // indirect
	mvdan.cc/editorconfig v0.2.1-0.20231228180347-1925077f8eb2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.3 // indirect
)

require (
	github.com/containers/common v0.60.4
	github.com/google/goterm v0.0.0-20190311235235-ce302be1d114 // indirect; indirect github.com/gophercloud/gophercloud v0.4.0 // indirect
)

replace (
	github.com/nxadm/tail => github.com/nxadm/tail v0.0.0-20211216163028-4472660a31a6
	github.com/openshift/api => github.com/openshift/api v0.0.0-20210105115604-44119421ec6b
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20210112165513-ebc401615f47
	github.com/operator-framework/operator-lifecycle-manager => github.com/operator-framework/operator-lifecycle-manager v0.0.0-20190128024246-5eb7ae5bdb7a

	k8s.io/api => k8s.io/api v0.32.5
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.32.5
	k8s.io/apimachinery => k8s.io/apimachinery v0.32.5
	k8s.io/apiserver => k8s.io/apiserver v0.32.5
	k8s.io/client-go => k8s.io/client-go v0.32.5
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.32.5
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.32.5
	k8s.io/code-generator => k8s.io/code-generator v0.32.5
	k8s.io/component-base => k8s.io/component-base v0.32.5
	k8s.io/cri-api => k8s.io/cri-api v0.32.5
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.32.5
	k8s.io/klog => k8s.io/klog v0.4.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.32.5
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.32.5
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20240430033511-f0e62f92d13f
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.32.5
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.32.5
	k8s.io/kubectl => k8s.io/kubectl v0.32.5
	k8s.io/kubelet => k8s.io/kubelet v0.32.5
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.32.5
	k8s.io/metrics => k8s.io/metrics v0.32.5
	k8s.io/node-api => k8s.io/node-api v0.32.5
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.32.5
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.32.5
	k8s.io/sample-controller => k8s.io/sample-controller v0.32.5

	kubevirt.io/api => ./staging/src/kubevirt.io/api
	kubevirt.io/client-go => ./staging/src/kubevirt.io/client-go
)
