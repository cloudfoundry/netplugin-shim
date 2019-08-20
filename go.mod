module code.cloudfoundry.org/netplugin-shim

go 1.12

require (
	code.cloudfoundry.org/guardian v0.0.0-00010101000000-000000000000
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/jessevdk/go-flags v1.4.0
	github.com/onsi/ginkgo v1.9.0
	github.com/onsi/gomega v1.5.0
	golang.org/x/sys v0.0.0-20190626221950-04f50cda93cb
)

replace (
	code.cloudfoundry.org/garden => ../garden
	code.cloudfoundry.org/grootfs => ../grootfs
	code.cloudfoundry.org/guardian => ../guardian
	code.cloudfoundry.org/idmapper => ../idmapper
)
