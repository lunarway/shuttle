package templates

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

var (
	input = make(map[string]interface{})
	data  = `
a: Easy!
b:
  c: 2
  h: 'ewff'
`
)

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

func TestTmplGetJsonValueByKeys(t *testing.T) {
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
				path: "user.name",
				data: nil,
			},
			output: nil,
		},
		{
			name: "empty JSON string",
			input: input{
				path: "user.name",
				data: "",
			},
			output: nil,
		},
		{
			name: "invalid JSON string",
			input: input{
				path: "user.name",
				data: "{invalid json}",
			},
			output: nil,
		},
		{
			name: "valid JSON string with nested object",
			input: input{
				path: "user.name",
				data: `{"user":{"name":"John","age":30}}`,
			},
			output: "John",
		},
		{
			name: "valid JSON string with array access",
			input: input{
				path: "users.0.name",
				data: `{"users":[{"name":"Alice","age":25},{"name":"Bob","age":35}]}`,
			},
			output: "Alice",
		},
		{
			name: "valid JSON string with deep nesting",
			input: input{
				path: "config.database.host",
				data: `{"config":{"database":{"host":"localhost","port":5432},"cache":{"enabled":true}}}`,
			},
			output: "localhost",
		},
		{
			name: "valid JSON string with number value",
			input: input{
				path: "user.age",
				data: `{"user":{"name":"John","age":30}}`,
			},
			output: float64(30), // JSON numbers are parsed as float64
		},
		{
			name: "valid JSON string with boolean value",
			input: input{
				path: "config.cache.enabled",
				data: `{"config":{"database":{"host":"localhost"},"cache":{"enabled":true}}}`,
			},
			output: true,
		},
		{
			name: "non-existent path in JSON string",
			input: input{
				path: "user.email",
				data: `{"user":{"name":"John","age":30}}`,
			},
			output: nil,
		},
		{
			name: "already parsed JSON data (map)",
			input: input{
				path: "user.name",
				data: map[string]interface{}{
					"user": map[string]interface{}{
						"name": "Jane",
						"age":  28,
					},
				},
			},
			output: "Jane",
		},
		{
			name: "already parsed JSON data with nested map",
			input: input{
				path: "config.database.port",
				data: map[string]interface{}{
					"config": map[string]interface{}{
						"database": map[string]interface{}{
							"host": "localhost",
							"port": 5432,
						},
					},
				},
			},
			output: 5432,
		},
		{
			name: "array access in parsed data",
			input: input{
				path: "items.1",
				data: map[string]interface{}{
					"items": []interface{}{"first", "second", "third"},
				},
			},
			output: "second",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output := TmplGetJsonValueByKeys(tc.input.path, tc.input.data)
			assert.Equal(t, tc.output, output, "output does not match the expected")
		})
	}
}

func TestTmplJsonPath(t *testing.T) {
	type input struct {
		jsonString string
		path       string
	}
	tt := []struct {
		name   string
		input  input
		output interface{}
	}{
		{
			name: "empty JSON string",
			input: input{
				jsonString: "",
				path:       "user.name",
			},
			output: nil,
		},
		{
			name: "invalid JSON string",
			input: input{
				jsonString: "{invalid json}",
				path:       "user.name",
			},
			output: nil,
		},
		{
			name: "valid JSON string with nested object",
			input: input{
				jsonString: `{"user":{"name":"John","age":30}}`,
				path:       "user.name",
			},
			output: "John",
		},
		{
			name: "valid JSON string with array access",
			input: input{
				jsonString: `{"users":[{"name":"Alice","age":25},{"name":"Bob","age":35}]}`,
				path:       "users.0.name",
			},
			output: "Alice",
		},
		{
			name: "valid JSON string with deep nesting",
			input: input{
				jsonString: `{"config":{"database":{"host":"localhost","port":5432},"cache":{"enabled":true}}}`,
				path:       "config.database.host",
			},
			output: "localhost",
		},
		{
			name: "valid JSON string with number value",
			input: input{
				jsonString: `{"user":{"name":"John","age":30}}`,
				path:       "user.age",
			},
			output: float64(30), // JSON numbers are parsed as float64
		},
		{
			name: "valid JSON string with boolean value",
			input: input{
				jsonString: `{"config":{"database":{"host":"localhost"},"cache":{"enabled":true}}}`,
				path:       "config.cache.enabled",
			},
			output: true,
		},
		{
			name: "valid JSON string with null value",
			input: input{
				jsonString: `{"user":{"name":"John","email":null}}`,
				path:       "user.email",
			},
			output: nil,
		},
		{
			name: "non-existent path in JSON string",
			input: input{
				jsonString: `{"user":{"name":"John","age":30}}`,
				path:       "user.email",
			},
			output: nil,
		},
		{
			name: "array access with index out of bounds",
			input: input{
				jsonString: `{"items":["first","second"]}`,
				path:       "items.5",
			},
			output: nil,
		},
		{
			name: "complex nested structure",
			input: input{
				jsonString: `{"data":{"users":[{"id":1,"profile":{"name":"Alice","settings":{"theme":"dark"}}},{"id":2,"profile":{"name":"Bob","settings":{"theme":"light"}}}]}}`,
				path:       "data.users.1.profile.settings.theme",
			},
			output: "light",
		},
		{
			name: "empty object",
			input: input{
				jsonString: `{}`,
				path:       "any.path",
			},
			output: nil,
		},
		{
			name: "array at root level",
			input: input{
				jsonString: `[{"name":"first"},{"name":"second"}]`,
				path:       "0.name",
			},
			output: "first",
		},
		{
			name: "string value at root level",
			input: input{
				jsonString: `"simple string"`,
				path:       "",
			},
			output: "simple string",
		},
		{
			name: "number value at root level",
			input: input{
				jsonString: `42`,
				path:       "",
			},
			output: float64(42),
		},
		{
			name: "boolean value at root level",
			input: input{
				jsonString: `true`,
				path:       "",
			},
			output: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output := TmplJsonPath(tc.input.jsonString, tc.input.path)
			assert.Equal(t, tc.output, output, "output does not match the expected")
		})
	}
}
