module code.cloudfoundry.org/netplugin-shim

go 1.16

require (
	code.cloudfoundry.org/guardian v0.0.0-20210527102652-2f945c09a983
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/jessevdk/go-flags v1.5.0
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.13.0
	golang.org/x/sys v0.0.0-20210603081109-ebe580a85c40
)

replace (
	code.cloudfoundry.org/garden => ../garden
	code.cloudfoundry.org/grootfs => ../grootfs
	code.cloudfoundry.org/guardian => ../guardian
	code.cloudfoundry.org/idmapper => ../idmapper
)
