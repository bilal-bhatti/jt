package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"

	jt "github.com/bilal-bhatti/jt/pkg"
	"github.com/google/subcommands"
)

type applyCmd struct {
	debug                bool
	dsl                  string
	input, template, out string
}

func (*applyCmd) Name() string { return "apply" }

func (*applyCmd) Synopsis() string {
	return "apply template to input"
}

func (*applyCmd) Usage() string {
	return `
apply jq tranformation template to input

	examples: 
		cat input.json | jt apply -t template.json 
		cat input.json | jt apply -t template.json -o template.json
		jt apply -i input.json -t template.json -o template.json

`
}

func (a *applyCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&a.input, "i", "", "read from STDIN or file, -i <file.json>")
	f.StringVar(&a.template, "t", "template.json", "required template file, -t <file.json>")
	f.StringVar(&a.out, "o", "", "write to STDOUT or file, -o <file.json>")
}

func (a *applyCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if _, err := os.Stat(a.template); errors.Is(err, os.ErrNotExist) {
		log.Println(a.Usage())
		log.Fatalln(err)
	}

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
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		err := jt.ReadFile(a.input, &input)
		if err != nil {
			log.Fatalln(err)
		}
	}

	err := jt.ReadFile(a.template, &template)
	if err != nil {
		log.Fatalln(err)
	}

	tmpl := jt.Template{
		Debug: a.debug,
		DSL:   a.dsl,
	}

	err = tmpl.Apply(input, template)
	if err != nil {
		log.Fatalln(err)
	}

	err = jt.WriteFile(a.out, template)
	if err != nil {
		log.Fatalln(err)
	}

	return subcommands.ExitSuccess
}
