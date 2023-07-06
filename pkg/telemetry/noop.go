package telemetry

import "context"

type NoopTelemetryClient struct{}

func (*NoopTelemetryClient) Trace(ctx context.Context, label string, options ...TelemetryOption) {
}

func (*NoopTelemetryClient) TraceError(
	ctx context.Context,
	label string,
	err error,
	options ...TelemetryOption,
) {
}

var _ TelemetryClient = &NoopTelemetryClient{}
