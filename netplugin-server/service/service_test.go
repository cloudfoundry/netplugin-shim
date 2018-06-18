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

			conn := dialUnix(address)
			defer conn.Close()
			writeString(conn, message)

			Eventually(buffer).Should(gbytes.Say(message))

			s.Stop()
		})
	})

	It("should fail when writing to the socket after stop", func() {
		message := "bbbbbbbbbb"
		s := service.New(handleFunc(buffer, int64(len(message)))).WithLogger(GinkgoWriter)
		go s.Serve(listener)
		conn := dialUnix(address)
		defer conn.Close()

		s.Stop()

		_, err := io.WriteString(conn, "bbbbbbbbbb")

		Expect(err).To(HaveOccurred())
	})
})

func handleFunc(buffer *gbytes.Buffer, limit int64) func(net.Conn) error {
	return func(conn net.Conn) error {
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
