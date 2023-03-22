package service_test

import (
	"io"
	"io/ioutil"
	"net"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}

func tempDir() string {
	name, err := ioutil.TempDir("", "")
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return name
}

func resolveUnixAddr(address string) *net.UnixAddr {
	addr, err := net.ResolveUnixAddr("unix", address)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return addr
}

func listenUnix(laddr *net.UnixAddr) *net.UnixListener {
	listener, err := net.ListenUnix("unix", laddr)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return listener
}

func dialUnix(raddr *net.UnixAddr) *net.UnixConn {
	conn, err := net.DialUnix("unix", nil, raddr)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return conn
}

func writeString(w io.Writer, data string) {
	_, err := io.WriteString(w, data)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
}
