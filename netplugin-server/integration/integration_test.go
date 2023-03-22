package integration_test

import (
	"encoding/json"
	"net"
	"os"
	"os/exec"
	"path/filepath"

	"code.cloudfoundry.org/netplugin-shim/garden-plugin/message"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"golang.org/x/sys/unix"
)

const responseJson = `{"foo": "bar"}`

var _ = Describe("Integration", func() {
	var (
		tmpDir               string
		socket               string
		session              *gexec.Session
		tmpFile              *os.File
		argsFilePath         string
		stdinFilePath        string
		extrafileOutFilePath string
	)

	BeforeEach(func() {
		tmpDir = tempDir("", "netplugin-server-integration")
		socket = filepath.Join(tmpDir, "sock.sock")
		tmpFile = tempFile(tmpDir)
		argsFilePath = filepath.Join(tmpDir, "args.txt")
		stdinFilePath = filepath.Join(tmpDir, "stdin.txt")
		extrafileOutFilePath = filepath.Join(tmpDir, "extrafile.txt")
	})

	JustBeforeEach(func() {
		session = gexecStart(exec.Command(
			netpluginServerPath,
			"--socket", socket,
			"--plugin-path", fakeNetpluginPath,
			"--plugin-arg", argsFilePath,
			"--plugin-arg", stdinFilePath,
			"--plugin-arg", responseJson,
			"--plugin-arg", extrafileOutFilePath,
		))
	})

	AfterEach(func() {
		Expect(session.Terminate().Wait()).To(gexec.Exit())
		Expect(os.RemoveAll(tmpDir)).To(Succeed())
	})

	It("creates the socket with proper permissions", func() {
		Eventually(socket).Should(BeAnExistingFile())
		info, err := os.Stat(socket)
		Expect(err).NotTo(HaveOccurred())

		Expect(info.Mode() & os.ModePerm).To(Equal(os.FileMode(0622)))
	})

	It("listens on the provided socket path", func() {
		var conn net.Conn
		dial := func() error {
			var err error
			conn, err = net.Dial("unix", filepath.Join(tmpDir, "sock.sock"))
			return err
		}
		Eventually(dial).Should(Succeed())
		Expect(conn.Close()).To(Succeed())
	})

	It("calls the netplugin and receives a response", func() {
		conn := dialAndSendUp(socket, tmpFile)

		response, err := readResponseFromSocket(conn)
		Expect(err).NotTo(HaveOccurred())

		Expect(response["foo"]).To(Equal("bar"))
	})

	When("the net plugin fails", func() {
		BeforeEach(func() {
			argsFilePath = "/does/not/exist"
		})

		It("returns the error over stderr", func() {
			conn := dialAndSendUp(socket, tmpFile)

			response, err := readResponseFromSocket(conn)
			Expect(err).NotTo(HaveOccurred())

			Expect(response["Error"]).NotTo(Equal(""))
		})
	})

	When("the server exits", func() {
		JustBeforeEach(func() {
			session.Terminate()
		})

		It("terminates successfully", func() {
			Eventually(session).Should(gexec.Exit())
		})

		It("cleans up the socket", func() {
			Eventually(socket).ShouldNot(BeAnExistingFile())
		})

		It("writes nothing to stderr", func() {
			Consistently(session.Err.Contents()).Should(BeEmpty())
		})
	})
})

func dialAndSendUp(socket string, tmpFile *os.File) *net.UnixConn {
	address, err := net.ResolveUnixAddr("unix", socket)
	Expect(err).NotTo(HaveOccurred())
	var conn *net.UnixConn

	dial := func() error {
		var err error
		conn, err = net.DialUnix("unix", nil, address)
		return err
	}
	Eventually(dial).Should(Succeed())

	message := message.New("up", "cake", []byte{})
	sendSocketMessage(conn, tmpFile.Fd(), message)

	return conn
}

func sendSocketMessage(conn *net.UnixConn, fd uintptr, msg message.Message) {
	socketControlMessage := unix.UnixRights(int(fd))
	_, _, err := conn.WriteMsgUnix(nil, socketControlMessage, nil)
	Expect(err).NotTo(HaveOccurred())

	encoder := json.NewEncoder(conn)
	Expect(encoder.Encode(msg)).To(Succeed())
}

func readResponseFromSocket(conn *net.UnixConn) (map[string]interface{}, error) {
	var output map[string]interface{}
	if err := json.NewDecoder(conn).Decode(&output); err != nil {
		return nil, err
	}

	return output, nil
}
