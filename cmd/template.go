package cmd

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

type context struct {
	Vars interface{}
	Args map[string]string
}

var templateCmd = &cobra.Command{
	Use:   "template [template]",
	Short: "Execute a template",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var templateName = args[0]
		loadedPlan := getPlan()

		namedArgs := map[string]string{}
		for _, arg := range args[1:] {
			parts := strings.SplitN(arg, "=", 2)
			namedArgs[parts[0]] = parts[1]
		}

		templatePath := path.Join(loadedPlan.PlanPath, templateName)

		tmpl, err := template.New(templateName).Funcs(template.FuncMap{
			"get":         templateGet,
			"string":      templateString,
			"array":       templateArray,
			"objectArray": templateObjectArray,
		}).ParseFiles(templatePath)

		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(os.Stdout, context{
			Args: namedArgs,
			Vars: loadedPlan.ShuttleConfig.Variables,
		})
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)
}

func templateGet(path string, input interface{}) interface{} {
	properties := strings.SplitN(path, ".", 2)

	if !strings.Contains(path, ".") {
		return templateGetInner(path, input)
	}

	property := properties[0]
	step := templateGetInner(property, input)

	if step != nil {
		return templateGetInner(properties[1], step)
	}
	return nil
}

func templateGetInner(property string, input interface{}) interface{} {
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

func templateString(path string, input interface{}) string {
	value := templateGet(path, input)
	if value == nil {
		return ""
	}
	return value.(string)
}

func templateArray(path string, input interface{}) []interface{} {
	value := templateGet(path, input)
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

func templateObjectArray(path string, input interface{}) []KeyValuePair {
	value := templateGet(path, input)
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

type KeyValuePair struct {
	Key   string
	Value interface{}
}
