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

func TestTmplArray(t *testing.T) {
	type input struct {
		path string
		data interface{}
	}
	tt := []struct {
		name   string
		input  input
		output []interface{}
	}{
		{
			name: "nil input",
			input: input{
				path: "b",
				data: nil,
			},
			output: nil,
		},
		{
			name: "empty input",
			input: input{
				path: "b",
				data: fromYaml(""),
			},
			output: nil,
		},
		{
			name: "empty array",
			input: input{
				path: "b",
				data: fromYaml(`
b:
`),
			},
			output: nil,
		},
		{
			name: "single value string array",
			input: input{
				path: "b",
				data: fromYaml(`
b:
- 'a'
`),
			},
			output: []interface{}{
				"a",
			},
		},
		{
			name: "single value object array",
			input: input{
				path: "b",
				data: fromYaml(`
b:
- name: 'a'
  field: 'b'
`),
			},
			output: []interface{}{
				map[interface{}]interface{}{
					"name":  "a",
					"field": "b",
				},
			},
		},
		{
			name: "large string array testing for correct order",
			input: input{
				path: "a",
				data: fromYaml(`a:
  - 'b'
  - 'd'
  - 'e'
  - 'g'
  - 'h'
  - 'f'
  - 'i'
  - 's'
  - 'k'
  - 'c'
  - 'l'
  - 'm'
  - 'n'
  - 'o'
  - 'r'
  - 'j'
  - 'p'
  - 'q'
`),
			},
			output: []interface{}{
				"b",
				"d",
				"e",
				"g",
				"h",
				"f",
				"i",
				"s",
				"k",
				"c",
				"l",
				"m",
				"n",
				"o",
				"r",
				"j",
				"p",
				"q",
			},
		},
		{
			name: "object testing for deterministic order",
			input: input{
				path: "a",
				data: fromYaml(`a:
  b: b
  d: d
  e: e
  g: g
  h: h
  f: f
  c: c
`),
			},
			output: []interface{}{
				"b",
				"c",
				"d",
				"e",
				"f",
				"g",
				"h",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output := TmplArray(tc.input.path, tc.input.data)

			assert.Equal(t, tc.output, output, "output does not match the expected")
		})
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
			output: nil,
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
			name: "large set testing for deterministic order",
			input: input{
				path: "a",
				data: fromYaml(`a:
  'b': 2
  'd': 4
  'e': 5
  'g': 7
  'h': 8
  'f': 6
  'i': 9
  's': 19
  'k': 11
  'c': 3
  'l': 12
  'm': 13
  'n': 14
  'o': 15
  'r': 18
  'j': 10
  'p': 16
  'q': 17
`),
			},
			output: []KeyValuePair{
				{Key: "b", Value: 2},
				{Key: "c", Value: 3},
				{Key: "d", Value: 4},
				{Key: "e", Value: 5},
				{Key: "f", Value: 6},
				{Key: "g", Value: 7},
				{Key: "h", Value: 8},
				{Key: "i", Value: 9},
				{Key: "j", Value: 10},
				{Key: "k", Value: 11},
				{Key: "l", Value: 12},
				{Key: "m", Value: 13},
				{Key: "n", Value: 14},
				{Key: "o", Value: 15},
				{Key: "p", Value: 16},
				{Key: "q", Value: 17},
				{Key: "r", Value: 18},
				{Key: "s", Value: 19},
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

			assert.Equal(t, tc.output, output, "output does not match the expected")
		})
	}
}

func TestTmplGet(t *testing.T) {
	type input struct {
		path string
		data interface{}
	}
	tt := []struct {
		name   string
		input  input
		output interface{}
	}{
		{
			name: "nil input",
			input: input{
				path: "b",
				data: nil,
			},
			output: nil,
		},
		{
			name: "empty input",
			input: input{
				path: "b",
				data: fromYaml(""),
			},
			output: nil,
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
			output: map[interface{}]interface{}{
				"c": 2,
				"h": "ewff",
			},
		},
		{
			name: "an array",
			input: input{
				path: "a",
				data: fromYaml(`
a:
- 4
- c`),
			},
			output: []interface{}{
				4,
				"c",
			},
		},
		{
			name: "a nested array",
			input: input{
				path: "a.b.c",
				data: fromYaml(`
a:
  b:
    c:
    - 4
    - c`),
			},
			output: []interface{}{
				4,
				"c",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output := TmplGet(tc.input.path, tc.input.data)
			assert.Equal(t, tc.output, output, "output does not match the expected")
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

func TestTmplGetFiles(t *testing.T) {
	tt := []struct {
		name      string
		directory string

		files []string
		err   error
	}{
		{
			name:      "existing directory",
			directory: "testdata/dir",

			files: []string{
				"file.test",
			},
			err: nil,
		},
		{
			name:      "non-existing directory",
			directory: "testdata/no-dir",

			files: nil,
			err:   nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			files := TmplGetFiles(tc.directory)

			var foundFiles []string
			for _, f := range files {
				foundFiles = append(foundFiles, f.Name())
			}
			assert.Equal(t, tc.files, foundFiles)
		})
	}
}
