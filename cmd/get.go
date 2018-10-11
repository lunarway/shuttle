package cmd

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/lunarway/shuttle/pkg/templates"
	"github.com/lunarway/shuttle/pkg/ui"

	"github.com/spf13/cobra"
)

type dynamicValue = interface{}

var getCmd = &cobra.Command{
	Use:   "get [variable]",
	Short: "Get a variable value",
	Args:  cobra.ExactArgs(1),
	//Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		uii = uii.SetContext(ui.LevelSilent)
		context := getProjectContext()
		path := args[0]
		value := templates.TmplGet(path, context.Config.Variables)
		switch value.(type) {
		case nil:
			// print nothing
		default:
			x, err := yaml.Marshal(value)
			if err != nil {
				uii.ExitWithErrorCode(9, "Could not yaml encoded %s\nError: %s", value, err)
			}
			fmt.Print(strings.TrimRight(string(x), "\n"))
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
