package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"code.cloudfoundry.org/lager"
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

	err = os.Chmod(args.SocketPath, 0622)
	exitOn(err)

	log := initLogger("netplugin-server")
	netpluginCaller := caller.New(log, args.NetpluginPath, args.NetpluginArgs)

	server := service.New(netpluginCaller.Handle).WithLogger(os.Stderr)

	go server.Serve(listener)

	done := make(chan os.Signal, 1)
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

func initLogger(component string) lager.Logger {
	internalSink := lager.NewPrettySink(os.Stdout, lager.DEBUG)
	logger := lager.NewLogger(component)
	sink := lager.NewReconfigurableSink(internalSink, lager.DEBUG)
	logger.RegisterSink(sink)

	return logger
}
