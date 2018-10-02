package templates

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

type KeyValuePair struct {
	Key   string
	Value interface{}
}

// GetFuncMap returns a Map of template functions
func GetFuncMap() template.FuncMap {
	f := sprig.TxtFuncMap()

	extra := template.FuncMap{
		"get":         TmplGet,
		"string":      TmplString,
		"array":       TmplArray,
		"objectArray": TmplObjectArray,
		"strConst":    TmplStrConst,
		"int":         TmplInt,
		"is":          TmplIs,
		"isnt":        TmplIsnt,
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
		return getInner(properties[1], step)
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
	switch value.(type) {
	case map[interface{}]interface{}:
		values := []KeyValuePair{}
		for k, v := range value.(map[interface{}]interface{}) {
			values = append(values, KeyValuePair{Key: k.(string), Value: v})
		}
		return values
	case map[string]interface{}:
		values := []KeyValuePair{}
		for k, v := range value.(map[string]interface{}) {
			values = append(values, KeyValuePair{Key: k, Value: v})
		}
		return values
	}
	return nil
}

func TmplIs(a interface{}, b interface{}) bool {
	return a == b
}

func TmplIsnt(a interface{}, b interface{}) bool {
	return a != b
}

// TODO: Add description
func getInner(property string, input interface{}) interface{} {
	switch t := input.(type) {
	default:
		fmt.Printf("unexpected type %T\n", t) // %T prints whatever type t has
		panic(fmt.Sprintf("unexpected type %T\n", t))
		//case config.DynamicYaml:
		//	return
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
