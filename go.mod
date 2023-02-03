module code.cloudfoundry.org/netplugin-shim

go 1.16

require (
	code.cloudfoundry.org/guardian v0.0.0-20220607160814-bbdc1696f4d2
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/jessevdk/go-flags v1.5.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.24.1
	golang.org/x/sys v0.4.0
)

replace (
	code.cloudfoundry.org/garden => ../garden
	code.cloudfoundry.org/grootfs => ../grootfs
	code.cloudfoundry.org/guardian => ../guardian
	code.cloudfoundry.org/idmapper => ../idmapper
	golang.org/x/text => golang.org/x/text v0.3.7
)
