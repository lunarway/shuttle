package cmd

import (
	"fmt"

	"github.com/lunarway/shuttle/pkg/templates"

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
