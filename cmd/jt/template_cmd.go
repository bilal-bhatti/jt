package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/pkg/errors"
	"log"
	"os"

	jt "github.com/bilal-bhatti/jt/pkg"
	"github.com/google/subcommands"
)

type templateCmd struct {
	debug      bool
	dsl        string
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

	examples: 
		cat input.json | jt template 
		jt template -i input.json -o template
		jt template -i input.json -o template.json

`
}

func (t *templateCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&t.input, "i", "", "read from STDIN or file, -i <file.json>")
	f.StringVar(&t.out, "o", "", "write to STDOUT or file, -o <file.json>")
}

func (t *templateCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	var input interface{}

	if t.debug {
		log.Println("executing with args")
		log.Println(" -i ", t.input)
		log.Println(" -o ", t.out)
	}

	if t.input == "" {
		err := json.NewDecoder(os.Stdin).Decode(&input)
		if err != nil {
			log.Fatal(errors.Errorf("failed to decode json from stdin, %v", err))
		}
	} else {
		err := jt.ReadFile(t.input, &input)
		if err != nil {
			log.Fatalln((err))
		}
	}

	tmpl := jt.Template{
		Debug: t.debug,
		DSL:   t.dsl,
	}

	err := tmpl.Templatize(input)
	if err != nil {
		log.Fatalln(err)
	}

	err = jt.WriteFile(t.out, input)
	if err != nil {
		log.Fatalln(err)
	}

	return subcommands.ExitSuccess
}
