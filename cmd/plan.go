package cmd

import (
	"github.com/lunarway/shuttle/cmd/utility"
	"github.com/lunarway/shuttle/pkg/ui"
	"github.com/spf13/cobra"
)

const planDefaultTempl = `{{.Plan}}`

var (
	planFlagTemplate string
)

func newPlan(uii *ui.UI, contextProvider utility.ContextProvider) *cobra.Command {
	planCmd := &cobra.Command{
		Use:   "plan",
		Short: "Output plan information to stdout",
		Long: `Output plan information to stdout.
By default the plan name is output. For projects without a plan (plan: false) an
empty string is written.

Configure the output with a template variable. The format is Go templates.
See http://golang.org/pkg/text/template/#pkg-overview for more details.

Available fields are:

  .LocalPlanPath     Path to the plan on the local file system.
  .Plan              Pretty plan string. Empty if no plan is set.
  .PlanRaw           Raw plan string as read from the configuration.
  .ProjectPath       Path to the current project.
  .TempDirectoryPath Path to the temporary files of the plan on the local file
                     system.
`,
		Example: `Get the raw plan string as it is written in the shuttle.yaml file:
  shuttle plan --template '{{.PlanRaw}}'`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			type templData struct {
				LocalPlanPath     string
				Plan              string
				PlanRaw           interface{}
				ProjectPath       string
				TempDirectoryPath string
			}
			uii.SetUserLevel(ui.LevelError)
			context, err := contextProvider()
			if err != nil {
				return err
			}

			var templ string
			if planFlagTemplate != "" {
				templ = planFlagTemplate
			} else {
				templ = planDefaultTempl
			}
			err = ui.Template(cmd.OutOrStdout(), "plan", templ, templData{
				Plan:              context.Config.Plan,
				PlanRaw:           context.Config.PlanRaw,
				LocalPlanPath:     context.LocalPlanPath,
				ProjectPath:       context.ProjectPath,
				TempDirectoryPath: context.TempDirectoryPath,
			})
			if err != nil {
				return err
			}

			return nil
		},
	}

	planCmd.Flags().StringVar(&planFlagTemplate, "template", "", "Template string to use. See --help for details.")

	return planCmd
}
