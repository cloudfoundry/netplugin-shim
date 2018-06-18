package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"code.cloudfoundry.org/netplugin-shim/netplugin-server/caller"
	"code.cloudfoundry.org/netplugin-shim/netplugin-server/service"
	"golang.org/x/sys/unix"
)

func main() {
	args, err := parseArgs()
	exitOn(err)

	addr, err := net.ResolveUnixAddr("unix", args.SocketPath)
	exitOn(err)

	listener, err := net.ListenUnix("unix", addr)
	exitOn(err)

	netpluginCaller := caller.New(args.NetpluginPath, args.NetpluginArgs)

	server := service.New(netpluginCaller.Handle).WithLogger(os.Stderr)

	go server.Serve(listener)

	done := make(chan os.Signal)
	signal.Notify(done, unix.SIGINT, unix.SIGTERM)
	<-done

	server.Stop()
}

func exitOn(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}
