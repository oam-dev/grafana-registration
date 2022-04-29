module github.com/zzxwill/grafana-configuration

go 1.16

require (
	github.com/go-logr/logr v0.4.0
	github.com/go-logr/zapr v0.4.0 // indirect
	github.com/grafana-tools/sdk v0.0.0-20210714133701-11b1efc100c9
	github.com/onsi/ginkgo v1.16.2
	github.com/onsi/gomega v1.12.0
	github.com/pkg/errors v0.9.1
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v0.18.8
	k8s.io/klog/v2 v2.10.0
	sigs.k8s.io/controller-runtime v0.6.2
)

replace github.com/grafana-tools/sdk v0.0.0-20210714133701-11b1efc100c9 => github.com/kubevela-contrib/grafana-sdk v0.9.5-0.20220428052213-cec5f5a9efd8
