package jt

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/ohler55/ojg/jp"
	"github.com/pkg/errors"

	"github.com/itchyny/gojq"
)

type Tool struct {
	Verbose bool
}

func (t Tool) Templatize(input interface{}) error {
	return t.templatize("", reflect.ValueOf(input), reflect.ValueOf(nil), reflect.ValueOf(input))
}

func (t Tool) templatize(prefix string, container, k, v reflect.Value) error {
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
		exp := fmt.Sprintf("$%s{%s}", "q", prefix)
		if k.Kind() == reflect.Int {
			container.Index(int(k.Int())).Set(reflect.ValueOf(exp))
		} else {
			container.SetMapIndex(k, reflect.ValueOf(exp))
		}
	}

	return nil
}

func (t Tool) Apply(input, template interface{}) error {
	return t.apply(input, reflect.ValueOf(template), reflect.ValueOf(nil), reflect.ValueOf(template))
}

var ep = regexp.MustCompile(`\$([p,q,e]{0,1}){(.+?)}`)

func (t Tool) apply(source interface{}, container, k, template reflect.Value) error {
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

		exp := template.String()

		// matches:
		// 0 : whole expression
		// 1 : dsl option
		// 2 : lookup expression

		all_matches := ep.FindAllStringSubmatch(exp, -1)

		var result interface{} = nil

		switch len(all_matches) {
		case 0:
			// nothing to interpolate
			return nil
		case 1:
			if len(all_matches[0][0]) == len(exp) {
				// full string match, json object replacement
				result = replace(all_matches[0], source)
			} else {
				// sub string match, string interpolation
				result = interpolate(exp, all_matches, source)
			}
		default:
			// multiple matches, string interpolation
			result = interpolate(exp, all_matches, source)
		}

		if k.Kind() == reflect.Int {
			t := container.Index(int(k.Int()))
			if result != nil {
				t.Set(reflect.ValueOf(result))
			}
		} else {
			container.SetMapIndex(k, reflect.ValueOf(result))
		}
	}

	return nil
}

func etype(matches []string) evaluator {
	if matches[1] == "" || matches[1] == "q" {
		return query
	} else if matches[1] == "p" {
		return path
	} else if matches[1] == "e" {
		return env
	}
	return query
}

func replace(match []string, source interface{}) interface{} {
	lookup := etype(match)

	results, err := lookup(match[2], source)
	if err != nil {
		return nil
	}

	// TODO: handle valid non-singular results
	// currently this will leave the expression
	// in the output
	// - 0 results
	// - more than 1 results
	if len(results) != 1 {
		return errors.Errorf("unexpected results, %v", results)
	}

	return results[0]
}

func interpolate(template string, matches [][]string, source interface{}) interface{} {
	out := string(template)

	for _, match := range matches {
		lookup := etype(match)

		results, err := lookup(match[2], source)
		if err != nil {
			return nil
		}

		// TODO: handle valid non-singular results
		// currently this will leave the expression
		// in the output
		// - 0 results
		// - more than 1 results
		if len(results) != 1 {
			return errors.Errorf("unexpected results, %v", results)
		}

		result := results[0]

		// TODO: result must be a string, fail otherwise
		if rt, ok := result.(string); ok {
			out = strings.ReplaceAll(out, match[0], rt)
		}
	}

	return out
}

type evaluator func(string, interface{}) ([]interface{}, error)

var query evaluator = func(exp string, source interface{}) ([]interface{}, error) {
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

var path evaluator = func(exp string, source interface{}) ([]interface{}, error) {
	jp_exp, err := jp.ParseString(exp)
	if err != nil {
		return nil, errors.Errorf("path expression parse error, %v", err)
	}

	results := jp_exp.Get(source)

	return results, nil
}

var env evaluator = func(exp string, source interface{}) ([]interface{}, error) {
	return []interface{}{os.Getenv(exp)}, nil
}
