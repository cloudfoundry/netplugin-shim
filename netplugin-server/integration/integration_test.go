package integration_test

import (
	"net"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Integration", func() {
	var (
		tmpDir  string
		socket  string
		session *gexec.Session
	)

	BeforeEach(func() {
		tmpDir = tempDir("", "netplugin-server-integration")
		socket = filepath.Join(tmpDir, "sock.sock")
	})

	JustBeforeEach(func() {
		session = gexecStart(exec.Command(
			netpluginServerPath,
			"--socket", socket,
			"--netpluginPath=/does/not/exist",
			`--netpluginArgs="--configFile"`,
			`--netpluginArgs="/not/a/config/file"`,
		))
	})

	AfterEach(func() {
		Expect(session.Terminate().Wait()).To(gexec.Exit())
		Expect(os.RemoveAll(tmpDir)).To(Succeed())
	})

	It("listens on the provided socket path", func() {
		dial := func() error {
			_, err := net.Dial("unix", filepath.Join(tmpDir, "sock.sock"))
			return err
		}
		Eventually(dial).Should(Succeed())
	})

	When("the server exits", func() {
		It("cleans up the socket", func() {
			Expect(session.Terminate().Wait()).To(gexec.Exit())
			Eventually(socket).ShouldNot(BeAnExistingFile())
		})

		It("writes nothing to stderr", func() {
			Expect(session.Terminate().Wait()).To(gexec.Exit())
			Expect(session.Err.Contents()).To(BeEmpty())
		})
	})
})
