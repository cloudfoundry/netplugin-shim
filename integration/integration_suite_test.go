package integration_test

import (
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var (
	pluginPath string
	daemonPath string
)

var _ = BeforeSuite(func() {
	var err error
	pluginPath, err = gexec.Build("code.cloudfoundry.org/netplugin-shim")
	Expect(err).NotTo(HaveOccurred())

	daemonPath, err = gexec.Build("code.cloudfoundry.org/netplugin-shim/fake-netplugin-daemon")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func gexecStart(cmd *exec.Cmd) *gexec.Session {
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	return session
}
