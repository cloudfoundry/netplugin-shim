package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"

	"code.cloudfoundry.org/guardian/netplugin"
	"code.cloudfoundry.org/netplugin-shim/message"
	"golang.org/x/sys/unix"
)

func main() {
	args, err := parseArgs()
	if err != nil {
		failWith(err)
	}

	conn, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: args.Socket, Net: "unix"})
	if err != nil {
		failWith(err)
	}
	defer conn.Close()

	data, pid, err := readDataAndPID(os.Stdin)
	if err != nil {
		failWith(err)
	}

	if err := writeNetNSFD(conn, pid); err != nil {
		failWith(err)
	}

	msg := message.Message{Command: args.Positional.Command, Data: string(data)}
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(&msg); err != nil {
		failWith(err)
	}

	if err := writeReply(conn); err != nil {
		failWith(err)
	}
}

func failWith(err error) {
	fmt.Fprint(os.Stderr, err.Error())
	os.Exit(1)
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
	file, err := os.Open(fmt.Sprintf("/proc/%d/ns/net", pid))
	if err != nil {
		return err
	}

	socketControlMessage := unix.UnixRights(int(file.Fd()))
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
