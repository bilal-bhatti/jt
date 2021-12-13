package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"

	jt "github.com/bilal-bhatti/jt/pkg"
	"github.com/google/subcommands"
)

type templateCmd struct {
	debug      bool
	language   string
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
	f.StringVar(&t.input, "i", "", "input file")
	f.StringVar(&t.out, "o", "", "out file")
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
		// TODO: use reader like above os.stdin instead of json.Unmarshal
		yf, err := ioutil.ReadFile(t.input)
		if err != nil {
			log.Fatal(errors.Errorf("failed to read input json file, %v", err))
		}

		err = json.Unmarshal(yf, &input)

		orPanic(err)
	}

	jt.Templatize(input)

	bites, err := json.MarshalIndent(input, "", "\t")
	if err != nil {
		log.Fatal(errors.Errorf("json marshal error, %v", err))
	}

	bites = append(bites, byte('\n'))

	if t.out != "" {
		err = ioutil.WriteFile(t.out, bites, 0644)
		if err != nil {
			log.Fatal(errors.Errorf("failed to write file, %v", err))
		}
	} else {
		_, err := os.Stdout.Write(bites)
		if err != nil {
			log.Fatal(errors.Errorf("failed to write to STDOUT, %v", err))
		}
	}

	return subcommands.ExitSuccess
}
