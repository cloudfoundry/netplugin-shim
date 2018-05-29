package main

import flags "github.com/jessevdk/go-flags"

type Args struct {
	Socket string `long:"socket"`
	Action string `long:"action"`
}

func parseArgs() (Args, error) {
	var args Args
	parser := flags.NewParser(&args, flags.IgnoreUnknown)
	_, err := parser.Parse()
	if err != nil {
		return Args{}, err
	}
	return args, nil
}
