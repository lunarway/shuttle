package telemetry

import "context"

type UploadTelemetryClient struct {
	url         string
	labelPrefix string
	properties  map[string]string
}

func (t *UploadTelemetryClient) Trace(
	ctx context.Context,
	label string,
	properties map[string]string,
) {
}

func (t *UploadTelemetryClient) TraceError(
	ctx context.Context,
	label string,
	err error,
	properties map[string]string,
) {
}

var _ TelemetryClient = &UploadTelemetryClient{}
