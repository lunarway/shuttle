package telemetry

import "context"

func Trace(ctx context.Context, label string, options ...TelemetryOption) {
	properties := setProperties(append(options, WithLabel(label))...)
	properties = mergeMaps(ctx, properties)
	client.Trace(ctx, properties)
}

func TraceError(ctx context.Context, label string, err error, options ...TelemetryOption) {
	properties := setProperties(append(options, WithLabel(label))...)
	properties = mergeMaps(ctx, properties)

	// TODO: consider enum for error (const list)
	properties["phase"] = "error"
	properties["error"] = err.Error()

	client.Trace(ctx, properties)
}

func setProperties(options ...TelemetryOption) map[string]string {
	properties := make(map[string]string)
	for _, o := range options {
		o(properties)
	}

	return properties
}
