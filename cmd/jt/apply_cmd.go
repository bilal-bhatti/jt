package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"

	jt "github.com/bilal-bhatti/jt/pkg"
	"github.com/google/subcommands"
)

type applyCmd struct {
	debug                bool
	input, template, out string
}

func (*applyCmd) Name() string { return "apply" }

func (*applyCmd) Synopsis() string {
	return "apply template to input"
}

func (*applyCmd) Usage() string {
	log.Println("version: ", Version)
	return `
apply jq tranformation template to input

	examples: 
		cat input.json | jt apply -t template.json 
		cat input.json | jt apply -t template.json -o template.json
		jt apply -i input.json -t template.json -o template.json

`
}

func (a *applyCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&a.debug, "d", false, "run with trace logging enabled")
	f.StringVar(&a.input, "i", "", "read from STDIN or file, -i <file.json>")
	f.StringVar(&a.template, "t", "template.json", "required template file, -t <file.json>")
	f.StringVar(&a.out, "o", "", "write to STDOUT or file, -o <file.json>")
}

func (a *applyCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if _, err := os.Stat(a.template); errors.Is(err, os.ErrNotExist) {
		log.Println(a.Usage())
		os.Exit(int(subcommands.ExitFailure))
	}

	var input interface{}
	var template interface{}

	if a.debug {
		log.Println("executing with args")
		log.Println(" -d ", a.debug)
		log.Println(" -i ", a.input)
		log.Println(" -t ", a.template)
		log.Println(" -o ", a.out)
	}

	if a.input == "" {
		err := json.NewDecoder(os.Stdin).Decode(&input)
		orPanic(err)
	} else {
		yf, err := ioutil.ReadFile(a.input)
		orPanic(err)

		err = json.Unmarshal(yf, &input)
		orPanic(err)
	}

	yf, err := ioutil.ReadFile(a.template)
	orPanic(err)

	err = json.Unmarshal(yf, &template)
	orPanic(err)

	jt.Apply(input, template)

	bites, err := json.MarshalIndent(template, "", "\t")
	orPanic(err)
	bites = append(bites, byte('\n'))

	if a.out != "" {
		err = ioutil.WriteFile(a.out, bites, 0644)
		orPanic(err)
	} else {
		_, err := os.Stdout.Write(bites)
		orPanic(err)
	}

	return subcommands.ExitSuccess
}
