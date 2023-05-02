module code.cloudfoundry.org/netplugin-shim

go 1.20

require (
	code.cloudfoundry.org/guardian v0.0.0-20230418163503-e44f10a60f42
	code.cloudfoundry.org/lager/v3 v3.0.1
	github.com/jessevdk/go-flags v1.5.1-0.20210607101731-3927b71304df
	github.com/onsi/ginkgo/v2 v2.9.2
	github.com/onsi/gomega v1.27.6
	golang.org/x/sys v0.7.0
)

require (
	code.cloudfoundry.org/commandrunner v0.0.0-20230406174834-03aed0535eb9 // indirect
	code.cloudfoundry.org/garden v0.0.0-20230418162742-6c127706d54f // indirect
	github.com/apoydence/eachers v0.0.0-20181020210610-23942921fe77 // indirect
	github.com/cloudfoundry/dropsonde v1.0.0 // indirect
	github.com/cloudfoundry/sonde-go v0.0.0-20230412182205-eaf74d09b55a // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/pprof v0.0.0-20230406165453-00490a63f317 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/opencontainers/runtime-spec v1.1.0-rc.2 // indirect
	github.com/openzipkin/zipkin-go v0.4.1 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	golang.org/x/tools v0.8.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	code.cloudfoundry.org/garden => ../garden
	code.cloudfoundry.org/grootfs => ../grootfs
	code.cloudfoundry.org/guardian => ../guardian
	code.cloudfoundry.org/idmapper => ../idmapper
)
