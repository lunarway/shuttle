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
		executors.Execute(context, commandName, args[1:])
	},
}

func init() {
	runCmd.SetHelpFunc(func(f *cobra.Command, args []string) {
		scripts := f.Flags().Args()
		if len(scripts) == 0 {
			runCmd.Usage()
			return
		}
		context := getProjectContext()
		err := executors.Help(context.Scripts, scripts[0], os.Stdout)
		if err != nil {
			context.UI.ExitWithError(err.Error())
			return
		}
	})
	rootCmd.AddCommand(runCmd)
}
