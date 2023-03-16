module code.cloudfoundry.org/netplugin-shim

go 1.19

require (
	code.cloudfoundry.org/guardian v0.0.0-20220607160814-bbdc1696f4d2
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/jessevdk/go-flags v1.5.1-0.20210607101731-3927b71304df
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.27.4
	golang.org/x/sys v0.6.0
)

require (
	code.cloudfoundry.org/commandrunner v0.0.0-20180212143422-501fd662150b // indirect
	code.cloudfoundry.org/garden v0.0.0-20230109170928-83a99bb745e3 // indirect
	github.com/apoydence/eachers v0.0.0-20181020210610-23942921fe77 // indirect
	github.com/cloudfoundry/dropsonde v1.0.0 // indirect
	github.com/cloudfoundry/sonde-go v0.0.0-20220627221915-ff36de9c3435 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/opencontainers/runtime-spec v1.1.0-rc.1 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	code.cloudfoundry.org/garden => ../garden
	code.cloudfoundry.org/grootfs => ../grootfs
	code.cloudfoundry.org/guardian => ../guardian
	code.cloudfoundry.org/idmapper => ../idmapper
)
