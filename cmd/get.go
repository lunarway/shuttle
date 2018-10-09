package cmd

import (
	"fmt"

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
		fmt.Print(templates.TmplGet(path, context.Config.Variables))
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
