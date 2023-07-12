package telemetry

import (
	"context"
	"testing"

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
}
