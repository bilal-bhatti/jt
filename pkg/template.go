package jt

import (
	"fmt"
	"log"
	"reflect"
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
		jq := prefix
		if k.Kind() == reflect.Int {
			container.Index(int(k.Int())).Set(reflect.ValueOf(jq))
		} else {
			container.SetMapIndex(k, reflect.ValueOf(jq))
		}
	}
}

func Apply(input, template interface{}) {
	apply(input, reflect.ValueOf(template), reflect.ValueOf(nil), reflect.ValueOf(template))
}

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
		query, err := gojq.Parse(template.String())
		if err != nil {
			log.Fatalln(err)
		}

		iter := query.Run(source)
		for { // TODO: clean up this for loop
			lookup, ok := iter.Next()
			if !ok {
				break
			}
			if err, ok := lookup.(error); ok {
				log.Fatalln(err)
			}

			if lookup == nil {
				log.Println(fmt.Errorf("query %s yielded nil result", template.String()))
			}

			if k.Kind() == reflect.Int {
				t := container.Index(int(k.Int()))
				if lookup != nil {
					t.Set(reflect.ValueOf(lookup))
				}
			} else {
				container.SetMapIndex(k, reflect.ValueOf(lookup))
			}
		}
	}
}
