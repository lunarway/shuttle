package cmd

import (
	"github.com/spf13/cobra"

	"github.com/lunarway/shuttle/pkg/ui"
)

func newUpload(uii *ui.UI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload shuttle telemetry",
		Run: func(cmd *cobra.Command, args []string) {
			uii.SetContext(ui.LevelSilent)
		},
	}

	return cmd
}
