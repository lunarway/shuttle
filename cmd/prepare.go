package cmd

import (
	"github.com/lunarway/shuttle/cmd/utility"
	"github.com/lunarway/shuttle/pkg/ui"
	"github.com/spf13/cobra"
)

func newPrepare(uii *ui.UI, contextProvider utility.ContextProvider) *cobra.Command {
	prepareCmd := &cobra.Command{
		Use:   "prepare",
		Short: "Load external resources",
		Long:  `Load external resources as a preparation step, before starting to use shuttle`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := contextProvider()
			if err != nil {
				return err
			}

			return nil
		},
	}

	return prepareCmd
}
