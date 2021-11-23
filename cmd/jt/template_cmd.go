package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	jt "github.com/bilal-bhatti/jt/pkg"
	"github.com/google/subcommands"
)

type templateCmd struct {
	debug      bool
	input, out string
}

func (*templateCmd) Name() string { return "template" }

func (*templateCmd) Synopsis() string {
	return "generate template from input"
}

func (*templateCmd) Usage() string {
	log.Println("version: ", Version)

	return `
generate template from input
	example: jt template -i input.json -o template.json

`
}

func (t *templateCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&t.debug, "d", false, "run with trace logging enabled")
	f.StringVar(&t.input, "i", "", "input file")
	f.StringVar(&t.out, "o", "", "out file")
}

func (t *templateCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	var input interface{}

	if t.debug {
		log.Println("executing with args")
		log.Println(" -d ", t.debug)
		log.Println(" -i ", t.input)
		log.Println(" -o ", t.out)
	}

	if t.input == "" {
		err := json.NewDecoder(os.Stdin).Decode(&input)
		orPanic(err)
	} else {
		yf, err := ioutil.ReadFile(t.input)
		orPanic(err)

		err = json.Unmarshal(yf, &input)
		orPanic(err)
	}

	jt.Templatize(input)

	bites, err := json.MarshalIndent(input, "", "\t")
	orPanic(err)
	bites = append(bites, byte('\n'))

	if t.out != "" {
		err = ioutil.WriteFile(t.out, bites, 0644)
		orPanic(err)
	} else {
		_, err := os.Stdout.Write(bites)
		orPanic(err)
	}

	return subcommands.ExitSuccess
}
