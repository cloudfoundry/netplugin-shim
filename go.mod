module code.cloudfoundry.org/netplugin-shim

go 1.21.0

toolchain go1.21.3

require (
	code.cloudfoundry.org/guardian v0.0.0-20231031182052-8789f016ea2d
	code.cloudfoundry.org/lager/v3 v3.0.2
	github.com/jessevdk/go-flags v1.5.1-0.20210607101731-3927b71304df
	github.com/onsi/ginkgo/v2 v2.13.0
	github.com/onsi/gomega v1.29.0
	golang.org/x/sys v0.14.0
)

require (
	code.cloudfoundry.org/commandrunner v0.0.0-20230612151827-2b11a2b4e9b8 // indirect
	code.cloudfoundry.org/garden v0.0.0-20231031181541-0ff6ac1ac49c // indirect
	github.com/cloudfoundry/dropsonde v1.1.0 // indirect
	github.com/cloudfoundry/sonde-go v0.0.0-20230911203642-fa89d986ae20 // indirect
	github.com/go-logr/logr v1.3.0 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/pprof v0.0.0-20231101202521-4ca4178f5c7a // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/opencontainers/runtime-spec v1.1.0 // indirect
	github.com/openzipkin/zipkin-go v0.4.2 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.14.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	code.cloudfoundry.org/garden => ../garden
	code.cloudfoundry.org/grootfs => ../grootfs
	code.cloudfoundry.org/guardian => ../guardian
	code.cloudfoundry.org/idmapper => ../idmapper
)
