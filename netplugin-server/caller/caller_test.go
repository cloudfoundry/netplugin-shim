package caller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"

	"code.cloudfoundry.org/lager/lagertest"
	"code.cloudfoundry.org/netplugin-shim/garden-plugin/message"
	"code.cloudfoundry.org/netplugin-shim/netplugin-server/caller"
	"code.cloudfoundry.org/netplugin-shim/netplugin-server/caller/callerfakes"
	"code.cloudfoundry.org/netplugin-shim/shimsocket"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NetpluginCaller", func() {
	var (
		netpluginCaller   *caller.NetpluginCaller
		cmdRun            *callerfakes.FakeCommandRunner
		listener          *net.UnixListener
		socketPath        string
		tmpFileToSend     *os.File
		sendingConnection *net.UnixConn
		handleConn        *net.UnixConn
		tmpDir            string
		logger            *lagertest.TestLogger
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("foo")
		cmdRun = new(callerfakes.FakeCommandRunner)

		var err error
		tmpDir, err = ioutil.TempDir("", "netplugin-caller")
		Expect(err).NotTo(HaveOccurred())

		tmpFileToSend, err = ioutil.TempFile(tmpDir, "garden-netns-fake")
		Expect(err).NotTo(HaveOccurred())

		_, err = tmpFileToSend.WriteString("potato")
		Expect(err).NotTo(HaveOccurred())

		socketPath = filepath.Join(tmpDir, "net-socket.sock")

		addr, err := net.ResolveUnixAddr("unix", socketPath)
		Expect(err).NotTo(HaveOccurred())

		listener, err = net.ListenUnix("unix", addr)
		Expect(err).NotTo(HaveOccurred())

		netpluginCaller = caller.New(logger, "/path/to/plugin", []string{"--configFile", "/path/to/config"}).WithCommandRunner(cmdRun.Spy)
	})

	AfterEach(func() {
		Expect(os.RemoveAll(tmpDir)).To(Succeed())
		sendingConnection.Close()
		handleConn.Close()
	})

	Context("when a message is received on the socket", func() {
		var (
			msg                 message.Message
			executedCommand     *exec.Cmd
			netpluginActionFunc func(cmd *exec.Cmd) error
			replyBuffer         *bytes.Buffer
		)

		BeforeEach(func() {
			msg = message.Message{
				Command: []byte("up"),
				Handle:  []byte("containery"),
				Data:    []byte(`{"Pid": 1001}`),
			}

			netpluginActionFunc = func(cmd *exec.Cmd) error {
				fmt.Fprint(cmd.Stdout, `{"Hey": "I succeeded"}`)
				return nil
			}

			replyBuffer = new(bytes.Buffer)
		})

		JustBeforeEach(func() {
			var err error
			sendingConnection, err = shimsocket.Send(socketPath, tmpFileToSend.Fd(), msg)
			Expect(err).NotTo(HaveOccurred())

			handleConn, err = listener.AcceptUnix()
			Expect(err).NotTo(HaveOccurred())

			cmdRun.Calls(func(cmd *exec.Cmd) error { return netpluginActionFunc(cmd) })

			netpluginCaller.Handle(handleConn)
			Expect(cmdRun.CallCount()).To(Equal(1))
			executedCommand = cmdRun.ArgsForCall(0)
		})

		It("calls the correct plugin", func() {
			Expect(executedCommand.Path).To(Equal("/path/to/plugin"))
		})

		It("calls the plugin with the correct args", func() {
			expectedArguments := []string{"--configFile", "/path/to/config", "--action", "up", "--handle", "containery"}
			Expect(executedCommand.Args[1:]).To(Equal(expectedArguments))
		})

		It("sends the fd from the socket as an extra file", func() {
			Expect(len(executedCommand.ExtraFiles)).To(Equal(1))
		})

		It("provides the data from the socket as stdin", func() {
			contents, err := ioutil.ReadAll(executedCommand.Stdin)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(contents)).To(MatchJSON(`{"Pid": 1001}`))
		})

		It("propagates the output of the netplugin to the socket", func() {
			Expect(shimsocket.PassReply(sendingConnection, replyBuffer)).To(Succeed())
			Expect(replyBuffer.String()).To(MatchJSON(`{"Hey": "I succeeded"}`))
		})

		When("the netplugin fails", func() {
			BeforeEach(func() {
				netpluginActionFunc = func(cmd *exec.Cmd) error {
					return errors.New("Error executing plugin")
				}
			})

			It("propagate the error back to the socket", func() {
				Expect(shimsocket.PassReply(sendingConnection, replyBuffer)).NotTo(Succeed())
				Expect(replyBuffer.String()).To(MatchJSON(`{"Error": "Error executing plugin"}`))
			})
		})

		When("the command is down", func() {
			BeforeEach(func() {
				msg = message.Message{
					Command: []byte("down"),
					Handle:  []byte("cake"),
					Data:    []byte{},
				}
			})

			It("calls the plugin with the down command as an argument", func() {
				expectedArguments := []string{"--configFile", "/path/to/config", "--action", "down", "--handle", "cake"}
				Expect(executedCommand.Args[1:]).To(Equal(expectedArguments))
			})
		})

		When("the network plugin outputs nothing", func() {
			BeforeEach(func() {
				netpluginActionFunc = func(cmd *exec.Cmd) error {
					return nil
				}
			})

			It("sends valid JSON to the socket", func() {
				var output map[string]interface{}
				handleConn.Close()
				Expect(json.NewDecoder(sendingConnection).Decode(&output)).To(Succeed())
			})
		})
	})
})
