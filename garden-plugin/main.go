package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"code.cloudfoundry.org/guardian/netplugin"
	"code.cloudfoundry.org/netplugin-shim/garden-plugin/message"
	"code.cloudfoundry.org/netplugin-shim/shimsocket"
)

func main() {
	args, err := parseArgs()
	exitOn(err)

	inputs, err := readData(os.Stdin)
	exitOn(err)

	netNSFile, err := os.Open(netNSFilepath(inputs.Pid))
	exitOn(err)
	defer netNSFile.Close()

	inputs.Pid = 0
	data, err := json.Marshal(inputs)
	exitOn(err)

	msg := message.Message{Command: []byte(args.Action), Handle: []byte(args.Handle), Data: data}
	conn, err := shimsocket.Send(args.Socket, netNSFile.Fd(), msg)
	exitOn(err)
	defer conn.Close()

	err = shimsocket.PassReply(conn, os.Stdout)
	exitOn(err)
}

func exitOn(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func readData(r io.Reader) (netplugin.UpInputs, error) {
	var upInputs netplugin.UpInputs

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return upInputs, err
	}

	err = json.NewDecoder(bytes.NewBuffer(data)).Decode(&upInputs)
	if err != nil {
		return upInputs, err
	}
	return upInputs, nil
}

func netNSFilepath(pid int) string {
	netNSFilepath := os.DevNull
	if pid != 0 {
		netNSFilepath = filepath.Join("/proc", strconv.Itoa(pid), "ns", "net")
	}

	return netNSFilepath
}
