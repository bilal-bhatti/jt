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

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&templateCmd{}, "")
	subcommands.Register(&applyCmd{}, "")

	flag.Parse()

	allCmds := map[string]bool{
		"commands": true, // builtin
		"help":     true, // builtin
		"flags":    true, // builtin
		"template": true,
		"apply":    true,
	}

	log.Println("version: ", Version)

	// Default to running the "template" command.
	if args := flag.Args(); len(args) == 0 || !allCmds[args[0]] {
		defCmd := &templateCmd{}
		os.Exit(int(defCmd.Execute(context.Background(), flag.CommandLine)))
	}
	os.Exit(int(subcommands.Execute(context.Background())))
}
