package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

type errResponse struct {
	Error string
}

func main() {
	args, err := parseArgs()
	if err != nil {
		failWith(err)
	}

	argsFile, err := os.Create(args.Positional.ArgsFile)
	if err != nil {
		failWith(err)
	}
	defer argsFile.Close()

	_, err = fmt.Fprintf(argsFile, "--action %s\n", args.Action)
	if err != nil {
		failWith(err)
	}

	_, err = fmt.Fprintf(argsFile, "--handle %s\n", args.Handle)
	if err != nil {
		failWith(err)
	}

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		failWith(err)
	}
	if err = ioutil.WriteFile(args.Positional.StdinFile, input, os.ModePerm); err != nil {
		failWith(err)
	}

	extrafileFile, err := os.Create(args.Positional.ExtrafileFile)
	if err != nil {
		failWith(err)
	}
	defer extrafileFile.Close()

	lsCmd := exec.Command("ls", "-l", "/proc/self/fd/3")
	lsCmd.Stdout = extrafileFile
	err = lsCmd.Run()
	if err != nil {
		failWith(err)
	}

	fmt.Println(args.Positional.OutputJSON)
}

func failWith(err error) {
	json.NewEncoder(os.Stderr).Encode(&errResponse{Error: err.Error()})
}
