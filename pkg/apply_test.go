package jt

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApply(t *testing.T) {
	var c interface{}

	yf, err := ioutil.ReadFile("apply_test.json")
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(yf, &c)
	if err != nil {
		t.Error(err)
	}

	tmpl := Template{}
	for _, data := range c.(map[string]interface{})["data"].([]interface{}) {
		input := data.(map[string]interface{})["input"]
		template := data.(map[string]interface{})["template"]
		expected := data.(map[string]interface{})["expected"]

		tmpl.Apply(input, template)

		assert.Equal(t, expected, template, "equal")
	}

}
