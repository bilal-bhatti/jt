package jt

import (
	"os"

	"github.com/spyzhov/ajson"
)

var json_str = `
{
	"object": {
		"foo": "bar"
	},
	"another": "anotherstring",
	"alist": [
		"one",
		"two",
		111
	],
	"anull": null
}
`

func json_test() {
	template, _ := ajson.Unmarshal([]byte(json_str))
	a_templatize(template)

	// bites, _ := ajson.Marshal(root)
	// os.Stdout.Write(bites)

	input, _ := ajson.Unmarshal([]byte(json_str))

	a_apply(input, template)
	bites, _ := ajson.Marshal(template)
	os.Stdout.Write(bites)
}

func a_templatize(input *ajson.Node) {
	switch input.Type() {
	case ajson.Object:
		for _, k := range input.Keys() {
			a_templatize(input.MustKey(k))
		}
	case ajson.Array:
		for _, v := range input.MustArray() {
			a_templatize(v)
		}
	default:
		input.Set(input.Path())
	}
}

func a_apply(input, template *ajson.Node) {
	switch template.Type() {
	case ajson.Object:
		for _, k := range template.Keys() {
			a_apply(input, template.MustKey(k))
		}
	case ajson.Array:
		for _, v := range template.MustArray() {
			a_apply(input, v)
		}
	default:
		res, _ := input.JSONPath(template.MustString())
		if len(res) == 1 {
			template.Set(res[0])
		} else {
			template.Set(res)
		}
	}
}
