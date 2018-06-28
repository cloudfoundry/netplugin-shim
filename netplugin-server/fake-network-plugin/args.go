package main

import flags "github.com/jessevdk/go-flags"

type PositionalArgs struct {
	ArgsFile      string
	StdinFile     string
	OutputJSON    string
	ExtrafileFile string
}

type Args struct {
	Positional PositionalArgs `positional-args:"true" required:"true"`
	Handle     string         `long:"handle" required:"true"`
	Action     string         `long:"action" required:"true"`
}

func parseArgs() (Args, error) {
	var args Args
	_, err := flags.Parse(&args)
	if err != nil {
		return Args{}, err
	}
	return args, nil
}
