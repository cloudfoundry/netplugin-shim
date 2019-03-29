package caller

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/netplugin-shim/shimsocket"
)

//go:generate counterfeiter . CommandRunner
type CommandRunner func(*exec.Cmd) error

type NetpluginCaller struct {
	path      string
	extraArgs []string
	cmdRun    CommandRunner
	logger    lager.Logger
}

func New(logger lager.Logger, path string, extraArgs []string) *NetpluginCaller {
	return &NetpluginCaller{
		path:      path,
		extraArgs: extraArgs,
		cmdRun:    func(cmd *exec.Cmd) error { return cmd.Run() },
		logger:    logger,
	}
}

func (c *NetpluginCaller) WithCommandRunner(run CommandRunner) *NetpluginCaller {
	c.cmdRun = run
	return c
}

func (c *NetpluginCaller) Handle(conn *net.UnixConn) error {
	logger := c.logger.Session("handle")
	logger.Info("start")
	defer logger.Info("end")

	procNSFile, msg, err := shimsocket.Receive(conn)
	if err != nil {
		logger.Error("error-receiving-message", err)
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

	logger.Debug("calling-external-network-plugin", lager.Data{"procNSFile": procNSFile.Fd(), "message": msg.String()})
	err = c.cmdRun(cmd)
	reply := stdout.Bytes()

	if len(reply) == 0 {
		reply = []byte("{}")
	}

	if err != nil {
		logger.Error("error-calling-external-network-plugin", err, lager.Data{"procNSFile": procNSFile.Fd(), "message": msg.String(), "stdout": stdout.String(), "stderr": stderr.String()})
		reply = []byte(fmt.Sprintf(`{"Error": "%v"}`, err))
	}

	if _, err := conn.Write(reply); err != nil {
		logger.Error("error-responding", err, lager.Data{"procNSFile": procNSFile.Fd(), "message": msg.String()})
		return err
	}

	return nil
}
