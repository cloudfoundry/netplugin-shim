package integration_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"code.cloudfoundry.org/guardian/netplugin"
	"code.cloudfoundry.org/netplugin-shim/garden-plugin/message"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"golang.org/x/sys/unix"
)

const waitDuration = time.Second * 10

var _ = Describe("Integration", func() {
	var (
		workDir     string
		socketPath  string
		fdPath      string
		messagePath string

		shimCmd     *exec.Cmd
		shimSession *gexec.Session

		initSession *gexec.Session

		daemonSession *gexec.Session
		daemonReply   []byte
	)

	BeforeEach(func() {
		workDir = tempDir("", "")
		socketPath = filepath.Join(workDir, "test.sock")
		fdPath = filepath.Join(workDir, "fd")
		messagePath = filepath.Join(workDir, "message")

		shimCmd = exec.Command(pluginPath, "--socket", socketPath)

		daemonReply = []byte(`{"Here":"Be Dragons"}`)

		initSession = gexecStart(initCommand())
	})

	JustBeforeEach(func() {
		daemonCmd := exec.Command(
			daemonPath,
			"--socket", socketPath,
			"--reply", string(daemonReply),
			"--fd-file", fdPath,
			"--message-file", messagePath,
		)
		daemonSession = gexecStart(daemonCmd)
		waitForFileToExist(socketPath)

		shimSession = gexecStart(shimCmd)
	})

	AfterEach(func() {
		Expect(initSession.Terminate().Wait(waitDuration)).To(gexec.Exit())
		Expect(shimSession.Wait(waitDuration)).To(gexec.Exit())
		Expect(daemonSession.Terminate().Wait(waitDuration)).To(gexec.Exit())
		Expect(os.RemoveAll(workDir)).To(Succeed())
	})

	Describe("Up", func() {
		var upInputs netplugin.UpInputs

		BeforeEach(func() {
			shimCmd.Args = append(shimCmd.Args, "--action", "up", "--handle", "potato")

			upInputs = netplugin.UpInputs{
				Pid: initSession.Command.Process.Pid,
			}
			shimCmd.Stdin = strings.NewReader(encode(upInputs))
		})

		It("exits successfully", func() {
			Expect(shimSession.Wait(waitDuration)).To(gexec.Exit(0))
		})

		It("sends the net ns fd of the provided pid to the socket", func() {
			waitForFileToExist(fdPath)
			fd := atoi(readFileAsString(fdPath))
			name, err := os.Readlink(fmt.Sprintf("/proc/%d/fd/%d", daemonSession.Command.Process.Pid, fd))
			Expect(err).NotTo(HaveOccurred())
			Expect(name).To(Equal(parseNetNS(initSession.Command.Process.Pid)))
		})

		It("sends the command to the provided socket", func() {
			waitForFileToExist(messagePath)
			message := decodeMessage(strings.NewReader(readFileAsString(messagePath)))
			Expect(string(message.Command)).To(Equal("up"))
		})

		It("sends the handle to the provided socket", func() {
			waitForFileToExist(messagePath)
			message := decodeMessage(strings.NewReader(readFileAsString(messagePath)))
			Expect(string(message.Handle)).To(Equal("potato"))
		})

		It("includes stdin contents with the pid set to 0 in the message sent to the socket", func() {
			waitForFileToExist(messagePath)
			message := decodeMessage(strings.NewReader(readFileAsString(messagePath)))
			Expect(string(message.Data)).To(MatchJSON(`{"Pid":0,"Properties":null}`))
		})

		It("writes JSON to stdout", func() {
			Expect(shimSession.Wait(waitDuration)).To(gexec.Exit())
			stdout := struct{}{}
			Expect(json.Unmarshal(shimSession.Out.Contents(), &stdout)).To(Succeed())
		})

		It("writes the network daemon's response to stdout", func() {
			Expect(shimSession.Wait(waitDuration)).To(gexec.Exit())
			stdout := string(shimSession.Out.Contents())
			Expect(strings.TrimSpace(stdout)).To(MatchJSON(`{"Here":"Be Dragons"}`))
		})

		Context("when the network daemon reports an error", func() {
			BeforeEach(func() {
				daemonReply = []byte(`{"Error":"no dragons received"}`)
			})

			It("writes the response to stderr", func() {
				Expect(shimSession.Wait(waitDuration)).To(gexec.Exit())
				stderr := string(shimSession.Err.Contents())
				Expect(stderr).To(ContainSubstring("no dragons received"))
			})

			It("exits non zero", func() {
				Expect(shimSession.Wait(waitDuration)).NotTo(gexec.Exit(0))
			})

		})
	})

	Describe("Down", func() {
		BeforeEach(func() {
			shimCmd.Args = append(shimCmd.Args, "--action", "down", "--handle", "potato")
			shimCmd.Stdin = strings.NewReader(encode(nil))
		})

		It("exits successfully", func() {
			Expect(shimSession.Wait(waitDuration)).To(gexec.Exit(0))
		})

		It("sends the command to the socket", func() {
			waitForFileToExist(messagePath)
			message := decodeMessage(strings.NewReader(readFileAsString(messagePath)))
			Expect(string(message.Command)).To(Equal("down"))
		})

		It("sends an fd to /dev/null", func() {
			waitForFileToExist(fdPath)
			fd := atoi(readFileAsString(fdPath))

			name, err := os.Readlink(fmt.Sprintf("/proc/%d/fd/%d", daemonSession.Command.Process.Pid, fd))
			Expect(err).NotTo(HaveOccurred())
			Expect(name).To(Equal("/dev/null"))
		})
	})
})

func encode(thing interface{}) string {
	bytes, err := json.Marshal(thing)
	Expect(err).NotTo(HaveOccurred())
	return string(bytes)
}

func tempDir(dir, prefix string) string {
	name, err := ioutil.TempDir(dir, prefix)
	Expect(err).NotTo(HaveOccurred())
	return name
}

func parseNetNS(pid int) string {
	netNS, err := os.Readlink(fmt.Sprintf("/proc/%d/ns/net", pid))
	Expect(err).NotTo(HaveOccurred())

	return strings.TrimSpace(netNS)
}

func initCommand() *exec.Cmd {
	cmd := exec.Command("sleep", "3600")
	cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: unix.CLONE_NEWUSER | unix.CLONE_NEWNET}
	return cmd
}

func readFileAsString(path string) string {
	content, err := ioutil.ReadFile(path)
	Expect(err).NotTo(HaveOccurred())
	return string(content)
}

func atoi(i string) int {
	s, err := strconv.Atoi(i)
	Expect(err).NotTo(HaveOccurred())
	return s
}

func decodeMessage(r io.Reader) message.Message {
	var content message.Message
	decoder := json.NewDecoder(r)
	Expect(decoder.Decode(&content)).To(Succeed())
	return content
}

func waitForFileToExist(filepath string) {
	for {
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			time.Sleep(time.Second)
		} else {
			break
		}
	}

	Expect(filepath).To(BeAnExistingFile())
}
