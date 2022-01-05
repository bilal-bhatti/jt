/*
Copyright © 2021 Bilal Bhatti
*/

package main

import (
	"encoding/json"
	"log"
	"os"

	jt "github.com/bilal-bhatti/jt/pkg"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jt [flags] -t <template.json> -i <input.json> -o <out.json>",
	Short: "apply jq tranformation template to input",
	Long: `apply jq tranformation template to input

examples: 
	cat input.json | jt -t template.json 
	cat input.json | jt -t template.json -o template.json
	jt -i input.json -t template.json -o template.json
`,
	// Args: cobra.MinimumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		if jt.Verbose {
			log.Printf("jt v%s\n", jt.Version)
		}
	},

	Run: run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.curly.yaml)")
	rootCmd.PersistentFlags().BoolVar(&jt.Verbose, "verbose", false, "run with verbose")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("input", "i", "", "read from STDIN or file, -i <file.json>")
	rootCmd.Flags().StringP("template", "t", "template.json", "required template file, -t <file.json>")
	rootCmd.Flags().StringP("out", "o", "", "write to STDOUT or file, -o <file.json>")
}

func run(cmd *cobra.Command, args []string) {
	var i, _ = cmd.Flags().GetString("input")
	var t, _ = cmd.Flags().GetString("template")
	var o, _ = cmd.Flags().GetString("out")

	if _, err := os.Stat(t); errors.Is(err, os.ErrNotExist) {
		log.Fatalln(err)
	}

	var input interface{}
	var template interface{}

	if i == "" {
		err := json.NewDecoder(os.Stdin).Decode(&input)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		err := jt.ReadFile(i, &input)
		if err != nil {
			log.Fatalln(err)
		}
	}

	err := jt.ReadFile(t, &template)
	if err != nil {
		log.Fatalln(err)
	}

	tmpl := jt.Tool{
		Verbose: jt.Verbose,
	}

	err = tmpl.Apply(input, template)
	if err != nil {
		log.Fatalln(err)
	}

	err = jt.WriteFile(o, template)
	if err != nil {
		log.Fatalln(err)
	}
}
