package cmd

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/lunarway/shuttle/pkg/templates"
	"github.com/lunarway/shuttle/pkg/ui"

	"github.com/spf13/cobra"
)

var (
	getFlagTemplate string
)

type dynamicValue = interface{}

var getCmd = &cobra.Command{
	Use:   "get [variable]",
	Short: "Get a variable value",
	Args:  cobra.ExactArgs(1),
	//Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		uii = uii.SetContext(ui.LevelError)
		context := getProjectContext()
		path := args[0]
		var templ string
		if getFlagTemplate != "" {
			templ = getFlagTemplate
		}
		value := templates.TmplGet(path, context.Config.Variables)
		if templ != "" {
			err := ui.Template(os.Stdout, "get", templ, value)
			context.UI.CheckIfError(err)
			return
		}
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
	getCmd.Flags().StringVar(&getFlagTemplate, "template", "", "Template string to use. See --help for details.")
	rootCmd.AddCommand(getCmd)
}
