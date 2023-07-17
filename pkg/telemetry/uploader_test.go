package telemetry

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTelemetryUploader(t *testing.T) {
	t.Parallel()

	t.Run("no telemetry files", func(t *testing.T) {
		t.Parallel()

		uploader := NewTelemetryUploader("some-url")

		files, err := uploader.getTelemetryFiles(
			context.Background(),
			"testdata/no-telemetry-files",
		)
		require.NoError(t, err)
		require.Empty(t, files)
	})

	t.Run("telemetry files", func(t *testing.T) {
		t.Parallel()

		uploader := NewTelemetryUploader("some-url")
		files, err := uploader.getTelemetryFiles(
			context.Background(),
			"testdata/telemetry-files",
		)
		require.NoError(t, err)
		require.Equal(t, []string{"testdata/telemetry-files/shuttle-telemetry.jsonl"}, files)
	})

	t.Run("get shuttle telemetry file", func(t *testing.T) {
		t.Parallel()
		uploader := NewTelemetryUploader("some-url")

		files, _, err := uploader.getTelemetryFile(
			context.Background(),
			"testdata/get-shuttle-telemetry-file/shuttle-telemetry.jsonl",
		)

		events := []UploadTraceEvent{
			{
				App:        "some-app",
				Timestamp:  time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC),
				Properties: map[string]string{},
			},
			{
				App:       "some-app",
				Timestamp: time.Date(2007, time.January, 2, 15, 4, 5, 0, time.UTC),
				Properties: map[string]string{
					"some-key":       "some-value",
					"some-other-key": "some-other-value",
				},
			},
		}

		require.NoError(t, err)
		require.Equal(t, events, files)
	})

	t.Run("full upload test", func(t *testing.T) {
		t.Parallel()

		uploader := NewTelemetryUploader(
			"some-url",
			WithUploadFunction(func(ctx context.Context, url string, event []UploadTraceEvent) error {
				assert.Equal(t, "some-url", url)
				assert.NotEmpty(t, event)

				return nil
			}),
			WithGetTelemetryFiles(func(ctx context.Context, location string) ([]string, error) {
				return getTelemetryFiles(ctx, "testdata/full-upload-test")
			}),
			WithCleanUp(false))

		err := uploader.Upload(context.Background())
		require.NoError(t, err)
	})
}
