package telemetry

import "context"

type NoopTelemetryClient struct{}

func (*NoopTelemetryClient) Trace(ctx context.Context, label string, properties map[string]string) {
}

func (*NoopTelemetryClient) TraceError(
	ctx context.Context,
	label string,
	err error,
	properties map[string]string,
) {
}

var _ TelemetryClient = &NoopTelemetryClient{}
