module code.cloudfoundry.org/netplugin-shim

go 1.14

require (
	code.cloudfoundry.org/guardian v0.0.0-00010101000000-000000000000
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/jessevdk/go-flags v1.4.0
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.1
	golang.org/x/sys v0.0.0-20201201145000-ef89a241ccb3
)

replace (
	code.cloudfoundry.org/garden => ../garden
	code.cloudfoundry.org/grootfs => ../grootfs
	code.cloudfoundry.org/guardian => ../guardian
	code.cloudfoundry.org/idmapper => ../idmapper
)
