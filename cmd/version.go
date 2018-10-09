package cmd

import (
	"fmt"

	"github.com/lunarway/shuttle/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	showCommit bool
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Info about version of shuttle",
		Run: func(cmd *cobra.Command, args []string) {
			uii = uii.SetContext(ui.LevelSilent)
			if showCommit {
				fmt.Println(commit)
			} else {
				fmt.Println(version)
			}
		},
	}
)

func init() {
	versionCmd.Flags().BoolVar(&showCommit, "commit", false, "Get git commit sha for current version")
	rootCmd.AddCommand(versionCmd)
}
