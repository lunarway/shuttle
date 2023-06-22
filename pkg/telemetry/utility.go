package telemetry

import (
	"context"

	"github.com/google/uuid"
)

const telemetryRunKey = "shuttle.runID"

func WithRunID(ctx context.Context) context.Context {
	return context.WithValue(ctx, telemetryRunKey, uuid.New().String())
}

func mergeMaps(ctx context.Context, properties map[string]string) map[string]string {
	if runID, ok := ctx.Value(telemetryRunKey).(string); ok && runID != "" {
		properties[telemetryRunKey] = runID
	}

	return properties
}
