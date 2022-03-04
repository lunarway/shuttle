package cmd

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/lunarway/shuttle/pkg/errors"
	"github.com/lunarway/shuttle/pkg/templates"
	"github.com/lunarway/shuttle/pkg/ui"

	"github.com/spf13/cobra"
)

func newGet(uii *ui.UI, contextProvider contextProvider) *cobra.Command {
	var getFlagTemplate string
	getCmd := &cobra.Command{
		Use:   "get [variable]",
		Short: "Get a variable value",
		Args:  cobra.ExactArgs(1),
		//Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			*uii = uii.SetContext(ui.LevelError)
			context, err := contextProvider()
			checkError(uii, err)
			path := args[0]
			var templ string
			if getFlagTemplate != "" {
				templ = getFlagTemplate
			}
			value := templates.TmplGet(path, context.Config.Variables)
			if templ != "" {
				err := ui.Template(cmd.OutOrStdout(), "get", templ, value)
				checkError(uii, err)
				return
			}
			switch value.(type) {
			case nil:
				// print nothing
			default:
				x, err := yaml.Marshal(value)
				if err != nil {
					checkError(uii, errors.NewExitCode(9, "Could not yaml encode value '%s'\nError: %s", value, err))
				}
				fmt.Fprint(cmd.OutOrStdout(), strings.TrimRight(string(x), "\n"))
			}
		},
	}

	getCmd.Flags().StringVar(&getFlagTemplate, "template", "", "Template string to use. See --help for details.")

	return getCmd
}
