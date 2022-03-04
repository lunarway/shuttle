package cmd

import (
	"fmt"

	"github.com/lunarway/shuttle/pkg/errors"
	"github.com/lunarway/shuttle/pkg/templates"
	"github.com/lunarway/shuttle/pkg/ui"
	"github.com/spf13/cobra"
)

func newHas(uii *ui.UI, contextProvider contextProvider) *cobra.Command {
	var (
		lookupInScripts bool
		outputAsStdout  bool
	)

	hasCmd := &cobra.Command{
		Use:           "has [variable]",
		Short:         "Check if a variable (or script) is defined",
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		//Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			uii.SetContext(ui.LevelSilent)

			context, err := contextProvider()
			if err != nil {
				return err
			}

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
				return nil
			} else {
				if found {
					return nil
				} else {
					return errors.NewExitCode(1, "")
				}
			}
		},
	}

	hasCmd.Flags().BoolVar(&lookupInScripts, "script", false, "Lookup existence in scripts instead of vars")
	hasCmd.Flags().BoolVarP(&outputAsStdout, "stdout", "o", false, "Print result to stdout instead of exit code as `true` or `false`")

	return hasCmd
}
