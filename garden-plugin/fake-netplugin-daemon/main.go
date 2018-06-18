package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"

	"code.cloudfoundry.org/netplugin-shim/garden-plugin/message"
	"golang.org/x/sys/unix"
)

type PluginInstruction struct {
	FD      uintptr
	Message message.Message
}

func main() {
	args, err := parseArgs()
	if err != nil {
		panic(err)
	}
	listener, err := net.Listen("unix", args.Socket)
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		fd, err := getFD(conn)
		if err != nil {
			panic(err)
		}
		if err = ioutil.WriteFile(args.FDFile, []byte(fmt.Sprintf("%d", fd)), os.ModePerm); err != nil {
			panic(err)
		}

		msg, err := decodeMsg(conn)
		if err != nil {
			panic(err)
		}

		jsonMessage, err := json.Marshal(msg)
		if err != nil {
			panic(err)
		}

		if err = ioutil.WriteFile(args.MessageFile, jsonMessage, os.ModePerm); err != nil {
			panic(err)
		}

		reply := []byte(args.Reply)
		var n int
		n, err = conn.Write(reply)
		if err != nil {
			panic(err)
		}
		if n != len(reply) {
			panic(err)
		}

		conn.Close()
	}
}

func getFD(conn net.Conn) (uintptr, error) {
	unixconn, ok := conn.(*net.UnixConn)
	if !ok {
		return 0, errors.New("failed to cast connection to unixconn")
	}

	return recvFD(unixconn)
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

func decodeMsg(r io.Reader) (message.Message, error) {
	var content message.Message
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&content); err != nil {
		return message.Message{}, err
	}
	return content, nil
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
