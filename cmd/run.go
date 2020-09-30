package cmd

import (
	"os"

	"github.com/lunarway/shuttle/pkg/executors"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [command]",
	Short: "Run a plan script",
	Long:  `Specify which plan script to run`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var commandName = args[0]
		context := getProjectContext()
		executors.Execute(context, commandName, args[1:], validateArgs)
	},
}

var (
	flagTemplate string
	validateArgs bool
)

func init() {
	runCmd.SetHelpFunc(func(f *cobra.Command, args []string) {
		scripts := f.Flags().Args()
		if len(scripts) == 0 {
			runCmd.Usage()
			return
		}
		context := getProjectContext()
		err := executors.Help(context.Scripts, scripts[0], os.Stdout, flagTemplate)
		if err != nil {
			context.UI.ExitWithError(err.Error())
			return
		}
	})
	runCmd.Flags().StringVar(&flagTemplate, "template", "", "Template string to use. The template format is golang templates [http://golang.org/pkg/text/template/#pkg-overview].")
	runCmd.Flags().BoolVar(&validateArgs, "validate", true, "Validate arguments against script definition in plan and exit with 1 on unknown or missing arguments")
	rootCmd.AddCommand(runCmd)
}
