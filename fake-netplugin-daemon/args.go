package main

import flags "github.com/jessevdk/go-flags"

type Args struct {
	Socket      string `long:"socket"`
	Reply       string `long:"reply"`
	FDFile      string `long:"fd-file"`
	MessageFile string `long:"message-file"`
}

func parseArgs() (Args, error) {
	var args Args
	_, err := flags.Parse(&args)
	if err != nil {
		return Args{}, err
	}
	return args, nil
}
