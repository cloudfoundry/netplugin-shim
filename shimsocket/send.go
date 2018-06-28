package shimsocket

import (
	"encoding/json"
	"fmt"
	"io"
	"net"

	"code.cloudfoundry.org/netplugin-shim/garden-plugin/message"
	"golang.org/x/sys/unix"
)

func Send(socketPath string, fd uintptr, msg message.Message) (*net.UnixConn, error) {
	conn, err := dial(socketPath)
	if err != nil {
		return nil, err
	}

	err = sendFD(conn, fd)
	if err != nil {
		conn.Close()
		return nil, err
	}

	err = sendMessage(conn, msg)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

func PassReply(conn net.Conn, writer io.Writer) error {
	var output map[string]interface{}
	if err := json.NewDecoder(conn).Decode(&output); err != nil {
		return err
	}

	err := json.NewEncoder(writer).Encode(output)
	if err != nil {
		return err
	}

	errString, found := output["Error"]
	if found {
		return fmt.Errorf("%v", errString)
	}

	return nil
}

func dial(socketPath string) (*net.UnixConn, error) {
	address, err := net.ResolveUnixAddr("unix", socketPath)
	if err != nil {
		return nil, err
	}

	return net.DialUnix("unix", nil, address)
}

func sendFD(conn *net.UnixConn, fd uintptr) error {
	socketControlMessage := unix.UnixRights(int(fd))
	_, _, err := conn.WriteMsgUnix(nil, socketControlMessage, nil)
	return err
}

func sendMessage(conn *net.UnixConn, msg message.Message) error {
	return json.NewEncoder(conn).Encode(&msg)
}
