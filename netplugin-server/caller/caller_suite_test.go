package caller_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestNetpluginCaller(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NetpluginCaller Suite")
}
