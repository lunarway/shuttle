package telemetry

import (
	"context"
	"os"
	"strings"
)

const appKey = "shuttle"

type TelemetryClient interface {
	Trace(ctx context.Context, label string, properties map[string]string)
	TraceError(ctx context.Context, label string, err error, properties map[string]string)
}

var (
	NoopClient TelemetryClient = &NoopTelemetryClient{}
	Client     TelemetryClient = NoopClient
)

func Setup() {
	if endpoint := os.Getenv("SHUTTLE_TRACING_ENDPOINT"); endpoint != "" {
		Client = &UploadTelemetryClient{
			labelPrefix: appKey,
			url:         endpoint,
			properties:  map[string]string{},
		}

		return
	}

	if logging_telemetry := os.Getenv("SHUTTLE_LOG_TRACING"); strings.ToLower(
		logging_telemetry,
	) == "true" {
		Client = &LoggingTelemetryClient{
			labelPrefix: appKey,
			properties:  map[string]string{},
		}

		return
	}
}
