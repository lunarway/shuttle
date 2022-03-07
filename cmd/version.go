package cmd

import (
	"fmt"

	"github.com/lunarway/shuttle/pkg/ui"
	"github.com/spf13/cobra"
)

func newVersion(uii *ui.UI) *cobra.Command {
	var showCommit bool

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Info about version of shuttle",
		Run: func(cmd *cobra.Command, args []string) {
			uii.SetContext(ui.LevelSilent)
			if showCommit {
				fmt.Fprintln(cmd.OutOrStdout(), commit)
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), version)
			}
		},
	}

	versionCmd.Flags().BoolVar(&showCommit, "commit", false, "Get git commit sha for current version")

	return versionCmd
}
