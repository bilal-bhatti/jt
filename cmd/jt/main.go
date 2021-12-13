package main

import (
	"context"
	"flag"
	"log"
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
var dsl string

func main() {
	log.Println("version: ", Version)

	flag.BoolVar(&debug, "d", false, "run with debug logging enabled")
	//flag.StringVar(&dsl, "l", "", "expression dsl [q (jq), p (json-path)], defaults to 'q'")

	flag.Parse()

	tCmd := &templateCmd{debug: debug, dsl: dsl}
	aCmd := &applyCmd{debug: debug, dsl: dsl}

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

	// Default to running the "apply" command.
	if args := flag.Args(); len(args) == 0 || !allCmds[args[0]] {
		os.Exit(int(tCmd.Execute(context.Background(), flag.CommandLine)))
	}
	os.Exit(int(subcommands.Execute(context.Background())))
}
