package templates

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

var input = make(map[string]interface{})
var data = `
a: Easy!
b:
  c: 2
  h: 'ewff'
`

func init() {
	err := yaml.Unmarshal([]byte(data), input)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func TestTmplStrConst(t *testing.T) {
	input := "log.console.as.json"
	output := TmplStrConst(input)
	expected := "LOG_CONSOLE_AS_JSON"

	if output != expected {
		t.Errorf("strCon conversion was incorrect, got: %s, want: %s", output, expected)
	}
}

func TestTmplStr(t *testing.T) {
	output := TmplString("b.h", input)

	if output != "ewff" {
		t.Errorf("string conversion was incorrect, got: %s, want: %v", output, "ewff")
	}
}

func TestTmplStr_with_int(t *testing.T) {
	output := TmplString("b.c", input)

	if output != "2" {
		t.Errorf("string conversion was incorrect, got: %s, want: %v", output, 2)
	}
}

func TestTmplInt(t *testing.T) {
	output := TmplInt("b.c", input)

	if output != 2 {
		t.Errorf("string conversion was incorrect, got: %d, want: %d", output, 2)
	}
}

func TestTmpObjArray(t *testing.T) {
	type input struct {
		path string
		data interface{}
	}
	tt := []struct {
		name   string
		input  input
		output []KeyValuePair
	}{
		{
			name: "nil input",
			input: input{
				path: "b",
				data: nil,
			},
			output: []KeyValuePair{},
		},
		{
			name: "empty input",
			input: input{
				path: "b",
				data: fromYaml(""),
			},
			output: []KeyValuePair{},
		},
		{
			name: "nested values",
			input: input{
				path: "b",
				data: fromYaml(`
b:
  c: 2
  h: 'ewff'
`),
			},
			output: []KeyValuePair{
				{Key: "c", Value: 2},
				{Key: "h", Value: "ewff"},
			},
		},
		{
			name: "non object value in path",
			input: input{
				path: "a",
				data: fromYaml(`a: b`),
			},
			output: []KeyValuePair{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output := TmplObjectArray(tc.input.path, tc.input.data)

			assert.ElementsMatch(t, tc.output, output, "output does not match the expected")
		})
	}
}

func fromYaml(data string) interface{} {
	m := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(data), m)
	if err != nil {
		log.Fatalf("faild to read yaml: %v", err)
	}
	return m
}
