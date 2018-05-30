package main

import flags "github.com/jessevdk/go-flags"

type Args struct {
	Socket string `long:"socket"`
	Handle string `long:"handle"`
	Action string `long:"action"`
}

func parseArgs() (Args, error) {
	var args Args
	_, err := flags.Parse(&args)
	if err != nil {
		return Args{}, err
	}
	return args, nil
}
