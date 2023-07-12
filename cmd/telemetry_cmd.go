package cmd

import (
	"github.com/spf13/cobra"

	"github.com/lunarway/shuttle/pkg/ui"
)

func newTelemetry(uii *ui.UI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "telemetry",
		Short: "Shuttle telemetry",
	}

	cmd.AddCommand(newTelemetryUploadCmd(uii))

	return cmd
}

func newTelemetryUploadCmd(uii *ui.UI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload shuttle telemetry",
		Run: func(cmd *cobra.Command, args []string) {
			uii.SetContext(ui.LevelSilent)
		},
	}

	return cmd
}
