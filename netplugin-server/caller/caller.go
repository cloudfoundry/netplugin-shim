package caller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"

	"code.cloudfoundry.org/commandrunner"
	"code.cloudfoundry.org/commandrunner/linux_command_runner"
	"code.cloudfoundry.org/netplugin-shim/garden-plugin/message"
	"golang.org/x/sys/unix"
)

type NetpluginCaller struct {
	path          string
	extraArgs     []string
	commandRunner commandrunner.CommandRunner
}

func New(path string, extraArgs []string) *NetpluginCaller {
	return &NetpluginCaller{
		path:          path,
		extraArgs:     extraArgs,
		commandRunner: linux_command_runner.New(),
	}
}

func (c *NetpluginCaller) WithCommandRunner(runner commandrunner.CommandRunner) *NetpluginCaller {
	c.commandRunner = runner
	return c
}

func (c *NetpluginCaller) Handle(conn net.Conn) error {
	procNSFile, err := readNsFileDescriptor(conn)
	if err != nil {
		return err
	}
	defer procNSFile.Close()

	msg, err := decodeMsg(conn)
	if err != nil {
		return err
	}

	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)

	args := append(c.extraArgs, "--action", "up", "--handle", string(msg.Handle))

	cmd := exec.Command(c.path, args...)
	cmd.ExtraFiles = []*os.File{procNSFile}
	cmd.Stdin = bytes.NewReader(msg.Data)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := c.commandRunner.Run(cmd); err != nil {
		return fmt.Errorf("Error while running netplugin: %s: %s", err.Error(), stderr.String())
	}

	reply := stdout.Bytes()
	if stderr.Bytes() != nil {
		reply = []byte(fmt.Sprintf(`{"Error": "%s"}`, stderr.String()))
	}

	if _, err := conn.Write(reply); err != nil {
		return err
	}

	return nil
}

func decodeMsg(r io.Reader) (message.Message, error) {
	var content message.Message
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&content); err != nil {
		return message.Message{}, err
	}
	return content, nil
}

func readNsFileDescriptor(conn net.Conn) (*os.File, error) {
	unixconn, ok := conn.(*net.UnixConn)
	if !ok {
		return nil, errors.New("failed to cast connection to unixconn")
	}

	fd, err := recvFD(unixconn)
	if err != nil {
		return nil, err
	}

	return os.NewFile(fd, fmt.Sprintf("fd%d", fd)), nil
}

func recvFD(conn *net.UnixConn) (uintptr, error) {
	controlMessageBytesSpace := unix.CmsgSpace(4)

	controlMessageBytes := make([]byte, controlMessageBytesSpace)
	_, readSocketControlMessageBytes, _, _, err := conn.ReadMsgUnix(nil, controlMessageBytes)
	if err != nil {
		return 0, err
	}

	if readSocketControlMessageBytes > controlMessageBytesSpace {
		return 0, errors.New("received too many things")
	}

	controlMessageBytes = controlMessageBytes[:readSocketControlMessageBytes]

	socketControlMessages, err := parseSocketControlMessage(controlMessageBytes)
	if err != nil {
		return 0, err
	}

	fds, err := parseUnixRights(&socketControlMessages[0])
	if err != nil {
		return 0, err
	}

	return uintptr(fds[0]), nil
}

func parseUnixRights(m *unix.SocketControlMessage) ([]int, error) {
	messages, err := unix.ParseUnixRights(m)
	if err != nil {
		return nil, err
	}
	if len(messages) != 1 {
		return nil, errors.New("no messages parsed")
	}
	return messages, nil
}

func parseSocketControlMessage(b []byte) ([]unix.SocketControlMessage, error) {
	messages, err := unix.ParseSocketControlMessage(b)
	if err != nil {
		return nil, err
	}
	if len(messages) != 1 {
		return nil, errors.New("no messages parsed")
	}
	return messages, nil
}
