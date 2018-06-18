package main

import flags "github.com/jessevdk/go-flags"

type Args struct {
	SocketPath    string   `long:"socket" required:"true"`
	NetpluginPath string   `long:"netpluginPath" required:"true"`
	NetpluginArgs []string `long:"netpluginArgs"`
}

func parseArgs() (Args, error) {
	var args Args
	_, err := flags.Parse(&args)
	if err != nil {
		return Args{}, err
	}
	return args, nil
}
