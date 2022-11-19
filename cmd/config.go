package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/lunarway/shuttle/pkg/git"
	"github.com/lunarway/shuttle/pkg/ui"
	"github.com/spf13/cobra"
)

func newConfig(uii *ui.UI, contextProvder contextProvider) *cobra.Command {
	var envVarsToExclude []string

	envCmd := &cobra.Command{
		Use:   "config",
		Short: "Display shuttle context information",
		RunE: func(cmd *cobra.Command, args []string) error {
			uii.SetContext(ui.LevelSilent)
			environmentVariables := os.Environ()
			shouldExclude := make(map[string]bool)

			for _, envVarToExclude := range envVarsToExclude {
				shouldExclude[envVarToExclude] = true
			}

			context, err := contextProvder()
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Version %v\n", version)
			breakLine(cmd)

			parsedPlan := git.ParsePlan(context.Config.Plan)
			if parsedPlan.IsGitPlan {
				fmt.Fprintf(cmd.OutOrStdout(), "Plan:\n%v %v", context.Config.Plan, parsedPlan.Head)
			} else {
				fmt.Printf("%+v\n", context.Config)
				plan := context.Config.Plan
				if plan == "" {
					plan = "false"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Plan:\n%v", plan)
			}
			breakLine(cmd)
			breakLine(cmd)

			fmt.Fprintf(cmd.OutOrStdout(), "Environment:\n")
			for _, envVar := range environmentVariables {
				if shouldExclude[extractEnvirontmentVariableName(envVar)] {
					continue
				}

				fmt.Fprintf(cmd.OutOrStdout(), "%v\n", envVar)
			}

			return nil
		},
	}

	envCmd.Flags().StringSliceVar(&envVarsToExclude, "exclude-env-vars", make([]string, 0), "Exclude environment variables from being displayed. Example: shuttle config --exclude-env-vars VAR1,VAR2,VAR3")
	return envCmd
}

func breakLine(cmd *cobra.Command) {
	fmt.Fprintln(cmd.OutOrStdout(), "")
}

func extractEnvirontmentVariableName(s string) string {
	return strings.Split(s, "=")[0]
}
