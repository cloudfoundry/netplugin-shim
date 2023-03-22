package integration_test

import (
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Plugin Integration Suite")
}

var (
	pluginPath string
	daemonPath string
)

var _ = BeforeSuite(func() {
	var err error
	pluginPath, err = gexec.Build("code.cloudfoundry.org/netplugin-shim/garden-plugin", "-mod=vendor")
	Expect(err).NotTo(HaveOccurred())

	daemonPath, err = gexec.Build("code.cloudfoundry.org/netplugin-shim/garden-plugin/fake-netplugin-daemon", "-mod=vendor")
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
