package cmd

import (
	"strings"

	"github.com/lunarway/shuttle/pkg/git"
	"github.com/lunarway/shuttle/pkg/ui"
	"github.com/spf13/cobra"
)

func newGitPlanCmd(uii *ui.UI, contextProvider contextProvider) *cobra.Command {
	gitPlanCmd := &cobra.Command{
		Use:   "git-plan [...git_args]",
		Short: "Run a git command for the plan",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: this is no longer possible to configure
			// skipGitPlanPulling = true
			context, err := contextProvider()
			checkError(uii, err)
			git.RunGitPlanCommand(strings.Join(args, " "), context.LocalPlanPath, context.UI)
		},
	}

	return gitPlanCmd
}
