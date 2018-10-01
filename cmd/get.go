package cmd

import (
	"fmt"

	"github.com/lunarway/shuttle/pkg/output"
	"github.com/lunarway/shuttle/pkg/templates"

	"github.com/lunarway/shuttle/pkg/config"
	"github.com/spf13/cobra"
)

type dynamicValue = interface{}

var getCmd = &cobra.Command{
	Use:   "get [variable]",
	Short: "Get a variable value",
	//Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		context := getProjectContext()
		path := args[0]
		fmt.Print(templates.TmplGet(path, context.Config.Variables))
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}

func resolve(dynyaml dynamicValue, properties []string) string {
	if len(properties) == 0 {
		return fmt.Sprintf("%s", dynyaml)
	}

	property := properties[0]
	step := get(dynyaml, property)

	if step != nil {
		return resolve(step, properties[1:])
	}
	return ""
}

func get(dynyaml dynamicValue, property string) dynamicValue {
	switch t := dynyaml.(type) {
	default:
		output.ExitWithErrorCode(2, fmt.Sprintf("unexpected type %T", t))
		return nil
	case map[interface{}]interface{}:
		values := dynyaml.(map[interface{}]interface{})
		return values[property]
	case map[string]interface{}:
		values := dynyaml.(config.DynamicYaml)
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
