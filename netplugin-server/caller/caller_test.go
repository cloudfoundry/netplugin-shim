package caller_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"

	"code.cloudfoundry.org/commandrunner/fake_command_runner"
	"code.cloudfoundry.org/netplugin-shim/garden-plugin/message"
	"code.cloudfoundry.org/netplugin-shim/netplugin-server/caller"
	"code.cloudfoundry.org/netplugin-shim/shimsocket"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NetpluginCaller", func() {
	var (
		netpluginCaller   *caller.NetpluginCaller
		commandRunner     *fake_command_runner.FakeCommandRunner
		listener          net.Listener
		socketPath        string
		tmpFileToSend     *os.File
		sendingConnection *net.UnixConn
		handleConn        net.Conn
		tmpDir            string
	)

	BeforeEach(func() {
		commandRunner = fake_command_runner.New()

		var err error
		tmpDir, err = ioutil.TempDir("", "netplugin-caller")
		Expect(err).NotTo(HaveOccurred())

		tmpFileToSend, err = ioutil.TempFile(tmpDir, "garden-netns-fake")
		Expect(err).NotTo(HaveOccurred())

		_, err = tmpFileToSend.WriteString("potato")
		Expect(err).NotTo(HaveOccurred())

		socketPath = filepath.Join(tmpDir, "net-socket.sock")

		listener, err = net.Listen("unix", socketPath)
		Expect(err).NotTo(HaveOccurred())

		netpluginCaller = caller.New("/path/to/plugin", []string{"--configFile", "/path/to/config"}).WithCommandRunner(commandRunner)
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

			handleConn, err = listener.Accept()
			Expect(err).NotTo(HaveOccurred())

			commandRunner.WhenRunning(fake_command_runner.CommandSpec{
				Path: "/path/to/plugin",
			}, func(cmd *exec.Cmd) error {
				return netpluginActionFunc(cmd)
			})

			netpluginCaller.Handle(handleConn)
			executedCommands := commandRunner.ExecutedCommands()
			Expect(len(executedCommands)).To(Equal(1))
			executedCommand = executedCommands[0]
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
	})
})
