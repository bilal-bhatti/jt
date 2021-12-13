package jt

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"reflect"
	"regexp"
	"strconv"

	"github.com/itchyny/gojq"
)

func Templatize(input interface{}) {
	templatize("", reflect.ValueOf(input), reflect.ValueOf(nil), reflect.ValueOf(input))
}

func templatize(prefix string, container, k, v reflect.Value) {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			templatize(prefix+"."+"["+strconv.Itoa(i)+"]", v, reflect.ValueOf(int(i)), v.Index(i))
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			templatize(prefix+"."+k.String(), v, k, v.MapIndex(k))
		}
	default:
		exp := fmt.Sprintf("$%s{%s}", "q", prefix)
		if k.Kind() == reflect.Int {
			container.Index(int(k.Int())).Set(reflect.ValueOf(exp))
		} else {
			container.SetMapIndex(k, reflect.ValueOf(exp))
		}
	}
}

func Apply(input, template interface{}) {
	apply(input, reflect.ValueOf(template), reflect.ValueOf(nil), reflect.ValueOf(template))
}

var ep = regexp.MustCompile("\\$([p,q]{0,1}){(.+)}")

func apply(source interface{}, container, k, template reflect.Value) {
	for template.Kind() == reflect.Ptr || template.Kind() == reflect.Interface {
		template = template.Elem()
	}
	switch template.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < template.Len(); i++ {
			apply(source, template, reflect.ValueOf(int(i)), template.Index(i))
		}
	case reflect.Map:
		for _, k := range template.MapKeys() {
			apply(source, template, k, template.MapIndex(k))
		}
	default:
		matches := ep.FindStringSubmatch(template.String())
		if len(matches) == 0 {
			return // not an expression, skip
		}

		log.Printf("processing: %v", matches)

		results := query(matches[2], source)

		if len(results) != 1 {
			log.Printf("unexpected results, %v", results)
			return
		}

		if results == nil {
			log.Println(errors.Errorf("%s yielded nil result", template.String()))
		}

		if k.Kind() == reflect.Int {
			t := container.Index(int(k.Int()))
			if results[0] != nil {
				t.Set(reflect.ValueOf(results[0]))
			}
		} else {
			container.SetMapIndex(k, reflect.ValueOf(results[0]))
		}
	}
}

func query(exp string, source interface{}) []interface{} {
	query, err := gojq.Parse(exp)
	if err != nil {
		log.Fatalln(errors.Errorf("expression parse error, %v", err))
	}

	results := make([]interface{}, 0)

	iter := query.Run(source)
	for {
		lookup, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := lookup.(error); ok {
			log.Fatalln(err) // TODO: handle error better
		}
		results = append(results, lookup)
	}
	return results
}
