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
	var (
		availabilityUrl string
		cleanUp         bool
		uploadUrl       string
	)

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload shuttle telemetry",
		Run: func(cmd *cobra.Command, args []string) {
			uii.SetContext(ui.LevelSilent)

			url := os.Getenv("SHUTTLE_REMOTE_TRACING_URL")
			if url == "" && uploadUrl == "" {
				log.Fatalln("SHUTTLE_REMOTE_TRACING_URL or upload-url is not set")
			}

			if uploadUrl != "" {
				url = uploadUrl
			}

			options := make([]telemetry.UploadOptions, 0)
			if availabilityUrl != "" {
				options = append(options, telemetry.WithAvailabilityCheck(availabilityUrl))
			}

			options = append(options, telemetry.WithCleanUp(cleanUp))

			uploader := telemetry.NewTelemetryUploader(url, options...)

			if err := uploader.Upload(cmd.Context()); err != nil {
				log.Fatalf("failed to upload traces: %s", err)
			}
		},
	}

	cmd.PersistentFlags().
		StringVar(&uploadUrl, "upload-url", "", "upload url is the url to which all the trace events will be uploaded to")
	cmd.PersistentFlags().
		StringVar(&availabilityUrl, "availability-url", "", "availability url is an address that needs to return a 200 http OK before continuing to upload, if anything else is returned, this command exits early")
	cmd.PersistentFlags().
		BoolVar(&cleanUp, "clean-up", true, "removes shuttle-telemetry files after upload")

	return cmd
}
