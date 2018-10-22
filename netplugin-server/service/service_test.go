package service_test

import (
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/netplugin-shim/netplugin-server/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Service", func() {
	var (
		workDir  string
		address  *net.UnixAddr
		listener *net.UnixListener
		buffer   *gbytes.Buffer
	)

	BeforeEach(func() {
		workDir = tempDir()
		address = resolveUnixAddr(filepath.Join(workDir, "socket.sock"))
		listener = listenUnix(address)
		buffer = gbytes.NewBuffer()
	})

	AfterEach(func() {
		Expect(os.RemoveAll(workDir)).To(Succeed())
	})

	Describe("Serve", func() {
		It("serves a listener", func() {
			message := "aaaaaaaaaa"
			s := service.New(handleFunc(buffer, int64(len(message)))).WithLogger(GinkgoWriter)
			go s.Serve(listener)

			conn, err := net.DialUnix("unix", nil, address)
			Expect(err).NotTo(HaveOccurred())
			defer conn.Close()
			writeString(conn, message)

			Eventually(buffer).Should(gbytes.Say(message))

			s.Stop()
		})
	})

	It("should fail to initiate new connections after stop", func() {
		s := service.New(func(*net.UnixConn) error { return nil }).WithLogger(GinkgoWriter)
		go s.Serve(listener)

		conn, err := net.DialUnix("unix", nil, address)
		Expect(err).NotTo(HaveOccurred())
		conn.Close()

		s.Stop()

		_, err = net.DialUnix("unix", nil, address)
		Expect(err).To(HaveOccurred())
	})
})

func handleFunc(buffer *gbytes.Buffer, limit int64) func(*net.UnixConn) error {
	return func(conn *net.UnixConn) error {
		reader := &io.LimitedReader{R: conn, N: limit}
		data, err := ioutil.ReadAll(reader)
		if err != nil {
			return err
		}

		_, err = buffer.Write(data)
		if err != nil {
			return err
		}

		return nil
	}
}
