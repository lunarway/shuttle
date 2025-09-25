package templates

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"text/template"

	sprig "github.com/Masterminds/sprig/v3"
	yaml "gopkg.in/yaml.v2"
)

type KeyValuePair struct {
	Key   string
	Value interface{}
}

// GetFuncMap returns a Map of template functions
func GetFuncMap() template.FuncMap {
	f := sprig.TxtFuncMap()

	extra := template.FuncMap{
		"get":                TmplGet,
		"string":             TmplString,
		"array":               TmplArray,
		"objectArray":         TmplObjectArray,
		"strConst":            TmplStrConst,
		"int":                 TmplInt,
		"is":                  TmplIs,
		"isnt":                TmplIsnt,
		"toYaml":              TmplToYaml,
		"fromYaml":            TmplFromYaml,
		"getFiles":            TmplGetFiles,
		"getFileContent":      TmplGetFileContent,
		"fileExists":          TmplFileExists,
		"trim":                strings.TrimSpace,
		"upperFirst":          TmplUpperFirst,
		"rightPad":            TmplRightPad,
		"getJsonValueByKeys":  TmplGetJsonValueByKeys,
		"jsonPath":            TmplJsonPath,
	}

	for k, v := range extra {
		f[k] = v
	}

	return f
}

func TmplGet(path string, input interface{}) interface{} {
	// If path is empty, return the input itself
	if path == "" {
		return input
	}

	if !strings.Contains(path, ".") {
		return getInner(path, input)
	}

	properties := strings.SplitN(path, ".", 2)
	property := properties[0]
	step := getInner(property, input)

	if step != nil {
		return TmplGet(properties[1], step)
	}
	return nil
}

// template function to convert from log.debug to LOG_DEBUG
func TmplStrConst(value string) string {
	value = strings.Replace(value, ".", "_", -1)
	value = strings.ToUpper(value)
	return value
}

func TmplString(path string, input interface{}) string {
	value := TmplGet(path, input)
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}

func TmplInt(path string, input interface{}) int {
	value := TmplGet(path, input)
	if value == nil {
		return 0
	}
	return value.(int)
}

// TmplArray parses a YAML path in the input parameter as an array.
// If variable found is an array, array values are returned
// If variable found is a map, the maps values are returned
// If none of above then nil is returned
func TmplArray(path string, input interface{}) []interface{} {
	value := TmplGet(path, input)
	switch actualValue := value.(type) {
	case map[interface{}]interface{}:
		var values []interface{}
		for _, v := range TmplObjectArray(path, input) {
			values = append(values, v.Value)
		}
		return values
	case map[string]interface{}:
		var values []interface{}
		for _, v := range TmplObjectArray(path, input) {
			values = append(values, v.Value)
		}
		return values
	case []interface{}:
		return actualValue
	}
	return nil
}

// TmplObjectArray parses a YAML path in the input parameter as an object and returns Key & Value pairs.
// If variable found is a map, a []KeyValuePair array is returned
// If variable found is a map, the maps values are returned
// If none of above then nil is returned
func TmplObjectArray(path string, input interface{}) []KeyValuePair {
	if input == nil {
		return nil
	}
	value := TmplGet(path, input)
	values := []KeyValuePair{}
	switch typedValue := value.(type) {
	case map[interface{}]interface{}:
		for k, v := range typedValue {
			values = append(values, KeyValuePair{Key: k.(string), Value: v})
		}
	case map[string]interface{}:
		for k, v := range typedValue {
			values = append(values, KeyValuePair{Key: k, Value: v})
		}
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i].Key < values[j].Key
	})
	return values
}

func TmplIs(a interface{}, b interface{}) bool {
	return a == b
}

func TmplIsnt(a interface{}, b interface{}) bool {
	return a != b
}

// TmplRightPad adds padding to the right of a string.
func TmplRightPad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(template, s)
}

// TmplUpperFirst upper cases the first letter in the string
func TmplUpperFirst(s string) string {
	if len(s) <= 1 {
		return strings.ToUpper(s)
	}
	return fmt.Sprintf("%s%s", strings.ToUpper(s[0:1]), s[1:])
}

// TmplToYaml takes an interface, marshals it to yaml, and returns a string.
//
// This is designed to be called from a template.
//
// Borrowed from github.com/helm/helm/pkg/chartutil
func TmplToYaml(v interface{}) (string, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// TmplFromYaml converts a YAML document into a map[string]interface{}.
//
// This is not a general-purpose YAML parser, and will not parse all valid
// YAML documents.
//
// Borrowed from github.com/helm/helm/pkg/chartutil
func TmplFromYaml(str string) (map[string]interface{}, error) {
	m := map[string]interface{}{}

	err := yaml.Unmarshal([]byte(str), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func TmplGetFileContent(filePath string) (string, error) {
	byteContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(byteContent), nil
}

func TmplFileExists(filePath string) bool {
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		return true
	}
	return false
}

func TmplGetFiles(directoryPath string) []os.FileInfo {
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return []os.FileInfo{}
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	return files
}

// TmplGetJsonValueByKeys extracts a value from JSON data using dot-notation keys.
// It works similarly to TmplGet but is specifically designed for JSON data structures.
func TmplGetJsonValueByKeys(path string, input interface{}) interface{} {
	if input == nil {
		return nil
	}

	// If input is a string, try to parse it as JSON
	var jsonData interface{}
	switch v := input.(type) {
	case string:
		err := json.Unmarshal([]byte(v), &jsonData)
		if err != nil {
			return nil
		}
	default:
		jsonData = input
	}

	return TmplGet(path, jsonData)
}

// TmplJsonPath parses a JSON string and extracts a value using a path expression.
// The path uses dot-notation to navigate through the JSON structure.
func TmplJsonPath(jsonString string, path string) interface{} {
	if jsonString == "" {
		return nil
	}

	var jsonData interface{}
	err := json.Unmarshal([]byte(jsonString), &jsonData)
	if err != nil {
		return nil
	}

	return TmplGet(path, jsonData)
}

func getInner(property string, input interface{}) interface{} {
	switch t := input.(type) {
	default:
		fmt.Printf("unexpected type %T\n", t) // %T prints whatever type t has
		panic(fmt.Sprintf("unexpected type %T\n", t))
	case nil:
		return nil
	case map[interface{}]interface{}:
		values := input.(map[interface{}]interface{})
		return values[property]
	case map[string]interface{}:
		values := input.(map[string]interface{})
		return values[property]
	case map[string]string:
		values := input.(map[string]string)
		return values[property]
	case []interface{}:
		// Handle array access by index
		if index, err := parseIndex(property); err == nil && index >= 0 && index < len(t) {
			return t[index]
		}
		return nil
	case string:
		// Handle string access by index
		if index, err := parseIndex(property); err == nil && index >= 0 && index < len(t) {
			return string(t[index])
		}
		return nil
	case bool:
		return nil
	case int:
		return nil
	case float64:
		return nil
	}
}

// parseIndex attempts to parse a string as an integer index
func parseIndex(s string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}
	
	// Simple integer parsing
	var index int
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, fmt.Errorf("not a number")
		}
		index = index*10 + int(r-'0')
	}
	return index, nil
}
