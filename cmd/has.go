package cmd

import (
	"fmt"
	"os"

	"github.com/lunarway/shuttle/pkg/templates"
	"github.com/spf13/cobra"
)

var (
	lookupInScripts bool
	outputAsStdout  bool
	hasCmd          = &cobra.Command{
		Use:   "has [variable]",
		Short: "Check if a variable (or script) is defined",
		Args:  cobra.ExactArgs(1),
		//Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			context := getProjectContext()
			variable := args[0]

			var found bool

			if lookupInScripts {
				_, found = context.Scripts[variable]
			} else {
				found = templates.TmplGet(variable, context.Config.Variables) != nil
			}

			if outputAsStdout {
				if found {
					fmt.Print("true")
				} else {
					fmt.Print("false")
				}
			} else {
				if found {
					os.Exit(0)
				} else {
					os.Exit(1)
				}
			}
		},
	}
)

func init() {
	hasCmd.Flags().BoolVar(&lookupInScripts, "script", false, "Lookup existence in scripts instead of vars")
	hasCmd.Flags().BoolVarP(&outputAsStdout, "stdout", "o", false, "Print result to stdout instead of exit code as `true` or `false`")
	rootCmd.AddCommand(hasCmd)
}
