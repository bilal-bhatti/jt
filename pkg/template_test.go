package jt

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplatization(t *testing.T) {
	var c interface{}

	yf, err := ioutil.ReadFile("template_test.json")
	t.Error(err)

	err = json.Unmarshal(yf, &c)
	t.Error(err)

	for _, data := range c.(map[string]interface{})["data"].([]interface{}) {
		input := data.(map[string]interface{})["input"]
		expected := data.(map[string]interface{})["expected"]

		Templatize(input)

		assert.Equal(t, expected, input)
	}
}
