package caller

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"

	"code.cloudfoundry.org/commandrunner"
	"code.cloudfoundry.org/commandrunner/linux_command_runner"
	"code.cloudfoundry.org/netplugin-shim/shimsocket"
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

func (c *NetpluginCaller) Handle(conn *net.UnixConn) error {
	procNSFile, msg, err := shimsocket.Receive(conn)
	if err != nil {
		return err
	}
	defer procNSFile.Close()

	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)

	args := append(c.extraArgs, "--action", string(msg.Command), "--handle", string(msg.Handle))

	cmd := exec.Command(c.path, args...)
	cmd.ExtraFiles = []*os.File{procNSFile}
	cmd.Stdin = bytes.NewReader(msg.Data)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err = c.commandRunner.Run(cmd)
	reply := stdout.Bytes()

	if err != nil {
		reply = []byte(fmt.Sprintf(`{"Error": "%v"}`, err))
	}

	if _, err := conn.Write(reply); err != nil {
		return err
	}

	return nil
}
