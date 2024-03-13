package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"code.cloudfoundry.org/netplugin-shim/garden-plugin/message"
	"code.cloudfoundry.org/netplugin-shim/shimsocket"
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

	addr, err := net.ResolveUnixAddr("unix", args.Socket)
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenUnix("unix", addr)
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	var procNSFile *os.File
	var msg message.Message

	defer procNSFile.Close()

	for {
		var err error
		conn, err := listener.AcceptUnix()
		if err != nil {
			panic(err)
		}

		procNSFile, msg, err = shimsocket.Receive(conn)
		if err != nil {
			panic(err)
		}

		if err = os.WriteFile(args.FDFile, []byte(fmt.Sprintf("%d", procNSFile.Fd())), os.ModePerm); err != nil {
			panic(err)
		}

		jsonMessage, err := json.Marshal(msg)
		if err != nil {
			panic(err)
		}

		if err = os.WriteFile(args.MessageFile, jsonMessage, os.ModePerm); err != nil {
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
