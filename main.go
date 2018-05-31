package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"code.cloudfoundry.org/guardian/netplugin"
	"code.cloudfoundry.org/netplugin-shim/message"
	"golang.org/x/sys/unix"
)

func main() {
	args, err := parseArgs()
	exitOn(err)

	conn, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: args.Socket, Net: "unix"})
	exitOn(err)
	defer conn.Close()

	data, pid, err := readDataAndPID(os.Stdin)
	exitOn(err)

	err = writeNetNSFD(conn, pid)
	exitOn(err)

	msg := message.Message{Command: args.Action, Handle: args.Handle, Data: string(data)}
	encoder := json.NewEncoder(conn)

	err = encoder.Encode(&msg)
	exitOn(err)

	err = writeReply(conn)
	exitOn(err)
}

func exitOn(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func readDataAndPID(r io.Reader) ([]byte, int, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, 0, err
	}

	var upInputs netplugin.UpInputs
	if err = json.NewDecoder(bytes.NewBuffer(data)).Decode(&upInputs); err != nil {
		return nil, 0, err
	}

	return data, upInputs.Pid, nil
}

func writeNetNSFD(socket *net.UnixConn, pid int) error {
	// Always send an FD over the socket, but it will only be
	// an FD to the net ns of process if the provided pid != 0.
	// This allows the same execution path for both "up" and "down" commands.
	netNSFilepath := os.DevNull
	if pid != 0 {
		netNSFilepath = filepath.Join("/proc", strconv.Itoa(pid), "ns", "net")
	}

	netNSFile, err := os.Open(netNSFilepath)
	if err != nil {
		return err
	}
	defer netNSFile.Close()

	socketControlMessage := unix.UnixRights(int(netNSFile.Fd()))
	_, _, err = socket.WriteMsgUnix(nil, socketControlMessage, nil)
	return err
}

func writeReply(conn net.Conn) error {
	var output map[string]interface{}
	if err := json.NewDecoder(conn).Decode(&output); err != nil {
		return err
	}

	errString, failed := output["Error"]
	if !failed {
		return json.NewEncoder(os.Stdout).Encode(output)
	}

	return fmt.Errorf("%v", errString)
}
