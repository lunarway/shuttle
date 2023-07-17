package telemetry

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
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

	t.Run("endpoint not available", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		ok, err := availabilityCheck(ctx, "http://some-url-which-doesnt-exist")
		assert.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("endpoint available bad response", func(t *testing.T) {
		t.Parallel()

		server := startServer(
			t,
			func() (string, func(w http.ResponseWriter, r *http.Request)) {
				return "/available", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(400)
				}
			},
		)
		defer server.Close()

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		ok, err := availabilityCheck(ctx, fmt.Sprintf("http://%s/available", server.Addr))
		require.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("endpoint available ok response", func(t *testing.T) {
		t.Parallel()

		server := startServer(
			t,
			func() (string, func(w http.ResponseWriter, r *http.Request)) {
				return "/available", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(200)
				}
			},
		)
		defer server.Close()

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		ok, err := availabilityCheck(ctx, fmt.Sprintf("http://%s/available", server.Addr))
		require.NoError(t, err)
		assert.True(t, ok)
	})
}

type serverFunc = func() (string, func(w http.ResponseWriter, r *http.Request))

func startServer(t *testing.T, serverFuncs ...serverFunc) *http.Server {
	t.Helper()
	server := &http.Server{}

	// Create a new mux for handling requests
	mux := http.NewServeMux()

	for _, f := range serverFuncs {
		path, f := f()
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			log.Printf("handling request for path: %s", path)
			f(w, r)
		})
	}

	server.Handler = mux

	// Start the server with a dynamically allocated port
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Set the listener address to the server
	server.Addr = listener.Addr().String()

	// Start serving requests
	go func(t *testing.T) {
		err := server.Serve(listener)
		if err != nil && err != http.ErrServerClosed {
			t.Fatalf("Server error: %v", err)
		}
	}(t)

	time.Sleep(time.Millisecond * 20)

	return server
}
