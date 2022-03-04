package cmd

import (
	"github.com/lunarway/shuttle/pkg/ui"
	"github.com/spf13/cobra"
)

func newPrepare(uii *ui.UI, contextProvider contextProvider) *cobra.Command {
	prepareCmd := &cobra.Command{
		Use:   "prepare",
		Short: "Load external resources",
		Long:  `Load external resources as a preparation step, before starting to use shuttle`,
		Run: func(cmd *cobra.Command, args []string) {
			_, err := contextProvider()
			checkError(uii, err)
		},
	}

	return prepareCmd
}
