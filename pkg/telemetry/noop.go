package telemetry

import "context"

type NoopTelemetryClient struct{}

func (*NoopTelemetryClient) Trace(ctx context.Context, properties map[string]string) {
}

var _ TelemetryClient = &NoopTelemetryClient{}
