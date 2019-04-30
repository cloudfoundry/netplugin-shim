package integration_test

import (
	"io/ioutil"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

var (
	netpluginServerPath string
	fakeNetpluginPath   string
)

var _ = BeforeSuite(func() {
	var err error
	netpluginServerPath, err = gexec.Build("code.cloudfoundry.org/netplugin-shim/netplugin-server", "-mod=vendor")
	Expect(err).NotTo(HaveOccurred())
	fakeNetpluginPath, err = gexec.Build("code.cloudfoundry.org/netplugin-shim/netplugin-server/fake-network-plugin")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Integration Suite")
}

func tempDir(dir, prefix string) string {
	name, err := ioutil.TempDir(dir, prefix)
	Expect(err).NotTo(HaveOccurred())
	return name
}

func tempFile(dir string) *os.File {
	file, err := ioutil.TempFile(dir, "")
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return file
}

func gexecStart(cmd *exec.Cmd) *gexec.Session {
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	return session
}
