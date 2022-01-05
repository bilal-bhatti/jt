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

// tmplCmd represents the swag command
var tmplCmd = &cobra.Command{
	Use:   "tmpl",
	Short: "generate template from input",
	Long: `generate template from input
	
examples: 
	cat input.json | jt tmpl
	jt tmpl -i input.json -o template
	jt tmpl -i input.json -o template.json
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if jt.Verbose {
			log.Printf("jt v%s\n", jt.Version)
		}
	},
	Run: templatize,
}

func init() {
	rootCmd.AddCommand(tmplCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tmplCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tmplCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	tmplCmd.Flags().StringP("input", "i", "", "read from STDIN or file, -i <file.json>")
	tmplCmd.Flags().StringP("out", "o", "", "write to STDOUT or file, -o <file.json>")
}

func templatize(cmd *cobra.Command, args []string) {
	var i, _ = cmd.Flags().GetString("input")
	var o, _ = cmd.Flags().GetString("out")

	var input interface{}

	if i == "" {
		err := json.NewDecoder(os.Stdin).Decode(&input)
		if err != nil {
			log.Fatal(errors.Errorf("failed to decode json from stdin, %v", err))
		}
	} else {
		err := jt.ReadFile(i, &input)
		if err != nil {
			log.Fatalln((err))
		}
	}

	tmpl := jt.Tool{
		Verbose: jt.Verbose,
	}

	err := tmpl.Templatize(input)
	if err != nil {
		log.Fatalln(err)
	}

	err = jt.WriteFile(o, input)
	if err != nil {
		log.Fatalln(err)
	}
}
