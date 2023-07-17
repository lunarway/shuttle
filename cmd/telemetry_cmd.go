package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/lunarway/shuttle/pkg/telemetry"
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

			url := os.Getenv("SHUTTLE_REMOTE_TRACING_URL")
			if url == "" {
				log.Fatalln("SHUTTLE_REMOTE_TRACING_URL is not set")
			}

			uploader := telemetry.NewTelemetryUploader(url)

			if err := uploader.Upload(cmd.Context()); err != nil {
				log.Fatalf("failed to upload traces: %s", err)
			}
		},
	}

	return cmd
}
