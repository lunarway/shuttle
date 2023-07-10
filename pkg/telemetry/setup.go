package telemetry

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"
	"strings"
)

const appKey = "shuttle"

type TelemetryClient interface {
	Trace(ctx context.Context, properties map[string]string)
}

type TelemetryOption func(properties map[string]string)

var (
	noopClient NoopTelemetryClient = NoopTelemetryClient{}
	client     TelemetryClient     = &noopClient
)

// Initializes the telemetry setup, if not called, NoopTelemetryClient will be used
func Setup() {
	if remoteTracing := os.Getenv("SHUTTLE_REMOTE_TRACING"); remoteTracing != "" {
		properties := make(map[string]string, 0)
		sysinfo := WithGoInfo()
		sysinfo(properties)

		logLocation := os.Getenv("SHUTTLE_REMOTE_LOG_LOCATION")
		if logLocation == "default" || logLocation == "" {
			usr, _ := user.Current()
			homeDir := usr.HomeDir
			logLocation = path.Join(
				homeDir,
				".local",
				"share",
				"shuttle",
				"telemetry",
			)

			if err := os.MkdirAll(logLocation, 0o755); err != nil {
				log.Fatal(err)
			}
		}
		client = &JsonLinesTelemetryClient{
			labelPrefix: appKey,
			properties:  properties,
			logLocation: logLocation,
			Client:      http.DefaultClient,
		}

		return
	}

	if logging_telemetry := os.Getenv("SHUTTLE_LOG_TRACING"); strings.ToLower(
		logging_telemetry,
	) == "true" {
		properties := make(map[string]string, 0)
		sysinfo := WithGoInfo()
		sysinfo(properties)
		client = &LoggingTelemetryClient{
			labelPrefix: appKey,
			properties:  properties,
		}

		return
	}
}
