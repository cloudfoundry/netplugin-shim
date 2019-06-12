module code.cloudfoundry.org/netplugin-shim

go 1.12

require (
	code.cloudfoundry.org/garden v0.0.0-00010101000000-000000000000 // indirect
	code.cloudfoundry.org/guardian v0.0.0-00010101000000-000000000000
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/apoydence/eachers v0.0.0-20181020210610-23942921fe77 // indirect
	github.com/cloudfoundry/sonde-go v0.0.0-20171206171820-b33733203bb4 // indirect
	github.com/jessevdk/go-flags v1.4.0
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/poy/eachers v0.0.0-20181020210610-23942921fe77 // indirect
	golang.org/x/sys v0.0.0-20190602015325-4c4f7f33c9ed
)

replace (
	code.cloudfoundry.org/garden => ../garden
	code.cloudfoundry.org/guardian => ../guardian
)
