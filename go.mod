module code.cloudfoundry.org/netplugin-shim

go 1.16

require (
	code.cloudfoundry.org/guardian v0.0.0-20210813144446-9d3aeb65f163
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/jessevdk/go-flags v1.5.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.17.0
	golang.org/x/sys v0.0.0-20211111213525-f221eed1c01e
)

replace (
	code.cloudfoundry.org/garden => ../garden
	code.cloudfoundry.org/grootfs => ../grootfs
	code.cloudfoundry.org/guardian => ../guardian
	code.cloudfoundry.org/idmapper => ../idmapper
)
