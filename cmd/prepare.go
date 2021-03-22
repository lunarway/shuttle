package cmd

import (
	"github.com/spf13/cobra"
)

var (
	prepareCmd = &cobra.Command{
		Use:   "prepare",
		Short: "Load external resources",
		Long:  `Load external resources as a preparation step, before starting to use shuttle`,
		Run: func(cmd *cobra.Command, args []string) {
			_, err := getProjectContext()
			checkError(err)
		},
	}
)

func init() {
	rootCmd.AddCommand(prepareCmd)
}
