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
	NoopClient NoopTelemetryClient = NoopTelemetryClient{}
	Client     TelemetryClient     = &NoopClient
)

func Setup() {
	properties := make(map[string]string, 0)
	sysinfo := WithGoInfo()
	sysinfo(properties)

	if remoteTracing := os.Getenv("SHUTTLE_REMOTE_TRACING"); remoteTracing != "" {
		logLocation := os.Getenv("SHUTTLE_REMOTE_LOG_LOCATION")
		if logLocation == "default" {
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
		endpoint := os.Getenv("SHUTTLE_REMOTE_TRACING_URL")

		Client = &JsonLinesTelemetryClient{
			labelPrefix: appKey,
			url:         endpoint,
			properties:  properties,
			logLocation: logLocation,
			Client:      http.DefaultClient,
		}

		return
	}

	if logging_telemetry := os.Getenv("SHUTTLE_LOG_TRACING"); strings.ToLower(
		logging_telemetry,
	) == "true" {
		Client = &LoggingTelemetryClient{
			labelPrefix: appKey,
			properties:  properties,
		}

		return
	}
}
