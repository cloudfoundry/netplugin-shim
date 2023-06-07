module code.cloudfoundry.org/netplugin-shim

go 1.20

require (
	code.cloudfoundry.org/guardian v0.0.0-20230602145355-4b998fcefc2a
	code.cloudfoundry.org/lager/v3 v3.0.1
	github.com/jessevdk/go-flags v1.5.1-0.20210607101731-3927b71304df
	github.com/onsi/ginkgo/v2 v2.10.0
	github.com/onsi/gomega v1.27.7
	golang.org/x/sys v0.8.0
)

require (
	code.cloudfoundry.org/commandrunner v0.0.0-20230427153105-c662e812fa6f // indirect
	code.cloudfoundry.org/garden v0.0.0-20230522165130-8444ff5a31d7 // indirect
	github.com/cloudfoundry/dropsonde v1.1.0 // indirect
	github.com/cloudfoundry/sonde-go v0.0.0-20230601214134-20e0b3c2ed5c // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/pprof v0.0.0-20230602150820-91b7bce49751 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/opencontainers/runtime-spec v1.1.0-rc.2 // indirect
	github.com/openzipkin/zipkin-go v0.4.1 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	golang.org/x/tools v0.9.3 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	code.cloudfoundry.org/garden => ../garden
	code.cloudfoundry.org/grootfs => ../grootfs
	code.cloudfoundry.org/guardian => ../guardian
	code.cloudfoundry.org/idmapper => ../idmapper
)
