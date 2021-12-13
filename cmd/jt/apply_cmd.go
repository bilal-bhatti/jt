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

type applyCmd struct {
	debug                bool
	language             string
	input, template, out string
}

func (*applyCmd) Name() string { return "apply" }

func (*applyCmd) Synopsis() string {
	return "apply template to input"
}

func (*applyCmd) Usage() string {
	log.Println("version: ", Version)
	return `
apply template to input:
	example: jt apply -i input.json -t template.json -o out.json

`
}

func (a *applyCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&a.input, "i", "", "input file")
	f.StringVar(&a.template, "t", "", "template file")
	f.StringVar(&a.out, "o", "", "out file")
}

func (a *applyCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	var input interface{}
	var template interface{}

	if a.debug {
		log.Println("executing with args")
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
