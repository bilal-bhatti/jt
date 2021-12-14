package jt

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/ohler55/ojg/jp"
	"github.com/pkg/errors"

	"github.com/itchyny/gojq"
)

type templatizer interface {
	Templatize(input interface{}) error
}

type applier interface {
	Apply(input, template interface{}) error
}

type Template struct {
	Debug bool
	DSL   string
}

func (t Template) Templatize(input interface{}) error {
	return t.templatize("", reflect.ValueOf(input), reflect.ValueOf(nil), reflect.ValueOf(input))
}

func (t Template) templatize(prefix string, container, k, v reflect.Value) error {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			err := t.templatize(prefix+"."+"["+strconv.Itoa(i)+"]", v, reflect.ValueOf(int(i)), v.Index(i))
			if err != nil {
				return err
			}
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			err := t.templatize(prefix+"."+k.String(), v, k, v.MapIndex(k))
			if err != nil {
				return err
			}
		}
	default:
		exp := fmt.Sprintf("$%s{%s}", t.DSL, prefix)
		if k.Kind() == reflect.Int {
			container.Index(int(k.Int())).Set(reflect.ValueOf(exp))
		} else {
			container.SetMapIndex(k, reflect.ValueOf(exp))
		}
	}

	return nil
}

func (t Template) Apply(input, template interface{}) error {
	return t.apply(input, reflect.ValueOf(template), reflect.ValueOf(nil), reflect.ValueOf(template))
}

var ep = regexp.MustCompile("^\\$([p,q]{0,1}){(.+)}$")

func (t Template) apply(source interface{}, container, k, template reflect.Value) error {
	for template.Kind() == reflect.Ptr || template.Kind() == reflect.Interface {
		template = template.Elem()
	}
	switch template.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < template.Len(); i++ {
			err := t.apply(source, template, reflect.ValueOf(int(i)), template.Index(i))
			if err != nil {
				return err
			}
		}
	case reflect.Map:
		for _, k := range template.MapKeys() {
			err := t.apply(source, template, k, template.MapIndex(k))
			if err != nil {
				return err
			}
		}
	default:
		if template.Kind() != reflect.String {
			return nil // not a lookup expression, skip
		}

		exp := strings.TrimSpace(template.String())

		// matches:
		// 0 : whole expression
		// 1 : dsl option
		// 2 : lookup expression
		matches := ep.FindStringSubmatch(exp)
		if len(matches) == 0 {
			return nil // not a lookup expression, skip
		}

		if t.Debug {
			log.Printf("processing: %v", matches)
		}

		var lookup func(string, interface{}) ([]interface{}, error)

		if matches[1] == "" || matches[1] == "q" {
			lookup = query
		} else {
			lookup = path
		}

		results, err := lookup(matches[2], source)
		if err != nil {
			return err
		}

		if t.Debug {
			if results == nil {
				log.Printf("nil result for expression: %s", exp)
			}
		}

		// TODO: handle valid non-singular results
		// currently this will leave the expression
		// in the output
		// - 0 results
		// - more than 1 results
		if len(results) != 1 {
			return errors.Errorf("unexpected results, %v", results)
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

	return nil
}

func query(exp string, source interface{}) ([]interface{}, error) {
	jq_exp, err := gojq.Parse(exp)
	if err != nil {
		return nil, errors.Errorf("query expression parse error, %v", err)
	}

	results := make([]interface{}, 0)

	iter := jq_exp.Run(source)
	for {
		lookup, ok := iter.Next()
		if !ok {
			break
		}

		if err, ok := lookup.(error); ok {
			return nil, errors.Errorf("jq lookup error, %v", err)
		}

		results = append(results, lookup)
	}
	return results, nil
}

func path(exp string, source interface{}) ([]interface{}, error) {
	jp_exp, err := jp.ParseString(exp)
	if err != nil {
		return nil, errors.Errorf("path expression parse error, %v", err)
	}

	results := jp_exp.Get(source)

	return results, nil
}
