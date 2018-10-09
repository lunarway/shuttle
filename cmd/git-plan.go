package cmd

import (
	"strings"

	"github.com/lunarway/shuttle/pkg/git"
	"github.com/spf13/cobra"
)

var (
	gitPlanCmd = &cobra.Command{
		Use:   "git-plan [...git_args]",
		Short: "Run a git command for the plan",
		Run: func(cmd *cobra.Command, args []string) {
			skipGitPlanPulling = true
			context := getProjectContext()
			git.RunGitPlanCommand(strings.Join(args, " "), context.LocalPlanPath, context.UI)
		},
	}
)

func init() {
	rootCmd.AddCommand(gitPlanCmd)
}
