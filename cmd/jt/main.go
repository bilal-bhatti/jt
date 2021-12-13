package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
)

func orPanic(err error) {
	if err != nil {
		panic(err)
	}
}

var Version = "DEV"
var debug bool
var language string

func main() {
	flag.BoolVar(&debug, "d", false, "run with trace logging enalbed")
	flag.StringVar(&language, "l", "q", "expression DSL [q (jq), p (json-path)], defaults to 'q'")

	flag.Parse()

	tCmd := &templateCmd{debug: debug, language: language}
	aCmd := &applyCmd{debug: debug, language: language}

	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(tCmd, "")
	subcommands.Register(aCmd, "")

	allCmds := map[string]bool{
		"commands": true, // builtin
		"help":     true, // builtin
		"flags":    true, // builtin
		"template": true,
		"apply":    true,
	}

	// Default to running the "template" command.
	if args := flag.Args(); len(args) == 0 || !allCmds[args[0]] {
		os.Exit(int(tCmd.Execute(context.Background(), flag.CommandLine)))
	}
	os.Exit(int(subcommands.Execute(context.Background())))
}
