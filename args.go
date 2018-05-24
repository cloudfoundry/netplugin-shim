package main

import flags "github.com/jessevdk/go-flags"

type PositionalArgs struct {
	Command string
}

type Args struct {
	Socket     string         `long:"socket"`
	Positional PositionalArgs `positional-args:"true"`
}

func parseArgs() (Args, error) {
	var args Args
	_, err := flags.Parse(&args)
	if err != nil {
		return Args{}, err
	}
	return args, nil
}
