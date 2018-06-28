package shimsocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"

	"code.cloudfoundry.org/netplugin-shim/garden-plugin/message"
	"golang.org/x/sys/unix"
)

func Receive(conn *net.UnixConn) (*os.File, message.Message, error) {
	fd, err := receiveFD(conn)
	if err != nil {
		return nil, message.Message{}, err
	}
	nsFile := os.NewFile(fd, fmt.Sprintf("fd%d", fd))

	msg, err := decodeMsg(conn)
	if err != nil {
		defer nsFile.Close()
		return nil, message.Message{}, err
	}

	return nsFile, msg, nil
}

func decodeMsg(r io.Reader) (message.Message, error) {
	var content message.Message
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&content); err != nil {
		return message.Message{}, err
	}
	return content, nil
}

func receiveFD(conn *net.UnixConn) (uintptr, error) {
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
