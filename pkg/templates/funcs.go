package templates

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
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
		"get":            TmplGet,
		"string":         TmplString,
		"array":          TmplArray,
		"objectArray":    TmplObjectArray,
		"strConst":       TmplStrConst,
		"int":            TmplInt,
		"is":             TmplIs,
		"isnt":           TmplIsnt,
		"toYaml":         TmplToYaml,
		"fromYaml":       TmplFromYaml,
		"getFiles":       TmplGetFiles,
		"getFileContent": TmplGetFileContent,
		"fileExists":     TmplFileExists,
	}

	for k, v := range extra {
		f[k] = v
	}

	return f
}

// TODO: Add description
func TmplGet(path string, input interface{}) interface{} {
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

// TODO: Add description
func TmplString(path string, input interface{}) string {
	value := TmplGet(path, input)
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}

// TODO: Add description
func TmplInt(path string, input interface{}) int {
	value := TmplGet(path, input)
	if value == nil {
		return 0
	}
	return value.(int)
}

// TODO: Add description
func TmplArray(path string, input interface{}) []interface{} {
	value := TmplGet(path, input)
	switch value.(type) {
	case map[interface{}]interface{}:
		values := []interface{}{}
		for _, v := range input.(map[interface{}]interface{}) {
			values = append(values, v)
		}
		return values
	case map[string]interface{}:
		values := []interface{}{}
		for _, v := range input.(map[string]interface{}) {
			values = append(values, v)
		}
		return values
	case []interface{}:
		return value.([]interface{})
	}
	return nil
}

// TODO: Add description
func TmplObjectArray(path string, input interface{}) []KeyValuePair {
	if input == nil {
		return nil
	}
	value := TmplGet(path, input)
	values := []KeyValuePair{}
	switch value.(type) {
	case map[interface{}]interface{}:
		for k, v := range value.(map[interface{}]interface{}) {
			values = append(values, KeyValuePair{Key: k.(string), Value: v})
		}
	case map[string]interface{}:
		for k, v := range value.(map[string]interface{}) {
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

// TmplToYaml takes an interface, marshals it to yaml, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// This is designed to be called from a template.
//
// Borrowed from github.com/helm/helm/pkg/chartutil
func TmplToYaml(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return string(data)
}

// TmplFromYaml converts a YAML document into a map[string]interface{}.
//
// This is not a general-purpose YAML parser, and will not parse all valid
// YAML documents. Additionally, because its intended use is within templates
// it tolerates errors. It will insert the returned error message string into
// m["Error"] in the returned map.
//
// Borrowed from github.com/helm/helm/pkg/chartutil
func TmplFromYaml(str string) map[string]interface{} {
	m := map[string]interface{}{}

	if err := yaml.Unmarshal([]byte(str), &m); err != nil {
		m["Error"] = err.Error()
	}
	return m
}

func TmplGetFileContent(filePath string) string {
	byteContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "" // TODO: Print error to shuttle output
	}
	return string(byteContent)
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
		return []os.FileInfo{} // TODO: Print error to shuttle output
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	return files
}

// TODO: Add description
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
		return nil
	case string:
		return nil
	case bool:
		return nil
	case int:
		return nil
	}
}
