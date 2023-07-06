package telemetry

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const appKey = "shuttle"

type TelemetryClient interface {
	Trace(ctx context.Context, label string, options ...TelemetryOption)
	TraceError(ctx context.Context, label string, err error, options ...TelemetryOption)
}

type TelemetryOption func(properties map[string]string)

var (
	NoopClient TelemetryClient = &NoopTelemetryClient{}
	Client     TelemetryClient = NoopClient
)

func Setup() {
	properties := make(map[string]string)
	sysinfo := WithGoInfo()
	sysinfo(properties)

	if endpoint := os.Getenv("SHUTTLE_TRACING_ENDPOINT"); endpoint != "" {
		fmt.Println("choosing remote telemetry")
		Client = &UploadTelemetryClient{
			labelPrefix: appKey,
			url:         endpoint,
			properties:  properties,
			Client:      http.DefaultClient,
		}

		return
	}

	if logging_telemetry := os.Getenv("SHUTTLE_LOG_TRACING"); strings.ToLower(
		logging_telemetry,
	) == "true" {
		fmt.Println("choosing logging telemetry")
		Client = &LoggingTelemetryClient{
			labelPrefix: appKey,
			properties:  properties,
		}

		return
	}

	fmt.Println("choosing noop telemetry")
}
